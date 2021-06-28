# wake-bot-go

This is a Telegram bot that helps to calculate sleep phases. 

Available at: https://t.me/wake_me_bot

## Run locally

To run this bot locally, you should get your own BOT_TOKEN, using BotFather - https://t.me/botfather
Then build docker image by running: 

```
docker build . -t wake-bot
```

Run docker container via: 

```
docker run --env BOT_TOKEN=<your bot token goes there> wake-bot
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
