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

List of authenticated and unauthenticated methods:

```
KeyServer:
~~~~~~~~~~
		Methods: map[string]rpc.Method{
			"Put": s.Put,
		},
		UnauthenticatedMethods: map[string]rpc.UnauthenticatedMethod{
			"Lookup": s.Lookup,
		},

DirServer:
~~~~~~~~~~
		Methods: map[string]rpc.Method{
			"Delete":      s.Delete,
			"Glob":        s.Glob,
			"Lookup":      s.Lookup,
			"Put":         s.Put,
			"WhichAccess": s.WhichAccess,
		},
        // No UnauthenticatedMethod

StoreServer:
~~~~~~~~~~~~
		Methods: map[string]rpc.Method{
			"Get":    s.Get,
			"Put":    s.Put,
			"Delete": s.Delete,
		},
        // No UnauthenticatedMethod
```

For these authenticated methods, the Dial method is used to create a
server dedicated to answering that user. With this mechanism we know
for sure that the server that is returned by Dial will only be
accessible from the legitimate user. Or is it? (TODO: check that)
