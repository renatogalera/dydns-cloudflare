package main

import (
	"log"
	"os"
	"strings"
	"time"

	externalip "github.com/GlenDC/go-external-ip"
	"github.com/cloudflare/cloudflare-go"
	"github.com/joho/godotenv"
)

var DOMAIN string
var CF_API_KEY string
var CF_API_EMAIL string

func loadConfig() error {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	DOMAIN = os.Getenv("DOMAIN")
	if DOMAIN == "" {
		log.Fatal("Need to define DOMAIN var")
	}
	CF_API_KEY = os.Getenv("CF_API_KEY")
	if CF_API_KEY == "" {
		log.Fatal("Need to define CF_API_KEY var")
	}
	CF_API_EMAIL = os.Getenv("CF_API_EMAIL")
	if CF_API_EMAIL == "" {
		log.Fatal("Need to define CF_API_EMAIL var")
	}
	return nil
}

func dynDNS() {
	PARTS := strings.SplitAfterN(DOMAIN, ".", 2)
	api, err := cloudflare.New(CF_API_KEY, CF_API_EMAIL)
	if err != nil {
		log.Println(err)
	}
	DOMAINSPLIT := PARTS[1]
	zoneID, err := api.ZoneIDByName(DOMAINSPLIT)
	if err != nil {
		log.Println(err)
		return
	}
	aIP := cloudflare.DNSRecord{Name: DOMAIN}
	recs, err := api.DNSRecords(zoneID, aIP)
	TARGETIP := getMyIP()
	for _, r := range recs {
		if TARGETIP == r.Content {
			log.Println("IP has not changed!")
			time.Sleep(120 * time.Second)
			main()
		}
	}
	newRecord := cloudflare.DNSRecord{
		Type:    "A",
		Name:    DOMAIN,
		Content: TARGETIP,
	}
	updateRecord(zoneID, api, &newRecord)
	log.Println("IP changed:", "\nDNS: ", newRecord.Name, "\nIP: ", newRecord.Content, "\n")
}

func updateRecord(zoneID string, api *cloudflare.API, newRecord *cloudflare.DNSRecord) {
	DNSRecordIP := cloudflare.DNSRecord{Type: newRecord.Type, Name: newRecord.Name}
	oldRecords, err := api.DNSRecords(zoneID, DNSRecordIP)
	if err != nil {
		log.Println(err)
		return
	}
	if len(oldRecords) == 1 {
		err := api.UpdateDNSRecord(zoneID, oldRecords[0].ID, *newRecord)
		if err != nil {
			log.Println(err)
			return
		}
		return
	}
	_, err = api.CreateDNSRecord(zoneID, *newRecord)
	if err != nil {
		log.Println(err)
		return
	}
}

func getMyIP() string {
	consensus := externalip.DefaultConsensus(nil, nil)
	currentIP, err := consensus.ExternalIP()
	if err != nil {
		log.Println("Error collecting external IP", err)
	}
	TARGETIP := currentIP.String()
	return TARGETIP
}

func main() {
	log.SetOutput(os.Stdout)
	loadConfig()
	for {
		dynDNS()
	}
}
