package main

import (
	"errors"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"
)

type User struct {
	Uid           int
	Mail          string
	EncryptedPass string
	Fname         string
	Lname         string
	Birth         *time.Time
	Gender        string
	Orientation   string
	Bio           string
	AvaID         int
	Avatar        string
	Latitude      *float64
	Longitude     *float64
	Interests     []string
	Status        string
	Rating        int
}

type Location struct {
	CountryCode int
	Country     string
	CityCode    int
	City        string
	Longitude   *float64
	Latitude    *float64
}

func (location *Location) CalculateLocation() {
	var latitude, longitude float64

	switch location.CountryCode {
	case 1:
		location.Country = "Russia"
		switch location.CityCode {
		case 1:
			location.City, latitude, longitude = "Moscow", 55.753960, 37.620393
		case 2:
			location.City, latitude, longitude = "Petersburg", 59.939095, 30.315868
		default:
			location.City, latitude, longitude = "Moscow", 55.753960, 37.620393
		}
	case 2:
		location.Country, location.City, latitude, longitude = "Ukraine", "Kiev", 50.4536, 30.5164
	default:
		randomized := rand.Int()
		if randomized%2 == 0 {
			location.Country, location.City, latitude, longitude = "Russia", "Moscow", 55.753960, 37.620393
		}
	}

	if location.Country != "" {
		latitude = randomizeNumber(latitude)
		longitude = randomizeNumber(longitude)
		location.Latitude = &latitude
		location.Longitude = &longitude
	}
}

func (user *User) ChooseGender() {
	if user.Gender != "" {
		return
	}
	if isFemale(user.Fname) {
		user.Gender = "female"
	} else if isMale(user.Fname) {
		user.Gender = "male"
	} else {
		if user.Fname == "Женя" || user.Fname == "Женя" {
			lname := []rune(user.Lname)
			if len(lname) > 0 && lname[len(lname)-1] == 'а' {
				user.Gender = "female"
			} else {
				user.Gender = "male"
			}
		}
	}
	if user.Gender == "" {
		randomized := rand.Int()
		if randomized%3 == 0 {
			user.Gender = ""
		} else if randomized%3 == 1 {
			user.Gender = "male"
		} else {
			user.Gender = "female"
		}
	}
}

func (user *User) ChooseOrientation() {
	randomized := rand.Int()
	if randomized%10 >= 0 && randomized%10 <= 6 {
		user.Orientation = "hetero"
	} else if randomized%10 == 7 || randomized%10 == 8 {
		user.Orientation = "homo"
	} else {
		user.Orientation = ""
	}
}

func (user *User) GenerateRating() {
	randomized := rand.Int()
	user.Rating = randomized % 12
}

func (user *User) FillNameFields(name string) error {
	var nameArr []string
	nameArrSrc := strings.Split(name, " ")

	for _, nameSrc := range nameArrSrc {
		if nameSrc != "" {
			nameArr = append(nameArr, nameSrc)
		}
	}
	if len(nameArr) < 1 {
		return errors.New("empty name")
	} else if len(nameArr) > 3 {
		return errors.New("too many words in name")
	}

	if len(nameArr) == 3 {
		lname := nameArr[2]
		if len(lname) < 2 {
			return errors.New("too short last name")
		}
		if lname[0] == '(' && lname[len(lname)-1] == ')' {
			//  Третье имя в скобках - это фамилия до ее замены.
			//	С высокой вероятностью это замужняя женщина
			user.Lname = string(([]rune(lname))[1 : utf8.RuneCountInString(lname)-1])
			user.Gender = "female"
			user.Orientation = "hetero"
		} else {
			user.Lname = lname
		}
	} else if len(nameArr) == 2 {
		user.Lname = nameArr[1]
	}

	user.Fname = nameArr[0]
	user.Bio += name + " "
	return nil
}

func (user *User) GenerateInterests() {
	var interestsArray []string = nil

	randomized := rand.Int()
	amounInterests := randomized % 5
	for i := 0; i < amounInterests; i++ {
		randomized := rand.Int()
		newInterest := gInterestsArr[randomized%10]
		if isArrayContainsString(interestsArray, newInterest) {
			continue
		}
		interestsArray = append(interestsArray, newInterest)
	}
	user.Interests = interestsArray
}
