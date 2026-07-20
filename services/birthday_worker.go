package services

import (
	"fmt"
	"log"
	"memories/discord"
	"memories/model"
	"strings"
	"time"
)

// CheckBirthdays checks for birthdays today and sends a Discord message if any are found.
func CheckBirthdays(dao model.LifeJournalDAO, webhookURL string) {
	people, err := dao.GetAllPeople()
	if err != nil {
		log.Printf("Error getting people for birthday check: %v\n", err)
		return
	}

	now := time.Now()
	todayMonth := int(now.Month())
	todayDay := now.Day()

	var bdayPeople []string
	for _, p := range people {
		if p.BirthMonth == todayMonth && p.BirthDay == todayDay {
			name := fmt.Sprintf("%s %s", p.First, p.Last)
			bdayPeople = append(bdayPeople, name)
		}
	}

	if len(bdayPeople) > 0 {
		message := "🎂 Birthdays today: \n" + strings.Join(bdayPeople, ", \n")
		err = discord.SendWebhookMessage(webhookURL, message)
		if err != nil {
			log.Printf("Error sending discord message: %v\n", err)
		} else {
			log.Println("Birthday notification sent successfully.")
		}
	} else {
		log.Println("No birthdays today. Sending status message to Discord.")
		err = discord.SendWebhookMessage(webhookURL, "No birthdays today.")
		if err != nil {
			log.Printf("Error sending discord message: %v\n", err)
		}
	}
}

// StartBirthdayWorker starts a background goroutine that runs CheckBirthdays daily at the specified time.
// runTime should be in "HH:MM" format (24-hour).
func StartBirthdayWorker(dao model.LifeJournalDAO, webhookURL string, runTime string) {
	go func() {
		for {
			now := time.Now()

			// Parse the runTime
			var hour, min int
			_, err := fmt.Sscanf(runTime, "%d:%d", &hour, &min)
			if err != nil {
				log.Printf("Invalid BIRTHDAY_CHECK_TIME format '%s', expected HH:MM. Defaulting to 09:00.\n", runTime)
				hour, min = 9, 0
			}

			// Calculate next run time
			nextRun := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())
			if nextRun.Before(now) {
				nextRun = nextRun.Add(24 * time.Hour)
			}

			log.Printf("Next birthday check scheduled for: %v\n", nextRun)

			// Sleep until nextRun
			timer := time.NewTimer(time.Until(nextRun))
			<-timer.C

			CheckBirthdays(dao, webhookURL)
		}
	}()
}
