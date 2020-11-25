package main

import (
	"database/sql"
	"fmt"
	"regexp"
)

const (
	// 6253950 брат
	passSalt  = "+++"
	urlString = "https://vk.com/id"
	filename  = "result.html"
)

var (
	nameFinder    = regexp.MustCompile(`<h1 class="page_name">(.*)</h1>`)
	photoFinder   = regexp.MustCompile(`<img class="page_avatar_img" src="(.*)" alt`)
	dayFinder     = regexp.MustCompile(`people&c\[bday\]=(\d|\d\d)&c\[bmonth\]=`)
	monthFinder   = regexp.MustCompile(`&c\[bmonth\]=(\d|\d\d)">`)
	yearFinder    = regexp.MustCompile(`=people&c\[byear\]=(\d\d\d\d)">`)
	countryFinder = regexp.MustCompile(`uni_country\]=(\d|\d\d|\d\d\d)&c\[uni_city`)
	cityFinder    = regexp.MustCompile(`uni_country\]=(\d|\d\d|\d\d\d)&c\[uni_city`)

	gInterestsArr = []string{
		"football", "politics", "video games", "cooking", "programming", "youtube",
		"find something new", "culture", "drink beer", "architecture", "VODKA!!",
	}

	gDbConn *sql.DB
)

func parser(vkUID, UsersAmount int, EncryptedPass string) {
	var unknownInterests []string

	println("userNbr vkUID             fname, lname, photo, birth, location")

	for userNumber := 1; userNumber <= UsersAmount; vkUID++ {
		result, err := runGetWithHeaders(vkUID)
		if err != nil {
			printParseFail(vkUID, err)
			continue
		}
		user, err := parseToUser(result, userNumber, EncryptedPass)
		if err != nil {
			printParseFail(vkUID, err)
			continue
		}
		err = user.SaveInDatabase()
		if err != nil {
			printParseFail(vkUID, err)
			fmt.Println("THIS IS SERIOUS ERROR. YOU MUST DROP YOUR DATABASE TABLES FIRST")
			return
		}
		unknownInterests = addToArrayIfNeeded(unknownInterests, user.Interests)
		printParseSuccess(userNumber, vkUID, *user)
		userNumber++
	}
	err := SaveInterestsToDB(unknownInterests)
	if err != nil {
		fmt.Println("Cannot save interests in its table: " + err.Error())
		return
	}
}

func main() {
	config, err := GetConfig()
	if err != nil {
		fmt.Println("Cannot read config file: ", err)
		return
	}
	err = ConnectDatabase(config.Sql)
	if err != nil {
		fmt.Println("Database connection error: ", err)
		return
	}
	fmt.Printf("config successfully readed\n%#v\n", config)
	parser(config.StartVkId, config.UsersAmount, PassHash(config.UsersPass))
}
