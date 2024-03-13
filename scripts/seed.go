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

func main() {
	repository, err := postgres.NewRepository(context.Background(), "")
	fail(err)

	ctx := context.Background()

	_, err = repository.Raw().Exec(ctx, "DELETE FROM environments")
	fail(err)

	fmt.Println("Seeding environments...")
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
}
