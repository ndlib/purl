# purl
Permanent Url Resource Locator

[x] Get a "hello world" HTTP server running.
[x] Add in the read-only routes from the repopurl spec, and get it mostly working using hard-coded data. We can use either the Gorilla Mux library or the httprouter library, like ndlib/bendo uses.
[] Add in unit tests and get a Travis CI set up. We may need to look at ndlib/bendo for bits and pieces and maybe a Makefile
[] Now add in the database. We will probably make an interface to decouple the database engine from everything else. Will also add fixity data and get all the tests set up to use that.
[] Add a utility to set up local testing databases, so we can run mysql locally and connect to it. We need a utility since the webapp will not handle database migrations, unlike most rails and python webapps.
[] Can we also make a utility to handle making records? Right now I edit them by entering SQL by hand. This is a bonus project.
[] Figure out how we can deploy this, and integration test it.
