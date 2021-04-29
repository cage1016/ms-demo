package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/kelseyhightower/envconfig"
	"github.com/opentracing/opentracing-go"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/uber/jaeger-client-go"
	jconfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"

	addservice "github.com/cage1016/ms-sample/internal/app/add/service"
	addtransportsgrpc "github.com/cage1016/ms-sample/internal/app/add/transports/grpc"
	"github.com/cage1016/ms-sample/internal/app/tictac/endpoints"
	"github.com/cage1016/ms-sample/internal/app/tictac/postgres"
	"github.com/cage1016/ms-sample/internal/app/tictac/service"
	transportsgrpc "github.com/cage1016/ms-sample/internal/app/tictac/transports/grpc"
	transportshttp "github.com/cage1016/ms-sample/internal/app/tictac/transports/http"
	"github.com/cage1016/ms-sample/internal/pkg/logconv"
	pb "github.com/cage1016/ms-sample/pb/tictac"
)

type Config struct {
	DbConfig    postgres.Config
	ServiceName string `envconfig:"QS_SERVICE_NAME" default:"tictac"`
	ServiceHost string `envconfig:"QS_SERVICE_HOST" default:"localhost"`
	LogLevel    string `envconfig:"QS_LOG_LEVEL" default:"error"`
	HttpPort    string `envconfig:"QS_HTTP_PORT" default:"8180"`
	GrpcPort    string `envconfig:"QS_GRPC_PORT" default:"8181"`
	ZipkinV2URL string `envconfig:"QS_ZIPKIN_V2_URL"`
	JaegerURL   string `envconfig:"QS_JAEGER_URL"`
	AddURL      string `envconfig:"QS_ADD_URL"`
}

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	}

	var cfg Config
	err := envconfig.Process("qs", &cfg)
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	logger = level.NewFilter(logger, level.AllowInfo())
	logger = log.With(logger, "service", cfg.ServiceName)
	logger = log.With(logger, "caller", log.DefaultCaller)
	level.Info(logger).Log("version", service.Version, "commitHash", service.CommitHash, "buildTimeStamp", service.BuildTimeStamp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := grpc.DialContext(ctx, cfg.AddURL, grpc.WithInsecure())
	if err != nil {
		level.Error(logger).Log("method", "grpc.DialContext", "err", err)
		os.Exit(1)
	}

	db := connectToDB(cfg.DbConfig, logger)
	if cfg.LogLevel == logconv.Debug {
		db = db.Debug()
	}

	tracer, closer := initJaeger(cfg.ServiceName, cfg.JaegerURL, logger)
	defer closer.Close()

	zipkinTracer := initZipkin(cfg.ServiceName, cfg.HttpPort, cfg.ZipkinV2URL, logger)

	addsvc := addtransportsgrpc.NewGRPCClient(conn, tracer, zipkinTracer, logger)
	service := NewServer(db, addsvc, logger)
	endpoints := endpoints.New(service, logger, tracer, zipkinTracer)

	hs := health.NewServer()
	hs.SetServingStatus(cfg.ServiceName, healthgrpc.HealthCheckResponse_SERVING)

	wg := &sync.WaitGroup{}

	go startHTTPServer(ctx, wg, endpoints, tracer, zipkinTracer, cfg.HttpPort, logger)
	go startGRPCServer(ctx, wg, endpoints, tracer, zipkinTracer, cfg.GrpcPort, hs, logger)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	cancel()
	wg.Wait()

	fmt.Println("main: all goroutines have told us they've finished")
}

func connectToDB(dbConfig postgres.Config, logger log.Logger) *gorm.DB {
	db, err := postgres.Connect(dbConfig)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}
	return db
}

func NewServer(db *gorm.DB, addsvc addservice.AddService, logger log.Logger) service.TictacService {
	repo := postgres.New(db, logger)
	service := service.New(repo, addsvc, logger)
	return service
}

func initJaeger(svcName string, url string, logger log.Logger) (opentracing.Tracer, io.Closer) {
	if url == "" {
		return opentracing.NoopTracer{}, ioutil.NopCloser(nil)
	}

	tracer, closer, err := jconfig.Configuration{
		ServiceName: svcName,
		Sampler: &jconfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jconfig.ReporterConfig{
			LocalAgentHostPort: url,
			LogSpans:           true,
		},
	}.NewTracer()
	if err != nil {
		level.Error(logger).Log("msg", fmt.Sprintf("Failed to init Jaeger: %s", err))
		os.Exit(1)
	}

	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

func initZipkin(serviceName, httpPort, zipkinV2URL string, logger log.Logger) (zipkinTracer *zipkin.Tracer) {
	if zipkinV2URL != "" {
		var (
			err           error
			hostPort      = fmt.Sprintf("localhost:%s", httpPort)
			useNoopTracer = (zipkinV2URL == "")
			reporter      = zipkinhttp.NewReporter(zipkinV2URL)
		)
		zEP, _ := zipkin.NewEndpoint(serviceName, hostPort)
		zipkinTracer, err = zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(useNoopTracer))
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		if !useNoopTracer {
			logger.Log("tracer", "Zipkin", "type", "Native", "URL", zipkinV2URL)
		}
	}
	return
}

func startHTTPServer(ctx context.Context, wg *sync.WaitGroup, endpoints endpoints.Endpoints, tracer stdopentracing.Tracer, zipkinTracer *zipkin.Tracer, port string, logger log.Logger) {
	wg.Add(1)
	defer wg.Done()

	if port == "" {
		level.Error(logger).Log("protocol", "HTTP", "exposed", port, "err", "port is not assigned exist")
		return
	}

	p := fmt.Sprintf(":%s", port)
	// create a server
	srv := &http.Server{Addr: p, Handler: transportshttp.NewHTTPHandler(endpoints, tracer, zipkinTracer, logger)}
	level.Info(logger).Log("protocol", "HTTP", "exposed", port)
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			level.Info(logger).Log("Listen", err)
		}
	}()

	<-ctx.Done()

	// shut down gracefully, but wait no longer than 5 seconds before halting
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ignore error since it will be "Err shutting down server : context canceled"
	srv.Shutdown(shutdownCtx)

	level.Info(logger).Log("protocol", "HTTP", "Shutdown", "http server gracefully stopped")
}

func startGRPCServer(ctx context.Context, wg *sync.WaitGroup, endpoints endpoints.Endpoints, tracer stdopentracing.Tracer, zipkinTracer *zipkin.Tracer, port string, hs *health.Server, logger log.Logger) {
	wg.Add(1)
	defer wg.Done()

	p := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		level.Error(logger).Log("protocol", "GRPC", "listen", port, "err", err)
		os.Exit(1)
	}

	var server *grpc.Server
	level.Info(logger).Log("protocol", "GRPC", "exposed", port)
	server = grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	pb.RegisterTictacServer(server, transportsgrpc.MakeGRPCServer(endpoints, tracer, zipkinTracer, logger))
	healthgrpc.RegisterHealthServer(server, hs)
	reflection.Register(server)

	go func() {
		// service connections
		err = server.Serve(listener)
		if err != nil {
			fmt.Printf("grpc serve : %s\n", err)
		}
	}()

	<-ctx.Done()

	// ignore error since it will be "Err shutting down server : context canceled"
	server.GracefulStop()

	fmt.Println("grpc server gracefully stopped")
}
