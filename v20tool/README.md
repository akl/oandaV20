To run this tool, create a INI file with the following
sections and fields:

    [oanda]
    account =
    access_token =
    env = fxpractice
    streaming = false


The account is the numeric decimal ID for the account in question.
The instruments/candles API doesn't actually use 'account', so it's
not imperative that you set this correctly.

The access token is the long hexadecimal string that is your REST API token
from OANDA.

"env" should be either fxpractice or fxtrade, depending on the OANDA
environment you are using.

"streaming" indicates whether to use the streaming API or not.
