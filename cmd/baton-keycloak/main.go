package main

// this connector should allow Conductor One to sync users and groups from Keycloak
// it should also allow for entitlement provisioning for users and groups
// this is mainly so that when someone at Weaviate tries to access a cluster, C1 can add them to the right group on a JIT basis.
import (
	"context"
	"fmt"
	"os"

	connectorSchema "github.com/spiros-spiros/baton-keycloak/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)
var version = "dev"
func main() {
		ctx := context.Background()
	
		_, cmd, err := config.DefineConfiguration(
			ctx,
			"baton-keycloak",
			getConnector,
			field.Configuration{
				Fields: ConfigurationFields,
			},
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	
		cmd.Version = version
	
		err = cmd.Execute()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

func getConnector(ctx context.Context, v *viper.Viper) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)
	if err := ValidateConfig(v); err != nil {
		return nil, err
	}

	keycloakServerURL := v.GetString(apiUrlField.FieldName)
	keycloakRealm := v.GetString(realmField.FieldName)
	keycloakClientID := v.GetString(clientField.FieldName)
	keycloakClientSecret := v.GetString(clientSecretField.FieldName)

	cb, err := connectorSchema.New(ctx, keycloakServerURL, keycloakRealm, keycloakClientID, keycloakClientSecret)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}
	connector, err := connectorbuilder.NewConnector(ctx, cb)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}
	return connector, nil
}
