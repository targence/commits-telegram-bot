Bot sends an every git push event from Github/Gitlab to Telegram chat.  

You need setup a webhook url in Github/Gitlab.   
For example `http://bot.your-site,com/github` and `http://bot.your-site,com/gitlab`   


```
## Gitlab setup for each your project
https://gitlab.com/your-company/your-project/settings/integrations

## Gitlab setup for all your company projects
https://github.com/organizations/your-company/settings/hooks
```

```
TG_TOKEN=xxxxxxxx CHAT_ID=xxxxxxxx go run *.go
```

```
docker run -d --restart=always -p "80:3000" -e TG_TOKEN=xxxxxxxx -e CHAT_ID=xxxxxxxx --name  commits targence/commits-telegram-bot

```

```
## TODO
- add tag_push event
- add https and and secret token support
- bypass a non-commits requests (spam, ect)
```

```
## NOTES
- Docker files I use mainly for testing and development, not for production
```