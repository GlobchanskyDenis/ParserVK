package main

import (
	"os"
	"os/signal"
	"syscall"
	"database/sql"
	"fmt"
	"regexp"
)

const (
	GREEN = "\033[32m"
	RED = "\033[31m"
	YELLOW = "\033[33m"
	NO_COLOR = "\033[m"
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
			fmt.Println(RED + "THIS IS SERIOUS ERROR. YOU MUST DROP YOUR DATABASE TABLES FIRST" + NO_COLOR)
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
		fmt.Println(RED+"Cannot read config file: ", err, NO_COLOR)
		return
	}
	fmt.Println(GREEN + "Config successfully readed" + NO_COLOR)
	err = ConnectDatabase(config.Sql)
	if err != nil {
		fmt.Println(RED+"Database connection error: ", err, NO_COLOR)
		return
	}
	fmt.Println(GREEN + "Database succussfully connected" + NO_COLOR)

	quit := make(chan os.Signal, 1)

	go func(config *Config, quit chan os.Signal) {
		println(YELLOW + "userNbr vkUID             fname, lname, photo, birth, location" + NO_COLOR)
		parser(config.StartVkId, config.UsersAmount, PassHash(config.UsersPass))
		quit <- syscall.SIGINT
	}(config, quit)

	
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	gDbConn.Close()
	fmt.Println(GREEN + "Database connection successfully closed" + NO_COLOR)
}
