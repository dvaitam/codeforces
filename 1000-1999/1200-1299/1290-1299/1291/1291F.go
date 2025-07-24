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

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	alive := make([]int, n)
	for i := range alive {
		alive[i] = i + 1
	}

	passes := 2 * n / k
	if passes < 1 {
		passes = 1
	}

	rand.Seed(time.Now().UnixNano())

	query := func(idx int) bool {
		fmt.Fprintf(writer, "? %d\n", idx)
		writer.Flush()
		var resp string
		if _, err := fmt.Fscan(reader, &resp); err != nil {
			os.Exit(0)
		}
		return resp == "Y"
	}

	reset := func() {
		fmt.Fprintln(writer, "R")
		writer.Flush()
	}

	for p := 0; p < passes && len(alive) > 0; p++ {
		rand.Shuffle(len(alive), func(i, j int) { alive[i], alive[j] = alive[j], alive[i] })
		newAlive := make([]int, 0, len(alive))
		for i, idx := range alive {
			if i%k == 0 {
				reset()
			}
			if !query(idx) {
				newAlive = append(newAlive, idx)
			}
		}
		alive = newAlive
	}

	fmt.Fprintf(writer, "! %d\n", len(alive))
}
