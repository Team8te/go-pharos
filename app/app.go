package app

import (
	"context"
	"syscall"
	"time"

	v1 "github.com/go-pharos/app/api/http/v1"
	"github.com/go-pharos/app/api/udp"
	"github.com/go-pharos/app/job"
	"github.com/go-pharos/app/job/worker"
	"github.com/go-pharos/app/repository"
	"github.com/go-pharos/app/service/device"
	"github.com/go-pharos/app/service/store"
	"github.com/go-pharos/pkg/platform/closer"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	v1device "github.com/go-pharos/app/api/http/v1/device"
	v1store "github.com/go-pharos/app/api/http/v1/store"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	*closer.Closer
	jobs []*job.Job
	e    *echo.Echo
}

func NewApp(config string) (*App, error) {
	viper.SetConfigFile(config)
	return &App{
		Closer: closer.New(syscall.SIGTERM, syscall.SIGINT),
	}, nil
}

func (a *App) Run() error {
	err := a.Init()
	if err != nil {
		return err
	}

	job.RunJobs(context.Background(), a.jobs...)
	defer job.StopJobs(a.jobs...)

	go a.e.Start(":3000")
	return a.Wait()
}

func (a *App) Init() error {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return err
	}

	deviceService := a.MakeDeviceService(db)
	jobQueue := job.NewQueue()
	port := 9910
	a.jobs = append(a.jobs,
		job.NewJob(
			300*time.Second, worker.NewScan(deviceService, port),
		),
		job.NewJob(
			0, worker.NewScanReceiver(udp.NewEndpoint(deviceService), port),
		),
		job.NewJob(
			0, jobQueue,
		),
	)

	storeService := a.MakeStoreKeeper(viper.GetString("root"))

	a.e = v1.Register(v1device.NewDeviceEndpoint(deviceService), v1store.NewStoreEndpoint(storeService))

	return nil
}

func (a *App) MakeDeviceService(db *sqlx.DB) *device.DeviceInformer {
	return device.NewDeviceInformer(
		uuid.New().String(),
		repository.NewRepository(db),
	)
}

func (a *App) MakeStoreKeeper(root string) *store.StoreKeeper {
	return store.NewStoreKeeper(root)
}

func (a *App) Wait() error {
	a.Closer.Wait()
	return a.e.Close()
}
