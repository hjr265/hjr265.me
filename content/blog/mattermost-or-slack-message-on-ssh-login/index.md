---
title: Mattermost (or Slack) Message on SSH Login
date: 2023-10-05T15:05:00+06:00
tags:
  - SSH
  - Mattermost
  - 100DaysToOffload
  - Slack
  - Curl
---

You have a server that you access over SSH. You have hardened it following the necessary best practices.

Now you can do one small thing for a little additional peace of mind: Set up Linux Pluggable Authentication Modules (PAM) to send a message to Mattermost (or Slack) on every successful SSH login.

{{< image src="example.png" >}}

Here is how you can do it:

1. Add a script to `/usr/local/bin/` to send the notification message. Name it `sshnotify.sh`. Make it executable.

    ``` sh
    #!/bin/sh
    MATTERMOST_HOST='...'
    MATTERMOST_WEBHOOK_KEY='...'
    MATTERMOST_CHANNEL='...'
    REMOTE_IP=`echo $SSH_CONNECTION | awk '{print $1}'`
    SERVER_HOSTNAME=`hostname`
    curl -i -X POST \
        -H 'Content-Type: application/json' \
        -d '{"channel":"'${MATTERMOST_CHANNEL}'","text": "Someone just logged in to your server '${SERVER_HOSTNAME}' from '${REMOTE_IP}'"}' \
        https://${MATTERMOST_HOST}/hooks/${MATTERMOST_WEBHOOK_KEY}
    ```

    ``` sh {linenos=false}
    chmod +x sshnotify.sh
    ```

    If you want to do the same but use Slack, then change the `curl` command accordingly.

    ``` sh
    curl -i -X POST \
        -H 'Content-Type: application/json' \
        -H "Authorization: Bearer ${SLACK_TOKEN}"
        -d '{"channel":"'${SLACK_CHANNEL}'","blocks":[{"type":"section","text":{"type":"mrkdwn","text":"Someone just logged in to your server '${SERVER_HOSTNAME}' from '${REMOTE_IP}'"}}]}' \
        https://slack.com/api/chat.postMessage
    ```

    Related documentations:

    - [Mattermost Webhooks](https://developers.mattermost.com/integrate/webhooks/)
    - [Posting Slack Messages Using Curl](https://api.slack.com/tutorials/tracks/posting-messages-with-curl)

2. Add the following line to the bottom of `/etc/pam.d/sshd`.

    ``` sh {linenos=false}
    session    optional     pam_exec.so /usr/local/bin/sshnotify.sh
    ```

    You can change `optional` to `require` to force PAM to allow the SSH connection only when the script runs successfully.

    Any time you are modifying `/etc/pam.d/sshd` be sure to keep a SSH connection active in a separate terminal window to avoid being locked out of your server because of bad configuration.

And that's it. 

Any time someone logs in to that server, you will get a message on Mattermost (or Slack) or anywhere the webhook points to.
