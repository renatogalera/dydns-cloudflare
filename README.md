## DynDNS Cloudflare update

Auto update your Dynamic IP in Cloudflare DNS.

Provide your cloudflare account data in the file **config.env**

Add your subdmain **home.example.com**

- Install local

```
#First edit/create conf.env first

cp config.env.example config.env

vim config.env

go mod download

go build main.go

./main
```

- Docker Image run in background

```
git clone https://github.com/renatogalera/dydns-cloudflare-update-go 

cd dydns-cloudflare-update-go

#First edit/create conf.env first

cp config.env.example config.env

vim config.env

docker build -t dyndns-cf-go .

docker run -d dyndns-cf-go
```
