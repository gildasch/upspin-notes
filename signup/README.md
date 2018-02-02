## Download the tools

- Go to [homepage](https://upspin.io/)

![homepage](01-homepage.png)

- Go to [download the Upspin client](https://upspin.io/dl/)

![download](02-download.png)

- `wget https://upspin.io/dl/upspin.linux_amd64.tar.gz`
- `tar -xf upspin.linux_amd64.tar.gz`

```
$ ls
cacheserver  README  upspin  upspin-audit  upspinfs  upspin.linux_amd64.tar.gz  upspin-ui
```

## Sign-up

- Go to [Signing up a new user](https://upspin.io/doc/signup.md)
- Binaries for Linux, MacOS and Windows desktops
- For others platforms, build it with:

```
$ go get upspin.io/cmd/...
$ go get augie.upspin.io/cmd/upspin-ui
```

- Start `./upspin-ui`
- Enter your email, here `upspin-fosdem@yopmail.com` (you need to be able to access it at least once)

- Write down your secret keys `pafik-susom-bajom-samum.josib-hituf-pakup-gumus`
- The key pair is stored at `$HOME/.ssh/upspin-fosdem@yopmail.com`

```
$ ls ~/.ssh/upspin-fosdem@yopmail.com
public.upspinkey  secret.upspinkey
```

- Verify your e-mail
- You can see it in the [key server logs](https://key.upspin.io/log)

```
2018-02-02 13:46:59.731832531 +0000 UTC: put attempt by "upspin-fosdem@yopmail.com": {"Name":"upspin-fosdem@yopmail.com","Dirs":null,"Stores":null,"PublicKey":"p256\n22684330827047887910595805123591373621548208819597934607991943487834538723113\n24420862709676394012547165196997970196674245372104689073549993353531700757159\n"}
SHA256:e3aa62f0e53e8956e7176a50ab1e572298dd0c58e1d9f7b740d00a7ca92844a4
2018-02-02 13:47:00.828355484 +0000 UTC: put success by "upspin-fosdem@yopmail.com": {"Name":"upspin-fosdem@yopmail.com","Dirs":null,"Stores":null,"PublicKey":"p256\n22684330827047887910595805123591373621548208819597934607991943487834538723113\n24420862709676394012547165196997970196674245372104689073549993353531700757159\n"}
SHA256:741474f8e4e4f8c1036d4dcc3024806177ff88159a10c0ec983b22cab8666fbd
```
