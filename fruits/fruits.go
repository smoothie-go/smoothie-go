package fruits

import (
	"math/rand/v2"
)

var fruits = []string{"Cucumber", "Durian", "Date", "Eggplant", "Fig",
	"Grape", "Guava", "Honeydew", "Kiwi", "Lemon",
	"Lime", "Lychee", "Mango", "Mirabelle", "Olive",
	"Orange", "Papaya", "Passion", "Peach", "Pear",
	"Pineapple", "Pitaya", "Plum", "Pomelo", "Quince",
	"Raspberry", "Starfruit", "Strawberry", "Tomato", "Watermelon"}

func GetRandomFruit() string {
	return fruits[rand.IntN(len(fruits))]
}
