package app

import (
	"context"
	"syscall"
	"time"

	v1 "github.com/Team8te/go-pharos/app/api/http/v1"
	"github.com/Team8te/go-pharos/app/api/udp"
	"github.com/Team8te/go-pharos/app/job"
	"github.com/Team8te/go-pharos/app/job/worker"
	"github.com/Team8te/go-pharos/app/repository"
	"github.com/Team8te/go-pharos/app/service/device"
	"github.com/Team8te/go-pharos/app/service/license"
	"github.com/Team8te/go-pharos/app/service/store"
	"github.com/Team8te/go-pharos/app/service/transfer"
	"github.com/Team8te/go-pharos/pkg/platform/closer"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	v1device "github.com/Team8te/go-pharos/app/api/http/v1/device"
	v1store "github.com/Team8te/go-pharos/app/api/http/v1/store"
	v1ws "github.com/Team8te/go-pharos/app/api/http/v1/ws"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	*closer.Closer
	jobs []*job.Job
	e    *echo.Echo
}

func NewApp(config string) (*App, error) {
	viper.SetConfigFile(config)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

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
	l, err := license.NewLicenser(viper.GetString("license"))
	if err != nil {
		return err
	}

	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return err
	}

	deviceService := a.MakeDeviceService(db, l)
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
	viper.SetDefault("chunkSize", 512)

	tr := transfer.NewTransfer(viper.GetInt("chunkSize"))
	storeService := a.MakeStoreKeeper(viper.GetString("root"), tr)

	a.e = v1.Register(
		v1device.NewDeviceEndpoint(deviceService),
		v1store.NewStoreEndpoint(storeService),
		v1ws.NewWSEndpoint(storeService),
	)

	return nil
}

func (a *App) MakeDeviceService(db *sqlx.DB, l *license.Licenser) *device.DeviceInformer {
	return device.NewDeviceInformer(
		l.GetUUID(),
		repository.NewRepository(db),
	)
}

func (a *App) MakeStoreKeeper(root string, tr *transfer.Transfer) *store.StoreKeeper {
	return store.NewStoreKeeper(root, tr)
}

func (a *App) Wait() error {
	a.Closer.Wait()
	return a.e.Close()
}
