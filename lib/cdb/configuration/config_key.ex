defmodule Cdb.Configuration.ConfigKey do
  use Cdb.Schema
  import Ecto.Changeset

  schema "config_keys" do
    field :name, :string

    timestamps()
  end

  @doc false
  def changeset(config_key, attrs) do
    config_key
    |> cast(attrs, [:name])
    |> validate_required([:name])
  end
end
