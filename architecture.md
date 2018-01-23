# The Architecture of Upspin

## 3 interfaces

```
// The KeyServer interface provides access to public information about users.
type KeyServer interface {
	Lookup(userName UserName) (*User, error)
	Put(user *User) error
}
```

```
// DirServer manages the name space for one or more users.
type DirServer interface {
	Lookup(name PathName) (*DirEntry, error)
	Put(entry *DirEntry) (*DirEntry, error)
	Glob(pattern string) ([]*DirEntry, error)
	Delete(name PathName) (*DirEntry, error)
	WhichAccess(name PathName) (*DirEntry, error)
	Watch(name PathName, sequence int64, done <-chan struct{}) (<-chan Event, error)
}
```

```
// The StoreServer saves and retrieves data without interpretation.
type StoreServer interface {
	Get(ref Reference) ([]byte, *Refdata, []Location, error)
	Put(data []byte) (*Refdata, error)
	Delete(ref Reference) error
}
```

## All authenticated methods

```
	rpc.NewServer(cfg, rpc.Service{
		Name: "Key",
		Methods: map[string]rpc.Method{
			"Put": s.Put,
		},
		UnauthenticatedMethods: map[string]rpc.UnauthenticatedMethod{
			"Lookup": s.Lookup,
		},
		Lookup: func(userName upspin.UserName) (upspin.PublicKey, error) {
			user, err := key.Lookup(userName)
			if err != nil {
				return "", err
			}
			return user.PublicKey, nil
		},
	})
```

```
	rpc.NewServer(cfg, rpc.Service{
		Name: "Dir",
		Methods: map[string]rpc.Method{
			"Delete":      s.Delete,
			"Glob":        s.Glob,
			"Lookup":      s.Lookup,
			"Put":         s.Put,
			"WhichAccess": s.WhichAccess,
		},
		Streams: map[string]rpc.Stream{
			"Watch": s.Watch,
		},
	})
```

```
	rpc.NewServer(cfg, rpc.Service{
		Name: "Store",
		Methods: map[string]rpc.Method{
			"Get":    s.Get,
			"Put":    s.Put,
			"Delete": s.Delete,
		},
	})
```

## Authentication check

### Doc from the rpc source

```
Authentication

The client authenticates itself to the server using special HTTP headers.

In its first request to a given Upspin server, the client presents a signed
authentication request as a series of HTTP headers with key
'Upspin-Auth-Request'. The header values are, in order:
       the user name,
       the server host name,
       the current time, and
       the R part of an upspin.Signature,
       the S part of an upspin.Signature.
The current time is formatted using Go's time.ANSIC presentation:
	"Mon Jan _2 15:04:05 2006"
To generate the upspin.Signature, the client concatenates the user name, host
name, and formatted time, each string prefixed by its length as a big
endian-encoded uint32:
	[len(username)][username][len(hostname)][hostname][len(time)][time]
The client then SHA256-hashes that string and signs it using its Factotum.

The server checks the signature and, if valid, returns an authentication token
that represents the current session in the 'Upspin-Auth-Token' response header.

In subsequent requests, the client presents that authentication token to the
server using the 'Upspin-Auth-Token' request header.

If there is an error validating an authentication request or token, the server
returns an error message in the 'Upspin-Auth-Error' response header.

TODO: document the 'Upspin-Proxy-Request' header.
```

### rpc.serverImpl.ServeHTTP (simplified)

```
func (s *serverImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    method := d.Methods[name]
	umethod := d.UnauthenticatedMethods[name]

	var session Session
	if umethod == nil {
		var err error
		session, err = s.SessionForRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}
}
```

### rpc.serverImpl.SessionForRequest (simplified)

```
func (s *serverImpl) SessionForRequest(w http.ResponseWriter, r *http.Request) (session Session, err error) {
	if tok, ok := r.Header[authTokenHeader]; ok && len(tok) == 1 {
		return s.validateToken(tok[0])
	}

    authRequest, ok := r.Header[authRequestHeader]

	return s.handleSessionRequest(w, authRequest, proxyRequest, r.Host)
}
```

