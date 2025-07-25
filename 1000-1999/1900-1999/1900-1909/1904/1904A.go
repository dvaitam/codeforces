package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct{ x, y int64 }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var a, b int64
		fmt.Fscan(reader, &a, &b)
		var xk, yk int64
		fmt.Fscan(reader, &xk, &yk)
		var xq, yq int64
		fmt.Fscan(reader, &xq, &yq)

		moves := [][2]int64{{a, b}, {a, -b}, {-a, b}, {-a, -b}, {b, a}, {b, -a}, {-b, a}, {-b, -a}}
		posK := make(map[pair]struct{})
		for _, m := range moves {
			posK[pair{xk + m[0], yk + m[1]}] = struct{}{}
		}
		posQ := make(map[pair]struct{})
		for _, m := range moves {
			posQ[pair{xq + m[0], yq + m[1]}] = struct{}{}
		}
		cnt := 0
		for p := range posK {
			if _, ok := posQ[p]; ok {
				cnt++
			}
		}
		fmt.Fprintln(writer, cnt)
	}
}
