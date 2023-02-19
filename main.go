package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type User struct {
	Nickname string
	Cooker   bool
	Age      uint16
	Gender   string
	Recipes  []string
}

func (u User) getAllInfo() string {
	return fmt.Sprintf("Nickname is: %s. Age is %d, and gender is %s", u.Nickname, u.Age, u.Gender)

}

func home(w http.ResponseWriter, r *http.Request) {
	data1 := User{"Nurly", true, 19, "Female", []string{"Deserts", "Rollton salad with sausage", "Waffles"}}
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, data1)
}
func account(w http.ResponseWriter, r *http.Request) {
	data1 := User{"Nurly", true, 19, "Female", []string{"Deserts", "Rollton salad with sausage", "Waffles"}}
	tmpl, _ := template.ParseFiles("templates/account.html")
	tmpl.Execute(w, data1)
}

func login_page(w http.ResponseWriter, r *http.Request) {
	data1 := User{"Nurly", true, 19, "Female", []string{"Deserts", "Rollton salad with sausage", "Waffles"}}
	tmpl, _ := template.ParseFiles("templates/login_page.html")
	tmpl.Execute(w, data1)
}

func handleRequest() {
	http.HandleFunc("/index/", home)
	http.HandleFunc("/", account)
	http.HandleFunc("/loginpage/", login_page)
	http.ListenAndServe(":8080", nil)

}
func main() {

	handleRequest()

	//installments()
}

func installments() {
	var basic, standard, curator int = 50000, 90000, 130000

	var m3, m6, m12, m24 int = 3, 6, 12, 24

	fmt.Println("iPhone prices:")
	fmt.Println("1. Basic package -", basic)
	fmt.Println("2. Standard package -", standard)
	fmt.Println("3. Curator package -", curator)

	fmt.Println("Напишите номер выбранного пакета (1,2,3):")
	var number int
	fmt.Scanln(&number)

	price := 0
	pack := ""

	switch number {
	case 1:
		pack = "Basic package"
		price = basic
	case 2:
		pack = "Standard package"
		price = standard
	case 3:
		pack = "Curator package"
		price = curator
	default:
		fmt.Println("Choose right number")

	}

	fmt.Println("Выберите, сколько месяцев вы хотите платить в рассрочку")
	fmt.Println(m3, m6, m12, m24)
	var month int
	fmt.Scanln(&month)

	fmt.Println("You have chosen a", pack, " package which costs", price)

	switch month {
	case 3:
		price = price / 3
	case 6:
		price = price / 6
	case 12:
		price = price / 12
	case 24:
		price = price / 24
	default:
		fmt.Println("Choose right number")

	}

	fmt.Println("You chose to pay in installments for", month, "months, it will be", price, "tg per month")

}
