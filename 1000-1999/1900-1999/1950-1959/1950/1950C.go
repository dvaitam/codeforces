package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		hour, _ := strconv.Atoi(s[:2])
		minute := s[3:]
		if hour < 12 {
			if hour == 0 {
				hour = 12
			}
			fmt.Fprintf(out, "%02d:%s AM\n", hour, minute)
		} else {
			if hour > 12 {
				hour -= 12
			}
			fmt.Fprintf(out, "%02d:%s PM\n", hour, minute)
		}
	}
}
