package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/internal/postgres"
	"github.com/rs/zerolog"
)

func fail(err error) {
	if err != nil {
		panic(err)
	}
}

func clearTable(repository *postgres.Repository, name string) {
	_, err := repository.Raw().Exec(context.Background(), fmt.Sprintf("DELETE FROM %s", name))
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

	repository, err := postgres.NewRepository(context.Background(), logger, "")
	fail(err)

	ctx := context.Background()

	fmt.Println("Seeding environments...")
	clearTable(repository, "environments")

	production, err := repository.CreateEnvironment(ctx, cdb.Environment{Name: "production"})
	fail(err)

	staging, err := repository.CreateEnvironment(ctx, cdb.Environment{Name: "staging", PromotesToID: &production.ID})
	fail(err)

	dev, err := repository.CreateEnvironment(ctx, cdb.Environment{Name: "dev", PromotesToID: &staging.ID})
	fail(err)

	fmt.Println("Done seeding environments.")

	fmt.Println("Seeding config keys...")
	clearTable(repository, "config_keys")

	owner, err := repository.CreateConfigKey(ctx, cdb.NewConfigKey("owner", cdb.TypeString))
	fail(err)

	maxReplicas, err := repository.CreateConfigKey(ctx, cdb.NewConfigKey("maxReplicas", cdb.TypeInteger))
	fail(err)

	minReplicas, err := repository.CreateConfigKey(ctx, cdb.NewConfigKey("minReplicas", cdb.TypeInteger))
	fail(err)

	sslEnabled, err := repository.CreateConfigKey(ctx, cdb.NewConfigKey("sslEnabled", cdb.TypeBoolean))
	fail(err)

	fmt.Println("Done seeding config keys.")

	fmt.Println("Seeding config values...")
	clearTable(repository, "config_values")

	_, err = repository.CreateConfigValue(ctx, cdb.NewStringConfigValue(
		production.ID,
		owner.ID,
		"SRE",
	))
	fail(err)

	_, err = repository.CreateConfigValue(ctx, cdb.NewIntConfigValue(
		production.ID,
		maxReplicas.ID,
		100,
	))
	fail(err)

	_, err = repository.CreateConfigValue(ctx, cdb.NewIntConfigValue(
		production.ID,
		minReplicas.ID,
		10,
	))
	fail(err)

	_, err = repository.CreateConfigValue(ctx, cdb.NewBoolConfigValue(
		production.ID,
		sslEnabled.ID,
		true,
	))
	fail(err)

	_, err = repository.CreateConfigValue(ctx, cdb.NewIntConfigValue(
		staging.ID,
		minReplicas.ID,
		1,
	))
	fail(err)

	_, err = repository.CreateConfigValue(ctx, cdb.NewIntConfigValue(
		dev.ID,
		maxReplicas.ID,
		10,
	))
	fail(err)

	fmt.Println("Done seeding config values.")

	fmt.Println("Seeding feature environments")

	for i := range 100 {
		fe, err := repository.CreateEnvironment(ctx, cdb.Environment{Name: fmt.Sprintf("feature-environment-%d", i+1), PromotesToID: &staging.ID})
		fail(err)

		_, err = repository.CreateConfigValue(ctx, cdb.NewBoolConfigValue(
			fe.ID,
			sslEnabled.ID,
			false,
		))
		fail(err)

		switch rand.Intn(3) {
		case 0:
			_, err = repository.CreateConfigValue(ctx, cdb.NewStringConfigValue(
				fe.ID,
				owner.ID,
				fmt.Sprintf("dev-team-%d", rand.Intn(10)),
			))
			fail(err)
		case 1:
			_, err = repository.CreateConfigValue(ctx, cdb.NewIntConfigValue(
				fe.ID,
				maxReplicas.ID,
				rand.Intn(30),
			))
			fail(err)
		case 2:
			_, err = repository.CreateConfigValue(ctx, cdb.NewIntConfigValue(
				fe.ID,
				minReplicas.ID,
				rand.Intn(9)+1,
			))
			fail(err)
		}
	}

	fmt.Println("Done seeding feature environments.")
}
