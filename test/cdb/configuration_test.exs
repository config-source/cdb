defmodule Cdb.ConfigurationTest do
  use Cdb.DataCase

  alias Cdb.Configuration

  describe "config_keys" do
    alias Cdb.Configuration.ConfigKey

    import Cdb.ConfigurationFixtures

    @invalid_attrs %{name: nil}

    test "list_config_keys/0 returns all config_keys" do
      config_key = config_key_fixture()
      assert Configuration.list_config_keys() == [config_key]
    end

    test "get_config_key!/1 returns the config_key with given id" do
      config_key = config_key_fixture()
      assert Configuration.get_config_key!(config_key.id) == config_key
    end

    test "create_config_key/1 with valid data creates a config_key" do
      valid_attrs = %{name: "some name", value_type: :string}

      assert {:ok, %ConfigKey{} = config_key} = Configuration.create_config_key(valid_attrs)
      assert config_key.name == "some name"
    end

    test "create_config_key/1 with invalid data returns error changeset" do
      assert {:error, %Ecto.Changeset{}} = Configuration.create_config_key(@invalid_attrs)
    end

    test "update_config_key/2 with valid data updates the config_key" do
      config_key = config_key_fixture()
      update_attrs = %{name: "some updated name"}

      assert {:ok, %ConfigKey{} = config_key} =
               Configuration.update_config_key(config_key, update_attrs)

      assert config_key.name == "some updated name"
    end

    test "update_config_key/2 with invalid data returns error changeset" do
      config_key = config_key_fixture()

      assert {:error, %Ecto.Changeset{}} =
               Configuration.update_config_key(config_key, @invalid_attrs)

      assert config_key == Configuration.get_config_key!(config_key.id)
    end

    test "delete_config_key/1 deletes the config_key" do
      config_key = config_key_fixture()
      assert {:ok, %ConfigKey{}} = Configuration.delete_config_key(config_key)
      assert_raise Ecto.NoResultsError, fn -> Configuration.get_config_key!(config_key.id) end
    end

    test "change_config_key/1 returns a config_key changeset" do
      config_key = config_key_fixture()
      assert %Ecto.Changeset{} = Configuration.change_config_key(config_key)
    end
  end

  describe "config_values" do
    alias Cdb.Configuration.ConfigValue

    import Cdb.ConfigurationFixtures
    import Cdb.EnvironmentsFixtures

    @invalid_attrs %{str_value: 1}

    test "list_config_values/0 returns all config_values" do
      config_value =
        Repo.preload(
          config_value_fixture(),
          [:environment, :config_key]
        )

      assert Configuration.list_config_values() == [config_value]
    end

    test "get_config_value!/1 returns the config_value with given id" do
      config_value =
        Repo.preload(
          config_value_fixture(),
          [:environment, :config_key]
        )

      assert Configuration.get_config_value!(config_value.id) == config_value
    end

    test "create_config_value/1 with valid data creates a config_value" do
      valid_attrs = %{
        str_value: "somevalue",
        environment_id: environment_fixture().id,
        config_key_id: config_key_fixture().id
      }

      assert {:ok, %ConfigValue{} = _config_value} =
               Configuration.create_config_value(valid_attrs)
    end

    test "create_config_value/1 with invalid data returns error changeset" do
      assert {:error, %Ecto.Changeset{}} = Configuration.create_config_value(@invalid_attrs)
    end

    test "update_config_value/2 with valid data updates the config_value" do
      config_value = config_value_fixture()
      update_attrs = %{}

      assert {:ok, %ConfigValue{} = _config_value} =
               Configuration.update_config_value(config_value, update_attrs)
    end

    test "update_config_value/2 with invalid data returns error changeset" do
      config_value = config_value_fixture()

      assert {:error, %Ecto.Changeset{}} =
               Configuration.update_config_value(config_value, @invalid_attrs)

      assert Repo.preload(
               config_value,
               [:environment, :config_key]
             ) == Configuration.get_config_value!(config_value.id)
    end

    test "delete_config_value/1 deletes the config_value" do
      config_value = config_value_fixture()
      assert {:ok, %ConfigValue{}} = Configuration.delete_config_value(config_value)
      assert_raise Ecto.NoResultsError, fn -> Configuration.get_config_value!(config_value.id) end
    end

    test "change_config_value/1 returns a config_value changeset" do
      config_value = config_value_fixture()
      assert %Ecto.Changeset{} = Configuration.change_config_value(config_value)
    end

    test "get_configuration/1 handles inheritance" do
      production = environment_fixture(%{name: "production"})
      staging = environment_fixture(%{name: "staging", promotes_to: production})
      dev = environment_fixture(%{name: "dev", promotes_to: staging})

      owner = config_key_fixture(%{name: "owner"})
      max_replicas = config_key_fixture(%{name: "max_replicas", value_type: :integer})
      min_replicas = config_key_fixture(%{name: "min_replicas", value_type: :integer})

      owner_value = config_value_fixture(%{
        environment_id: production.id,
        config_key_id: owner.id,
        str_value: "SRE"
      })

      min_replicas_value_prod = config_value_fixture(%{
        environment_id: production.id,
        config_key_id: min_replicas.id,
        str_value: nil,
        int_value: 10,
      })
      max_replicas_value_prod = config_value_fixture(%{
        environment_id: production.id,
        config_key_id: max_replicas.id,
        str_value: nil,
        int_value: 100,
      })

      min_replicas_value_staging = config_value_fixture(%{
        environment_id: staging.id,
        config_key_id: min_replicas.id,
        str_value: nil,
        int_value: 5,
      })
      max_replicas_value_staging = config_value_fixture(%{
        environment_id: staging.id,
        config_key_id: max_replicas.id,
        str_value: nil,
        int_value: 50,
      })

      min_replicas_value_dev = config_value_fixture(%{
        environment_id: dev.id,
        config_key_id: min_replicas.id,
        str_value: nil,
        int_value: 1,
      })
      max_replicas_value_dev = config_value_fixture(%{
        environment_id: dev.id,
        config_key_id: max_replicas.id,
        str_value: nil,
        int_value: 5,
      })

      configs = Configuration.get_configuration(dev)
      expected = Enum.sort(Enum.map([owner_value, min_replicas_value_dev, max_replicas_value_dev], fn cv -> cv |> Repo.preload(:config_key) end))
      assert Enum.sort(configs) == expected
    end
  end
end
