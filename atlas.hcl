variable "database_url" {
  type    = string
  default = getenv("DATABASE_URL")
}

env "local" {
  migration {
    dir = "file://migrations"
  }

  # The Atlas container runs with --network host, so it reaches the host-mode
  # Postgres at the same localhost address the app uses.
  url = var.database_url
  src = "file://schema.sql"
  dev = urlsetpath(var.database_url, "atlas_dev")

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
