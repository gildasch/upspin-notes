The officiel RPC client will make authenticated requests for every
method.

There are two kinds of authentifications:

1. The first authentication:

```
-> Upspin-Auth-Request: gildaschbt@gmail.com,upspin.gildas.ch:443,Sat Jan  6 16:08:41 2018,XXXXXXXXXXXXXXX,XXXXXXXXXXXXXXX

<- Upspin-Auth-Token:[13DA28092C3224920B29CA9834811407]
```

2. The following authentications with the auth token:

```
-> Upspin-Auth-Token:[13DA28092C3224920B29CA9834811407]
```

About access right, the official directory server will check the
Access right of the client and return a `withheld information` error
if he does not have the right to perform the requested action. The
official store server however doesn't check any access right and will
serve the stored bytes to anybody with a valid token. The security is
provided here by the fact that the data that is served by the store is
encrypted.
