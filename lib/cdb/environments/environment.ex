defmodule Cdb.Environments.Environment do
  use Cdb.Schema
  import Ecto.Changeset

  schema "environments" do
    field :name, :string

    belongs_to :promotes_to, Cdb.Environments.Environment, foreign_key: :promotes_to_id
    has_many :promotes_from, Cdb.Environments.Environment, foreign_key: :promotes_to_id

    has_many :config_values, Cdb.Configuration.ConfigValue

    timestamps()
  end

  @doc false
  def changeset(environment, attrs) do
    environment
    |> cast(attrs, [:name])
    |> validate_required([:name])
  end
end
