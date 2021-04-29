// +build integration

package tictac

import (
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	addendpoints "github.com/cage1016/ms-sample/internal/app/add/endpoints"
	addservice "github.com/cage1016/ms-sample/internal/app/add/service"
	addtransports "github.com/cage1016/ms-sample/internal/app/add/transports/grpc"
	"github.com/cage1016/ms-sample/internal/app/tictac/endpoints"
	"github.com/cage1016/ms-sample/internal/app/tictac/postgres"
	"github.com/cage1016/ms-sample/internal/app/tictac/service"
	transports "github.com/cage1016/ms-sample/internal/app/tictac/transports/http"
	"github.com/cage1016/ms-sample/internal/pkg/errors"
	pb "github.com/cage1016/ms-sample/pb/add"
)

const (
	// databaseHost is the host name of the test database.
	databaseHost = "db"
	// databaseHost = "localhost"

	// databasePort is the port that the test database is listening on.
	databasePort = "5432"

	// databaseUser is the user for the test database.
	databaseUser = "postgres"

	// databasePass is the password of the user for the test database.
	databasePass = "password"

	// databaseName is the name of the test database.
	databaseName = "tictac"
)

const (
	hostPort string = "localhost:8002"
)

var a *Application

type Application struct {
	DB      *gorm.DB
	handler http.Handler
}

func Truncate(dbc *gorm.DB) error {
	stmt := "TRUNCATE TABLE tictac;"

	if err := dbc.Raw(stmt).Error; err != nil {
		return errors.Wrap(errors.New("truncate test database tables"), err)
	}

	return nil
}

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	logger := log.NewLogfmtLogger(os.Stderr)

	db, err := postgres.Connect(postgres.Config{
		Host:        databaseHost,
		Port:        databasePort,
		User:        databaseUser,
		Pass:        databasePass,
		Name:        databaseName,
		SSLMode:     "disable",
		SSLCert:     "",
		SSLKey:      "",
		SSLRootCert: "",
	})
	if err != nil {
		logger.Log("err", err)
		return 1
	}

	zkt, _ := zipkin.NewTracer(nil, zipkin.WithNoopTracer(true))
	tracer := opentracing.GlobalTracer()

	// add server
	server := grpc.NewServer()
	eps1 := addendpoints.New(addservice.New(logger), logger, tracer, zkt)
	sc, err := net.Listen("tcp", hostPort)
	if err != nil {
		logger.Log("unable to listen: %+v", err)
		return 1
	}
	defer server.GracefulStop()

	go func() {
		pb.RegisterAddServer(server, addtransports.MakeGRPCServer(eps1, tracer, zkt, logger))
		_ = server.Serve(sc)
	}()

	// add client
	cc, err := grpc.Dial(hostPort, grpc.WithInsecure())
	if err != nil {
		logger.Log("unable to Dial: %+v", err)
		return 1
	}
	repo := postgres.New(db, logger)
	addsvcclient := addtransports.NewGRPCClient(cc, tracer, zkt, logger)
	svc := service.New(repo, addsvcclient, logger)
	eps2 := endpoints.New(svc, logger, tracer, zkt)

	a = &Application{
		DB:      db,
		handler: transports.NewHTTPHandler(eps2, tracer, zkt, logger),
	}

	return m.Run()
}
