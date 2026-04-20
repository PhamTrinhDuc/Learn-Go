package main

import "fmt"

type BankAccount struct {
	Owner   string
	Balance float64
}

func (b *BankAccount) Deposit(amount float64) {
	if amount > 0 {
		b.Balance += amount
	}
}

func (b *BankAccount) Withdraw(amount float64) bool {
	if amount <= 0 || b.Balance < amount {
		return false
	}
	b.Balance -= amount
	return true
}

func Transfer(b1 *BankAccount, b2 *BankAccount, amount float64) {
	if b1.Withdraw(amount) {
		b2.Deposit(amount)
	} else {
		fmt.Println("Transfer failed. Not enough money")
	}

}
func (b BankAccount) Display() {
	fmt.Println("Owner: ", b.Owner)
	fmt.Println("Balance: ", b.Balance)
}

func main() {
	b := BankAccount{Owner: "John", Balance: 1000}
	b.Deposit(100)
	b.Withdraw(200)
	b.Display()
	
	b2 := BankAccount{Owner: "Jane", Balance: 2000}
	Transfer(&b, &b2, 500)
	b.Display()
	b2.Display()
}
