defmodule CdbWeb.EnvironmentControllerTest do
  use CdbWeb.ConnCase

  import Cdb.EnvironmentsFixtures

  @create_attrs %{name: "some name"}
  @update_attrs %{name: "some updated name"}
  @invalid_attrs %{name: nil}

  describe "index" do
    test "lists all environments", %{conn: conn} do
      conn = get(conn, ~p"/environments")
      assert html_response(conn, 200) =~ "Listing Environments"
    end
  end

  describe "new environment" do
    test "renders form", %{conn: conn} do
      conn = get(conn, ~p"/environments/new")
      assert html_response(conn, 200) =~ "New Environment"
    end
  end

  describe "create environment" do
    test "redirects to show when data is valid", %{conn: conn} do
      conn = post(conn, ~p"/environments", environment: @create_attrs)

      assert %{id: id} = redirected_params(conn)
      assert redirected_to(conn) == ~p"/environments/#{id}"

      conn = get(conn, ~p"/environments/#{id}")
      assert html_response(conn, 200) =~ "Environment #{id}"
    end

    test "renders errors when data is invalid", %{conn: conn} do
      conn = post(conn, ~p"/environments", environment: @invalid_attrs)
      assert html_response(conn, 200) =~ "New Environment"
    end
  end

  describe "edit environment" do
    setup [:create_environment]

    test "renders form for editing chosen environment", %{conn: conn, environment: environment} do
      conn = get(conn, ~p"/environments/#{environment}/edit")
      assert html_response(conn, 200) =~ "Edit Environment"
    end
  end

  describe "update environment" do
    setup [:create_environment]

    test "redirects when data is valid", %{conn: conn, environment: environment} do
      conn = put(conn, ~p"/environments/#{environment}", environment: @update_attrs)
      assert redirected_to(conn) == ~p"/environments/#{environment}"

      conn = get(conn, ~p"/environments/#{environment}")
      assert html_response(conn, 200) =~ "some updated name"
    end

    test "renders errors when data is invalid", %{conn: conn, environment: environment} do
      conn = put(conn, ~p"/environments/#{environment}", environment: @invalid_attrs)
      assert html_response(conn, 200) =~ "Edit Environment"
    end
  end

  describe "delete environment" do
    setup [:create_environment]

    test "deletes chosen environment", %{conn: conn, environment: environment} do
      conn = delete(conn, ~p"/environments/#{environment}")
      assert redirected_to(conn) == ~p"/environments"

      assert_error_sent 404, fn ->
        get(conn, ~p"/environments/#{environment}")
      end
    end
  end

  defp create_environment(_) do
    environment = environment_fixture()
    %{environment: environment}
  end
end
