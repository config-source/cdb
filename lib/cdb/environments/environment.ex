defmodule Cdb.Environments.Environment do
  use Cdb.Schema
  import Ecto.Changeset

  schema "environments" do
    field :name, :string

    belongs_to :parent, Cdb.Environments.Environment, foreign_key: :promotes_to
    has_many :children, Cdb.Environments.Environment, foreign_key: :promotes_to

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
