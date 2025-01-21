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

	infeedpb "github.com/goverland-labs/goverland-core-web-api/protocol/feed"

	instopb "github.com/goverland-labs/goverland-core-web-api/protocol/storage"

	"github.com/goverland-labs/goverland-core-web-api/internal/config"
	ingrpc "github.com/goverland-labs/goverland-core-web-api/internal/grpc"
	"github.com/goverland-labs/goverland-core-web-api/internal/rest"
	apihandlers "github.com/goverland-labs/goverland-core-web-api/internal/rest/handlers"
	"github.com/goverland-labs/goverland-core-web-api/pkg/grpcsrv"
	"github.com/goverland-labs/goverland-core-web-api/pkg/health"
	"github.com/goverland-labs/goverland-core-web-api/pkg/prometheus"
)

type Application struct {
	sigChan <-chan os.Signal
	manager *process.Manager
	cfg     config.App

	cdc  storagepb.DaoClient
	cpc  storagepb.ProposalClient
	cefc feedpb.FeedEventsClient
	csfc storagepb.VoteClient
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

		a.initGRPCServer,

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
	storageConn, err := grpc.NewClient(a.cfg.InternalAPI.CoreStorageAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("create connection with core storage server: %v", err)
	}

	a.cdc = storagepb.NewDaoClient(storageConn)
	a.cpc = storagepb.NewProposalClient(storageConn)
	a.csfc = storagepb.NewVoteClient(storageConn)
	vc := storagepb.NewVoteClient(storageConn)
	ec := storagepb.NewEnsClient(storageConn)
	sc := storagepb.NewStatsClient(storageConn)
	delegateClient := storagepb.NewDelegateClient(storageConn)

	feedConn, err := grpc.NewClient(a.cfg.InternalAPI.CoreFeedAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("create connection with core feed server: %v", err)
	}

	subscriberClient := feedpb.NewSubscriberClient(feedConn)
	subscriptionClient := feedpb.NewSubscriptionClient(feedConn)
	fc := feedpb.NewFeedClient(feedConn)
	a.cefc = feedpb.NewFeedEventsClient(feedConn)

	handlers := []apihandlers.APIHandler{
		apihandlers.NewDaoHandler(a.cdc, fc, delegateClient),
		apihandlers.NewProposalHandler(a.cpc, vc),
		apihandlers.NewSubscribeHandler(subscriberClient, subscriptionClient),
		apihandlers.NewFeedHandler(fc),
		apihandlers.NewVotesHandler(vc),
		apihandlers.NewEnsHandler(ec),
		apihandlers.NewStatsHandler(sc),
		apihandlers.NewDelegateHandler(delegateClient),
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

func (a *Application) initGRPCServer() error {
	authInterceptor := grpcsrv.NewAuthInterceptor()
	srv := grpcsrv.NewGrpcServer(
		[]string{
			"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
		},
		authInterceptor.AuthAndIdentifyTickerFunc,
	)

	instopb.RegisterDaoServer(srv, ingrpc.NewDaoServer(a.cdc))
	instopb.RegisterProposalServer(srv, ingrpc.NewProposalServer(a.cpc))
	infeedpb.RegisterFeedEventsServer(srv, ingrpc.NewFeedServer(ingrpc.NewService(a.cefc, a.csfc)))

	a.manager.AddWorker(grpcsrv.NewGrpcServerWorker("API", srv, a.cfg.InternalAPI.Bind))

	return nil
}
