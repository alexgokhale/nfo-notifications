package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func stripEmpty(elems []string) []string {
	elemsToReturn := make([]string, 0)

	for _, elem := range elems {
		elem = strings.Trim(elem, " \n")
		if elem != "" {
			elemsToReturn = append(elemsToReturn, elem)
		}
	}

	return elemsToReturn
}

func getEventById(ID string, events []*Event) *Event {
	for _, event := range events {
		if event.ID == ID {
			return event
		}
	}

	return nil
}

func getTokenFromCookies(cookies []*http.Cookie) (string, error) {
	for _, cookie := range cookies {
		if cookie.Name == "cookietoken" {
			return cookie.Value, nil
		}
	}

	return "", errors.New("token cookie missing")
}

func findNewEvents(old []*Event, new []*Event) []Event {
	oldMap := map[string]bool{}
	newMap := map[string]bool{}

	for _, event := range old {
		oldMap[event.ID] = true
	}

	for _, event := range new {
		newMap[event.ID] = true
	}

	eventsToReturn := make([]Event, 0)

	for k := range newMap {
		if !oldMap[k] {
			eventsToReturn = append(eventsToReturn, *getEventById(k, new))
		}
	}

	return eventsToReturn
}

func removeLastOctet(ip string) string {
	ipParts := strings.Split(ip, ".")
	ipParts[len(ipParts)-1] = "xxx"

	newIp := ""

	for i, part := range ipParts {
		newIp += part

		if i != len(ipParts)-1 {
			newIp += "."
		}
	}

	return newIp
}

func sendNewEvents(events []Event) {
	for _, event := range events {
		err := sendDiscordWebhook(WebhookURL, Webhook{
			Embeds: []Embed{
				{
					Title:       "New Event",
					Description: "A DDoS attack was filtered by NFO's firewall",
					Color:       9772083,
					Timestamp:   event.Time,
					Fields: []Field{
						{
							Name:  "Attack Type",
							Value: event.AttackType,
						},
						{
							Name:   "Router",
							Value:  event.Router,
							Inline: true,
						},
						{
							Name:   "Filter Duration",
							Value:  event.FilterDuration,
							Inline: true,
						},
						{
							Name:   "Target Address",
							Value:  removeLastOctet(event.TargetAddress),
							Inline: true,
						},
					},
				},
			},
		})

		if err != nil {
			fmt.Println(err)
		}
	}
}
