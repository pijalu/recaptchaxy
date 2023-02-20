# recaptchaxy

A basic golang proxy with added recaptcha enterprise calls
If will perform a call to recaptcha before proxying - if the recaptcha fails, an error will be returned

It expect, in the request to the proxy the following header:
* x-recaptcha-action: the action (not used)
* x-recaptcha-site : the site key
* x-recaptcha-token : the recaptcha token

the proxy expect the following env variable:
* RC_LISTEN : the proxy server bind address
* RC_TARGET: the proxied end point (protocol://host)
* RC_PROJECT_ID : recaptcha enterpise project id
* RC_APIKEY : google recaptcha enterprise API key
* RC_MIN_SCORE: Min score to accept the request , 0.5 as default