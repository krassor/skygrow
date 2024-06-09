package dto

type ResponseCalendarStatus struct {
	CalendarId string `json:"calendarId"`
	Status     string `json:"status"`
}

type ResponseCalendar struct {
	CalendarId       string `json:"calendar_id"`
	CalendarOwnerId  string `json:"calendar_owner_id"`
	GoogleCalendarId string `json:"google_calendar_id"`
	Description      string `json:"description"`
	Etag             string `json:"etag"`
	Summary          string `json:"summary"`
	TimeZone         string `json:"time_zone"`
	Status           string `json:"status"`
}
