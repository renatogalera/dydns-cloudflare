package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/joho/godotenv"
)

var DOMAIN string
var CF_API_KEY string
var CF_API_EMAIL string
var SUBDOMAIN string
var NEWIPADDR string

func argParse() error {

	// Checa se arquivo config.env existe
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Coletando variáveis
	DOMAIN = os.Getenv("DOMAIN")
	if DOMAIN == "" {
		msg := fmt.Sprintf("É necessário configurar variável DOMAIN")
		return errors.New(msg)
	}
	CF_API_KEY = os.Getenv("CF_API_KEY")
	if CF_API_KEY == "" {
		msg := fmt.Sprintf("É necessário configurar variável CF_API_KEY")
		return errors.New(msg)
	}
	CF_API_EMAIL = os.Getenv("CF_API_EMAIL")
	if CF_API_EMAIL == "" {
		msg := fmt.Sprintf("É necessário configurar variável CF_API_EMAIL")
		return errors.New(msg)
	}
	SUBDOMAIN = os.Getenv("SUBDOMAIN")
	if SUBDOMAIN == "" {
		msg := fmt.Sprintf("É necessário configurar variável SUBDOMAIN")
		return errors.New(msg)
	}
	NEWIPADDR = os.Getenv("NEWIPADDR")
	if NEWIPADDR == "" {
		msg := fmt.Sprintf("É necessário configurar variável NEWIPADDR")
		return errors.New(msg)
	}
	// Criando arquivo novoip caso não exista
	if _, err := os.Stat(NEWIPADDR); os.IsNotExist(err) {
		os.OpenFile(NEWIPADDR, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	}

	return nil
}

func checkIP() {

	log.Printf("Checando IP...\n")
	IPV4 := getMyIP(4)
	// IPV6 := getMyIP(6)

	// Pesquisa IP no arquivo
	OLDIPTMP, _ := ioutil.ReadFile(NEWIPADDR)
	OLDIP := string(OLDIPTMP)

	if OLDIP != IPV4 {
		log.Printf("IP Mudou! Alterando: %s -> %s", OLDIP, IPV4)
		dynDNS(IPV4)
	} else {
		log.Printf("IP Não mudou!\n")
		os.Exit(0)
	}
}

func dynDNS(IPV4 string) {

	// Removendo referência antiga
	os.Remove(NEWIPADDR)

	//Salvando ipatual no arquivo
	saveip, err := os.OpenFile(NEWIPADDR, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	defer saveip.Close()
	if _, err = saveip.WriteString(IPV4); err != nil {
		panic(err)
	}

	// API Cloudflare
	api, err := cloudflare.New(CF_API_KEY, CF_API_EMAIL)
	if err != nil {
		log.Fatal(err)

	}

	// API Cloudflare
	zoneID, err := api.ZoneIDByName(DOMAIN)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Criando entrada DNS
	newRecord := cloudflare.DNSRecord{
		Type:    "A",
		Name:    SUBDOMAIN + "." + DOMAIN,
		Content: IPV4,
	}

	updateRecord(zoneID, api, &newRecord)
	log.Println("Setando entrada DNS:", newRecord.Name, newRecord.Content, "\n")
}

func updateRecord(zoneID string, api *cloudflare.API, newRecord *cloudflare.DNSRecord) {

	dns := cloudflare.DNSRecord{Type: newRecord.Type, Name: newRecord.Name}
	OLDRECORDS, err := api.DNSRecords(zoneID, dns)
	if err != nil {
		log.Fatal(err)
		return
	}

	if len(OLDRECORDS) == 1 {
		// Update
		err := api.UpdateDNSRecord(zoneID, OLDRECORDS[0].ID, *newRecord)
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	}

	_, err = api.CreateDNSRecord(zoneID, *newRecord)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func getMyIP(protocol int) string {
	var target string
	if protocol == 4 {
		target = "http://ifconfig.me/ip"

	} else if protocol == 6 {
		target = "http://ifconfig.me/ip" //Alterar para fonte ipv6

	} else {
		return ""

	}
	resp, err := http.Get(target)

	if err == nil {
		contents, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			defer resp.Body.Close()
			return strings.TrimSpace(string(contents))

		}

	}
	return ""
}

func main() {

	err := argParse()
	if err != nil {
		log.Fatal(err)
	}

	checkIP()

}
