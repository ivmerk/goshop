package app

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	inventoryAPI "github.com/ivmerk/goshop/inventory/internal/api/inventory/v1"
	"github.com/ivmerk/goshop/inventory/internal/config"
	"github.com/ivmerk/goshop/inventory/internal/repository"
	partRepo "github.com/ivmerk/goshop/inventory/internal/repository/part"
	"github.com/ivmerk/goshop/inventory/internal/service"
	partService "github.com/ivmerk/goshop/inventory/internal/service/part"
	"github.com/ivmerk/goshop/platform/pkg/closer"
	inventoryV1 "github.com/ivmerk/goshop/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1API inventoryV1.InventoryServiceServer
	partService    service.PartService
	partRepository repository.PartRepository

	mongoClient *mongo.Client
	partsColl   *mongo.Collection
}

func NewDiContainer() *diContainer { return &diContainer{} }

func (d *diContainer) InventoryV1API(ctx context.Context) (inventoryV1.InventoryServiceServer, error) {
	if d.inventoryV1API == nil {
		svc, err := d.PartService(ctx)
		if err != nil {
			return nil, err
		}
		d.inventoryV1API = inventoryAPI.NewAPI(svc)
	}
	return d.inventoryV1API, nil
}

func (d *diContainer) PartService(ctx context.Context) (service.PartService, error) {
	if d.partService == nil {
		repo, err := d.PartRepository(ctx)
		if err != nil {
			return nil, err
		}
		d.partService = partService.NewPartService(repo)
	}
	return d.partService, nil
}

func (d *diContainer) PartRepository(ctx context.Context) (repository.PartRepository, error) {
	if d.partRepository == nil {
		coll, err := d.PartsCollection(ctx)
		if err != nil {
			return nil, err
		}
		repo := partRepo.NewPartRepository(coll)
		if err := repo.Seed(ctx); err != nil {
			return nil, err
		}
		d.partRepository = repo
	}
	return d.partRepository, nil
}

func (d *diContainer) MongoClient(ctx context.Context) (*mongo.Client, error) {
	if d.mongoClient == nil {
		client, err := mongo.Connect(options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			return nil, err
		}
		if err := client.Ping(ctx, nil); err != nil {
			return nil, err
		}
		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})
		d.mongoClient = client
	}
	return d.mongoClient, nil
}

func (d *diContainer) PartsCollection(ctx context.Context) (*mongo.Collection, error) {
	if d.partsColl == nil {
		client, err := d.MongoClient(ctx)
		if err != nil {
			return nil, err
		}
		d.partsColl = client.Database(config.AppConfig().Mongo.DatabaseName()).Collection("parts")
	}
	return d.partsColl, nil
}
