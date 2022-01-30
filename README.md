# wake-bot-go

This is a Telegram bot that helps to calculate sleep phases. 

Available at: https://t.me/wake_me_bot

## Run locally

To run this bot locally, you should get your own `BOT_TOKEN`, using [BotFather](https://t.me/botfather).

There are two ways to get updates from this bot

- get updates from a chan
- get updates via API call, by setting up a webhook

### Updates from chan

1) Modify `compose.env` file, to store yours valid `BOT_TOKEN`
2) Run the following command to start datastore emulator and bot automatically:
    ```
   docker compose -f chan-compose.yml up --build
   ```
3) Try to use your bot from a telegram client

### Updates from tg Webhook

1) Install [ngrok](https://ngrok.com/download)
2) Set up a tunnel from `:8080` port using the following command:
    ```
    ngrok http 8080
    ```
    As the result you will see something similar to this message:
    ```
    ngrok by @inconshreveable                                                                                                                                   (Ctrl+C to quit)
    
    Session Status                online                                                                                                                                        
    Account                       ____________ (Plan: Free)                                                                                                              
    Update                        update available (version 2.3.40, Ctrl-U to update)                                                                                           
    Version                       2.3.35                                                                                                                                        
    Region                        United States (us)                                                                                                                            
    Web Interface                 http://127.0.0.1:4040                                                                                                                         
    Forwarding                    http://ab88-134-17-149-60.ngrok.io -> http://localhost:8080                                                                                   
    Forwarding                    https://ab88-134-17-149-60.ngrok.io -> http://localhost:8080
    
    Connections                   ttl     opn     rt1     rt5     p50
    ```
3) Copy a forwarding URL from the ngrok result screen. Choose the one, that starts with `https://`.
Paste it into `compose.env` file as a `WEBHOOK_URL`.

4) Modify `compose.env` file, to store yours valid `BOT_TOKEN`.

5) Run the following command to start datastore emulator and bot automatically.
    ```
    docker compose -f webhook-compose.yml up --build
    ```
6) Try to use your bot from a telegram client

## Defect Management

We use Trello board as a public issues tracker. To leave an issue to the project, please open a new github issue, we will migrate it to trello, then you will be able to track a progress on it.

https://trello.com/b/iycB60bT/wakebot

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
