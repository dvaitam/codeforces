package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		v := make([]uint64, n)
		for i := 0; i < m; i++ {
			var a, b int
			fmt.Fscan(reader, &a, &b)
			k := rnd.Uint64()
			v[a-1] ^= k
			v[b-1] ^= k
		}
		count := make(map[uint64]int)
		var prefix uint64
		best := 0
		for _, r := range v {
			prefix ^= r
			count[prefix]++
			if count[prefix] > best {
				best = count[prefix]
			}
		}
		fmt.Fprintln(writer, n-best)
	}
}
