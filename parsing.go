package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

func parseName(src []byte) ([]byte, error) {

	nameSlice := nameFinder.FindAllSubmatch(src, -1)
	if len(nameSlice) > 0 {
		nameArr := nameSlice[0]
		if len(nameArr) > 1 {
			if len(nameArr[1]) > 60 {
				return nil, errors.New("too long user name")
			}
			return nameArr[1], nil
		} else {
			return nil, errors.New("invalid user name array")
		}
	}
	return nil, errors.New("invalid user name slice")
}

func parseAvatar(src []byte) (string, error) {
	photoSlice := photoFinder.FindAllSubmatch(src, -1)
	if len(photoSlice) > 0 {
		photoArr := photoSlice[0]
		if len(photoArr) > 1 {
			return string(photoArr[1]), nil
		} else {
			return "", errors.New("invalid photo array")
		}
	}
	return "", errors.New("invalid photo slice")
}

func parseBirth(src []byte) *time.Time {
	var day, month, year string
	var t time.Time
	var err error

	daySlice := dayFinder.FindAllSubmatch(src, -1)
	if len(daySlice) > 0 {
		dayArr := daySlice[0]
		if len(dayArr) > 1 {
			day = string(dayArr[1])
		} else {
			return nil
		}
	} else {
		return nil
	}
	if len(day) == 1 {
		day = "0" + day
	}

	monthSlice := monthFinder.FindAllSubmatch(src, -1)
	if len(monthSlice) > 0 {
		monthArr := monthSlice[0]
		if len(monthArr) > 1 {
			month = string(monthArr[1])
		} else {
			return nil
		}
	} else {
		return nil
	}
	if len(month) == 1 {
		month = "0" + month
	}

	yearSlice := yearFinder.FindAllSubmatch(src, -1)
	if len(yearSlice) > 0 {
		yearArr := yearSlice[0]
		if len(yearArr) > 1 {
			year = string(yearArr[1])
		} else {
			return nil
		}
	} else {
		return nil
	}

	t, err = time.Parse("2006-01-02", year+"-"+month+"-"+day)
	if err != nil {
		fmt.Println("time parse error:", err)
		return nil
	}
	return &t
}

func parseToLocation(src []byte) (*Location, error) {
	var location Location
	var err error

	countrySlice := countryFinder.FindAllSubmatch(src, -1)
	if len(countrySlice) > 0 {
		countryArr := countrySlice[0]
		if len(countryArr) > 1 {
			location.CountryCode, err = strconv.Atoi(string(countryArr[1]))
			if err != nil {
				return &location, errors.New("country parse error")
			}
		} else {
			return &location, errors.New("invalid country array")
		}
	} else {
		return &location, errors.New("invalid country slice")
	}

	citySlice := cityFinder.FindAllSubmatch(src, -1)
	if len(citySlice) > 0 {
		cityArr := citySlice[0]
		if len(cityArr) > 1 {
			location.CityCode, err = strconv.Atoi(string(cityArr[1]))
			if err != nil {
				return &location, errors.New("city parse error")
			}
		} else {
			return &location, errors.New("invalid city array")
		}
	} else {
		return &location, errors.New("invalid city slice")
	}
	return &location, nil
}

func parseToUser(src []byte, userNumber int, encryptedPass string) (*User, error) {
	var nameString string
	var nameByte []byte
	var user User
	var err error

	user.Mail = "user" + strconv.Itoa(userNumber) + "@gmail.com"
	user.EncryptedPass = encryptedPass

	nameByte, err = parseName(src)
	if err != nil {
		return nil, err
	}
	nameString, err = myConverter(nameByte)
	if err != nil {
		return nil, err
	}

	err = user.FillNameFields(nameString)
	if err != nil {
		return nil, err
	}

	user.Avatar, err = parseAvatar(src)
	if err != nil {
		return nil, err
	}

	location, _ := parseToLocation(src)
	location.CalculateLocation()
	if location.Longitude != nil && location.Latitude != nil {
		user.Longitude = location.Longitude
		user.Latitude = location.Latitude
		user.Bio += location.Country + " " + location.City + " "
	}

	user.Birth = parseBirth(src)

	user.ChooseGender()
	user.ChooseOrientation()
	user.GenerateRating()
	user.GenerateInterests()
	user.Status = "confirmed"
	return &user, nil
}

