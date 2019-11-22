package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	externalip "github.com/GlenDC/go-external-ip"
	"github.com/cloudflare/cloudflare-go"
	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
)

var DOMAIN string
var CF_API_KEY string
var CF_API_EMAIL string
var DOMAINSPLIT string
var TARGETIPIP string

func writeConfig(result string) {

	f, err := os.OpenFile(checkConfig(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	cfdata := map[string]string{
		"Set domain name": "DOMAIN",
		"Set API Key":     "CF_API_KEY",
		"Set EMAIL":       "CF_API_EMAIL",
	}

	for k, v := range cfdata {
		//fmt.Printf("%s %s\n", k, v)
		labelprompt = k
		prompt := promptui.Prompt{
			Label: labelprompt,
		}
		result, err := prompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}
		writeConfig(v + "=" + result + "\n")
	}

	return

}

func checkConfig() string {

	currDir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		log.Fatal(err)
	}

	configDir := (currDir + "/config.env")

	return configDir

}

func loadConfig() error {

	err := godotenv.Load(checkConfig())

	if err != nil {
		fmt.Println("Error loading .env file, I would like to create?")
		result := yesNo()
		if result == false {
			os.Exit(0)
		} else {
			menuAPIdata()
		}
	}

	_ = godotenv.Load(checkConfig())

	DOMAIN = os.Getenv("DOMAIN")
	CF_API_KEY = os.Getenv("CF_API_KEY")
	CF_API_EMAIL = os.Getenv("CF_API_EMAIL")

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

func dynDNS() {

	PARTS := strings.SplitAfterN(DOMAIN, ".", 2)

	api, err := cloudflare.New(CF_API_KEY, CF_API_EMAIL)

	if err != nil {
		log.Fatal(err)
	}

	DOMAINSPLIT := PARTS[1]

	zoneID, err := api.ZoneIDByName(DOMAINSPLIT)

	if err != nil {
		log.Fatal(err)
		return
	}

	//Check ip change
	aIP := cloudflare.DNSRecord{Name: DOMAIN}

	recs, err := api.DNSRecords(zoneID, aIP)

	TARGETIP := getMyIP()

	for _, r := range recs {
		if TARGETIP == r.Content {
			log.Fatal("IP has not changed!")
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
		log.Fatal(err)
		return
	}

	if len(oldRecords) == 1 {
		err := api.UpdateDNSRecord(zoneID, oldRecords[0].ID, *newRecord)
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

	currentIP, err := consensus.ExternalIP()

	if err != nil {
		log.Fatal("Error collecting external IP", err)
	}

	TARGETIP := currentIP.String()

	return TARGETIP

}

func main() {

	err := loadConfig()

	if err != nil {
		log.Fatal("Error check ./config file", err)
	}

	dynDNS()

}
