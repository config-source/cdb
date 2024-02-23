defmodule Cdb.Configuration.ConfigKey do
  use Cdb.Schema
  import Ecto.Changeset

  schema "config_keys" do
    field :name, :string
    field :value_type, Ecto.Enum, values: [:string, :integer, :float, :boolean]
    field :can_propagate, :boolean, default: true

    timestamps()
  end

  @doc false
  def changeset(config_key, attrs) do
    config_key
    |> cast(attrs, [:name])
    |> validate_required([:name])
  end
end
