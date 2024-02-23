defmodule Cdb.Repo.Migrations.AddValueColumnsToConfigValue do
  use Ecto.Migration

  def change do
    alter table("config_values") do
      add :str_value, :text
      add :int_value, :integer
      add :float_value, :float
      add :bool_value, :boolean, default: false
    end
  end
end
