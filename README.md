# Tool Belt

Tool Belt is a module to wrap up a number of smaller utility tools into a single application & deployment. Tools can
share HTTP routers, datastores and config to make it faster to deploy small tools quickly.

Tool Belt has been created to bring together a number of utility tools I have been running separately. Tools I have
built for toolbelt:

* [tool-webhook-rss](https://github.com/charlieegan3/tool-webhook-rss)
* [tool-twitter-rss](https://github.com/charlieegan3/tool-twitter-rss)
* [tool-activities-rss](https://github.com/charlieegan3/tool-activities-rss)

As you can see, the idea is to bring together a number of other small to _tiny_ tools under a single monolithic
deployment.

This is good because:

* I can run a single long-running HTTP server for many tools and _pay_ for a single server
* I can share things like a SQL database which takes time to configure
* Where the workload is low, I can run jobs on the same server instance
