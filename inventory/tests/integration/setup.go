//go:build integration

package integration

import (
	"context"
	"io"

	"github.com/ivmerk/goshop/platform/pkg/testcontainers"
	"github.com/ivmerk/goshop/platform/pkg/testcontainers/app"
	"github.com/ivmerk/goshop/platform/pkg/testcontainers/mongo"
	"github.com/ivmerk/goshop/platform/pkg/testcontainers/network"
	"github.com/ivmerk/goshop/platform/pkg/testcontainers/path"
)

func setupTestEnvironment(ctx context.Context) (*TestEnvironment, error) {
	env := &TestEnvironment{}

	var err error
	defer func() {
		if err != nil {
			teardownTestEnvironment(ctx, env)
		}
	}()

	env.Network, err = network.NewNetwork(ctx, projectName)
	if err != nil {
		return nil, err
	}

	env.Mongo, err = mongo.NewContainer(ctx,
		mongo.WithNetworkName(env.Network.Name()),
		mongo.WithContainerName(mongoContainerName),
		mongo.WithImageName(mongoImage),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithAuthDB(mongoAuthDB),
	)
	if err != nil {
		return nil, err
	}

	env.App, err = app.NewContainer(ctx,
		app.WithName(appName),
		app.WithPort(appPort),
		app.WithDockerfile(path.GetProjectRoot(), dockerfilePath),
		app.WithNetwork(env.Network.Name()),
		app.WithLogOutput(io.Discard), // для отладки замените на os.Stdout
		app.WithEnv(map[string]string{
			"GRPC_HOST":                     "0.0.0.0",
			"GRPC_PORT":                     appPort,
			"LOGGER_LEVEL":                  "info",
			"LOGGER_AS_JSON":                "true",
			testcontainers.MongoHostKey:     mongoContainerName, // адрес внутри docker-сети
			testcontainers.MongoPortKey:     mongoInternalPort,
			testcontainers.MongoDatabaseKey: mongoDatabase,
			testcontainers.MongoAuthDBKey:   mongoAuthDB,
			testcontainers.MongoUsernameKey: mongoUsername,
			testcontainers.MongoPasswordKey: mongoPassword,
		}),
	)
	if err != nil {
		return nil, err
	}

	return env, nil
}
