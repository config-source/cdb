defmodule Cdb.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  @impl true
  def start(_type, _args) do
    children = [
      # Start the Telemetry supervisor
      CdbWeb.Telemetry,
      # Start the Ecto repository
      Cdb.Repo,
      # Start the PubSub system
      {Phoenix.PubSub, name: Cdb.PubSub},
      # Start Finch
      {Finch, name: Cdb.Finch},
      # Start the Endpoint (http/https)
      CdbWeb.Endpoint
      # Start a worker by calling: Cdb.Worker.start_link(arg)
      # {Cdb.Worker, arg}
    ]

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: Cdb.Supervisor]
    Supervisor.start_link(children, opts)
  end

  # Tell Phoenix to update the endpoint configuration
  # whenever the application is updated.
  @impl true
  def config_change(changed, _new, removed) do
    CdbWeb.Endpoint.config_change(changed, removed)
    :ok
  end
end
