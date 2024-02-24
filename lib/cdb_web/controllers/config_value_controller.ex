defmodule CdbWeb.ConfigValueController do
  use CdbWeb, :controller

  alias Cdb.Configuration
  alias Cdb.Configuration.ConfigValue

  def index(conn, _params) do
    config_values = Configuration.list_config_values()
    render(conn, :index, config_values: config_values)
  end

  def new(conn, _params) do
    changeset = Configuration.change_config_value(%ConfigValue{})
    render(conn, :new, changeset: changeset)
  end

  def create(conn, %{"config_value" => config_value_params}) do
    case Configuration.create_config_value(config_value_params) do
      {:ok, config_value} ->
        conn
        |> put_flash(:info, "Config value created successfully.")
        |> redirect(to: ~p"/config-values/#{config_value}")

      {:error, %Ecto.Changeset{} = changeset} ->
        render(conn, :new, changeset: changeset)
    end
  end

  def show(conn, %{"id" => id}) do
    config_value = Configuration.get_config_value!(id)
    render(conn, :show, config_value: config_value)
  end

  def edit(conn, %{"id" => id}) do
    config_value = Configuration.get_config_value!(id)
    changeset = Configuration.change_config_value(config_value)
    render(conn, :edit, config_value: config_value, changeset: changeset)
  end

  def update(conn, %{"id" => id, "config_value" => config_value_params}) do
    config_value = Configuration.get_config_value!(id)

    case Configuration.update_config_value(config_value, config_value_params) do
      {:ok, config_value} ->
        conn
        |> put_flash(:info, "Config value updated successfully.")
        |> redirect(to: ~p"/config-values/#{config_value}")

      {:error, %Ecto.Changeset{} = changeset} ->
        render(conn, :edit, config_value: config_value, changeset: changeset)
    end
  end

  def delete(conn, %{"id" => id}) do
    config_value = Configuration.get_config_value!(id)
    {:ok, _config_value} = Configuration.delete_config_value(config_value)

    conn
    |> put_flash(:info, "Config value deleted successfully.")
    |> redirect(to: ~p"/config-values")
  end
end
