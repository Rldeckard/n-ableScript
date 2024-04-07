package main

import (
	"encoding/csv"
	"fmt"
	"os"
    "flag"

     "github.com/neteng-tools/cliPrompt"
     n "github.com/neteng-tools/n-ableScraper"

)

type Device struct {
    Address string
    Name string
}

func readCSV(filename string) ([]*Device, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []*Device{}, err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return []*Device{}, err
	}
	deviceMap := []*Device{}
	for _, element := range lines {
		deviceMap = append(deviceMap, &Device{
			Address: element[0],
			Name:  element[1],
		})
	}
	return deviceMap, err
}

func main() {
    username := flag.String("u", "", "Define username as an argument instead of through the prompt")
    password := flag.String("p", "", "Define password in plaintext. If not provided you'll get a secure prompt")
    server := flag.String("s", "", "Provide the https:// link to your N-Able tenant/server")
    fileLocation := flag.String("f", "data.csv", "Provide name of csv file.")
    flag.Parse()

    deviceList, err := readCSV(*fileLocation)
    if err != nil {
        panic(err)
    }
    if *server == "" {
        *server = prompt.Scan("Enter N-Able Server URL with https:// prefix:")
    }
    if *username == "" {
        *username = prompt.Credentials("Username: ")
    }
    if *password == "" {
        *password = prompt.Credentials("Password: ")
    }
    
    var nable n.NewPage
    nable.Connect(*server)
    nable.Page.MustWaitStable()
    defer nable.Page.MustClose()

    nable.Login(*username, *password)
    fmt.Println("Running")

    nable.AllDevicesPage()
    for _, device := range deviceList {
        fmt.Print("Searching tenant for " + device.Address)
        nable.Search(device.Address)
        nable.SelectAll().Edit()
        name, ok := nable.GetDeviceName()
        if !ok {
                fmt.Println("\nWARNING: Skipping " + device.Address + " due to navigation error.")
                nable.AllDevicesPage()
                continue
        }
        fmt.Print("...Current device name is '" + name + "'")
        nable.DeviceProps()
        nable.ChangeDeviceName(device.Name)
        nable.InputOsName("Other Operating System")
        nable.SaveChanges()
        newName, ok := nable.GetDeviceName()
        if ok {
            fmt.Println("...New device name is '" + newName + "'")
        } else {
            fmt.Println("...Getting updated device name failed")
        }
        nable.AllDevicesPage()
    }
    os.Exit(0)
}