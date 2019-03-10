## DynDNS Cloudflare update

Altere seu IP Dinâmico automaticamente em uma entrada DNS de um domínio hospedado na Cloudflare.

Forneça dados de sua conta cloudflare no arquivo **config.env**

Exemplo de configuração utilizando subdomínio **home.example.com**

```
CF_API_KEY=SUAAPIKEY
CF_API_EMAIL=SEUEMAIL
SUBDOMAIN=home
DOMAIN=example.com
```

- Instalação

```
git clone https://github.com/renatoguilhermini/dydns-cloudflare-update-go

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
