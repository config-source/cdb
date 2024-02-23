defmodule Cdb.Repo.Migrations.AddValueColumnsToConfigKeys do
  use Ecto.Migration

  def change do
    alter table("config_keys") do
      add :value_type, :string
      add :can_propagate, :boolean, default: false, null: false
    end
  end
end
