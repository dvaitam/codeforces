package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	// read number of nodes and k-ary parameter
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	if n <= 1 {
		fmt.Fprintf(writer, "! 1\n")
		return
	}
	// votes counts positive path checks for each node
	votes := make([]int, n+1)
	t := 50 // number of samples per node
	for b := 1; b <= n; b++ {
		for i := 0; i < t; i++ {
			// pick random a, c distinct from b and each other
			a := rand.Intn(n) + 1
			for a == b {
				a = rand.Intn(n) + 1
			}
			c := rand.Intn(n) + 1
			for c == b || c == a {
				c = rand.Intn(n) + 1
			}
			// query if b is on path a-c
			fmt.Fprintf(writer, "? %d %d %d\n", a, b, c)
			writer.Flush()
			resp, _ := reader.ReadString('\n')
			resp = strings.TrimSpace(resp)
			if len(resp) > 0 && (resp[0] == 'Y' || resp[0] == 'y') {
				votes[b]++
			}
		}
	}
	// find node with maximum votes
	root := 1
	maxv := votes[1]
	for i := 2; i <= n; i++ {
		if votes[i] > maxv {
			maxv = votes[i]
			root = i
		}
	}
	// report root
	fmt.Fprintf(writer, "! %d\n", root)
}
