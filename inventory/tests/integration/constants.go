//go:build integration

package integration

const (
	projectName = "goshop-inventory-e2e"

	mongoImage         = "mongo:7.0.45"
	mongoContainerName = "inventory-e2e-mongo"
	mongoInternalPort  = "27017"
	mongoDatabase      = "inventory"
	mongoAuthDB        = "admin"
	mongoUsername      = "inventory_admin"
	mongoPassword      = "inventory_secret"

	// testSessionUUID сессия, которую заглушка IAM считает валидной
	testSessionUUID = "11111111-1111-1111-1111-111111111111"
	hostGatewayName = "host.docker.internal"

	appName        = "inventory-e2e-app"
	appPort        = "50051"
	dockerfilePath = "deploy/docker/inventory/Dockerfile"
)
