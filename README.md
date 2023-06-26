# Mattermost Bot Webhook

This plugin registers a configurable webhook for a specific Bot ID. The webhook contains the post made to the bot, other messages or messaged by the bot are ignored.

## Installation

Clone this repository and build the plugin using `make dist`. After that upload the plugin to your Mattermost instance and enable it.

If the upload fails make sure that Mattermost is configured to handle large file uploads and your webserver also allows uploading files of the same size.
