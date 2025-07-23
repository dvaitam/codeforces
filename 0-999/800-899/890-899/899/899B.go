package main

import (
	"bufio"
	"fmt"
	"os"
)

func isLeap(year int) bool {
	if year%400 == 0 {
		return true
	}
	if year%100 == 0 {
		return false
	}
	return year%4 == 0
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	// build months for a range of years covering at least two 400-year cycles
	months := make([]int, 0, 9600)
	for y := 2000; y < 2800; y++ {
		feb := 28
		if isLeap(y) {
			feb = 29
		}
		months = append(months, 31, feb, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31)
	}

	found := false
	for i := 0; i+len(a) <= len(months); i++ {
		match := true
		for j := 0; j < n; j++ {
			if months[i+j] != a[j] {
				match = false
				break
			}
		}
		if match {
			found = true
			break
		}
	}

	if found {
		fmt.Fprintln(out, "YES")
	} else {
		fmt.Fprintln(out, "NO")
	}
}
