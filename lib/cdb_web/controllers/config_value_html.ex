defmodule CdbWeb.ConfigValueHTML do
  use CdbWeb, :html

  embed_templates "config_value_html/*"

  @doc """
  Renders a config_value form.
  """
  attr :changeset, Ecto.Changeset, required: true
  attr :action, :string, required: true

  def config_value_form(assigns)
end
