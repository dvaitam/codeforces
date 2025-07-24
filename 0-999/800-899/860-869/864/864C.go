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

	var a, b, f, k int
	if _, err := fmt.Fscan(reader, &a, &b, &f, &k); err != nil {
		return
	}

	fuel := b
	refuels := 0

	for i := 1; i <= k; i++ {
		if i%2 == 1 { // 0 -> a
			// move from start to station
			if fuel < f {
				fmt.Fprintln(writer, -1)
				return
			}
			fuel -= f
			// distance needed after station
			var dist int
			if i == k {
				dist = a - f
			} else {
				dist = 2 * (a - f)
			}
			if b < dist {
				fmt.Fprintln(writer, -1)
				return
			}
			if fuel < dist {
				refuels++
				fuel = b
			}
			fuel -= a - f
		} else { // a -> 0
			if fuel < a-f {
				fmt.Fprintln(writer, -1)
				return
			}
			fuel -= a - f
			var dist int
			if i == k {
				dist = f
			} else {
				dist = 2 * f
			}
			if b < dist {
				fmt.Fprintln(writer, -1)
				return
			}
			if fuel < dist {
				refuels++
				fuel = b
			}
			fuel -= f
		}
	}

	fmt.Fprintln(writer, refuels)
}
