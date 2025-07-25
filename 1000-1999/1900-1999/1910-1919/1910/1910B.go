package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		firstMinus := -1
		lastPlus := -1
		bal := 0
		minPref := 0
		for i, ch := range s {
			if ch == '+' {
				bal++
				lastPlus = i
			} else {
				bal--
				if firstMinus == -1 {
					firstMinus = i
				}
			}
			if bal < minPref {
				minPref = bal
			}
		}
		if bal < 0 {
			fmt.Fprintln(writer, -1)
			continue
		}
		if minPref >= 0 {
			fmt.Fprintln(writer, 1, 1)
			continue
		}
		if minPref < -2 || firstMinus == -1 || lastPlus == -1 || firstMinus >= lastPlus {
			fmt.Fprintln(writer, -1)
			continue
		}
		arr := []rune(s)
		arr[firstMinus], arr[lastPlus] = arr[lastPlus], arr[firstMinus]
		bal = 0
		ok := true
		for _, ch := range arr {
			if ch == '+' {
				bal++
			} else {
				bal--
			}
			if bal < 0 {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(writer, firstMinus+1, lastPlus+1)
		} else {
			fmt.Fprintln(writer, -1)
		}
	}
}
