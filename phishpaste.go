package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gophish/gophish/models"

	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql" // swap for sqlite below if you want mysql
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DB is a global database connection for use by all functions
var DB *gorm.DB

func main() {

	listUsers := flag.Bool("list-users", false, "List all users")

	copyLandingPages := flag.Bool("landing-pages", false, "Copy Landing Pages from source user to desintation user")
	copyEmailTemplates := flag.Bool("email-templates", false, "Copy Email Templates from source user to desintation user")
	copySendingProfiles := flag.Bool("sending-profiles", false, "Copy Sending Profiles from source user to desintation user")
	dryRun := flag.Bool("dry-run", false, "Don't actually copy the data, just show what would be copied")

	sourceUser := flag.String("source", "", "User to copy data from")
	desintationUser := flag.String("destination", "", "User to copy data to")

	flag.Parse()

	if !*listUsers && *sourceUser == "" && *desintationUser == "" {
		fmt.Println("Welcome to Phish Paste. This utility will allow you to copy landing pages, email templates, and sending profiles between Gophish users.")
		fmt.Println("\nUsage: ")
		flag.PrintDefaults()
		os.Exit(1)
	}

	//Connect to the database
	err := initDB()
	if err != nil {
		fmt.Println("[!] Error connecting to database: ", err.Error())
		os.Exit(-1)
	}
	defer DB.Close()

	users, err := getUsers() // map from username to id
	if err != nil {
		fmt.Println("[!] Error getting users: ", err.Error())
		os.Exit(-1)
	}

	if *listUsers {

		fmt.Println("[+] Available users: ")
		for u := range users {
			fmt.Printf("\t%s\n", u)
		}

	} else {
		if *dryRun {
			fmt.Println("[-] Dry run mode enabled, no data will actually be copied")
		}

		// Check users exist
		if _, ok := users[*sourceUser]; !ok {
			fmt.Println("[!] No such source user: ", *sourceUser)
			os.Exit(-1)
		}
		if _, ok := users[*desintationUser]; !ok {
			fmt.Println("[!] No such source user: ", *desintationUser)
			os.Exit(-1)
		}

		// Copy each of the requried tables
		if *copyLandingPages {
			var pagesSource []models.Page
			var pagesDestination []models.Page
			pdMap := make(map[string]bool) // Used to check we don't have duplicate named pages
			var count int
			DB.Where("user_id = ?", users[*sourceUser]).Find(&pagesSource).Count(&count) // Source
			DB.Where("user_id = ?", users[*desintationUser]).Find(&pagesDestination)     // Destination
			for _, pd := range pagesDestination {
				pdMap[pd.Name] = true
			}

			fmt.Printf("[+] Copying %d Landing Pages from '%s' to '%s'\n", count, *sourceUser, *desintationUser)
			for _, p := range pagesSource {
				if _, ok := pdMap[p.Name]; ok {
					fmt.Printf("\t❌ %s (destination user already has page of the same name)\n", p.Name)
				} else {
					newPage := p
					newPage.Id = 0
					newPage.UserId = users[*desintationUser]
					if !*dryRun {
						DB.Create(&newPage)
					}
					fmt.Printf("\t✅ %s \n", p.Name)
				}
			}

		}
		if *copyEmailTemplates {
			var templatesSource []models.Template
			var templatesDestination []models.Template

			pdMap := make(map[string]bool) // Used to check we don't have duplicate template pages
			var count int
			var attachmentCount int
			DB.Where("user_id = ?", users[*sourceUser]).Find(&templatesSource).Count(&count) // Source
			DB.Where("user_id = ?", users[*desintationUser]).Find(&templatesDestination)     // Destination
			for _, pd := range templatesDestination {
				pdMap[pd.Name] = true
			}

			fmt.Printf("[+] Copying %d Email Templates from '%s' to '%s'...\n", count, *sourceUser, *desintationUser)
			for _, p := range templatesSource {
				if _, ok := pdMap[p.Name]; ok {
					fmt.Printf("\t❌ %s (destination user already has template of the same name)\n", p.Name)
				} else {

					// Grab attachments
					var attachments []models.Attachment
					DB.Where("template_id = ?", p.Id).Find(&attachments).Count(&attachmentCount)

					newTemplate := p
					newTemplate.Id = 0
					newTemplate.UserId = users[*desintationUser]
					//newTemplate.Attachments = newAttachments		// This does an update on the old attachment, can't figure out why. Solution is to do manual inserts below.

					if !*dryRun {
						DB.Create(&newTemplate)

						//Insert attachments
						newAttachments := attachments
						for _, a := range newAttachments {
							a.Id = 0
							a.TemplateId = newTemplate.Id // We grab the ID of the freshly inserted new template
							DB.Create(&a)                 // And insert each attachment
						}

					}
					if attachmentCount > 0 {
						fmt.Printf("\t✅ %s (including %d attachments) \n", p.Name, attachmentCount)
					} else {
						fmt.Printf("\t✅ %s \n", p.Name)
					}
				}
			}

		}
		if *copySendingProfiles {
			var smtpsSource []models.SMTP
			var smtpsDestination []models.SMTP
			pdMap := make(map[string]bool) // Used to check we don't have duplicate template pages
			var count int
			DB.Where("user_id = ?", users[*sourceUser]).Find(&smtpsSource).Count(&count) // Source
			DB.Where("user_id = ?", users[*desintationUser]).Find(&smtpsDestination)     // Destination
			for _, pd := range smtpsDestination {
				pdMap[pd.Name] = true
			}

			fmt.Printf("[+] Copying %d Sending Profiles from '%s' to '%s'...\n", count, *sourceUser, *desintationUser)
			for _, p := range smtpsSource {
				if _, ok := pdMap[p.Name]; ok {
					fmt.Printf("\t❌ %s (destination user already has sending profile of the same name)\n", p.Name)
				} else {
					newSMTP := p
					newSMTP.Id = 0
					newSMTP.UserId = users[*desintationUser]
					if !*dryRun {
						DB.Create(&newSMTP)
					}
					fmt.Printf("\t✅ %s \n", p.Name)
				}
			}

		}

	}

}

// initDB setups the database connection
func initDB() error {

	var err error
	if os.Getenv("DATABASE_URL") == "" {
		fmt.Println("[+] Connecting to gophish.db (set $DATABASE_URL to use MySQL)")
		DB, err = gorm.Open("sqlite3", "gophish.db")
	} else {
		fmt.Println("[+] Connecting to datbase")
		DB, err = gorm.Open("mysql", os.Getenv("DATABASE_URL"))
	}

	return err
}

func getUsers() (map[string]int64, error) {

	userMap := make(map[string]int64)
	var users []models.User
	if err := DB.Find(&users).Error; err != nil {
		return userMap, err
	}

	for _, u := range users {
		userMap[u.Username] = u.Id
	}

	return userMap, nil
}
