# Mattermost Bot Webhook

This plugin provides the possibility to configure a webhook to call when a specific Bot ID is messaged.  
The webhook contains the post made to the bot.  
Other messages or messages by the bot itself are ignored.

## Installation

Download the most current release from the [Releases](https://github.com/timkley/mattermost-plugin-bot-webhook/releases) page.

> [!IMPORTANT]
> If the upload fails make sure that Mattermost is configured to handle large file uploads and your webserver also allows uploading files of the same size.
> 
> **Mattermost:** Search the System Console for "Maximum File Size".  
> **Nginx:** The server config can be found at `/etc/nginx/conf.d/mattermost.conf`. Edit the setting `client_max_body_size` under `location /`.

## Manual Build / Installation

Clone this repository and build the plugin using `make dist`. The bundled plugin can be found in the `dist` folder.

## Configuration

After successful installation you may configure the Bot Webhook from the System Console → Plugins → Bot Webhook.
