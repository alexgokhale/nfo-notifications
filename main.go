package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	Token           string
	Email           string
	Password        string
	WebhookURL      string
	Hostname        string
	pacificLocation *time.Location
)

type Event struct {
	ID             string
	Time           time.Time
	Router         string
	TargetAddress  string
	AttackType     string
	FilterDuration string
}

func handleError(s string, err error) {
	if err != nil {
		log.Fatalf("💔 %s, %v\n", s, err)
	}
}

func init() {
	flag.StringVar(&Token, "t", "", "An existing authentication token")
	flag.StringVar(&Email, "e", "", "Your account email")
	flag.StringVar(&Password, "p", "", "Your account password")
	flag.StringVar(&WebhookURL, "w", "", "Discord webhook to send new events to")
	flag.StringVar(&Hostname, "h", "", "The hostname to fetch events for")
	flag.Parse()

	if Email == "" || Password == "" {
		log.Fatalln("💔 Email or Password argument missing")
	}

	if WebhookURL == "" {
		log.Fatalln("💔 Webhook URL argument missing")
	}

	if Hostname == "" {
		log.Fatalln("💔 Hostname argument missing")
	}

	loc, err := time.LoadLocation("America/Los_Angeles")

	if err != nil {
		log.Fatalln("💔 Couldn't load Pacific Time:", err)
	}

	pacificLocation = loc
}

func main() {
	httpClient := http.Client{
		CheckRedirect: nil,
	}

	if Token == "" {
		loginRes, err := httpClient.PostForm("https://www.nfoservers.com/control/status.pl", url.Values{
			"email":       {Email},
			"password":    {Password},
			"permacookie": {"on"},
			"login":       {"Log in to your control panel"},
		})

		handleError("Login request failed", err)

		if loginRes.StatusCode != http.StatusOK {
			log.Fatalln("💔 Couldn't login to account, are your credentials correct?")
		}

		token, err := getTokenFromCookies(loginRes.Cookies())

		handleError("Couldn't get token", err)

		Token = token
	}

	eventsUrl := fmt.Sprintf("https://www.nfoservers.com/control/events.pl?name=%s&typeofserver=virtual", Hostname)
	eventsReq, _ := http.NewRequest("GET", eventsUrl, nil)
	eventsReq.AddCookie(&http.Cookie{
		Name:  "cookietoken",
		Value: Token,
	})
	eventsReq.AddCookie(&http.Cookie{
		Name:  "email",
		Value: Email,
	})
	eventsReq.AddCookie(&http.Cookie{
		Name:  "password",
		Value: Password,
	})

	var oldEvents []*Event

	for {
		fmt.Println("Getting events...")
		eventsRes, err := httpClient.Do(eventsReq)

		handleError("Events request failed", err)

		if eventsRes.StatusCode != http.StatusOK {
			log.Fatalln("💔 Couldn't fetch event log")
		}

		fmt.Println("Events fetched successfully")
		doc, err := goquery.NewDocumentFromReader(eventsRes.Body)

		handleError("Couldn't parse response body", err)

		logTables := doc.Find(`.logtable`)

		if logTables.Length() < 2 {
			log.Fatalln("💔 Can't find event table on page")
		}

		events := make([]*Event, 0)
		eventTable := goquery.NewDocumentFromNode(logTables.Get(1))

		eventTable.Find(`tbody > tr:not(.logheading)`).Each(func(_ int, selection *goquery.Selection) {
			selection.Find(`td`).Each(func(_ int, selection *goquery.Selection) {
				if len(selection.Children().Nodes) == 2 {
					subjectElement := selection.Find(`span[id^="event_standard_"]`)

					if subjectElement.Text() != "\n(D)DoS attack against your service\n" {
						return
					}

					id, ok := subjectElement.Attr("id")

					if !ok {
						return
					}

					id = strings.TrimPrefix(id, "event_standard_")
					id = strings.TrimSuffix(id, "_subj")

					event := Event{
						ID: id,
					}

					t := strings.Replace(selection.Find(`span i`).Text(), "PT", "PST", 1)
					parsedTime, _ := time.ParseInLocation("Jan 02 2006 03:04:05 PM MST", t, pacificLocation)
					event.Time = parsedTime

					events = append(events, &event)
				} else {
					infoElement := selection.Find(`span`)

					id, ok := infoElement.Attr("id")

					if !ok {
						return
					}

					id = strings.TrimPrefix(id, "event_standard_")

					event := getEventById(id, events)

					if event == nil {
						return
					}

					infoElement.Children().RemoveFiltered(`a`)
					infoElement.Children().RemoveFiltered(`div`)
					infoHtml, err := infoElement.Html()

					if err != nil {
						return
					}

					infoParts := strings.Split(infoHtml, "<br/>")
					infoParts = stripEmpty(infoParts)

					if len(infoParts) != 4 {
						return
					}

					event.Router = strings.TrimPrefix(infoParts[0], "Our system responded to a (D)DoS against your service with a filter on ")
					event.Router = strings.Trim(event.Router, ".\n")

					event.TargetAddress = strings.TrimPrefix(infoParts[1], "Target address: ")
					event.TargetAddress = strings.Trim(event.TargetAddress, ".\n")

					event.AttackType = strings.TrimPrefix(infoParts[2], "Attack: ")
					event.AttackType = strings.Trim(event.AttackType, ".\n")

					event.FilterDuration = strings.TrimPrefix(infoParts[3], "Filter duration: ")
					event.FilterDuration = strings.Trim(event.FilterDuration, ".\n")
				}
			})
		})

		fmt.Println("Parsed", len(events), "events")

		if len(oldEvents) > 0 {
			newEvents := findNewEvents(oldEvents, events)

			if len(newEvents) != 0 {
				fmt.Println(len(newEvents), "new events, sending webhook")
				sendNewEvents(newEvents)
			} else {
				fmt.Println("No new events, continuing")
			}
		}

		oldEvents = events

		fmt.Println("Pausing for 1 minute")
		time.Sleep(time.Minute * 1)
	}
}