func myConverter(src []byte) (string, error) {
	var dst []rune
	for _, letter := range src {
		switch letter {
		case 0x41:
			dst = append(dst, 'A')
		case 0x42:
			dst = append(dst, 'B')
		case 0x43:
			dst = append(dst, 'C')
		case 0x44:
			dst = append(dst, 'D')
		case 0x45:
			dst = append(dst, 'E')
		case 0x46:
			dst = append(dst, 'F')
		case 0x47:
			dst = append(dst, 'G')
		case 0x48:
			dst = append(dst, 'H')
		case 0x49:
			dst = append(dst, 'I')
		case 0x4a:
			dst = append(dst, 'J')
		case 0x4b:
			dst = append(dst, 'K')
		case 0x4c:
			dst = append(dst, 'L')
		case 0x4d:
			dst = append(dst, 'M')
		case 0x4e:
			dst = append(dst, 'N')
		case 0x4f:
			dst = append(dst, 'O')
		case 0x50:
			dst = append(dst, 'P')
		case 0x51:
			dst = append(dst, 'Q')
		case 0x52:
			dst = append(dst, 'R')
		case 0x53:
			dst = append(dst, 'S')
		case 0x54:
			dst = append(dst, 'T')
		case 0x55:
			dst = append(dst, 'U')
		case 0x56:
			dst = append(dst, 'V')
		case 0x57:
			dst = append(dst, 'W')
		case 0x58:
			dst = append(dst, 'X')
		case 0x59:
			dst = append(dst, 'Y')
		case 0x5a:
			dst = append(dst, 'Z')
		/////////////
		case 0x61:
			dst = append(dst, 'a')
		case 0x62:
			dst = append(dst, 'b')
		case 0x63:
			dst = append(dst, 'c')
		case 0x64:
			dst = append(dst, 'd')
		case 0x65:
			dst = append(dst, 'e')
		case 0x66:
			dst = append(dst, 'f')
		case 0x67:
			dst = append(dst, 'g')
		case 0x68:
			dst = append(dst, 'h')
		case 0x69:
			dst = append(dst, 'i')
		case 0x6a:
			dst = append(dst, 'j')
		case 0x6b:
			dst = append(dst, 'k')
		case 0x6c:
			dst = append(dst, 'l')
		case 0x6d:
			dst = append(dst, 'm')
		case 0x6e:
			dst = append(dst, 'n')
		case 0x6f:
			dst = append(dst, 'o')
		case 0x70:
			dst = append(dst, 'p')
		case 0x71:
			dst = append(dst, 'q')
		case 0x72:
			dst = append(dst, 'r')
		case 0x73:
			dst = append(dst, 's')
		case 0x74:
			dst = append(dst, 't')
		case 0x75:
			dst = append(dst, 'u')
		case 0x76:
			dst = append(dst, 'v')
		case 0x77:
			dst = append(dst, 'w')
		case 0x78:
			dst = append(dst, 'x')
		case 0x79:
			dst = append(dst, 'y')
		case 0x7a:
			dst = append(dst, 'z')
		/////////////
		case 0xe0:
			dst = append(dst, 'а')
		case 0xe1:
			dst = append(dst, 'б')
		case 0xe2:
			dst = append(dst, 'в')
		case 0xe3:
			dst = append(dst, 'г')
		case 0xe4:
			dst = append(dst, 'д')
		case 0xe5:
			dst = append(dst, 'е')
		case 0xe6:
			dst = append(dst, 'ж')
		case 0xe7:
			dst = append(dst, 'з')
		case 0xe8:
			dst = append(dst, 'и')
		case 0xe9:
			dst = append(dst, 'й')
		case 0xea:
			dst = append(dst, 'к')
		case 0xeb:
			dst = append(dst, 'л')
		case 0xec:
			dst = append(dst, 'м')
		case 0xed:
			dst = append(dst, 'н')
		case 0xee:
			dst = append(dst, 'о')
		case 0xef:
			dst = append(dst, 'п')
		case 0xf0:
			dst = append(dst, 'р')
		case 0xf1:
			dst = append(dst, 'с')
		case 0xf2:
			dst = append(dst, 'т')
		case 0xf3:
			dst = append(dst, 'у')
		case 0xf4:
			dst = append(dst, 'ф')
		case 0xf5:
			dst = append(dst, 'х')
		case 0xf6:
			dst = append(dst, 'ц')
		case 0xf7:
			dst = append(dst, 'ч')
		case 0xf8:
			dst = append(dst, 'ш')
		case 0xf9:
			dst = append(dst, 'щ')
		case 0xfa:
			dst = append(dst, 'ь')
		case 0xfb:
			dst = append(dst, 'ы')
		case 0xfc:
			dst = append(dst, 'ь') /// Был знак 'ъ'
		case 0xfd:
			dst = append(dst, 'э')
		case 0xfe:
			dst = append(dst, 'ю')
		case 0xff:
			dst = append(dst, 'я')
		//////////////
		case 0xc0:
			dst = append(dst, 'А')
		case 0xc1:
			dst = append(dst, 'Б')
		case 0xc2:
			dst = append(dst, 'В')
		case 0xc3:
			dst = append(dst, 'Г')
		case 0xc4:
			dst = append(dst, 'Д')
		case 0xc5:
			dst = append(dst, 'Е')
		case 0xc6:
			dst = append(dst, 'Ж')
		case 0xc7:
			dst = append(dst, 'З')
		case 0xc8:
			dst = append(dst, 'И')
		case 0xc9:
			dst = append(dst, 'Й')
		case 0xca:
			dst = append(dst, 'К')
		case 0xcb:
			dst = append(dst, 'Л')
		case 0xcc:
			dst = append(dst, 'М')
		case 0xcd:
			dst = append(dst, 'Н')
		case 0xce:
			dst = append(dst, 'О')
		case 0xcf:
			dst = append(dst, 'П')
		case 0xd0:
			dst = append(dst, 'Р')
		case 0xd1:
			dst = append(dst, 'С')
		case 0xd2:
			dst = append(dst, 'Т')
		case 0xd3:
			dst = append(dst, 'У')
		case 0xd4:
			dst = append(dst, 'Ф')
		case 0xd5:
			dst = append(dst, 'Х')
		case 0xd6:
			dst = append(dst, 'Ц')
		case 0xd7:
			dst = append(dst, 'Ч')
		case 0xd8:
			dst = append(dst, 'Ш')
		case 0xd9:
			dst = append(dst, 'Щ')
		case 0xda:
			dst = append(dst, 'Ь')
		case 0xdb:
			dst = append(dst, 'Ы')
		case 0xdc:
			dst = append(dst, 'Ъ')
		case 0xdd:
			dst = append(dst, 'Э')
		case 0xde:
			dst = append(dst, 'Ю')
		case 0xdf:
			dst = append(dst, 'Я')
		case 0x20:
			dst = append(dst, ' ')
		case 0x28:
			dst = append(dst, '(')
		case 0x29:
			dst = append(dst, ')')
		default:
			return "", errors.New("unexpected character found")
		}
	}
	return string(dst), nil
}
