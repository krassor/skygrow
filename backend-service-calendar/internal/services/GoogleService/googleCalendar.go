package GoogleService

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"net/http"
	"os"
)

var (
	ErrorEmptyCalendarField = errors.New("calendar field is empty")
)

type GoogleCalendar struct {
	googleService *calendar.Service
	googleClient  *http.Client
}

//func (gc *GoogleCalendar) CreateCalendar(
//	description string, etag string, summary string, timezone string) (string, error) {
//	op := "GoogleService.CreateCalendar()"
//	if summary == "" {
//		return "", fmt.Errorf("%s : %w", op, ErrorEmptyCalendarField)
//	}
//	if etag == "" {
//		return "", fmt.Errorf("%s : %w", op, ErrorEmptyCalendarField)
//	}
//	if timezone == "" {
//		return "", fmt.Errorf("%s : %w", op, ErrorEmptyCalendarField)
//	}
//	newCalendar := &calendar.Calendar{
//		ConferenceProperties: &calendar.ConferenceProperties{
//			AllowedConferenceSolutionTypes: []string{"eventHangout", "eventNamedHangout", "hangoutsMeet"},
//		},
//		Description: description,
//		Etag:        etag,
//		Summary:     summary,
//		TimeZone:    timezone,
//	}
//	cal, err := gc.googleService.Calendars.Insert(newCalendar).Do()
//	if err != nil {
//		return "", fmt.Errorf("%s : %w", op, err)
//	}
//
//	return cal.Id, nil
//}

// CreateCalendar() return Google calendar ID. Return non nil error if function cannot create calendar with Google API
func (gc *GoogleCalendar) CreateCalendar(
	description string, summary string, timezone string) (string, error) {
	op := "GoogleService.CreateCalendar()"
	if summary == "" {
		return "", fmt.Errorf("%s : %w", op, ErrorEmptyCalendarField)
	}
	if timezone == "" {
		return "", fmt.Errorf("%s : %w", op, ErrorEmptyCalendarField)
	}
	newCalendar := &calendar.Calendar{
		ConferenceProperties: &calendar.ConferenceProperties{
			AllowedConferenceSolutionTypes: []string{"eventHangout", "eventNamedHangout", "hangoutsMeet"},
		},
		Description: description,
		Summary:     summary,
		TimeZone:    timezone,
	}
	cal, err := gc.googleService.Calendars.Insert(newCalendar).Do()
	if err != nil {
		return "", fmt.Errorf("%s : %w", op, err)
	}

	return cal.Id, nil
}

func NewGoogleCalendar() *GoogleCalendar {
	b := getClientSecret("credentials/client_secret.json")

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatal().Msgf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	srv, err := calendar.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatal().Msgf("unable to retrieve Calendar client: %v", err)
	}

	return &GoogleCalendar{
		googleService: srv,
		googleClient:  client,
	}

}

func getClientSecret(filepath string) []byte {
	if filepath == "" {
		log.Fatal().Msgf("google client_secret path is empty")
	}

	// check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Fatal().Msgf("google client_secret.json file does not exist: " + filepath)
	}
	//"credentials/client_secret.json"
	b, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal().Msgf("unable to read client secret file: %v", err)
	}

	return b
}

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
		log.Fatal().Msgf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatal().Msgf("Unable to retrieve token from web: %v", err)
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
		log.Fatal().Msgf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		log.Fatal().Msgf("cannot encode token to save file: %v", err)
	}
}
