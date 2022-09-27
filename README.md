# Tool Belt

**Tool Belt is a work in progress.**

Tool Belt is a module to wrap up a number of smaller utility tools into a single application. Tools can share
HTTP routers, datastores and config to make it faster to deploy small tools quickly.

Tool Belt has been created to bring together a number of utility tools I have been running separately. Tools of mine to
be wrapped in Tool Belt for deployment include:

* a tool to read the contents of a URL and make string replacements on the fly. Used to reformat Calendar URLs.
* a recurring task to sync an AirTable contact database to CardDav.
* running of [json-charlieegan3](https://github.com/charlieegan3/json-charlieegan3)
* RSS feed substitution and filtering
* Twitter and Strava RSS feed generation
* Webhook to RSS feed server
* etc

As you can see the idea is to bring together a number of other small-tiny tools under a single monolithic deployment.
This is good because:

* I can run a single HTTP server and *pay* for a single server
* I can share things like a SQL database which takes time to configure
