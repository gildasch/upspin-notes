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

- Click on `Proceed`

## Configure your Upspin servers

- 3 choices:
  1. I will use existing Upspin servers.
  1. I will deploy new Upspin servers to the Google Cloud Platform.
  1. Skip configuring my servers; I'll use Upspin in read-only mode for now.

- GCP
  - Server user
  - New keys to `$HOME/.ssh/upspin-fosdem+server@yopmail.com`
  - New secret seed `hfeiw-hejei-nnbds-eesfn.lfmfp-hdgss-gowwp-rnej`
  - Host name: you can leave it blank
    - Generated for you: 8cda9311ce4bed564f1004cf4dd864b7.upspin.services
  - Wait for DNS propagation
  - (Optional) Add other users for root read & write rights
  - Done

- Once it is done:
  - After clicking continue, upload files to your tree by dragging them
    into an upspin-ui navigation panel.
  - Configure a cacheserver to speed up Upspin operations.
  - Try using the upspin command-line tool.
  - If you run Linux or macOS, you can set up upspinfs to browse
    Upspin files as if they were on your local machine.

- In the key log:

```
2018-02-02 14:20:31.118590695 +0000 UTC: put attempt by "upspin-fosdem@yopmail.com": {"Name":"upspin-fosdem+server@yopmail.com","Dirs":["unassigned"],"Stores":["unassigned"],"PublicKey":"p256\n114944207301036971755855933023701478199414145030435701152657383838778212475297\n82252520421664963245805888324805718385377964184829806800713681790767939963428\n"}
SHA256:4ddf22ebf0a32bc9f7b0bf7ed311f5e22f1e551b3734e9af2789bcec00e8f4ec
2018-02-02 14:20:31.891360604 +0000 UTC: put success by "upspin-fosdem@yopmail.com": {"Name":"upspin-fosdem+server@yopmail.com","Dirs":["unassigned"],"Stores":["unassigned"],"PublicKey":"p256\n114944207301036971755855933023701478199414145030435701152657383838778212475297\n82252520421664963245805888324805718385377964184829806800713681790767939963428\n"}
SHA256:7b0fe5aa0d53393da312a680a866784420d41d26a9ebba8d65f67e12a9219a9b
2018-02-02 14:33:25.235354889 +0000 UTC: put attempt by "upspin-fosdem@yopmail.com": {"Name":"upspin-fosdem@yopmail.com","Dirs":["remote,8cda9311ce4bed564f1004cf4dd864b7.upspin.services:443"],"Stores":["remote,8cda9311ce4bed564f1004cf4dd864b7.upspin.services:443"],"PublicKey":"p256\n22684330827047887910595805123591373621548208819597934607991943487834538723113\n24420862709676394012547165196997970196674245372104689073549993353531700757159\n"}
SHA256:a88dff22313f015cd337c72c6b8aabc3af2e3d2ddb177b467c05f7a58a442a86
2018-02-02 14:33:25.979881229 +0000 UTC: put success by "upspin-fosdem@yopmail.com": {"Name":"upspin-fosdem@yopmail.com","Dirs":["remote,8cda9311ce4bed564f1004cf4dd864b7.upspin.services:443"],"Stores":["remote,8cda9311ce4bed564f1004cf4dd864b7.upspin.services:443"],"PublicKey":"p256\n22684330827047887910595805123591373621548208819597934607991943487834538723113\n24420862709676394012547165196997970196674245372104689073549993353531700757159\n"}
SHA256:42dfd2757ed75ec327dbdcc6eeae2d5028c6deacf7f1d6bb6983fe35e8ea111e
2018-02-02 14:33:26.471262025 +0000 UTC: put attempt by "upspin-fosdem@yopmail.com": {"Name":"upspin-fosdem+snapshot@yopmail.com","Dirs":["remote,8cda9311ce4bed564f1004cf4dd864b7.upspin.services:443"],"Stores":["remote,8cda9311ce4bed564f1004cf4dd864b7.upspin.services:443"],"PublicKey":"p256\n22684330827047887910595805123591373621548208819597934607991943487834538723113\n24420862709676394012547165196997970196674245372104689073549993353531700757159\n"}
SHA256:dbc629a9356212caaa71890524db48c195065bd45cd5142152d580429819ae82
2018-02-02 14:33:27.059761056 +0000 UTC: put success by "upspin-fosdem@yopmail.com": {"Name":"upspin-fosdem+snapshot@yopmail.com","Dirs":["remote,8cda9311ce4bed564f1004cf4dd864b7.upspin.services:443"],"Stores":["remote,8cda9311ce4bed564f1004cf4dd864b7.upspin.services:443"],"PublicKey":"p256\n22684330827047887910595805123591373621548208819597934607991943487834538723113\n24420862709676394012547165196997970196674245372104689073549993353531700757159\n"}
SHA256:f00344402a7a424de67427ca998b915cfab991327392aaa1a63127a0d1043b00
2018-02-02 14:33:27.651042608 +0000 UTC: put attempt by "upspin-fosdem@yopmail.com": {"Name":"upspin-fosdem+server@yopmail.com","Dirs":["remote,8cda9311ce4bed564f1004cf4dd864b7.upspin.services:443"],"Stores":["remote,8cda9311ce4bed564f1004cf4dd864b7.upspin.services:443"],"PublicKey":"p256\n114944207301036971755855933023701478199414145030435701152657383838778212475297\n82252520421664963245805888324805718385377964184829806800713681790767939963428\n"}
SHA256:ceeeef15111477b8f282247dc01b3ac0563a9b7327e6f732af0ce4b15e6eacc7
2018-02-02 14:33:28.31606501 +0000 UTC: put success by "upspin-fosdem@yopmail.com": {"Name":"upspin-fosdem+server@yopmail.com","Dirs":["remote,8cda9311ce4bed564f1004cf4dd864b7.upspin.services:443"],"Stores":["remote,8cda9311ce4bed564f1004cf4dd864b7.upspin.services:443"],"PublicKey":"p256\n114944207301036971755855933023701478199414145030435701152657383838778212475297\n82252520421664963245805888324805718385377964184829806800713681790767939963428\n"}
SHA256:fa6b336005b7667bf8f55d1ac584e65431c1da4378cbb6a21aa966419e80d3c0
```
