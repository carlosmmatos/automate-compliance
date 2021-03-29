package main

import (
	"encoding/json"
	"fmt"
	"github.com/carlosmmatos/automate-compliance/internal/parser"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func normalizeControl(control string) {

}

func main() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Accessing the data in the spreadsheet
	// Using the following NIST RHACM 800-53 Example sheet:
	// https://docs.google.com/spreadsheets/d/12883Aj3eK3O0mgOesZMVnoVf8UmEPf1kPMyqFP7cp68/edit
	spreadsheetId := "12883Aj3eK3O0mgOesZMVnoVf8UmEPf1kPMyqFP7cp68"
	readRange := "800-53-controls-new!A2:M75"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).ValueRenderOption("FORMATTED_VALUE").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	// ***NOTE**** // The output of the Values.Get will mess up the length of
	// the array if the last cell in the index is empty. I created a dummy
	// column in the spreadsheet that will always be the last column in my
	// readRange and have a '.' in it so that it's not empty. Unless there is a
	// better way to ensure a consistent length of the row of values??
	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		p := parser.NewParser()
		fmt.Println("Control, Implementation Status, Narrative:")
		for _, row := range resp.Values {
			family := row[0].(string)
			control := row[1].(string)
			err := p.ParseEntry(family, control)
			if err != nil {
				fmt.Printf("Found error in control %s: %v\n", control, err)
				os.Exit(1)
			}
			//fmt.Printf("Raw Entry: %s, %s, %s\n", row[1], row[8], row[11])
		}

		fmt.Printf("Parsed data\n")
		fmt.Printf("===========\n\n")

		for family, controls := range p.GetData() {
			fmt.Printf("> %s\n", family)
			for _, ctrl := range controls {
				fmt.Printf("- %v\n", ctrl)
			}
		}
	}
}
