defmodule Cdb.Repo.Migrations.CreateEnvironments do
  use Ecto.Migration

  def change do
    create table(:environments, primary_key: false) do
      add :id, :uuid, primary_key: true
      add :name, :string
      add :promotes_to_id, references(:environments, on_delete: :nilify_all, type: :uuid)

      timestamps()
    end

    create index(:environments, [:promotes_to_id])
  end
end
