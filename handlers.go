package main

import (
	"errors"
	"hash/crc32"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

func addToArrayIfNeeded(src []string, add []string) []string {
	for _, addItem := range add {
		if !isArrayContainsString(src, addItem) {
			src = append(src, addItem)
		}
	}
	return src
}

func isArrayContainsString(array []string, str string) bool {
	for _, line := range array {
		if strings.Compare(line, str) == 0 {
			return true
		}
	}
	return false
}

func randomizeNumber(number float64) float64 {
	randomized := rand.Int()
	delta := float64(randomized%41-20) / 40
	return number + delta
}

func isMale(Name string) bool {
	var nameArr = []string{
		"Андрей", "Антон", "Артур", "Александр", "Алексей", "Анатолий", "Анвар", "Аркадий", "Артем", "Ахмад", "Ашот",
		"Богдан", "Борис",
		"Вадим", "Валерий", "Василий", "Виктор", "Виталий", "Владимир", "Владислав", "Вячеслав",
		"Геннадий", "Георгий", "Глеб", "Григорий", "Гриша",
		"Давид", "Данил", "Даниил", "Денис", "Джамал", "Дмитрий",
		"Евгений", "Егор",
		"Иван", "Игорь", "Илья",
		"Карен", "Кирилл", "Кирил", "Константин", "Костя",
		"Максим", "Михаил", "Мурад", "Макс",
		"Никита", "Николай", "Коля",
		"Олег",
		"Павел",
		"Рашид", "Рашит", "Ренат", "Рома", "Роман", "Руслан", "Рустам",
		"Святослав", "Славик", "Слава", "Сергей", "Серега", "Станислав",
		"Тимур", "Тимофей", "Тарас",
		"Фёдор",
		"Юрий",
		"Ярослав",
	}

	for _, name := range nameArr {
		if strings.Compare(name, Name) == 0 {
			return true
		}
	}
	return false
}

func isFemale(Name string) bool {
	var nameArr = []string{
		"Анна", "Аня", "Алёна", "Ася", "Алла", "Альбина", "Анастасия", "Настя", "Алина", "Алеся",
		"Валентина", "Валя", "Варвара", "Вера", "Василина", "Вероника", "Виктория", "Вита",
		"Галина", "Галя",
		"Дарья", "Даша",
		"Ева", "Екатерина", "Катя", "Елена", "Лена", "Евгения",
		"Зоя",
		"Ира", "Инна", "Ирина",
		"Карина", "Ксения", "Ксюша",
		"Леся", "Лада", "Лиза", "Лилия", "Лиля", "Люся", "Любовь", "Люба", "Лариса",
		"Марина", "Маргарита", "Рита", "Мария", "Маша",
		"Надежда", "Надя", "Нина", "Наталья",
		"Оксана", "Ксюша", "Ольга", "Оля", "Олеся",
		"Полина",
		"Раиса", "Рая", "Руслана",
		"Светлана", "Софья", "София",
		"Тамара", "Тома", "Тина", "Татьяна",
		"Элиза", "Эля",
		"Юлия", "Юля", "Юлианна",
		"Ярослава", "Яна",
	}

	for _, name := range nameArr {
		if strings.Compare(name, Name) == 0 {
			return true
		}
	}
	return false
}

func runGet() string {
	resp, err := http.Get(urlString)
	if err != nil {
		println("Error: " + err.Error())
		return ""
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println("Error: " + err.Error())
		return ""
	}

	return string(respBody)
}

func runGetWithHeaders(vkUID int) ([]byte, error) {
	var err error

	request := &http.Request{
		Method: http.MethodGet,
		Header: http.Header{
			"User-Agent":     {"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36"},
			"Accept":         {"image/avif,image/webp,image/apng,image/*,*/*;q=0.8"},
			"Connection":     {"keep-alive"},
			"Referer":        {"http://localhost:8080/"},
			"Sec-Fetch-Dest": {"image"},
			"Sec-Fetch-Mode": {"no-cors"},
			"Sec-Fetch-Site": {"same-origin"},
		},
	}

	request.URL, err = url.Parse(urlString + strconv.Itoa(vkUID))
	if err != nil {
		return nil, errors.New("url parse error")
	}
	// request.URL.Query().Set("user", "Denis')

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.New("request vk error")
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("response body error")
	}

	return respBody, nil
}

func writeToFile(src []byte) {
	file, err := os.Create(filename)
	if err != nil {
		println("Unable to create file: " + err.Error())
		return
	}
	defer file.Close()

	file.Write(src)

	println("Done.")
}

func printParseSuccess(userNumber int, vkUID int, user User) {
	var OK = GREEN + "+" + NO_COLOR
	var KO = RED + "-" + NO_COLOR
	var dst = strconv.Itoa(userNumber) + "\t" + strconv.Itoa(vkUID) + "\t"

	dst += user.Fname + " " + user.Lname + "\t"

	if utf8.RuneCountInString(user.Fname+" "+user.Lname) < 16 {
		dst += "\t"
	}

	dst += OK + OK

	if user.Avatar != "" {
		dst += OK
	} else {
		dst += KO
	}

	if user.Birth != nil {
		dst += OK
	} else {
		dst += KO
	}

	if user.Latitude != nil && user.Longitude != nil {
		dst += OK
	} else {
		dst += KO
	}

	dst += "\t" + user.Gender + "\t" + user.Orientation
	println(dst)
}

func printParseFail(vkUID int, err error) {
	var dst = "\t" + strconv.Itoa(vkUID) + "\t" + err.Error()
	println(dst)
}

func PassHash(pass string) string {
	pass += passSalt
	crcH := crc32.ChecksumIEEE([]byte(pass))
	return strconv.FormatUint(uint64(crcH), 20)
}
