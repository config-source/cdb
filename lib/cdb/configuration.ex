defmodule Cdb.Configuration do
  @moduledoc """
  The Configuration context.
  """

  import Ecto.Query, warn: false
  alias Cdb.Repo

  alias Cdb.Configuration.ConfigKey

  @doc """
  Returns the list of config_keys.

  ## Examples

      iex> list_config_keys()
      [%ConfigKey{}, ...]

  """
  def list_config_keys do
    Repo.all(ConfigKey)
  end

  @doc """
  Gets a single config_key.

  Raises `Ecto.NoResultsError` if the Config key does not exist.

  ## Examples

      iex> get_config_key!(123)
      %ConfigKey{}

      iex> get_config_key!(456)
      ** (Ecto.NoResultsError)

  """
  def get_config_key!(id), do: Repo.get!(ConfigKey, id)

  @doc """
  Creates a config_key.

  ## Examples

      iex> create_config_key(%{field: value})
      {:ok, %ConfigKey{}}

      iex> create_config_key(%{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def create_config_key(attrs \\ %{}) do
    %ConfigKey{}
    |> ConfigKey.changeset(attrs)
    |> Repo.insert()
  end

  @doc """
  Updates a config_key.

  ## Examples

      iex> update_config_key(config_key, %{field: new_value})
      {:ok, %ConfigKey{}}

      iex> update_config_key(config_key, %{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def update_config_key(%ConfigKey{} = config_key, attrs) do
    config_key
    |> ConfigKey.changeset(attrs)
    |> Repo.update()
  end

  @doc """
  Deletes a config_key.

  ## Examples

      iex> delete_config_key(config_key)
      {:ok, %ConfigKey{}}

      iex> delete_config_key(config_key)
      {:error, %Ecto.Changeset{}}

  """
  def delete_config_key(%ConfigKey{} = config_key) do
    Repo.delete(config_key)
  end

  @doc """
  Returns an `%Ecto.Changeset{}` for tracking config_key changes.

  ## Examples

      iex> change_config_key(config_key)
      %Ecto.Changeset{data: %ConfigKey{}}

  """
  def change_config_key(%ConfigKey{} = config_key, attrs \\ %{}) do
    ConfigKey.changeset(config_key, attrs)
  end

  alias Cdb.Configuration.ConfigValue

  @doc """
  Returns the list of config_values.

  ## Examples

      iex> list_config_values()
      [%ConfigValue{}, ...]

  """
  def list_config_values do
    Repo.all(ConfigValue)
  end

  @doc """
  Gets a single config_value.

  Raises `Ecto.NoResultsError` if the Config value does not exist.

  ## Examples

      iex> get_config_value!(123)
      %ConfigValue{}

      iex> get_config_value!(456)
      ** (Ecto.NoResultsError)

  """
  def get_config_value!(id), do: Repo.get!(ConfigValue, id)

  @doc """
  Creates a config_value.

  ## Examples

      iex> create_config_value(%{field: value})
      {:ok, %ConfigValue{}}

      iex> create_config_value(%{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def create_config_value(attrs \\ %{}) do
    %ConfigValue{}
    |> ConfigValue.changeset(attrs)
    |> Repo.insert()
  end

  @doc """
  Updates a config_value.

  ## Examples

      iex> update_config_value(config_value, %{field: new_value})
      {:ok, %ConfigValue{}}

      iex> update_config_value(config_value, %{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def update_config_value(%ConfigValue{} = config_value, attrs) do
    config_value
    |> ConfigValue.changeset(attrs)
    |> Repo.update()
  end

  @doc """
  Deletes a config_value.

  ## Examples

      iex> delete_config_value(config_value)
      {:ok, %ConfigValue{}}

      iex> delete_config_value(config_value)
      {:error, %Ecto.Changeset{}}

  """
  def delete_config_value(%ConfigValue{} = config_value) do
    Repo.delete(config_value)
  end

  @doc """
  Returns an `%Ecto.Changeset{}` for tracking config_value changes.

  ## Examples

      iex> change_config_value(config_value)
      %Ecto.Changeset{data: %ConfigValue{}}

  """
  def change_config_value(%ConfigValue{} = config_value, attrs \\ %{}) do
    ConfigValue.changeset(config_value, attrs)
  end
end