### rpc.serverImpl.validateToken  (simplified)

```
func (s *serverImpl) validateToken(authToken string) (Session, error) {
	session := GetSession(authToken)
	if session == nil {
		return nil, errors.E(errors.Permission, errUnauthenticated)
    }

	if session.Expires().Before(time.Now()) {
		return nil, errors.E(errors.Permission, errExpired)
	}

	return session, nil
}
```

### rpc.serverImpl.handleSessionRequest (simplified)

```
func (s *serverImpl) handleSessionRequest(w http.ResponseWriter, authRequest []string, proxyRequest []string, host string) (Session, error) {
	// Validate the username.
	user := upspin.UserName(authRequest[0])
	if err := valid.UserName(user); err != nil {
		return nil, errors.E(user, err)
	}

	// Get user's public key.
	key, err := s.lookup(user)
	if err != nil {
		return nil, errors.E(user, err)
	}

	// If this is a proxy request, extract the endpoint and
	// set the signed host to that endpoint.
	ep := &upspin.Endpoint{}
	if len(proxyRequest) == 1 {
		if pUser := s.config.UserName(); user != pUser {
			return nil, errors.E(errors.Permission, errors.Errorf("client %q and proxy %q users mismatched", user, pUser))
		}
		ep, err = upspin.ParseEndpoint(proxyRequest[0])
		if err != nil {
			return nil, errors.E(errors.Invalid, errors.Errorf("invalid proxy endpoint: %v", err))
		}
		host = string(ep.NetAddr)
	}

	now := time.Now()

	// Validate signature.
	if err := verifyUser(key, authRequest, clientAuthMagic, host, now); err != nil {
		return nil, errors.E(errors.Permission, user, errors.Errorf("invalid signature: %v", err))
	}

	// Generate an auth token and bind it to a session for the client.
	expiration := now.Add(authTokenDuration)
	authToken, err := generateRandomToken()
	if err != nil {
		return nil, err
	}
	w.Header().Set(authTokenHeader, authToken)

	// If there is a proxy request, authenticate server to client.
	if len(proxyRequest) == 1 {
		// Authenticate the server to the user.
		authMsg, err := signUser(s.config, serverAuthMagic, "[localproxy]")
		if err != nil {
			return nil, errors.E(errors.Permission, err)
		}
		w.Header()[authRequestHeader] = authMsg
	}

	return NewSession(user, expiration, authToken, ep, nil), nil
}
```

### dirserver.Lookup (simplified)

```
func (s *server) Lookup(session rpc.Session, reqBytes []byte) (pb.Message, error) {
	dir, err := s.serverFor(session, reqBytes, &req)
	return op.entryError(dir.Lookup(upspin.PathName(req.Name)))
}
```

### dir/server.Lookup

```
func (s *server) Lookup(name upspin.PathName) (*upspin.DirEntry, error) {
	const op errors.Op = "dir/server.Lookup"
	o, m := newOptMetric(op)
	defer m.Done()
	return s.lookupWithPermissions(op, name, o)
}

func (s *server) lookupWithPermissions(op errors.Op, name upspin.PathName, opts ...options) (*upspin.DirEntry, error) {
	p, err := path.Parse(name)
	if err != nil {
		return nil, errors.E(op, name, err)
	}

	entry, err := s.lookup(p, entryMustBeClean, opts...)

	// Check if the user can know about the file at all. If not, to prevent
	// leaking its existence, return Private.
	if err == upspin.ErrFollowLink {
		return s.errLink(op, entry, opts...)
	}
	if err != nil {
		if errors.Is(errors.NotExist, err) {
			if canAny, _, err := s.hasRight(access.AnyRight, p, opts...); err != nil {
				return nil, errors.E(op, err)
			} else if !canAny {
				return nil, errors.E(op, name, errors.Private)
			}
		}
		return nil, errors.E(op, err)
	}

	// Check for Read access permission.
	canRead, _, err := s.hasRight(access.Read, p, opts...)
	if err == upspin.ErrFollowLink {
		return nil, errors.E(op, errors.Internal, p.Path(), "can't be link at this point")
	}
	if err != nil {
		return nil, errors.E(op, err)
	}
	if !canRead {
		canAny, _, err := s.hasRight(access.AnyRight, p, opts...)
		if err != nil {
			return nil, errors.E(op, err)
		}
		if !canAny {
			return nil, s.errPerm(op, p, opts...)
		}
		if !access.IsAccessControlFile(entry.SignedName) {
			entry.MarkIncomplete()
		}
	}
	return entry, nil
}
```

