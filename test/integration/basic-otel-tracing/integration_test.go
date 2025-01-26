package basic_otel_tracing

import (
	"context"
	"fmt"
	"github.com/Trendyol/go-dcp"
	"github.com/couchbase/gocbcore/v10"
	"github.com/testcontainers/testcontainers-go"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"

	"github.com/Trendyol/go-dcp/models"

	"github.com/Trendyol/go-dcp/logger"

	"github.com/Trendyol/go-dcp/config"

	"github.com/Trendyol/go-dcp/couchbase"
	"github.com/testcontainers/testcontainers-go/wait"

	// You know the drill
	_ "github.com/Trendyol/otel-go-dcp"
)

func panicVersion(version string) {
	panic(fmt.Sprintf("invalid version: %v", version))
}

func parseVersion(version string) (int, int, int) {
	parse := strings.Split(version, ".")
	if len(parse) < 3 {
		panicVersion(version)
	}

	major, err := strconv.Atoi(parse[0])
	if err != nil {
		panicVersion(version)
	}

	minor, err := strconv.Atoi(parse[1])
	if err != nil {
		panicVersion(version)
	}

	patch, err := strconv.Atoi(parse[2])
	if err != nil {
		panicVersion(version)
	}

	return major, minor, patch
}

func isVersion5xx(version string) bool {
	major, _, _ := parseVersion(version)
	return major == 5
}

func getConfig() *config.Dcp {
	return &config.Dcp{
		Hosts:      []string{"localhost:8091"},
		Username:   "user",
		Password:   "123456",
		BucketName: "dcp-test",
		Dcp: config.ExternalDcp{
			Group: config.DCPGroup{
				Name: "groupName",
				Membership: config.DCPGroupMembership{
					RebalanceDelay: 3 * time.Second,
				},
			},
		},
	}
}

func setupContainer(c *config.Dcp, ctx context.Context, version string) (testcontainers.Container, error) {
	var entrypoint string
	if isVersion5xx(version) {
		entrypoint = "../../../scripts/entrypoint_5.sh"
	} else {
		entrypoint = "../../../scripts/entrypoint.sh"
	}

	req := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("couchbase/server:%v", version),
		ExposedPorts: []string{"8091/tcp", "8093/tcp", "11210/tcp"},
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.PortBindings = map[nat.Port][]nat.PortBinding{
				"8091/tcp":  {{HostIP: "0.0.0.0", HostPort: "8091"}},
				"8093/tcp":  {{HostIP: "0.0.0.0", HostPort: "8093"}},
				"11210/tcp": {{HostIP: "0.0.0.0", HostPort: "11210"}},
			}
		},
		WaitingFor: wait.ForLog("/entrypoint.sh couchbase-server").WithStartupTimeout(30 * time.Second),
		Env: map[string]string{
			"USERNAME":                  c.Username,
			"PASSWORD":                  c.Password,
			"BUCKET_NAME":               c.BucketName,
			"BUCKET_TYPE":               "couchbase",
			"BUCKET_RAMSIZE":            "1024",
			"CLUSTER_RAMSIZE":           "1024",
			"CLUSTER_INDEX_RAMSIZE":     "512",
			"CLUSTER_EVENTING_RAMSIZE":  "256",
			"CLUSTER_FTS_RAMSIZE":       "256",
			"CLUSTER_ANALYTICS_RAMSIZE": "1024",
			"INDEX_STORAGE_SETTING":     "memopt",
			"REST_PORT":                 "8091",
			"CAPI_PORT":                 "8092",
			"QUERY_PORT":                "8093",
			"FTS_PORT":                  "8094",
			"MEMCACHED_SSL_PORT":        "11207",
			"MEMCACHED_PORT":            "11210",
			"SSL_REST_PORT":             "18091",
		},
		Entrypoint: []string{
			"/config-entrypoint.sh",
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      entrypoint,
				ContainerFilePath: "/config-entrypoint.sh",
				FileMode:          600,
			},
		},
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func setupJeagerAllInOneContainer(ctx context.Context) func() {
	req := testcontainers.ContainerRequest{
		Image:        "jaegertracing/all-in-one:latest",
		ExposedPorts: []string{"5775/udp", "6831/udp", "6832/udp", "5778/tcp", "16686/tcp", "14268/tcp", "9411/tcp"},
		WaitingFor:   wait.ForLog("Starting HTTP server").WithStartupTimeout(30 * time.Second),
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.PortBindings = nat.PortMap{
				"5775/udp":  []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "5775"}},
				"6831/udp":  []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "6831"}},
				"6832/udp":  []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "6832"}},
				"5778/tcp":  []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "5778"}},
				"16686/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "16686"}},
				"14268/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "14268"}},
				"9411/tcp":  []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "9411"}},
			}
		},
	}

	// Start the container
	jaegerContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	// Get the mapped ports
	port5775, _ := jaegerContainer.MappedPort(ctx, "5775/udp")
	port6831, _ := jaegerContainer.MappedPort(ctx, "6831/udp")
	port6832, _ := jaegerContainer.MappedPort(ctx, "6832/udp")
	port5778, _ := jaegerContainer.MappedPort(ctx, "5778/tcp")
	port16686, _ := jaegerContainer.MappedPort(ctx, "16686/tcp")
	port14268, _ := jaegerContainer.MappedPort(ctx, "14268/tcp")
	port9411, _ := jaegerContainer.MappedPort(ctx, "9411/tcp")

	// Print the mapped ports
	fmt.Printf("Jaeger container started!\n")
	fmt.Printf("Port 5775/udp mapped to: %s\n", port5775)
	fmt.Printf("Port 6831/udp mapped to: %s\n", port6831)
	fmt.Printf("Port 6832/udp mapped to: %s\n", port6832)
	fmt.Printf("Port 5778/tcp mapped to: %s\n", port5778)
	fmt.Printf("Port 16686/tcp mapped to: %s\n", port16686)
	fmt.Printf("Port 14268/tcp mapped to: %s\n", port14268)
	fmt.Printf("Port 9411/tcp mapped to: %s\n", port9411)

	// Clean up the container when done
	return func() {
		if err := jaegerContainer.Terminate(ctx); err != nil {
			panic(err)
		}
	}
}

