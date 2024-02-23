defmodule Cdb.Repo.Migrations.CreateConfigValues do
  use Ecto.Migration

  def change do
    create table(:config_values, primary_key: false) do
      add :id, :uuid, primary_key: true
      add :config_key_id, references(:config_keys, on_delete: :nothing, type: :uuid)
      add :environment_id, references(:environments, on_delete: :nothing, type: :uuid)

      timestamps()
    end

    create index(:config_values, [:config_key_id])
    create index(:config_values, [:environment_id])
  end
end
