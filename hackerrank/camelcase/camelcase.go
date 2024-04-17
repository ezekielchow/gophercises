package main

import (
	"fmt"
	"strings"
)

func main() {

	s := "頁設是"

	noOfWords := 1

	for k, v := range s {
		fmt.Println(v, k)

		c := fmt.Sprintf("%c", v)

		upper := strings.ToUpper(c)

		if strings.Compare(upper, c) == 0 {
			noOfWords += 1
		}
	}

	fmt.Println("amount of words:", noOfWords)
}
