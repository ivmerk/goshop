//go:build integration

package integration

import (
	"github.com/ivmerk/goshop/platform/pkg/testcontainers/app"
	"github.com/ivmerk/goshop/platform/pkg/testcontainers/mongo"
	"github.com/ivmerk/goshop/platform/pkg/testcontainers/network"
)

type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}
