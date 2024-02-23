defmodule CdbWeb.ConfigKeyHTML do
  use CdbWeb, :html

  embed_templates "config_key_html/*"

  @doc """
  Renders a config_key form.
  """
  attr :changeset, Ecto.Changeset, required: true
  attr :action, :string, required: true

  def config_key_form(assigns)
end
