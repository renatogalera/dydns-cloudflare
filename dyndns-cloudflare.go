package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/bogdanovich/dns_resolver"
	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/joho/godotenv"
)

var DOMAIN string
var CF_API_KEY string
var CF_API_EMAIL string
var SUBDOMAIN string
var NEWIPADDR string

func loadConfig() error {
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
	return nil
}

func checaIPDNS(target string) {
	resolver := dns_resolver.New([]string{"1.1.1.1"})
	resolver.RetryTimes = 5
	ip, err := resolver.LookupHost(SUBDOMAIN + "." + DOMAIN)
	
	if err != nil {
		log.Fatal(err.Error())
	}
	
	for _, ips := range ip {
		dnsip := ips.String()
		if target == dnsip {
			fmt.Println("IP não mudou!")
			os.Exit(0)
		}else{
			dynDNS(target)
		}
	}
}

func dynDNS(target string) {
	api, err := cloudflare.New(CF_API_KEY, CF_API_EMAIL)
	if err != nil {
		log.Fatal(err)
	}
	
	zoneID, err := api.ZoneIDByName(DOMAIN)
	if err != nil {
		log.Fatal(err)
		return
	}
	
	newRecord := cloudflare.DNSRecord{
		Type:    "A",
		Name:    SUBDOMAIN + "." + DOMAIN,
		Content: target,
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
		//} else if protocol == 6 {
		//	target = "http://ifconfig.me/ip" //Alterar para fonte ipv6
	}else{
		os.Exit(0)
	}

	resp, err := http.Get(target)

	if err == nil {
		contents, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			defer resp.Body.Close()
			return strings.TrimSpace(string(contents))
		}

	}
	return target
}

func main() {
	err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	
	target := getMyIP(4)
	checaIPDNS(target)
}
