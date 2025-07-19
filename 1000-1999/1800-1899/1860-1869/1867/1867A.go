package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

func readInt() (int, error) {
	var x int
	var neg bool
	// read first non-space
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}
		if b == '-' {
			neg = true
			break
		}
		if b >= '0' && b <= '9' {
			x = int(b - '0')
			break
		}
	}
	// read rest digits
	for {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}
		if b < '0' || b > '9' {
			break
		}
		x = x*10 + int(b-'0')
	}
	if neg {
		x = -x
	}
	return x, nil
}

func main() {
	defer writer.Flush()
	t, err := readInt()
	if err != nil {
		return
	}
	for ti := 0; ti < t; ti++ {
		n, _ := readInt()
		qs := make([]struct{ val, id int }, n)
		for i := 0; i < n; i++ {
			v, _ := readInt()
			qs[i].val = v
			qs[i].id = i
		}
		// sort by val descending
		sort.Slice(qs, func(i, j int) bool {
			return qs[i].val > qs[j].val
		})
		// assign ranks
		for i := 0; i < n; i++ {
			qs[i].val = i + 1
		}
		// restore original order
		sort.Slice(qs, func(i, j int) bool {
			return qs[i].id < qs[j].id
		})
		// output
		for i := 0; i < n; i++ {
			if i > 0 {
				writer.WriteByte(' ')
			}
			writer.WriteString(fmt.Sprintf("%d", qs[i].val))
		}
		writer.WriteByte('\n')
	}
}
