defmodule CdbWeb.EnvironmentHTML do
  use CdbWeb, :html

  embed_templates "environment_html/*"

  @doc """
  Renders a environment form.
  """
  attr :changeset, Ecto.Changeset, required: true
  attr :action, :string, required: true

  def environment_form(assigns)
end
