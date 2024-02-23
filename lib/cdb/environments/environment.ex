defmodule Cdb.Environments.Environment do
  use Cdb.Schema
  import Ecto.Changeset

  schema "environments" do
    field :name, :string
    field :promotes_to, :id

    timestamps()
  end

  @doc false
  def changeset(environment, attrs) do
    environment
    |> cast(attrs, [:name])
    |> validate_required([:name])
  end
end
