package main

import (
	"errors"
	"fmt"
	"github.com/GlenDC/go-external-ip"
	"github.com/cloudflare/cloudflare-go"
	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
	"log"
	"os"
	"path/filepath"
)

var DOMAIN string
var CF_API_KEY string
var CF_API_EMAIL string
var SUBDOMAIN string

func writeConfig(result string) {
	f, err := os.OpenFile(dir(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(result)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func menuAPIdata() (labelprompt string, result string) {
	apidata := map[string]string{
		"Set domain name": "DOMAIN",
		"Set subdomain": "SUBDOMAIN",
		"Set API Key": "CF_API_KEY",
		"Set EMAIL": "CF_API_EMAIL",
	}
	for k, v := range apidata {
		//fmt.Printf("%s %s\n", k, v)
		labelprompt = k
		prompt := promptui.Prompt{
			Label:    labelprompt,
		}
		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
		}
		writeConfig(v + "=" + result + "\n")
	}

	return
}

func dir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
	log.Fatal(err)
	}
	configdir := (dir + "/config.env")
	return configdir
}

func loadConfig() error {
	err := godotenv.Load(dir())
	if err != nil {
		fmt.Println("Error loading .env file, I would like to create?")
		result := yesNo()
		if result ==false {
			os.Exit(0)
		}else{
			menuAPIdata()
		}
	}

	_ = godotenv.Load(dir())

	DOMAIN = os.Getenv("DOMAIN")
	if DOMAIN == "" {
		msg := fmt.Sprintf("Please, set variable DOMAIN")
		return errors.New(msg)
	}

	CF_API_KEY = os.Getenv("CF_API_KEY")
	if CF_API_KEY == "" {
		msg := fmt.Sprintf("Please, set variable CF_API_KEY")
		return errors.New(msg)
	}

	CF_API_EMAIL = os.Getenv("CF_API_EMAIL")
	if CF_API_EMAIL == "" {
		msg := fmt.Sprintf("Please, set variable CF_API_EMAIL")
		return errors.New(msg)
	}

	SUBDOMAIN = os.Getenv("SUBDOMAIN")
	if SUBDOMAIN == "" {
		msg := fmt.Sprintf("Please, set variable SUBDOMAIN")
		return errors.New(msg)
	}
	return nil
}

func yesNo() bool {
	prompt := promptui.Select{
		Label: "Select [Yes/No]",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "Yes"
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
	log.Println("IP changed, setting DNS entry:", newRecord.Name, newRecord.Content, "\n")
}

func updateRecord(zoneID string, api *cloudflare.API, newRecord *cloudflare.DNSRecord) {

	dns := cloudflare.DNSRecord{Type: newRecord.Type, Name: newRecord.Name}
	OLDRECORDS, err := api.DNSRecords(zoneID, dns)

	if err != nil {
		log.Fatal(err)
		return
	}

	if len(OLDRECORDS) == 1 {
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

func getMyIP() string {

	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()

	if err != nil {
		fmt.Println("Error collecting external IP") // print IPv4/IPv6 in string format
	}

	target := ip.String()
	return target
}

func main() {

	err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	api, err := cloudflare.New(CF_API_KEY, CF_API_EMAIL)
	if err != nil {
		log.Fatal(err)
	}

	zoneID, err := api.ZoneIDByName(DOMAIN)
	if err != nil {
		log.Fatal(err)
		return
	}

	ipdns := cloudflare.DNSRecord{Name: SUBDOMAIN + "." + DOMAIN}
	recs, err := api.DNSRecords(zoneID, ipdns)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range recs {
		if getMyIP() == r.Content {
			fmt.Println("IP has not changed!")
			os.Exit(0)
		}else{
			dynDNS(getMyIP())
		}
	}
}
