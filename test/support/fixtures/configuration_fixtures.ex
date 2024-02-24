defmodule Cdb.ConfigurationFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `Cdb.Configuration` context.
  """

  import Cdb.EnvironmentsFixtures

  @doc """
  Generate a config_key.
  """
  def config_key_fixture(attrs \\ %{}) do
    {:ok, config_key} =
      attrs
      |> Enum.into(%{
        name: "some name",
        value_type: :string
      })
      |> Cdb.Configuration.create_config_key()

    config_key
  end

  @doc """
  Generate a config_value.
  """
  def config_value_fixture(attrs \\ %{}) do
    {:ok, config_value} =
      attrs
      |> Enum.into(%{
        environment_id: environment_fixture().id,
        config_key_id: config_key_fixture().id,
        str_value: "some value"
      })
      |> Cdb.Configuration.create_config_value()

    config_value
  end
end
