# Project Manager

Project manager is a terminal client application to manage multiple projects and configurations, in a easy and organized way.

# Features

## Create projects

Create a new folder with a `.project.hcl` file, if the folder already exist initialize the folder with the `.project.hcl` file.

The project environments are stored in a `.env` file also created in the project folder.

The project manager also managed the user git configuration updating the `.gitconfig` to store the path and a custom `%s.gitconfig` file to manage the git user, email, signkeys...

## Add Environments to a Project

## Update Projects and environments

## Delete Projects

# TUI and Cli

The project manager allow the user interact with a Terminal User Interface or with a set of cli commands.

## TUI

`pm [path]` run a new TUI to select start a new session

## CLI

-   `pm [project] [env|[path]]` start a new terminal session
-   `pm new [project]` add a new project
-   `pm new [project] [env]` add a new environment to the selected project
-   `pm config` open the configuration to edit global configs
-   `pm edit [project]` edit the selected project
-   `pm edit [project] [env]` edit the selected environment of a project
-   `pm delete [project]` delete a project and all associated environments
-   `pm delete [project] [env]` delete the selected environment of a project

# Commands

-   build: `make build`
-   lint: `make lint`
-   fix linter errors: `make lint-fix`
-   test: `make test`
