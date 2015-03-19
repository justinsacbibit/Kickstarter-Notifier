# Kickstarter-Notifier
Interested in a Kickstarter project's rewards, but too indecisive to take the plunge? As you think about whether or not you'd like to back it, you can be notified of how many of the reward is left as others purchase it.

Running this project requires a few environment variables to be set up:
``TWILIO_SID````TWILIO_TOKEN````TWILIO_NUMBER````SEND_TO_NUMBER````PORT``

``TWILIO_SID`` and ``TWILIO_TOKEN`` are obtained by creating a Twilio account, which is used to send text messages through API calls.
``TWILIO_NUMBER`` is your Twilio phone number, and ``SEND_TO_NUMBER`` is your cell phone number.
``PORT`` is used to host the web server.

Example:
```
$ TWILIO_SID=<sid> TWILIO_TOKEN=<token> TWILIO_NUMBER=+11234567890 SEND_TO_NUMBER=+10987654321 PORT=3000 go run main.go
```

Scrapes are triggered by a GET request to /ping. This is due to what seems to be Heroku killing off goroutines after an hour. See https://github.com/thoughtpolice/heroku-ping for a worker that can ping your endpoint every so often.
