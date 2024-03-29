defmodule Cdb.Repo.Migrations.CreateConfigKeys do
  use Ecto.Migration

  def change do
    create table(:config_keys, primary_key: false) do
      add :id, :uuid, primary_key: true
      add :name, :string

      timestamps()
    end
  end
end
