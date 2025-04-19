package utils

import (
	"fmt"
	"math"
)

func DepCalc() {
	var amount float64
	var months int
	var capitalized bool

	fmt.Println("=== Депозитный калькулятор «Тапочек» (Активбанк) ===")
	fmt.Print("Введите сумму вклада (сомони): ")
	fmt.Scanln(&amount)

	fmt.Print("Введите срок вклада (в месяцах, от 3 до 24): ")
	fmt.Scanln(&months)

	if months < 3 || months > 24 {
		fmt.Println("Ошибка: срок должен быть от 3 до 24 месяцев.")
		return
	}

	fmt.Print("Капитализировать проценты? (1 — да, 0 — нет): ")
	var capInput int
	fmt.Scanln(&capInput)
	capitalized = capInput == 1

	// Получаем процентную ставку
	rate := getInterestRate(months)

	if capitalized {
		finalAmount := calculateWithCapitalization(amount, rate, months)
		income := finalAmount - amount
		tax := income * 0.14
		netIncome := income - tax
		netTotal := amount + netIncome

		fmt.Printf("\nСумма с капитализацией: %.2f сомони\n", finalAmount)
		fmt.Printf("Доход до налогообложения: %.2f сомони\n", income)
		fmt.Printf("Налог (14%%): %.2f сомони\n", tax)
		fmt.Printf("Чистый доход: %.2f сомони\n", netIncome)
		fmt.Printf("Общая сумма к получению: %.2f сомони\n", netTotal)

	} else {
		monthlyIncome := amount * (rate / 100 / 12)
		totalIncome := monthlyIncome * float64(months)
		tax := totalIncome * 0.14
		netIncome := totalIncome - tax
		netTotal := amount + netIncome

		fmt.Printf("\nПроценты выплачиваются ежемесячно.\n")
		fmt.Printf("Доход до налогообложения: %.2f сомони\n", totalIncome)
		fmt.Printf("Налог (14%%): %.2f сомони\n", tax)
		fmt.Printf("Чистый доход: %.2f сомони\n", netIncome)
		fmt.Printf("Общая сумма к получению: %.2f сомони\n", netTotal)
	}
}

// Определение процентной ставки по сроку
func getInterestRate(months int) float64 {
	switch {
	case months >= 3 && months <= 3:
		return 8.0
	case months >= 4 && months <= 6:
		return 10.0
	case months >= 7 && months <= 12:
		return 14.0
	case months >= 13 && months <= 18:
		return 16.0
	case months >= 19 && months <= 24:
		return 17.0
	default:
		return 0.0
	}
}

// Расчет с капитализацией процентов
func calculateWithCapitalization(principal float64, annualRate float64, months int) float64 {
	monthlyRate := annualRate / 12 / 100
	return principal * math.Pow(1+monthlyRate, float64(months))
}
