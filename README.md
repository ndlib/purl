# Purl
Permanent Url Resource Locator

## About
A Go web application keeping
all urls and serving them
statically to an admin with
access to the application

## TODO
- [x] Get a "hello world" HTTP server running.
- [x] Add in the read-only routes from the repopurl spec, and get it mostly working using hard-coded data. We can use either the Gorilla Mux library or the httprouter library, like ndlib/bendo uses.
- [x] Add in unit tests and get a Travis CI set up. We may need to look at ndlib/bendo for bits and pieces and maybe a Makefile
- [x] Now add in the database. We will probably make an interface to decouple the database engine from everything else. Will also add fixity data and get all the tests set up to use that.
- [x] Add a utility to set up local testing databases, so we can run mysql locally and connect to it. We need a utility since the webapp will not handle database migrations, unlike most rails and python webapps.
- [ ] Can we also make a utility to handle making records? Right now I edit them by entering SQL by hand. This is a bonus project.
- [x] Figure out how we can deploy this, and integration test it.
