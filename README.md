<p align="center">
  <img alt="golangci-lint logo" src="docs/res/gfm.png" height="150" />
  <h3 align="center">GoForMail</h3>
  <p align="center">The Mailing List Manager written in Go</p>
</p>

---

GoForMail hooks up to your Mail Transfer Agent and reroutes any appropriate traffic.

It can be managed through a web UI or CLI, for convenience.

## Usage

Usage instructions can be found on the docs pdf

You can find additional information on how to use the [docs](docs) directory.

## Installing

<h3 style="color:#d4514e">GoForMail requires a separate MTA and DB to be running! It will not function standalone.</h3>

1. Install Postfix on your machine
2. Clone this repo
3. Go into /deploy
4. Feel free to use the configs provided! (adapt them to your domain)
5. docker compose up -d

## To do

- [X] Email Forwarding
- [X] List Management
- [X] List Locking
- [X] Web UI
- [X] User Management
- [X] Email Scheduling
- [X] Email Archiving
- [X] CLI interface

The project is still in early development. Please be patient as we develop it.

## Maintainers

This project is currently maintained by @fonseca3 and @dagohos2 (GitLab) or @CostDeath and @sdk194 (GitHub).

## Additional resources

- [Postfix](https://www.postfix.org/download.html) (Supported MTA)
- [PostgreSQL](https://www.postgresql.org/download/) (Supported DB)
- [Project Logo Source](https://mailtrap.io/blog/golang-send-email/)

