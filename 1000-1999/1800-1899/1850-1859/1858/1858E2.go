package main

import (
	"bufio"
	"fmt"
	"os"
)

// This solution maintains the current array as a stack and keeps
// track of the number of distinct elements using a frequency array.
// To support rollbacks, every change (+ or - query) is stored in a
// history stack so it can be undone when a '!' query appears.

// operation represents a change to the array that can be rolled back.
type operation struct {
	typ  byte  // '+' for push, '-' for pop
	vals []int // affected values
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}

	const maxVal = 1000000
	freq := make([]int, maxVal+1)
	distinct := 0

	arr := make([]int, 0)
	history := make([]operation, 0)

	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(reader, &op)
		switch op {
		case "+":
			var x int
			fmt.Fscan(reader, &x)
			arr = append(arr, x)
			history = append(history, operation{typ: '+', vals: []int{x}})
			if freq[x] == 0 {
				distinct++
			}
			freq[x]++
		case "-":
			var k int
			fmt.Fscan(reader, &k)
			removed := make([]int, k)
			for j := k - 1; j >= 0; j-- {
				v := arr[len(arr)-1]
				arr = arr[:len(arr)-1]
				removed[j] = v
				freq[v]--
				if freq[v] == 0 {
					distinct--
				}
			}
			history = append(history, operation{typ: '-', vals: removed})
		case "!":
			if len(history) == 0 {
				continue
			}
			last := history[len(history)-1]
			history = history[:len(history)-1]
			if last.typ == '+' {
				v := last.vals[0]
				arr = arr[:len(arr)-1]
				freq[v]--
				if freq[v] == 0 {
					distinct--
				}
			} else { // '-'
				for _, v := range last.vals {
					arr = append(arr, v)
					if freq[v] == 0 {
						distinct++
					}
					freq[v]++
				}
			}
		case "?":
			fmt.Fprintln(writer, distinct)
			writer.Flush() // flush to satisfy online requirement
		}
	}
}
