defmodule Cdb.Configuration.ConfigValue do
  use Cdb.Schema
  import Ecto.Changeset

  schema "config_values" do
    field :config_key_id, :id
    field :environment_id, :id

    timestamps()
  end

  @doc false
  def changeset(config_value, attrs) do
    config_value
    |> cast(attrs, [])
    |> validate_required([])
  end
end
