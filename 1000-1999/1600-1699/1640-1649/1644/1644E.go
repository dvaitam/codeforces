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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for i := 0; i < t; i++ {
		var n int64
		var s string
		if _, err := fmt.Fscan(reader, &n, &s); err != nil {
			break
		}
		solve(n, s, writer)
	}
}

func solve(n int64, s string, writer *bufio.Writer) {
	cntD := 0
	cntR := 0
	kD := -1
	kR := -1

	for i, r := range s {
		if r == 'D' {
			cntD++
			if kD == -1 {
				kD = i
			}
		} else {
			cntR++
			if kR == -1 {
				kR = i
			}
		}
	}

	if cntD == 0 || cntR == 0 {
		fmt.Fprintln(writer, n)
		return
	}

	Sx := n - int64(cntD) - 1
	Sy := n - int64(cntR) - 1

	var ans int64

	if kD < kR {
		// Starts with D (Vertical moves first)
		// k is number of D moves before the first R
		k := int64(kR)
		
		// Vertical slack and Horizontal slack
		Sv := Sx
		Sh := Sy
		
		// Base area covers the prefix of D's and the transition R
		ans = k + (Sv+1)*(Sh+2)
		
		// Remaining moves
		remD := int64(cntD) - k
		remR := int64(cntR) - 1
		
		// Each remaining D adds a slice of width 1 and height (Sh + 1)
		ans += remD * (Sh + 1)
		// Each remaining R adds a slice of height 1 and width (Sv + 1)
		ans += remR * (Sv + 1)
	} else {
		// Starts with R (Horizontal moves first)
		// k is number of R moves before the first D
		k := int64(kD)
		
		Sv := Sy
		Sh := Sx
		
		ans = k + (Sv+1)*(Sh+2)
		
		remR := int64(cntR) - k
		remD := int64(cntD) - 1
		
		ans += remR * (Sh + 1)
		ans += remD * (Sv + 1)
	}

	fmt.Fprintln(writer, ans)
}