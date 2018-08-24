# Purl
Permanent URL Resource Locator

[![Go Report
Card](https://goreportcard.com/badge/github.com/ndlib/repopurl)](https://goreportcard.com/report/github.com/ndlib/repopurl)
![image](https://travis-ci.org/ndlib/repopurl.svg?branch=master)


## About

A simple daemon to resolve the legacy repository PURLs. For each purl the
upstream content is proxied back through this app, or a 302 redirect is
returned pointing to the content URL. The choice depends on the contents of the
`Information` field in the database. (refer to the database schema in [SPEC.md](SPEC.md)).
This application does not create or modify the database with the exception of
logging access usage.

## Configuration

All configuration is via environment variables.

 * `PORT` - the port number to listen on. Defaults to 8000
 * `MYSQL_CONNECTION` - how to connect to the database (see more below). Defaults to "/repopurl"
 * `FEDORA_USER` - the Fedora user name to use. Defaults to ""
 * `FEDORA_PASS` - the Fedora password to use. Defaults to ""
 * `TEMPLATE_PATH` - the directory to load our HTML files from. Defaults to "./templates"
 * `ROOT_REDIRECT` - the URL to redirect the root route `/` to. (Maybe this should be hardcoded?)

The only database supported is MySQL (since that is where the legacy database is).
The database connection string has the form

    username:password@tcp(hostname)/databasename

Where `username`, `password`, and `hostname` are optional.
(See https://github.com/go-sql-driver/mysql#dsn-data-source-name)
Examples of connection strings are

    /repopurl
    tcp(dbhost.example.com)/test
    jdoe:12345@tcp(localhost)/test_db

The Fedora user name and password are passed to an upstream Fedora server
using standard authentication if the `FEDORA_*` environment variables are
set, and the upstream content URL contains the string `fedora` somewhere.
This is to support legacy behavior.

## Testing

The unit tests are run with the `go test` command. These include basic integration tests on the web handlers.

To test the database connection, set up a test database and use
the `MYSQL_CONNECTION` environment variable to point to the test
database. Run the integration tests by passing the tag `mysql`, e.g.:

    env MYSQL_CONNECTION=/repopurl go test -tags=mysql

### Local Database Setup

To set up a local test database, install MySQL.
Then from a command line prompt run the following:

    $ mysql -u root
    mysql> create database repopurl;
    mysql> grant all on repopurl to ''@'localhost';
    mysql> use repopurl;
    mysql> source test/seed_data.sql;
    mysql> exit;

## Other

* [SPEC.md](SPEC.md) is the specification for the legacy service (that this one replaced).

