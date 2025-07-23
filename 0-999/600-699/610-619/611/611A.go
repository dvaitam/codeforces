package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// It calculates how many times a certain weekday or day of the month
// occurs in the year 2016 (a leap year starting on Friday).
func main() {
	reader := bufio.NewReader(os.Stdin)
	var x int
	var of string
	var typ string
	if _, err := fmt.Fscan(reader, &x, &of, &typ); err != nil {
		return
	}

	if typ == "week" {
		if x == 5 || x == 6 {
			fmt.Println(53)
		} else {
			fmt.Println(52)
		}
		return
	}

	// typ == "month"
	days := []int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	cnt := 0
	for _, d := range days {
		if x <= d {
			cnt++
		}
	}
	fmt.Println(cnt)
}