## Annexes

### Complete interfaces

```
// The KeyServer interface provides access to public information about users.
type KeyServer interface {
	Dialer
	Service

	// Lookup returns all public information about a user.
	Lookup(userName UserName) (*User, error)

	// Put sets or updates information about a user. The user's name must
	// match the authenticated user. The call can update any field except
	// the user name.
	// To add new users, see the signup subcommand of cmd/upspin.
	Put(user *User) error
}
```

```
// DirServer manages the name space for one or more users.
type DirServer interface {
	Dialer
	Service

	// Lookup returns the directory entry for the named file.
	//
	// If the returned error is ErrFollowLink, the caller should
	// retry the operation as outlined in the description for
	// ErrFollowLink. Otherwise in the case of error the
	// returned DirEntry will be nil.
	Lookup(name PathName) (*DirEntry, error)

	// Put stores the DirEntry in the directory server. The entry
	// may be a plain file, a link, or a directory. (Only one of
	// these attributes may be set.)
	// In practice the data for the file should be stored in
	// a StoreServer as specified by the blocks in the entry,
	// all of which should be stored with the same packing.
	//
	// Within the DirEntry, several fields have special properties.
	// Time represents a timestamp for the item. It is advisory only
	// but is included in the packing signature and so should usually
	// be set to a non-zero value.
	//
	// Sequence represents a sequence number that is incremented
	// after each Put. If it is neither 0 nor -1, the DirServer will
	// reject the Put operation if the file does not exist or, for an
	// existing item, if the Sequence is not the same as that
	// stored in the metadata. If it is -1, Put will fail if there
	// is already an item with that name.
	//
	// The Name field of the DirEntry identifies where in the directory
	// tree the entry belongs. The SignedName field, which usually has the
	// same value, is the name used to sign the DirEntry to guarantee its
	// security. They may differ if an entry appears in multiple locations,
	// such as in its original location plus within a second tree holding
	// a snapshot of the original tree but starting from a different root.
	//
	// Most software will concern itself only with the Name field unless
	// generating or validating the entry's signature.
	//
	// All but the last element of the path name must already exist
	// and be directories or links. The final element, if it exists,
	// must not be a directory. If something is already stored under
	// the path, the new location and packdata replace the old.
	//
	// If the returned error is ErrFollowLink, the caller should
	// retry the operation as outlined in the description for
	// ErrFollowLink (with the added step of updating the
	// Name field of the argument DirEntry). For any other error,
	// the return DirEntry will be nil.
	//
	// A successful Put returns an incomplete DirEntry (see the
	// description of AttrIncomplete) containing only the
	// new sequence number.
	Put(entry *DirEntry) (*DirEntry, error)

	// Glob matches the pattern against the file names of the full
	// rooted tree. That is, the pattern must look like a full path
	// name, but elements of the path may contain metacharacters.
	// Matching is done using Go's path.Match elementwise. The user
	// name must be present in the pattern and is treated as a literal
	// even if it contains metacharacters.
	// If the caller has no read permission for the items named in the
	// DirEntries, the returned Location and Packdata fields are cleared.
	//
	// If the returned error is ErrFollowLink, one or more of the
	// returned DirEntries is a link (the others are completely
	// evaluated). The caller should retry the operation for those
	// DirEntries as outlined in the description for ErrFollowLink,
	// updating the pattern as appropriate. Note that any returned
	// links may only partially match the original argument pattern.
	//
	// If the pattern evaluates to one or more name that identifies
	// a link, the DirEntry for the link is returned, not the target.
	// This is analogous to passing false as the second argument
	// to Client.Lookup.
	Glob(pattern string) ([]*DirEntry, error)

	// Delete deletes the DirEntry for a name from the directory service.
	// It does not delete the data it references; use StoreServer.Delete
	// for that. If the name identifies a link, Delete will delete the
	// link itself, not its target.
	//
	// If the returned error is ErrFollowLink, the caller should
	// retry the operation as outlined in the description for
	// ErrFollowLink. (And in that case, the DirEntry will never
	// represent the full path name of the argument.) Otherwise, the
	// returned DirEntry will be nil whether the operation succeeded
	// or not.
	Delete(name PathName) (*DirEntry, error)

	// WhichAccess returns the DirEntry of the Access file that is
	// responsible for the access rights defined for the named item.
	// WhichAccess requires that the calling user have at least one access
	// right granted for the argument name. If not, WhichAccess will return
	// a "does not exist" error, even if the item and/or the Access file
	// exist.
	//
	// If the returned error is ErrFollowLink, the caller should
	// retry the operation as outlined in the description for
	// ErrFollowLink. Otherwise, in the case of error the returned
	// DirEntry will be nil.
	WhichAccess(name PathName) (*DirEntry, error)

	// Watch returns a channel of Events that describe operations that
	// affect the specified path and any of its descendants, beginning
	// at the specified sequence number for the corresponding user root.
	//
	// If sequence is 0, all events known to the DirServer are sent.
	//
	// If sequence is WatchCurrent, the server first sends a sequence
	// of events describing the entire tree rooted at name. The Events are
	// sent in sequence such that a directory is sent before its contents.
	// After the full tree has been sent, the operation proceeds as normal.
	//
	// If sequence is WatchNew, the server sends only new events.
	//
	// If the sequence is otherwise invalid, this is reported by the
	// server sending a single event with a non-nil Error field with
	// Kind=errors.Invalid. The events channel is then closed.
	//
	// When the provided done channel is closed the event channel
	// is closed by the server.
	//
	// To receive an event for a given path under name, the caller must have
	// one or more of the Upspin access rights to that path. Events for
	// which the caller does not have enough rights to watch will not be
	// sent. If the caller has rights but not Read, the entry will be
	// present but incomplete (see the description of AttrIncomplete). If
	// the name does not exist, Watch will succeed and report events if and
	// when it is created.
	//
	// If the caller does not consume events in a timely fashion
	// the server will close the event channel.
	//
	// If this server does not support this method it returns
	// ErrNotSupported.
	//
	// The only errors returned by the Watch method itself are
	// to report that the name is invalid or refers to a non-existent
	// root, or that the operation is not supported.
	Watch(name PathName, sequence int64, done <-chan struct{}) (<-chan Event, error)
}
```

// The StoreServer saves and retrieves data without interpretation.
type StoreServer interface {
	Dialer
	Service

	// Get attempts to retrieve the data identified by the reference.
	// Three things might happen:
	// 1. The data is in this StoreServer. It is returned. The Location slice
	// and error are nil. Refdata contains information about the data.
	// 2. The data is not in this StoreServer, but may be in one or more
	// other locations known to the store. The slice of Locations
	// is returned. The data, Refdata, Locations, and error are nil.
	// 3. An error occurs. The data, Locations and Refdata are nil
	// and the error describes the problem.
	Get(ref Reference) ([]byte, *Refdata, []Location, error)

	// Put puts the data into the store and returns the reference
	// to be used to retrieve it.
	Put(data []byte) (*Refdata, error)

	// Delete permanently removes all storage space associated
	// with the reference. After a successful Delete, calls to Get with the
	// same reference will fail. If the reference is not found, an error is
	// returned. Implementations may disable this method except for
	// privileged users.
	Delete(ref Reference) error
}
```
