defmodule Cdb.EnvironmentsFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `Cdb.Environments` context.
  """

  @doc """
  Generate a environment.
  """
  def environment_fixture(attrs \\ %{}) do
    {:ok, environment} =
      attrs
      |> Enum.into(%{
        name: "some name"
      })
      |> Cdb.Environments.create_environment()

    environment |> Cdb.Repo.preload([:promotes_to])
  end
end
