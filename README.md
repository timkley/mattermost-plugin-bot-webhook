# Mattermost Bot Webhook

This plugin provides the possibility to configure a webhook to call when a specific Bot ID is messaged.  
The webhook contains the post made to the bot.  
Other messages or messages by the bot itself are ignored.

## Installation

Clone this repository and build the plugin using `make dist`. After that upload the plugin to your Mattermost instance and enable it.

If the upload fails make sure that Mattermost is configured to handle large file uploads and your webserver also allows uploading files of the same size.

## Configuration

After successful installation you may configure the Bot Webhook from the System Console → Plugins → Bot Webhook.
