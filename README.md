# toolbelt

toolbelt is a module to wrap up a number of smaller utility tools into a single application & deployment. Tools can
share HTTP routers, datastores and config to make it faster to deploy small tools quickly.

toolbelt has been created to bring together a number of utility tools I have been running separately. See
[tools](#tools) for a list of tools. The idea is to bring together a number of other small tools under
a single monolithic deployment - a toolbelt.

This is good because:

* I can run a single long-running HTTP server for many tools and _pay_ for a single server
* I can share things like a SQL database which takes time to configure
* Where the workload is low, I can run jobs on the same server instance

I have written in more detail about why I created toolbelt [here on my blog](https://charlieegan3.com/posts/2022-10-10-toolbelt-building-a-personal-side-project-platform/).

## Tools

Below is a list I keep up to date of all the tools I've created for use on my toolbelt deployment:

* [tool-activities-rss](https://github.com/charlieegan3/tool-activities-rss)
* [tool-airtable-contacts](https://github.com/charlieegan3/tool-airtable-contacts)
* [tool-inoreader-github-actions-trigger](https://github.com/charlieegan3/tool-inoreader-github-actions-trigger)
* [tool-json-status](https://github.com/charlieegan3/tool-json-status)
* [tool-subpub](https://github.com/charlieegan3/tool-subpub)
* [tool-twitter-rss](https://github.com/charlieegan3/tool-twitter-rss)
* [tool-webhook-rss](https://github.com/charlieegan3/tool-webhook-rss)
* [tool-dropbox-backup](https://github.com/charlieegan3/tool-dropbox-backup)
* [food](https://github.com/charlieegan3/food/blob/main/pkg/tool/tool.go) website refresh

# Extenstions

Currently, it's possible to extend toolbelts with a means to external jobs, so far I have the built one integration I need:

* [toolbelt-external-job-runner-northflank](https://github.com/charlieegan3/toolbelt-external-job-runner-northflank) for running external jobs on [Northflank](https://northflank.com).

