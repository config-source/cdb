defmodule CdbWeb.ConfigValueControllerTest do
  use CdbWeb.ConnCase

  import Cdb.ConfigurationFixtures
  import Cdb.EnvironmentsFixtures

  describe "index" do
    test "lists all config-values", %{conn: conn} do
      conn = get(conn, ~p"/config-values")
      assert html_response(conn, 200) =~ "Listing Config values"
    end
  end

  describe "new config_value" do
    test "renders form", %{conn: conn} do
      conn = get(conn, ~p"/config-values/new")
      assert html_response(conn, 200) =~ "New Config value"
    end
  end

  describe "create config_value" do
    test "redirects to show when data is valid", %{conn: conn} do
      conn =
        post(conn, ~p"/config-values",
          config_value: %{
            str_value: "some value",
            environment_id: environment_fixture().id,
            config_key_id: config_key_fixture().id
          }
        )

      assert %{id: id} = redirected_params(conn)
      assert redirected_to(conn) == ~p"/config-values/#{id}"

      conn = get(conn, ~p"/config-values/#{id}")
      assert html_response(conn, 200) =~ "Config value #{id}"
    end

    test "renders errors when data is invalid", %{conn: conn} do
      conn =
        post(conn, ~p"/config-values",
          config_value: %{
            str_value: 1
          }
        )

      assert html_response(conn, 200) =~ "New Config value"
    end
  end

  describe "edit config_value" do
    setup [:create_config_value]

    test "renders form for editing chosen config_value", %{conn: conn, config_value: config_value} do
      conn = get(conn, ~p"/config-values/#{config_value}/edit")
      assert html_response(conn, 200) =~ "Edit Config value"
    end
  end

  describe "update config_value" do
    setup [:create_config_value]

    test "redirects when data is valid", %{conn: conn, config_value: config_value} do
      conn =
        put(conn, ~p"/config-values/#{config_value}",
          config_value: %{
            str_value: "updated"
          }
        )

      assert redirected_to(conn) == ~p"/config-values/#{config_value}"

      conn = get(conn, ~p"/config-values/#{config_value}")
      assert html_response(conn, 200)
    end

    test "renders errors when data is invalid", %{conn: conn, config_value: config_value} do
      conn =
        put(conn, ~p"/config-values/#{config_value}",
          config_value: %{
            str_value: 1
          }
        )

      assert html_response(conn, 200) =~ "Edit Config value"
    end
  end

  describe "delete config_value" do
    setup [:create_config_value]

    test "deletes chosen config_value", %{conn: conn, config_value: config_value} do
      conn = delete(conn, ~p"/config-values/#{config_value}")
      assert redirected_to(conn) == ~p"/config-values"

      assert_error_sent 404, fn ->
        get(conn, ~p"/config-values/#{config_value}")
      end
    end
  end

  defp create_config_value(_) do
    config_value = config_value_fixture()
    %{config_value: config_value}
  end
end
