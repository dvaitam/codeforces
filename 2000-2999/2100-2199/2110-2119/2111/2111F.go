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
	const limit = 50000
	for ; t > 0; t-- {
		var p, s int64
		fmt.Fscan(reader, &p, &s)
		if p > 4*s {
			fmt.Fprintln(writer, -1)
			continue
		}
		if p == 4*s {
			fmt.Fprintln(writer, 1)
			fmt.Fprintln(writer, "0 0")
			continue
		}
		found := false
		var wFound, hFound int64
		for w := int64(1); w <= limit && !found; w++ {
			den := p*w - 2*s
			if den <= 0 {
				continue
			}
			num := 2 * s * w
			if num%den != 0 {
				continue
			}
			h := num / den
			if h <= 0 || w*h > limit {
				continue
			}
			found = true
			wFound = w
			hFound = h
		}
		if !found {
			fmt.Fprintln(writer, -1)
			continue
		}
		area := wFound * hFound
		fmt.Fprintln(writer, area)
		for y := int64(0); y < hFound; y++ {
			for x := int64(0); x < wFound; x++ {
				fmt.Fprintf(writer, "%d %d\n", x, y)
			}
		}
	}
}
