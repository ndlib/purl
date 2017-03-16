# Purl
Permanent URL Resource Locator

[![Go Report
Card](https://goreportcard.com/badge/github.com/ndlib/repopurl)](https://goreportcard.com/report/github.com/ndlib/repopurl)
![image](https://travis-ci.org/ndlib/repopurl.svg?branch=master)


## About

A simple daemon to resolve the legacy repository PURLs. For each purl, a
database is consulted and either the upstream content is proxied back through
this app, or a 302 redirect is returned pointing to the content URL. The choice
depends on the contents of the `Information` field in the database. (refer to
the database schema in SPEC.md). This application does not create or modify the
database with the exception of logging access usage.

## Configuration

All configuration is via environment variables.

 * `PORT` - the port number to listen on. Defaults to 8000
 * `MYSQL_CONNECTION` - how to connect to the database (see more below). Defaults to "/repopurl"
 * `FEDORA_USER` - the Fedora user name to use. Defaults to ""
 * `FEDORA_PASS` - the Fedora password to use. Defaults to ""

The only database supported is MySQL. The database connection string has the form

    username:password@tcp(hostname)/databasename

Where the user name, password, and host name are optional.
(See https://github.com/go-sql-driver/mysql#dsn-data-source-name)
Examples of connection strings are

    /repopurl
    tcp(dbhost.example.com)/test
    jdoe:12345@tcp(localhost)/test_db

The Fedora user name and password are passed to the upstream server using HTTP
standard auth only if 1) the envrionment variables are set, and 2) the upstream
content URL contains the string `fedora` somewhere. This is to support legacy
behavior.

## Testing
The unit tests are run with the `go test` command.

To run the database integration tests, the database needs to be seeded with the
test data in the file `test/seed_data.sql`. For example, the following
command would do:

    mysql --user=jdoe -h hostname databasename -p < ./test/seed_data.sql

To run the tests, set the `MYSQL_CONNECTION` environment variable and tell `go
test` to run the integration tests:

    env MYSQL_CONNECTION=/repopurl go test -tags=mysql

## Demo
Currently live at https://repopurl.herokuapp.com
