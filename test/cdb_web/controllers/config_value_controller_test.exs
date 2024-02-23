defmodule CdbWeb.ConfigValueControllerTest do
  use CdbWeb.ConnCase

  import Cdb.ConfigurationFixtures

  @create_attrs %{}
  @update_attrs %{}
  @invalid_attrs %{}

  describe "index" do
    test "lists all config_values", %{conn: conn} do
      conn = get(conn, ~p"/config_values")
      assert html_response(conn, 200) =~ "Listing Config values"
    end
  end

  describe "new config_value" do
    test "renders form", %{conn: conn} do
      conn = get(conn, ~p"/config_values/new")
      assert html_response(conn, 200) =~ "New Config value"
    end
  end

  describe "create config_value" do
    test "redirects to show when data is valid", %{conn: conn} do
      conn = post(conn, ~p"/config_values", config_value: @create_attrs)

      assert %{id: id} = redirected_params(conn)
      assert redirected_to(conn) == ~p"/config_values/#{id}"

      conn = get(conn, ~p"/config_values/#{id}")
      assert html_response(conn, 200) =~ "Config value #{id}"
    end

    test "renders errors when data is invalid", %{conn: conn} do
      conn = post(conn, ~p"/config_values", config_value: @invalid_attrs)
      assert html_response(conn, 200) =~ "New Config value"
    end
  end

  describe "edit config_value" do
    setup [:create_config_value]

    test "renders form for editing chosen config_value", %{conn: conn, config_value: config_value} do
      conn = get(conn, ~p"/config_values/#{config_value}/edit")
      assert html_response(conn, 200) =~ "Edit Config value"
    end
  end

  describe "update config_value" do
    setup [:create_config_value]

    test "redirects when data is valid", %{conn: conn, config_value: config_value} do
      conn = put(conn, ~p"/config_values/#{config_value}", config_value: @update_attrs)
      assert redirected_to(conn) == ~p"/config_values/#{config_value}"

      conn = get(conn, ~p"/config_values/#{config_value}")
      assert html_response(conn, 200)
    end

    test "renders errors when data is invalid", %{conn: conn, config_value: config_value} do
      conn = put(conn, ~p"/config_values/#{config_value}", config_value: @invalid_attrs)
      assert html_response(conn, 200) =~ "Edit Config value"
    end
  end

  describe "delete config_value" do
    setup [:create_config_value]

    test "deletes chosen config_value", %{conn: conn, config_value: config_value} do
      conn = delete(conn, ~p"/config_values/#{config_value}")
      assert redirected_to(conn) == ~p"/config_values"

      assert_error_sent 404, fn ->
        get(conn, ~p"/config_values/#{config_value}")
      end
    end
  end

  defp create_config_value(_) do
    config_value = config_value_fixture()
    %{config_value: config_value}
  end
end
