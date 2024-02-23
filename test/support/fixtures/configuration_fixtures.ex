defmodule Cdb.ConfigurationFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `Cdb.Configuration` context.
  """

  @doc """
  Generate a config_key.
  """
  def config_key_fixture(attrs \\ %{}) do
    {:ok, config_key} =
      attrs
      |> Enum.into(%{
        name: "some name"
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

      })
      |> Cdb.Configuration.create_config_value()

    config_value
  end
end
