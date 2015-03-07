# Kickstarter-Notifier
Interested in a Kickstarter project's rewards, but too indecisive to take the plunge? As you think about whether or not you'd like to back it, you can be notified of how many of the reward is left as others purchase it.

### Developing
Running this project requires a few environment variables to be set up:
``twilioSid````twilioToken````fromNum````toNum````PORT``

``twilioSid`` and ``twilioToken`` are obtained by creating a Twilio account, which is used to send text messages through API calls.
``fromNum`` is your Twilio phone number, and ``toNum`` is your cell phone number.
``PORT`` is used to host the web server.

Example:
```
$ twilioSid=<sid> twilioToken=<token> fromNum=+11234567890 toNum=+10987654321 PORT=3000 go run main.go
```
