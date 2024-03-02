package internal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/goverland-labs/goverland-core-feed/protocol/feedpb"
	"github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
	"github.com/s-larionov/process-manager"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"

	"github.com/goverland-labs/goverland-core-web-api/internal/config"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest"
	apihandlers "github.com/goverland-labs/goverland-core-web-api/internal/rest/handlers"
	"github.com/goverland-labs/goverland-core-web-api/pkg/health"
	"github.com/goverland-labs/goverland-core-web-api/pkg/prometheus"
)

type Application struct {
	sigChan <-chan os.Signal
	manager *process.Manager
	cfg     config.App
}

func NewApplication(cfg config.App) (*Application, error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	a := &Application{
		sigChan: sigChan,
		cfg:     cfg,
		manager: process.NewManager(),
	}

	err := a.bootstrap()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Application) Run() {
	a.manager.StartAll()
	a.registerShutdown()
}

func (a *Application) bootstrap() error {
	initializers := []func() error{
		// Init Dependencies
		a.initServices,

		// Init Workers: Application
		a.initRestAPI,

		// Init Workers: System
		a.initPrometheusWorker,
		a.initHealthWorker,
	}

	for _, initializer := range initializers {
		if err := initializer(); err != nil {
			return err
		}
	}

	return nil
}

func (a *Application) initServices() error {
	// TODO

	return nil
}

func (a *Application) initRestAPI() error {
	storageConn, err := grpc.Dial(a.cfg.InternalAPI.CoreStorageAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("create connection with core storage server: %v", err)
	}

	dc := storagepb.NewDaoClient(storageConn)
	pc := storagepb.NewProposalClient(storageConn)
	vc := storagepb.NewVoteClient(storageConn)

	feedConn, err := grpc.Dial(a.cfg.InternalAPI.CoreFeedAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("create connection with core feed server: %v", err)
	}

	subscriberClient := feedpb.NewSubscriberClient(feedConn)
	subscriptionClient := feedpb.NewSubscriptionClient(feedConn)
	fc := feedpb.NewFeedClient(feedConn)

	handlers := []apihandlers.APIHandler{
		apihandlers.NewDaoHandler(dc, fc),
		apihandlers.NewProposalHandler(pc, vc),
		apihandlers.NewSubscribeHandler(subscriberClient, subscriptionClient),
		apihandlers.NewFeedHandler(fc),
		apihandlers.NewVotesHandler(vc),
	}

	a.manager.AddWorker(process.NewServerWorker("rest-API", rest.NewRestServer(a.cfg.REST, handlers)))

	return nil
}

func (a *Application) initPrometheusWorker() error {
	srv := prometheus.NewServer(a.cfg.Prometheus.Listen, "/metrics")
	a.manager.AddWorker(process.NewServerWorker("prometheus", srv))

	return nil
}

func (a *Application) initHealthWorker() error {
	srv := health.NewHealthCheckServer(a.cfg.Health.Listen, "/status", health.DefaultHandler(a.manager))
	a.manager.AddWorker(process.NewServerWorker("health", srv))

	return nil
}

func (a *Application) registerShutdown() {
	go func(manager *process.Manager) {
		<-a.sigChan

		manager.StopAll()
	}(a.manager)

	a.manager.AwaitAll()
}
