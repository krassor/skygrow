/**
 * @license
 * Copyright Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
// [START calendar_quickstart]
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "credentials/token.json"
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

func main() {
	ctx := context.Background()
	b, err := os.ReadFile("credentials/client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	yogaTestCalendar := "cebd29e02309dea02f13bd2acdd06aa666be7bcab4a98f956b64a98ea4aeb0b1@group.calendar.google.com"

	//t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List(yogaTestCalendar).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("\t\t%v (%v)\n\n", item.Summary, date)
			fmt.Printf("%v\n%v\n%v\n\n", item.ConferenceData, item.EventType, item.HangoutLink)
			//fmt.Printf("%v\n", item.ConferenceData.EntryPoints[0])

		}
	}

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day()+1, 12, 0, 0, 0, time.Local)
	end := time.Date(now.Year(), now.Month(), now.Day()+1, 13, 0, 0, 0, time.Local)

	meetId := "7ex-a0wv-y0t"
	// entryPoint := &calendar.EntryPoint{
	// 	EntryPointType: "video",
	// 	Label:          fmt.Sprintf("meet.google.com/%s", meetId),
	// 	Uri:            fmt.Sprintf("https://meet.google.com/%s", meetId),
	// }
	// entryPoints := make([]*calendar.EntryPoint, 1)
	// entryPoints[0] = entryPoint

	event := &calendar.Event{
		ConferenceData: &calendar.ConferenceData{
			ConferenceId: meetId,
			// ConferenceSolution: &calendar.ConferenceSolution{
			// 	IconUri: "https://fonts.gstatic.com/s/i/productlogos/meet_2020q4/v6/web-512dp/logo_meet_2020q4_color_2x_web_512dp.png",
			// 	Key: &calendar.ConferenceSolutionKey{
			// 		Type: "hangoutsMeet",
			// 	},
			// 	Name: "Meet name",
			// },
			// EntryPoints: entryPoints,
			CreateRequest: &calendar.CreateConferenceRequest{
				ConferenceSolutionKey: &calendar.ConferenceSolutionKey{
					Type: "hangoutsMeet",
				},
				RequestId: "7qxalsvy0e",
			},
		},
		Description: "API_TEST",
		End: &calendar.EventDateTime{
			DateTime: end.Format(time.RFC3339),
			TimeZone: "Europe/Moscow",
		},
		Id: fmt.Sprintf("%v", now.Unix()),
		Start: &calendar.EventDateTime{
			DateTime: start.Format(time.RFC3339),
			TimeZone: "Europe/Moscow",
		},
		Status:       "tentative",
		Summary:      fmt.Sprintf("API_TEST_Summary %v", now.Unix()),
		Transparency: "opaque",
	}

	_, err = srv.Events.Insert(yogaTestCalendar, event).ConferenceDataVersion(1).Do()
	if err != nil {
		fmt.Printf("Error insert: %v", err)
		return
	}

	calendar := &calendar.Calendar {
		ConferenceProperties: &calendar.ConferenceProperties{
			AllowedConferenceSolutionTypes: []string {"eventHangout", "eventNamedHangout", "hangoutsMeet"},
		},
		Description: "Description 1",
		Etag: "yogi 1",
		Summary: "Summary 1",
		TimeZone: "Europe/Moscow",

	}	
	_, err = srv.Calendars.Insert(calendar).Do()
	if err != nil {
		fmt.Printf("Error insert: %v", err)
		return
	}

}