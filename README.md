## DynDNS Cloudflare update

Utilize IP Dinâmico e atualize automaticamente entrada DNS de um domínio hospedado na Cloudflare.

- Instalação

```
git clone https://github.com/renatoguilhermini/dydns-cloudflare-update-go
```

Altere config.env

Exemplo subdomínio Cloudflare home.example.com

```
CF_API_KEY=SUAAPIKEY
CF_API_EMAIL=SEUEMAIL
SUBDOMAIN=home
DOMAIN=example.com
```

Com Go instalado

```
cd dydns-cloudflare-update-go

go get github.com/cloudflare/cloudflare-go

go get github.com/joho/godotenv

go build dyndns-cloudflare.go

chmod +x dyndns-cloudflare
```

Em linux, criar tarefa no crontab. Obs: Não esqueça de apontar para diretório correto

```
crontab -l | { cat; echo "*/3 * * * * dir/dydns-cloudflare-update-go/dyndns-cloudflare"; } | crontab -
```
