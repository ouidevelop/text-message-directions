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

# Cost breakdown
I use twilio for the messaging. They charge .75 cents per segment (which I believe is 160 characters or so). On average, requests have enough segments to cost me about 4 cents. That includes both the incoming message and outgoing message.

The business information and directions come from google api's which also charge. In my last month of providing this for free, that cost was about the same as the messaging fees. So together that's up to about 8 cents per request.

I also don't charge for every request, like when there is a user error or an error in the application.

The payment system (stripe) also takes a cut of each payment.

And then there's paying for the server and two phone numbers (one for production and one for testing).

# Payment
This project costs me about 10 cents per message (on average). As I'm not trying to make a profit, that is what I charge. Here's how to get set up:
1. Buy as many messages as you want (you can adjust the quantity in the link). 100 messages would be 10 dollars, 50 would be 5 dollars etc.
   [PAY HERE](https://buy.stripe.com/4gwbLO8j6dA4g24fZ0)
2. This will save your credit card information so that when you run out of messages, you will be given the option to buy more from your phone.
3. If you have questions about payment, or would like to pay over the phone, please call or text support at (805) 423-4224.

Built by Oui Develop (at http://ouidevelop.org)
