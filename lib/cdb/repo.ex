defmodule Cdb.Repo do
  use Ecto.Repo,
    otp_app: :cdb,
    adapter: Ecto.Adapters.Postgres
end
