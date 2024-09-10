package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/config-source/cdb/auth"
	authpg "github.com/config-source/cdb/auth/postgres"
	"github.com/config-source/cdb/configkeys"
	"github.com/config-source/cdb/configvalues"
	"github.com/config-source/cdb/environments"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

func fail(err error) {
	if err != nil {
		panic(err)
	}
}

func clearTable(pool *pgxpool.Pool, name string) {
	_, err := pool.Exec(context.Background(), fmt.Sprintf("DELETE FROM %s", name))
	fail(err)
}

func main() {
	logger := zerolog.New(os.Stdout).
		Level(zerolog.ErrorLevel).
		With().
		Timestamp().
		Logger().
		Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})

	pool, err := pgxpool.New(context.Background(), "")
	fail(err)

	authGw := authpg.NewGateway(logger, pool)
	envRepo := environments.NewRepository(logger, pool)
	keyRepo := configkeys.NewRepository(logger, pool)
	valueRepo := configvalues.NewRepository(logger, pool)

	ctx := context.Background()

	fmt.Println("Seeding environments...")
	clearTable(pool, "environments")

	production, err := envRepo.CreateEnvironment(ctx, environments.Environment{Name: "production", Sensitive: true})
	fail(err)

	staging, err := envRepo.CreateEnvironment(ctx, environments.Environment{Name: "staging", PromotesToID: &production.ID, Sensitive: true})
	fail(err)

	dev, err := envRepo.CreateEnvironment(ctx, environments.Environment{Name: "dev", PromotesToID: &staging.ID})
	fail(err)

	fmt.Println("Done seeding environments.")

	fmt.Println("Seeding config keys...")
	clearTable(pool, "config_keys")

	owner, err := keyRepo.CreateConfigKey(ctx, configkeys.NewConfigKey("owner", configkeys.TypeString))
	fail(err)

	maxReplicas, err := keyRepo.CreateConfigKey(ctx, configkeys.NewConfigKey("maxReplicas", configkeys.TypeInteger))
	fail(err)

	minReplicas, err := keyRepo.CreateConfigKey(ctx, configkeys.NewConfigKey("minReplicas", configkeys.TypeInteger))
	fail(err)

	sslEnabled, err := keyRepo.CreateConfigKey(ctx, configkeys.NewConfigKey("sslEnabled", configkeys.TypeBoolean))
	fail(err)

	// Add an unconfigured config key for testing those features which require it.
	_, err = keyRepo.CreateConfigKey(ctx, configkeys.NewConfigKey("readyForReaping", configkeys.TypeBoolean))
	fail(err)

	fmt.Println("Done seeding config keys.")

	fmt.Println("Seeding config values...")
	clearTable(pool, "config_values")

	_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewStringConfigValue(
		production.ID,
		owner.ID,
		"SRE",
	))
	fail(err)

	_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewIntConfigValue(
		production.ID,
		maxReplicas.ID,
		100,
	))
	fail(err)

	_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewIntConfigValue(
		production.ID,
		minReplicas.ID,
		10,
	))
	fail(err)

	_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewBoolConfigValue(
		production.ID,
		sslEnabled.ID,
		true,
	))
	fail(err)

	_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewIntConfigValue(
		staging.ID,
		minReplicas.ID,
		1,
	))
	fail(err)

	_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewIntConfigValue(
		dev.ID,
		maxReplicas.ID,
		10,
	))
	fail(err)

	fmt.Println("Done seeding config values.")

	fmt.Println("Seeding feature environments")

	for i := range 100 {
		fe, err := envRepo.CreateEnvironment(ctx, environments.Environment{Name: fmt.Sprintf("feature-environment-%d", i+1), PromotesToID: &staging.ID})
		fail(err)

		_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewBoolConfigValue(
			fe.ID,
			sslEnabled.ID,
			false,
		))
		fail(err)

		switch rand.Intn(3) {
		case 0:
			_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewStringConfigValue(
				fe.ID,
				owner.ID,
				fmt.Sprintf("dev-team-%d", rand.Intn(10)),
			))
			fail(err)
		case 1:
			_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewIntConfigValue(
				fe.ID,
				maxReplicas.ID,
				rand.Intn(30),
			))
			fail(err)
		case 2:
			_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewIntConfigValue(
				fe.ID,
				minReplicas.ID,
				rand.Intn(9)+1,
			))
			fail(err)
		}
	}

	fmt.Println("Done seeding feature environments.")

	fmt.Println("Seeding users")
	clearTable(pool, "users")
	clearTable(pool, "users_to_roles")

	adminUser, err := authGw.CreateUser(ctx, auth.User{
		Email:    "admin@example.com",
		Password: "password",
	})
	fail(err)
	fail(authGw.AssignRoleToUserNoAuth(ctx, adminUser, "Administrator"))

	operatorUser, err := authGw.CreateUser(ctx, auth.User{
		Email:    "operator@example.com",
		Password: "password",
	})
	fail(err)
	fail(authGw.AssignRoleToUserNoAuth(ctx, operatorUser, "Operator"))

	fmt.Println("Done seeding users")
}
