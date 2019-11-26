## DynDNS Cloudflare update

Auto update your Dynamic IP in Cloudflare DNS.

Provide your cloudflare account data in the file **config.env**

Add your subdmain **home.example.com**

- Docker Image run in background

```
git clone https://github.com/renatogalera/dydns-cloudflare 

cd dydns-cloudflare

#First edit/create config.env 

cp config.env.example config.env

vim config.env

docker build -t dyndns-cf-go .

docker run -dit --restart unless-stopped --name dyndns-cf-go dyndns-cf-go

#Check logs

docker logs -f dyndns-cf-go
```

- Install local

```
git clone https://github.com/renatogalera/dydns-cloudflare

cd dydns-cloudflare

#First edit/create config.env

cp config.env.example config.env

vim config.env

go mod download

go build main.go

./main
```


