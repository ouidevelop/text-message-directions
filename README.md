This is a service that provides text message based info for people who don't have a smart phone or don't have data. I know this works for USA and Canada.

You have to send a text message to 1 (855) 651-4636 with the right format depending on what you want

# Directions
Use the form `[drive or walk or bike or transit] from [origin] to [destination]`.

For example:

`drive from UC berkeley to Oakland airport`

`transit from los angeles airport to santa monica`

`bike from santa monica to venice beach`

# Contact info
Use the form `info [place]` to get the phone number and address of a business.

For example:

`info uc berkeley`

Use the form `find [type of place] near [specific place]` to get up to 5 nearby places.

For example:

`find grocery near uc berkeley`

# Weather
Use the form `weather [zip or postal code]`

For example:

`weather 10017`

# Definitions
To define a word, use the form `define [word]`.

For example:

`define car`

# Payment 

This project costs me about 10 cents per message (on average). As I'm not trying to make a profit, that is what I charge. Here's how to get set up:
1. Buy as many messages as you want. 100 messages would be 10 dollars, 50 would be 5 dollars etc.
[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/cgi-bin/webscr?cmd=_donations&business=YLQFA7GD6GZYG&lc=US&currency_code=USD&bn=PP%2dDonationsBF%3abtn_donate_LG%2egif%3aNonHosted)
2. Send me an email at ouidevelop@gmail.com telling me your phone number and how much you paid and I'll add your messages to your number.

# Cost breakdown
I use twilio for the messaging. They charge .75 cents per segment (which I believe is 160 characters or so). On average, requests have enough segments to cost me about 4 cents. That includes both the incoming message and outgoing message.

The business information and directions come from google api's which also charge. In my last month of providing this for free, that cost was about the same as the messaging fees. So together that's up to about 8 cents per request. And then there's paying for the server and two phone numbers (one for production and one for testing).

Built by Oui Develop (at http://ouidevelop.org)
