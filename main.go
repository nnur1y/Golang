package main

import "fmt"

func main() {

	var iPhone14, iPhone13, iPhoneSE int = 799, 599, 429

	var m3, m6, m12, m24 int = 3, 6, 12, 24

	fmt.Println("iPhone prices:")
	fmt.Println("1. iPhone 14 - $", iPhone14)
	fmt.Println("2. iPhone 13 - $", iPhone13)
	fmt.Println("3. iPhone SE - $", iPhoneSE)

	fmt.Println("Напишите номер выбранного iPhone (1,2,3): ")
	var number int
	fmt.Scanln(&number)

	price := 0
	phone := ""

	switch number {
	case 1:
		phone = "iPhone 14"
		price = iPhone14
	case 2:
		phone = "iPhone 13"
		price = iPhone13
	case 3:
		phone = "iPhone SE"
		price = iPhoneSE
	default:
		fmt.Println("Choose right number")

	}

	fmt.Println("Выберите, сколько месяцев вы хотите платить в рассрочку")
	fmt.Println(m3, m6, m12, m24)
	var month int
	fmt.Scanln(&month)

	fmt.Println("You have chosen a ", phone, " which costs $", price)

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

	fmt.Println("You chose to pay in installments for", month, " months, it will be $", price, " per month")

}
