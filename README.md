## DynDNS Cloudflare update

Change your Dynamic IP automatically in a DNS entry for a domain hosted on Cloudflare.

Provide your cloudflare account data in the file **config.env**

Add your subdomain **home.example.com**

```
CF_API_KEY=SUAAPIKEY
CF_API_EMAIL=SEUEMAIL
DOMAIN=example.com
```

- Install

```
go get github.com/manifoldco/promptui

go get github.com/cloudflare/cloudflare-go

go get github.com/joho/godotenv

go get github.com/GlenDC/go-external-ip

go build dyndns-cloudflare.go

chmod +x dyndns-cloudflare
```

In linux, create task in crontab. Note: Add your correct directory.

```
crontab -l | { cat; echo "*/3 * * * * dir/dydns-cloudflare-update-go/dyndns-cloudflare"; } | crontab -
```
