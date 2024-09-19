package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/config-source/cdb/pkg/auth"
	authpg "github.com/config-source/cdb/pkg/auth/postgres"
	"github.com/config-source/cdb/pkg/configkeys"
	"github.com/config-source/cdb/pkg/configvalues"
	"github.com/config-source/cdb/pkg/environments"
	"github.com/config-source/cdb/pkg/services"
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

func resetIDs(pool *pgxpool.Pool, tableName string) {
	_, err := pool.Exec(context.Background(), fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART", tableName))
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
	svcRepo := services.NewRepository(logger, pool)
	valueRepo := configvalues.NewRepository(logger, pool)

	ctx := context.Background()

	fmt.Println("Seeding services")
	clearTable(pool, "services")

	makeSvc := func(name string) services.Service {
		svc, err := svcRepo.CreateService(ctx, services.Service{Name: name})
		if err != nil {
			panic(err)
		}

		return svc
	}

	services := []services.Service{
		makeSvc("legacy-service"),
		makeSvc("deployment-service"),
		makeSvc("delivery-service"),
	}

	clearTable(pool, "environments")
	clearTable(pool, "config_values")
	clearTable(pool, "config_keys")

	for _, svc := range services {
		fmt.Printf("Seeding config keys for %s...\n", svc.Name)

		owner, err := keyRepo.CreateConfigKey(ctx, configkeys.New(svc.ID, "owner", configkeys.TypeString))
		fail(err)

		maxReplicas, err := keyRepo.CreateConfigKey(ctx, configkeys.New(svc.ID, "maxReplicas", configkeys.TypeInteger))
		fail(err)

		minReplicas, err := keyRepo.CreateConfigKey(ctx, configkeys.New(svc.ID, "minReplicas", configkeys.TypeInteger))
		fail(err)

		sslEnabled, err := keyRepo.CreateConfigKey(ctx, configkeys.New(svc.ID, "sslEnabled", configkeys.TypeBoolean))
		fail(err)

		// Add an unconfigured config key for testing those features which require it.
		_, err = keyRepo.CreateConfigKey(ctx, configkeys.New(svc.ID, "readyForReaping", configkeys.TypeBoolean))
		fail(err)

		fmt.Println("Done seeding config keys.")

		fmt.Printf("Seeding environments for %s...\n", svc.Name)

		production, err := envRepo.CreateEnvironment(
			ctx,
			environments.Environment{
				Name:      "production",
				Sensitive: true,
				ServiceID: svc.ID,
			},
		)
		fail(err)

		staging, err := envRepo.CreateEnvironment(
			ctx,
			environments.Environment{
				Name:         "staging",
				PromotesToID: &production.ID,
				Sensitive:    true,
				ServiceID:    svc.ID,
			},
		)
		fail(err)

		dev, err := envRepo.CreateEnvironment(
			ctx,
			environments.Environment{
				Name:         "dev",
				PromotesToID: &staging.ID,
				ServiceID:    svc.ID,
			},
		)
		fail(err)

		fmt.Println("Done seeding environments.")

		fmt.Printf("Seeding config values for %s...\n", svc.Name)

		_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewString(
			production.ID,
			owner.ID,
			"SRE",
		))
		fail(err)

		_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewInt(
			production.ID,
			maxReplicas.ID,
			100,
		))
		fail(err)

		_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewInt(
			production.ID,
			minReplicas.ID,
			10,
		))
		fail(err)

		_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewBool(
			production.ID,
			sslEnabled.ID,
			true,
		))
		fail(err)

		_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewInt(
			staging.ID,
			minReplicas.ID,
			1,
		))
		fail(err)

		_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewInt(
			dev.ID,
			maxReplicas.ID,
			10,
		))
		fail(err)

		fmt.Println("Done seeding config values.")

		featureEnvCount := rand.Intn(100)
		fmt.Printf("Seeding %d feature environments for %s...\n", featureEnvCount, svc.Name)

		for i := range featureEnvCount {
			fe, err := envRepo.CreateEnvironment(
				ctx,
				environments.Environment{
					Name:         fmt.Sprintf("feature-environment-%d", i+1),
					PromotesToID: &staging.ID,
					ServiceID:    svc.ID,
				},
			)
			fail(err)

			_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewBool(
				fe.ID,
				sslEnabled.ID,
				false,
			))
			fail(err)

			switch rand.Intn(3) {
			case 0:
				_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewString(
					fe.ID,
					owner.ID,
					fmt.Sprintf("dev-team-%d", rand.Intn(10)),
				))
				fail(err)
			case 1:
				_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewInt(
					fe.ID,
					maxReplicas.ID,
					rand.Intn(30),
				))
				fail(err)
			case 2:
				_, err = valueRepo.CreateConfigValue(ctx, configvalues.NewInt(
					fe.ID,
					minReplicas.ID,
					rand.Intn(9)+1,
				))
				fail(err)
			}
		}

		fmt.Println("Done seeding feature environments.")
	}

	fmt.Println("Seeding users")
	clearTable(pool, "users_to_roles")
	clearTable(pool, "users")
	resetIDs(pool, "users")

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

	clearTable(pool, "api_tokens")

	_, err = pool.Query(
		ctx,
		"INSERT INTO api_tokens (user_id, token, created_at) VALUES ($1, $2, $3)",
		strconv.Itoa(int(adminUser.ID)),
		"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiRW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsImlzcyI6ImNkYiIsImlhdCI6MTcyNjA2Njk0M30.QkpwWcOXxYOhWvbo8sf0F-xlAfQ59X84ZlQaDHERxWc06nWtjjkh5vqTMl8haZ9mSaOK_-FwZXABrtGoV3nRuA",
		time.Now(),
	)
	fail(err)

	_, err = pool.Query(
		ctx,
		"INSERT INTO api_tokens (user_id, token, created_at) VALUES ($1, $2, $3)",
		strconv.Itoa(int(operatorUser.ID)),
		"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJJRCI6MiwiRW1haWwiOiJvcGVyYXRvckBleGFtcGxlLmNvbSIsImlzcyI6ImNkYiIsImlhdCI6MTcyNjA2Njk0M30.CzI0bmdz6MCJPPmcHGD_pltgCkZMTXxweX96Ejy4R499IqPSNqdS-igX9pvxyeJvtgn8jIuBnMIUUxjGKIia5A",
		time.Now(),
	)
	fail(err)
}
