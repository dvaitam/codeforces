package main

import (
	"bufio"
	"fmt"
	"os"
)

type Op struct {
	typ     byte
	val     int
	removed []int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}

	arr := make([]int, 0, q)
	freq := make([]int, 1000001)
	distinct := 0

	history := make([]Op, 0, q)

	for i := 0; i < q; i++ {
		var cmd string
		if _, err := fmt.Fscan(reader, &cmd); err != nil {
			return
		}
		switch cmd[0] {
		case '+':
			var x int
			fmt.Fscan(reader, &x)
			arr = append(arr, x)
			freq[x]++
			if freq[x] == 1 {
				distinct++
			}
			history = append(history, Op{typ: '+', val: x})
		case '-':
			var k int
			fmt.Fscan(reader, &k)
			removed := make([]int, k)
			for j := 0; j < k; j++ {
				idx := len(arr) - 1
				val := arr[idx]
				arr = arr[:idx]
				freq[val]--
				if freq[val] == 0 {
					distinct--
				}
				removed[k-j-1] = val
			}
			history = append(history, Op{typ: '-', removed: removed})
		case '!':
			if len(history) == 0 {
				continue
			}
			last := history[len(history)-1]
			history = history[:len(history)-1]
			if last.typ == '+' {
				val := last.val
				arr = arr[:len(arr)-1]
				freq[val]--
				if freq[val] == 0 {
					distinct--
				}
			} else {
				for _, val := range last.removed {
					arr = append(arr, val)
					freq[val]++
					if freq[val] == 1 {
						distinct++
					}
				}
			}
		case '?':
			fmt.Fprintln(writer, distinct)
		}
	}
}
