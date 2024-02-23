defmodule Cdb.Environments do
  @moduledoc """
  The Environments context.
  """

  import Ecto.Query, warn: false
  alias Cdb.Repo

  alias Cdb.Environments.Environment

  @doc """
  Returns the list of environments.

  ## Examples

      iex> list_environments()
      [%Environment{}, ...]

  """
  def list_environments do
    Repo.all(Environment)
  end

  @doc """
  Gets a single environment.

  Raises `Ecto.NoResultsError` if the Environment does not exist.

  ## Examples

      iex> get_environment!(123)
      %Environment{}

      iex> get_environment!(456)
      ** (Ecto.NoResultsError)

  """
  def get_environment!(id), do: Repo.get!(Environment, id)

  @doc """
  Creates a environment.

  ## Examples

      iex> create_environment(%{field: value})
      {:ok, %Environment{}}

      iex> create_environment(%{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def create_environment(attrs \\ %{}) do
    %Environment{}
    |> Environment.changeset(attrs)
    |> Repo.insert()
  end

  @doc """
  Updates a environment.

  ## Examples

      iex> update_environment(environment, %{field: new_value})
      {:ok, %Environment{}}

      iex> update_environment(environment, %{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def update_environment(%Environment{} = environment, attrs) do
    environment
    |> Environment.changeset(attrs)
    |> Repo.update()
  end

  @doc """
  Deletes a environment.

  ## Examples

      iex> delete_environment(environment)
      {:ok, %Environment{}}

      iex> delete_environment(environment)
      {:error, %Ecto.Changeset{}}

  """
  def delete_environment(%Environment{} = environment) do
    Repo.delete(environment)
  end

  @doc """
  Returns an `%Ecto.Changeset{}` for tracking environment changes.

  ## Examples

      iex> change_environment(environment)
      %Ecto.Changeset{data: %Environment{}}

  """
  def change_environment(%Environment{} = environment, attrs \\ %{}) do
    Environment.changeset(environment, attrs)
  end
end
