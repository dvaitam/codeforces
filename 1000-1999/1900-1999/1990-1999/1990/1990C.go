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

	const maxN = 200005

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int64
		fmt.Fscan(reader, &n)
		var b, c int64
		sum := int64(0)
		p1 := make([]bool, maxN)
		p2 := make([]bool, maxN)
		for n > 0 {
			n--
			var x int64
			fmt.Fscan(reader, &x)
			if p1[x] {
				if x > b {
					b = x
				}
			} else {
				p1[x] = true
			}
			if p2[b] {
				if b > c {
					c = b
				}
			} else if b != 0 {
				p2[b] = true
			}
			sum += x + b + (n+1)*c
		}
		fmt.Fprintln(writer, sum)
	}
}
