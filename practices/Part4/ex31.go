package main

import "fmt"

type Rule func(price float64) float64

func discount10(price float64) float64 {
	return price * 0.1
}

func discount20(price float64) float64 {
	return price * 0.2
}

func ComposeRules(mode string, rules ...Rule) Rule {
	return func(price float64) float64 {
		if len(rules) == 0 {
			return price
		}

		price_discounted := make([]float64, len(rules))
		for i, rule := range rules {
			res := rule(price)
			if res < 0 {
				res = 0
			}
			price_discounted[i] = res
		}

		finalPrice := price_discounted[0]
		if mode == "MIN" {
			for _, value := range price_discounted {
				if float64(value) < finalPrice {
					finalPrice = float64(value)
				}
			}
		} else if mode == "MAX" {
			for value := range price_discounted {
				if float64(value) > finalPrice {
					finalPrice = float64(value)
				}
			}
		} else {
			curr_price := price
			for _, rule := range rules {
				curr_price = rule(curr_price)
			}
			finalPrice = curr_price
		}

		if finalPrice < 0 {
			return 0
		}
		return finalPrice
	}
}

func main() {
	var r1 Rule = discount10
	var r2 Rule = discount20

	price := ComposeRules("MIN", r1, r2)(100)
	fmt.Println(price)
}
