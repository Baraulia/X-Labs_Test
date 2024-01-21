package main

//nolint:depguard
import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"

	"github.com/Baraulia/X-Labs_Test/internal/api/grpc/grpcserver"
	"github.com/Baraulia/X-Labs_Test/internal/api/grpc/pb"
	"github.com/Baraulia/X-Labs_Test/internal/app"
	memorystorage "github.com/Baraulia/X-Labs_Test/internal/storage/inMemory"
	"github.com/Baraulia/X-Labs_Test/pkg/logger"
	"github.com/Baraulia/X-Labs_Test/pkg/validation"
	"google.golang.org/grpc"
)

var (
	initAdminName     string
	initAdminPassword string
	secretKey         string
	configPath        string
)

func init() {
	flag.StringVar(&initAdminName, "admin_name", "admin", "Name for initialization the first admin in database")
	flag.StringVar(&initAdminPassword, "admin_password", "admin", "Password for initialization the first admin in database")
	flag.StringVar(&secretKey, "secret_key", "secret", "Secret key for password hashing")
	flag.StringVar(&configPath, "config", "./configs/config.yaml", "Path to config file")
}

func main() {
	flag.Parse()

	config, err := NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.GetLogger(config.Logger.Level)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage := memorystorage.NewUserStorage(logg)

	if err = storage.InitAdmin(initAdminName, initAdminPassword, secretKey); err != nil {
		logg.Fatal(err.Error(), map[string]interface{}{"initAdminName": initAdminName, "initAdminPassword": initAdminPassword, "secretKey": secretKey})
	}

	validator := validation.New()

	service := app.NewApp(logg, storage, validator)
	grpcService := grpcserver.NewServer(service, logg)

	server := grpc.NewServer(grpc.UnaryInterceptor(
		grpcService.BasicAuthInterceptor,
	))

	pb.RegisterUserServiceServer(server, grpcService)

	go func() {
		<-ctx.Done()
		logg.Info("stopping grpc server...", nil)
		server.GracefulStop()
	}()

	lsn, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GRPC.Port))
	if err != nil {
		logg.Fatal(err.Error(), nil)
	}

	logg.Info("starting server on "+lsn.Addr().String(), nil)
	if err := server.Serve(lsn); err != nil {
		logg.Fatal("failed to start grpc server", map[string]interface{}{"error": err})
	}
}
