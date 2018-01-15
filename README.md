# GitLab NotBot

This service is a workaround for the lack of "NOT" filter in GitLab issues.
There is an [issue (feature)](https://gitlab.com/gitlab-org/gitlab-ce/issues/27747)
reported to GitLab requesting to introduce this, but it has been planned,
rephrased and delayed multiple times.

This service automatically adds "negative" labels for all labels missing from issues.
It does this by listening for issue update webhook and updating the issue labels
whenever labels change.

For example, if a project has labels "bug", "design", "testing" and an issue is
created with label "bug", this bot will add labels "~design" and "~testing". This
way it'll be possible, for example, to find all issues that do not have "testing"
tag, as they will have tag "~testing".

## Installation
Download and compile by running
```
go get github.com/astax-t/gitlab-notbot
``` 
This creates the executable file `$GOPATH/bin/gitlab-notbot[.exe]`. Follow the
[Cross Compiling guide](http://golangcookbook.com/chapters/running/cross-compiling/)
if you need to build the binary for a different platform.

## Configuration
The service uses environment variables and/or `.env` file for configuration. Take
a look at the example file `.env-example` for the list of available variables. You
can copy this file to `.env`, update the values and put it next to the executable.
Note, variables defined in `.env` file _do not_ overwrite the variables already
defined in environment.

So the installation steps are the following:
  * Go to your settings in GitLab and click "Access Tokens". Create a new token and
    allow "api" scope access. If you're running a standalone instance of GitLab,
	it's not a bad idea to create a separate user for this - call it "NotBot",
	for example and give it "Reporter" right for the GitLab project.
  * Put the GitLab URL and the access token into `.env` file.
  * Also you'll probably want to change the value of `LISTEN_HOST` unless you're
    running the bot on the same server as GitLab. Set it to "*" in this case.
  * Alternatively, put all these values into environment variables.
  * Run the service executable.
  * Go to GitLab again, Project Settings, Integrations, and add a new webhook
    pointing to the NotBot service, such as `http://localhost:8085`. Set
	`Issue events` checkbox and unset all other.
  * Test the webhook for issue event and check the output of the NotBot service.
    You may set higher debug level for initial tests.

## TODO / Known problems
  * Currently list of all available labels is fetched for every issue update. This
    works ok for not very busy instances of GitLab, but may be slow on higher load.
  * Only up to 100 project labels are fetched, which leaves only 50 "real" labels
    which can be effectively used (due to 50 more "negative" labels). But if you
	have more than this, probably this bot shouldn't be used as each issue will
	have too many labels.
  * Perhaps it would be good to have a filter/regexp to restrict which labels are
    included into processing. Some labels may not need "NOT" filter.
