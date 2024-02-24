defmodule CdbWeb.ConfigKeyControllerTest do
  use CdbWeb.ConnCase

  import Cdb.ConfigurationFixtures

  @create_attrs %{name: "some name", value_type: :string}
  @update_attrs %{name: "some updated name", value_type: :string}
  @invalid_attrs %{name: nil, value_type: nil}

  describe "index" do
    test "lists all config-keys", %{conn: conn} do
      conn = get(conn, ~p"/config-keys")
      assert html_response(conn, 200) =~ "Listing Config keys"
    end
  end

  describe "new config_key" do
    test "renders form", %{conn: conn} do
      conn = get(conn, ~p"/config-keys/new")
      assert html_response(conn, 200) =~ "New Config key"
    end
  end

  describe "create config_key" do
    test "redirects to show when data is valid", %{conn: conn} do
      conn = post(conn, ~p"/config-keys", config_key: @create_attrs)

      assert %{id: id} = redirected_params(conn)
      assert redirected_to(conn) == ~p"/config-keys/#{id}"

      conn = get(conn, ~p"/config-keys/#{id}")
      assert html_response(conn, 200) =~ "Config key #{id}"
    end

    test "renders errors when data is invalid", %{conn: conn} do
      conn = post(conn, ~p"/config-keys", config_key: @invalid_attrs)
      assert html_response(conn, 200) =~ "New Config key"
    end
  end

  describe "edit config_key" do
    setup [:create_config_key]

    test "renders form for editing chosen config_key", %{conn: conn, config_key: config_key} do
      conn = get(conn, ~p"/config-keys/#{config_key}/edit")
      assert html_response(conn, 200) =~ "Edit Config key"
    end
  end

  describe "update config_key" do
    setup [:create_config_key]

    test "redirects when data is valid", %{conn: conn, config_key: config_key} do
      conn = put(conn, ~p"/config-keys/#{config_key}", config_key: @update_attrs)
      assert redirected_to(conn) == ~p"/config-keys/#{config_key}"

      conn = get(conn, ~p"/config-keys/#{config_key}")
      assert html_response(conn, 200) =~ "some updated name"
    end

    test "renders errors when data is invalid", %{conn: conn, config_key: config_key} do
      conn = put(conn, ~p"/config-keys/#{config_key}", config_key: @invalid_attrs)
      assert html_response(conn, 200) =~ "Edit Config key"
    end
  end

  describe "delete config_key" do
    setup [:create_config_key]

    test "deletes chosen config_key", %{conn: conn, config_key: config_key} do
      conn = delete(conn, ~p"/config-keys/#{config_key}")
      assert redirected_to(conn) == ~p"/config-keys"

      assert_error_sent 404, fn ->
        get(conn, ~p"/config-keys/#{config_key}")
      end
    end
  end

  defp create_config_key(_) do
    config_key = config_key_fixture()
    %{config_key: config_key}
  end
end
