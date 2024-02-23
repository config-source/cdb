defmodule CdbWeb.ConfigKeyController do
  use CdbWeb, :controller

  alias Cdb.Configuration
  alias Cdb.Configuration.ConfigKey

  def index(conn, _params) do
    config_keys = Configuration.list_config_keys()
    render(conn, :index, config_keys: config_keys)
  end

  def new(conn, _params) do
    changeset = Configuration.change_config_key(%ConfigKey{})
    render(conn, :new, changeset: changeset)
  end

  def create(conn, %{"config_key" => config_key_params}) do
    case Configuration.create_config_key(config_key_params) do
      {:ok, config_key} ->
        conn
        |> put_flash(:info, "Config key created successfully.")
        |> redirect(to: ~p"/config_keys/#{config_key}")

      {:error, %Ecto.Changeset{} = changeset} ->
        render(conn, :new, changeset: changeset)
    end
  end

  def show(conn, %{"id" => id}) do
    config_key = Configuration.get_config_key!(id)
    render(conn, :show, config_key: config_key)
  end

  def edit(conn, %{"id" => id}) do
    config_key = Configuration.get_config_key!(id)
    changeset = Configuration.change_config_key(config_key)
    render(conn, :edit, config_key: config_key, changeset: changeset)
  end

  def update(conn, %{"id" => id, "config_key" => config_key_params}) do
    config_key = Configuration.get_config_key!(id)

    case Configuration.update_config_key(config_key, config_key_params) do
      {:ok, config_key} ->
        conn
        |> put_flash(:info, "Config key updated successfully.")
        |> redirect(to: ~p"/config_keys/#{config_key}")

      {:error, %Ecto.Changeset{} = changeset} ->
        render(conn, :edit, config_key: config_key, changeset: changeset)
    end
  end

  def delete(conn, %{"id" => id}) do
    config_key = Configuration.get_config_key!(id)
    {:ok, _config_key} = Configuration.delete_config_key(config_key)

    conn
    |> put_flash(:info, "Config key deleted successfully.")
    |> redirect(to: ~p"/config_keys")
  end
end
