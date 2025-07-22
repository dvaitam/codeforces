package main

import (
	"bufio"
	"fmt"
	"os"
)

type row struct{ a, b, c byte }

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(in, &s)

	rows := []row{{'.', '.', '.'}}          // ensure top neighbor of first row is '.'
	rows = append(rows, row{'.', 'X', 'X'}) // move pointer to cell1
	cur := 0
	for i := 0; i < len(s); i++ {
		target := int(s[i])
		diff := (cur - target + 256) % 256
		for j := 0; j < diff; j++ {
			rows = append(rows, row{'.', '.', '.'}) // top '.' for '-'
			rows = append(rows, row{'.', 'X', '.'}) // '-'
		}
		rows = append(rows, row{'.', 'X', '.'}) // '.'
		cur = target
	}
	rows = append(rows, row{'.', '.', '.'}) // closing blank row for toroidal topology

	w := bufio.NewWriter(os.Stdout)
	for _, r := range rows {
		fmt.Fprintf(w, "%c%c%c\n", r.a, r.b, r.c)
	}
	w.Flush()
}
