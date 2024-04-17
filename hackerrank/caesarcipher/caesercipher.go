package main

import "fmt"

func main() {
	str := "i-will_be$hiddenz"
	rotateBy := 1

	min := int('a')
	max := int('z')
	minCapital := int('A')
	maxCapital := int('Z')

	encrypted := ""

	for _, v := range str {

		letter := fmt.Sprintf("%c", v)
		letterRune := int(v)

		if !(letterRune >= min && letterRune <= max) && !(letterRune >= minCapital && letterRune <= maxCapital) {
			encrypted = encrypted + letter
			continue
		}

		newRune := letterRune + rotateBy

		if newRune > maxCapital && newRune < min {
			encrypted = encrypted + string(rune(newRune-26))
			continue
		}

		if newRune > max {
			steps := (newRune - max) % 26

			if steps > 0 {
				newRune = min + steps - 1
			} else {
				newRune = max
			}
		}

		encrypted = encrypted + string(rune(newRune))
	}

	fmt.Println("O:", str, "\nE:", encrypted)
}
