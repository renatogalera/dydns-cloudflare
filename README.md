## DynDNS Cloudflare update

Change your Dynamic IP automatically in a DNS entry for a domain hosted on Cloudflare.

Provide your cloudflare account data in the file **config.env**

Add your subdmain **home.example.com**

- Install

```
go mod download

go build main.go
```

- Docker Image run in background

```
git clone https://github.com/renatogalera/dydns-cloudflare-update-go 

cd dydns-cloudflare-update-go

#First create edit/create conf.env first

cp config.env.example config.env

vim config.env

docker build -t go-docker .

docker run -d go-docker
```
