# Horus

:warning: non-functional - early brain dump

## Purpose
Horus takes a secret, a repository and a key, encrypt it and store it.

When retrieved, the encrypted string will be returned, so whoever requests it
must do the decryption.

The reason it takes a repository is that the intent is to store secrets needed for
build-time in a CI system, and have that system know how to decrypt the secret as well.

It does not implement its own encryption strategy. Initially it will only
rely on Google Cloud KMS, but more options may be added.

Submission of keys should always be done over TLS.

## Backend
Backend is written in Go, and is very simple.

Build with `go build`.
Requires Go 1.12 (or at the very least vgo).

### Configuration

See the [backend config_test.toml](backend/config/config_test.toml) file for available configuration settings.

### API
Examples use [httpie](https://httpie.org/)

#### Adding a secret
```bash
http --json POST "http://0.0.0.0:7654/api/v1/secret" "repo=github.com/alde/horus" "key=MY_SECRET" "secret=a-totally-secret-secret"
```
#### Getting a secret
Note: repo should be url-encoded
```bash
http "http://0.0.0.0:7654/api/v1/secret?repo=github.com%2Falde%2Fhorus&key=MY_SECRET"
```

#### Listing available secrets for a repository
Note the s
```bash
http "http://0.0.0.0:7654/api/v1/secrets?repo=github.com%2Falde%2Fhorus"
```

## CLI
Written in Go. Needs to be configured to talk to the backend, and what KMS to use to decrypt secrets.

You can either specify a config file with `-c/--config`, or place it in `${HOME}/.config/horus/cli_config.toml`

Build with `go build`.
Requires Go 1.12 (or at the very least vgo).

### Configuration
See the [cli config_test.toml](cli/config/config_test.toml) file for available configuration settings.

### Usage
```bash
horus -c config.toml download "http://github.com/alde/horus" MY_SECRET
```

## Frontend
Not implemented yet, but it's a react app.

# TODO
* Tests
* Frontend
* Probably much more
