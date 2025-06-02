## go-fitbit-get

The Fitbit API is pretty simple to use.  You can even do it from a shell
script using `curl`

eg

`https://api.fitbit.com/1/user/-/activities/heart/date/2025-06-02/2025-06-02/1sec/time/00:00/23:59.json` will give your heart rate data as fine grained as
possible for a whole day as a JSON file.  We can then use `jq` on that.

The problem is the authentication is done via Oauth2, and that makes things
harder.

So what this program does is provide a wrapper around the URL to get and
manage the oauth2 token.  It will refresh it as needed and save the
access/refresh tokens in the filesystem (eg in `$HOME/.fitbit_get`)

*WARNING: these tokens won't be protected so if anyone else can read your directory then they could steal your token and access your data*


## Register an APP wit Fitbit

To access the API you need to create an application at https://dev.fitbit.com/

The following fields need to be filled out

| Field | What |
| - | - |
| Name | Put whatever you like |
| Description | Put whatever you like |
| App URL | I use http://localhost - it's just got to look right |
| Organisation | Put whatever you like |
| Org URL | Again I use http://localhost |
| ToS URL | And again http://localhost |
| Privacy URL | And again http://localhost |
| Type | Personal |
| Redirect URL | http://localhost:16601 |
| Access | Read Only |

The Type and Redirect URL are important, the rest just have to be there
to keep the fitbit registration happy.

This will create a Client ID and a Client Secret.  You need these values.


## Configure this program

Now run this program, it will moan that the config file is missing; e.g

```
Error parsing /home/sweh/.fitbit_get
   open /home/sweh/.fitbit_get: no such file or directory

```

You need to create this file.  It's in JSON format and needs the two
fields defined:

```
{
        "ClientID": "111111",
        "ClientSecret": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}
````

Obviously use your values!

## Get your access tokens.

Rerun the program.  This will spin up a temporary server on localhost:16601
and then will tell you to go to a URL

```
% fitbit_get
No access token found so we need to authenticate this application

Go to the following URL and authorize the application:

https://www.fitbit.com/oauth2/authorize?access_type=offline&client_id=....
```

Copy/Paste that into your browser

You may be asked to login to fitbit again.  Select all the scope values
and accept everything.  This should automatically redirect back to
the temporary server.  At that point you should see `Saving new config`
and if you look in the config file you'll see a lot more data:

```
{
        "ClientID": "111111",
        "ClientSecret": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
        "Port": 16601,
        "AccessToken": "eyJh....",
        "RefreshToken": "108...",
        "Expiry": "2025-06-03T01:19:36.890001256-04:00"
}
```

And now we can do what we really wanted:
```
% fitbit_get https://api.fitbit.com/1/user/-/activities/heart/date/2025-06-02/2025-06-02/1sec/time/00:00/23:59.json
{
  "activities-heart": [
    {
      "customHeartRateZones": [],
      "heartRateZones": [
        {
...
```

## The redirect failed; eg the browser said "can't establish a connection".

This can happen if your webbrowser runs on a different machine to where
you're running this program.  In that case it should be simple enough to
copy the URL and then "curl" it from the machine

```
% curl 'http://localhost:16601/?code=d3axxx&state=7194a0xxx#_=_'
```

## I don't like port 16601

Just add a `Port` setting to the config to pick the port you want.  Remember
to make the RedirectURL match.

## My token no longer works

Remove the `AccessToken` from the config file and the code will go through
the authentication process again.