func insertDataToContainer(c *config.Dcp, t *testing.T, iteration int, chunkSize int, bulkSize int) {
	logger.Log.Info("mock data stream started with iteration=%v", iteration)

	client := couchbase.NewClient(c)

	err := client.Connect()
	if err != nil {
		t.Fatal(err)
	}

	var iter int

	for iteration > iter {
		for chunk := 0; chunk < chunkSize; chunk++ {
			wg := &sync.WaitGroup{}
			wg.Add(bulkSize)

			for id := 0; id < bulkSize; id++ {
				go func(id int, chunk int) {
					ch := make(chan error, 1)

					opm := couchbase.NewAsyncOp(context.Background())

					op, err := client.GetAgent().Set(gocbcore.SetOptions{
						Key:           []byte(fmt.Sprintf("%v_%v_%v", iter, chunk, id)),
						Value:         []byte(fmt.Sprintf("%v_%v_%v", iter, chunk, id)),
						Deadline:      time.Now().Add(time.Second * 5),
						RetryStrategy: gocbcore.NewBestEffortRetryStrategy(nil),
					}, func(result *gocbcore.StoreResult, err error) {
						opm.Resolve()

						ch <- err
					})

					err = opm.Wait(op, err)
					if err != nil {
						t.Error(err)
					}

					err = <-ch
					if err != nil {
						t.Error(err)
					}

					wg.Done()
				}(id, chunk)
			}

			wg.Wait()
		}

		iter++
	}

	client.Close()

	logger.Log.Info("mock data stream finished with totalSize=%v", iteration)
}

func testTraces(ctx *models.ListenerContext) {
	// Traces
	lt1 := ctx.ListenerTracerComponent.InitializeListenerTrace("test1", map[string]interface{}{})
	time.Sleep(time.Second * 1)
	lt11 := ctx.ListenerTracerComponent.CreateListenerTrace(lt1, "test1-1", map[string]interface{}{})
	time.Sleep(time.Second * 1)
	lt11.Finish()
	lt12 := ctx.ListenerTracerComponent.CreateListenerTrace(lt1, "test1-2", map[string]interface{}{
		"test1-2": "This is a test metadata",
	})
	time.Sleep(time.Second * 1)
	lt121 := ctx.ListenerTracerComponent.CreateListenerTrace(lt12, "test1-2-1", map[string]interface{}{})
	time.Sleep(time.Millisecond * 100)
	lt121.Finish()
	time.Sleep(time.Millisecond * 300)
	lt12.Finish()
	lt1.Finish()
}

func testWithTraces(t *testing.T, version string) {
	chunkSize := 4
	bulkSize := 1024
	iteration := 512
	mockDataSize := iteration * bulkSize * chunkSize
	totalNotify := 10
	notifySize := mockDataSize / totalNotify

	c := getConfig()
	c.ApplyDefaults()

	ctx := context.Background()

	container, err := setupContainer(c, ctx, version)
	if err != nil {
		t.Fatal(err)
	}

	var counter atomic.Int32
	finish := make(chan struct{}, 1)

	dcp, err := dcp.NewDcp(c, func(ctx *models.ListenerContext) {
		if _, ok := ctx.Event.(models.DcpMutation); ok {
			ctx.Ack()
			testTraces(ctx)

			val := int(counter.Add(1))

			if val%notifySize == 0 {
				logger.Log.Info("%v/%v processed", val/notifySize, totalNotify)
			}

			if val == mockDataSize {
				finish <- struct{}{}
			}
		}
	})
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		<-dcp.WaitUntilReady()
		insertDataToContainer(c, t, iteration, chunkSize, bulkSize)
	}()

	go func() {
		<-finish
		dcp.Close()
	}()

	dcp.Start()

	err = container.Terminate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	logger.Log.Info("mock data stream finished with totalSize=%v", counter.Load())
}

func TestDcpWithTraces(t *testing.T) {
	terminateFunc := setupJeagerAllInOneContainer(context.Background())
	defer terminateFunc()

	version := "7.6.3"

	t.Run(version, func(t *testing.T) {
		testWithTraces(t, version)
	})
}
