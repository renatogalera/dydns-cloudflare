## DynDNS Cloudflare update

Change automatically your Dynamic IP in a DNS hosted on Cloudflare. 

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
