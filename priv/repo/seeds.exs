# Script for populating the database. You can run it as:
#
#     mix run priv/repo/seeds.exs
#
# Inside the script, you can read and write to any of your
# repositories directly:
#
#     Cdb.Repo.insert!(%Cdb.SomeSchema{})
#
# We recommend using the bang functions (`insert!`, `update!`
# and so on) as they will fail if something goes wrong.

production = Cdb.Repo.insert!(%Cdb.Environments.Environment{name: "production"})

staging =
  Cdb.Repo.insert!(%Cdb.Environments.Environment{name: "staging", promotes_to: production})

dev1 = Cdb.Repo.insert!(%Cdb.Environments.Environment{name: "dev1", promotes_to: staging})
dev2 = Cdb.Repo.insert!(%Cdb.Environments.Environment{name: "dev2", promotes_to: staging})

owner = Cdb.Repo.insert!(%Cdb.Configuration.ConfigKey{name: "owner", value_type: :string})

min_replicas =
  Cdb.Repo.insert!(%Cdb.Configuration.ConfigKey{name: "min_replicas", value_type: :integer})

max_replicas =
  Cdb.Repo.insert!(%Cdb.Configuration.ConfigKey{name: "max_replicas", value_type: :integer})

Cdb.Repo.insert!(%Cdb.Configuration.ConfigValue{
  str_value: "SRE",
  environment: production,
  config_key: owner
})

Cdb.Repo.insert!(%Cdb.Configuration.ConfigValue{
  int_value: 10,
  environment: production,
  config_key: min_replicas
})

Cdb.Repo.insert!(%Cdb.Configuration.ConfigValue{
  int_value: 100,
  environment: production,
  config_key: max_replicas
})

Cdb.Repo.insert!(%Cdb.Configuration.ConfigValue{
  int_value: 1,
  environment: staging,
  config_key: min_replicas
})

Cdb.Repo.insert!(%Cdb.Configuration.ConfigValue{
  int_value: 10,
  environment: staging,
  config_key: max_replicas
})
