# Purl
Permanent Url Resource Locator

[![Go Report
Card](https://goreportcard.com/badge/github.com/ndlib/repopurl)](https://goreportcard.com/report/github.com/ndlib/repopurl)

## About
A Go web application keeping
all urls and serving them
statically to an admin with
access to the application

## Testing
To run test that need the database execute the sql file in
the test directory. In order to setup the connection to any database
set the environment variables "MYSQL_USER", "MYSQL_PASSWORD", "MYSQL_HOST",
"MYSQL_PORT" and "MYSQL_DB" to the correct values. Execute the file in test to
initialize the seed data with the correct values and the test database. Command
shown below:

`mysql --user=(someuser) --password=(somepassword) < ./test/seed_data.sql `


## Demo
currently live at https://repopurl.herokuapp.com

### TravisCI Build status:
![image](https://travis-ci.org/ndlib/repopurl.svg?branch=master)

## TODO
- [x] Get a "hello world" HTTP server running.
- [x] Add in the read-only routes from the repopurl spec, and get it mostly working using hard-coded data. We can use either the Gorilla Mux library or the httprouter library, like ndlib/bendo uses.
- [x] Add in unit tests and get a Travis CI set up. We may need to look at ndlib/bendo for bits and pieces and maybe a Makefile
- [x] Now add in the database. We will probably make an interface to decouple the database engine from everything else. Will also add fixity data and get all the tests set up to use that.
- [x] Add a utility to set up local testing databases, so we can run mysql locally and connect to it. We need a utility since the webapp will not handle database migrations, unlike most rails and python webapps.
- [x] Can we also make a utility to handle making records? Right now I edit them by entering SQL by hand. This is a bonus project.
- [x] Figure out how we can deploy this, and integration test it.
