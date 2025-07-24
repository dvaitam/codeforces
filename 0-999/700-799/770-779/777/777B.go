package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var sherlock, moriarty string
	fmt.Fscan(in, &sherlock)
	fmt.Fscan(in, &moriarty)

	var freq1 [10]int
	var freq2 [10]int
	for i := 0; i < n; i++ {
		d := moriarty[i] - '0'
		freq1[d]++
		freq2[d]++
	}

	minFlicks := 0
	for i := 0; i < n; i++ {
		d := int(sherlock[i] - '0')
		j := d
		for j < 10 && freq1[j] == 0 {
			j++
		}
		if j < 10 {
			freq1[j]--
		} else {
			minFlicks++
			for k := 0; k < 10; k++ {
				if freq1[k] > 0 {
					freq1[k]--
					break
				}
			}
		}
	}

	maxFlickSherlock := 0
	for i := 0; i < n; i++ {
		d := int(sherlock[i] - '0')
		j := d + 1
		for j < 10 && freq2[j] == 0 {
			j++
		}
		if j < 10 {
			freq2[j]--
			maxFlickSherlock++
		} else {
			for k := 0; k < 10; k++ {
				if freq2[k] > 0 {
					freq2[k]--
					break
				}
			}
		}
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, minFlicks)
	fmt.Fprintln(out, maxFlickSherlock)
	out.Flush()
}
