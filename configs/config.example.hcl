
theme = "default"

root_folder = "$HOME"

project "example" {
    theme = "default" # overwrite global
    path = "./example"

    env_vars = {
        "FOO" = "BAR"
    }

    env_vars_file = ".env"

    environment "dev" {
        theme = "default" # overwrite default and project
        
        env_vars = {
            "FOO" = "BAR_DEV"
        }

        env_vars_file = ".dev.env"
    }

    environment "pre" {
        theme = "default" # overwrite default and project
        
        env_vars = {
            "FOO" = "BAR_PRE"
        }

        env_vars_file = ".pre.env"
    }
}
