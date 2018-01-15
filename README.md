GitLab NotBot
-------------------------

This service is a workaround for the lack of "NOT" filter in GitLab issues.
There is an issue (feature) reported to GitLab requesting to introduce this,
but it has been planned, rephrased and delayed multiple times - https://gitlab.com/gitlab-org/gitlab-ce/issues/27747

This service automatically adds "negative" labels for all labels missing from issues.
It does it by subscribing to issue update webhook and updating the issue labels
whenever labels change.
