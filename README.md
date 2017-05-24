```

```
go run commits.go types.go
```
docker run -it --rm -p "3000:3000" -e TG_TOKEN=xxxxxxxx -e CHAT_ID=xxxxxxx --name  commits-tmp  targence/commits-telegram-bot go run commits.go types.go
```

```
## TODO
- add github support
- add tag_push event
- add https and secret token support
```