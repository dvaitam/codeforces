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
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		aWhite, aBlack := 0, 0
		bWhite, bBlack := 0, 0
		pos := 1
		step := 1
		remaining := n
		for remaining > 0 {
			take := step
			if take > remaining {
				take = remaining
			}
			l := pos
			r := pos + take - 1
			white := (r+1)/2 - (l)/2
			black := take - white
			// determine current player
			var alice bool
			if step == 1 {
				alice = true
			} else {
				pair := (step - 2) / 2
				if pair%2 == 1 {
					alice = true
				} else {
					alice = false
				}
			}
			if alice {
				aWhite += white
				aBlack += black
			} else {
				bWhite += white
				bBlack += black
			}
			pos += take
			remaining -= take
			step++
		}
		fmt.Fprintf(writer, "%d %d %d %d\n", aWhite, aBlack, bWhite, bBlack)
	}
}
