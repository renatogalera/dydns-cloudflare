## DynDNS Cloudflare update

Auto update your Dynamic IP in Cloudflare DNS.

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

docker build -t dyndns-cf-go .

docker run -d dyndns-cf-go
```
