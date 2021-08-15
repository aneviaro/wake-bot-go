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

## Defect Management

We use Trello board as a public issues tracker. To leave an issue to the project, please open a new github issue, we will migrate it to trello, then you will be able to track a progress on it.

https://trello.com/b/iycB60bT/wakebot

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
