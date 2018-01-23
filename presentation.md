# Upspin - Fosdem 2018

## Intro

## The problem

### It's about _your_ data

- Your personal data
- What is _your_ data?
- The photos that...
  - ... you have on your phone?
  - ... you saved on Google photo?
  - ... you liked?
  - ... you saved?
  - ... you saw?
  - ... you created?

### It's about your _power_ (control?)

- You should be able to:
  - access your data from everywhere
  - use them with the app of your choice
  - share them with your friends and family
  - have a good overview of what they are
  - be sure you won't loose them
  - a lot of other things
- We are missing out on multiple points...

### How's that?

- The apps and data come together as a "service"
  - TODO: Insert drawing of user talking to a service
- This means:
  - No control, only what the service is willing to give you
  - Mostly shitty app as they as re-done for every new data

### What about this?

- Some services provide data
- Some other services provide apps which can access the data that
  _you_ can access
  - TODO: insert drawing of these 2 kinds of services, intacting with
    each other
- That looks like working on a _regular_ file system

## Upspin: the idea

### Upspin is...

- Upspin is a protocol for accessing all of your data remotely with
  access control
  - TODO: check the definitions at upspin.io as Rob Pike usually have
    nice chosen words

### Basic usage

- A shared Key server at key.upspin.io
- Your Dir server running somewhere
- Your Store server running somewhere and using one of the available
  storage providers (gcp, s3, dropbox, drive, local filesystem)

- On your local machine, you run `upspinfs` which mount all the files
  you can access
- You use your favorite native app to play with your data

- TODO: insert schema of this

### What I hear

- A lot of different apps that access your data for different precise
  functions (remind of the Unix philosophy)
- Apps chaining each other
  - An app can transform the data by reading it somewhere and serving
    its transformation somewhere else
- External apps that can produce new data or trigger operations

- TODO: insert schema of this

### Example (simple): photos

- Synchronization on your phone: write-only
- Photo gallery: read-only
- Sharing organizer: read & write (copy to shared folder)
- Photo editor: read & write
- ...

### Example (hardcore): pacemaker

- Your pacemaker is a read & write upspin server!
- A service that you trust can access it directly
- A security app exposes some not-too-dangerous functions so that
  other apps can change some settings
- A read-only monitoring app gives data to its state to you and your
  doctor

## Upspin: the reality

### The protocol

- Upspin gives the basic data access operation:
  - Open, Create, Glob (list), Lookup, Put, Delete
- Looks like:
  - REST, plus some others
  - filesystem

### 3 kinds of servers

- Key (only one: key.upspin.io)
- Dir
- Store

- TODO: insert schema with 1 example operation of each server

### The KeyServer interface

```
// The KeyServer interface provides access to public information about users.
type KeyServer interface {
	Lookup(userName UserName) (*User, error)
	Put(user *User) error
}
```

(simplified from upspin.io/upspin/upspin.go)

### The DirServer interface

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

(simplified from upspin.io/upspin/upspin.go)

### The StoreServer interface

```
// The StoreServer saves and retrieves data without interpretation.
type StoreServer interface {
	Get(ref Reference) ([]byte, *Refdata, []Location, error)
	Put(data []byte) (*Refdata, error)
	Delete(ref Reference) error
}
```

(simplified from upspin.io/upspin/upspin.go)

### Access control

- Plus, Upspin gives access control over classes of operations:
  - For files: Read, Write
  - For directories: Create, List, Delete (items from the directory)
  - No concept of execute
- At directory level (for now?)
- They are defined in special `Access` files with a nice synthax
- Give rights to you, users or groups...

### 2 kinds of access control

1. Who can operate with the Dir (and Store) server of the files
2. Who can decrypt the files i.e. get the decryption key

### Encryption

- End-to-end encryption
- Nice exotic algorithms
- Create and store symetric key encrypted with public key
- To share, create a copy of the symetric key, encrypted with the
  target user private key
- TODO: check if it's 1 key per directory; where is the key actually
  stored (dir or store?)

### Authorization

- All methods are _authenticated_ except KeyServer.Lookup
  - Need to be a legitimate Upspin user
  - Doesn't mean your rights are checked
- Real-time
- The server performs the authentication of the user based on
  signatures and tokens.
- For operations like DirServer.Lookup, the rights are checked
- In standard implementation, StoreServer.Get doesn't check the rights
  of the user
  - The encryption of the file is sufficient
  - You'd better not use the `plain` encryption then...
