package main

import (
	"context"
	"fmt"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/internal/postgres"
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
	repository, err := postgres.NewRepository(context.Background(), "")
	fail(err)

	ctx := context.Background()

	fmt.Println("Seeding environments...")
	clearTable(repository, "environments")

	production, err := repository.CreateEnvironment(ctx, cdb.Environment{Name: "production"})
	fail(err)

	staging, err := repository.CreateEnvironment(ctx, cdb.Environment{Name: "staging", PromotesToID: &production.ID})
	fail(err)

	dev1, err := repository.CreateEnvironment(ctx, cdb.Environment{Name: "dev1", PromotesToID: &staging.ID})
	fail(err)

	dev2, err := repository.CreateEnvironment(ctx, cdb.Environment{Name: "dev2", PromotesToID: &staging.ID})
	fail(err)

	fmt.Println("Done seeding environments.")

	// here for a bit of "testing" and to keep go from freaking out about dev1
	// and dev2 being unused for now.
	for _, env := range []cdb.Environment{production, staging, dev1, dev2} {
		_, err = repository.GetEnvironmentByName(ctx, env.Name)
		fail(err)

		_, err = repository.GetEnvironment(ctx, env.ID)
		fail(err)
	}

	fmt.Println("Seeding config keys...")
	clearTable(repository, "config_keys")

	owner, err := repository.CreateConfigKey(ctx, cdb.NewConfigKey("owner", cdb.TypeString))
	fail(err)

	maxReplicas, err := repository.CreateConfigKey(ctx, cdb.NewConfigKey("maxReplicas", cdb.TypeInteger))
	fail(err)

	minReplicas, err := repository.CreateConfigKey(ctx, cdb.NewConfigKey("minReplicas", cdb.TypeInteger))
	fail(err)

	// here for a bit of "testing" and to keep go from freaking out about the
	// config key vars being unused for now.
	allKeys, err := repository.ListConfigKeys(ctx)
	fail(err)

	for idx, ck := range []cdb.ConfigKey{owner, maxReplicas, minReplicas} {
		if allKeys[idx].ID != ck.ID && allKeys[idx].Name != ck.Name && *allKeys[idx].CanPropagate != *ck.CanPropagate {
			fail(fmt.Errorf("Expected: %v Got: %v", ck, allKeys[idx]))
		}

		_, err = repository.GetConfigKey(ctx, ck.ID)
		fail(err)
	}

	fmt.Println("Done seeding config keys.")

}
