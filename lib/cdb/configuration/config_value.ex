defmodule Cdb.Configuration.ConfigValue do
  use Cdb.Schema
  import Ecto.Changeset

  schema "config_values" do
    field :str_value, :string
    field :int_value, :integer
    field :float_value, :float
    field :bool_value, :boolean, default: false

    belongs_to :environment, Cdb.Environments.Environment
    belongs_to :config_key, Cdb.Configuration.ConfigKey

    timestamps()
  end

  @doc false
  def changeset(config_value, attrs) do
    config_value
    |> cast(attrs, [])
    |> validate_required([])
  end
end
