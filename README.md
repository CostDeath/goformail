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

### Installing GoForMail

This is very easy, due to the provided dockerfile. Just download [these files](deploy) and run `docker compose up`.

Make sure you set email_domain to the domain you're setting your mailing lists to. If there is a mismatch, all emails
will be rejected.

### Configuring Postfix

Firstly, you will need to set up a hash db using postfix. First, create a text file. Below is an example of one:
```
lists.costwynn.xyz	lmtp:inet:127.0.0.1:8024
```
You should replace `lists.costwynn.xyz` with your list domain, and 8024 with the port you have allocated GoForMail
(it is set to 8024 on the dockerfile by default).

To create it into a proper hash db, run `postmap <filename>` on your db file.

Additionally, you should set the following properties in your `main.cf` file:
```
inet_interfaces = all
mynetworks = 172.20.0.0/16 127.0.0.0/8
relay_domains = lists.$mydomain
transport_maps = hash:/etc/postfix/goformail_lmtp
local_recipient_maps = hash:/etc/postfix/goformail_lmtp, proxy:unix:passwd.byname 
```
Update the path `/etc/postfix/goformail_lmtp` to your created db file. **DO NOT INCLUDE THE .db EXTENSION**.

These properties will redirect anything directed at `@lists.<your_domain>` to GoForMail.

### Configuring GoForMail

All configurations can be done through the `configs.cf` file. The dockerfile automatically mounts it, so restarting the
docker compose with `docker compose stop` followed by `docker compose up` will read the new values.

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

