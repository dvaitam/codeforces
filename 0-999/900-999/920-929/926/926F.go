package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var p, m int64
	if _, err := fmt.Fscan(in, &n, &p, &m); err != nil {
		return
	}

	d := make([]int64, n)
	t := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &d[i], &t[i])
	}

	var balance int64
	var res int64
	currDay := int64(1)

	for i := 0; i < n; i++ {
		nextDay := d[i]
		if nextDay > currDay {
			length := nextDay - currDay
			if balance >= 0 {
				nonNegDays := balance / p
				if nonNegDays < length {
					res += length - nonNegDays
				}
			} else {
				res += length
			}
			balance -= p * length
		}
		balance += t[i]
		balance -= p
		if balance < 0 {
			res++
		}
		currDay = nextDay + 1
	}

	if m >= currDay {
		length := m - currDay + 1
		if balance >= 0 {
			nonNegDays := balance / p
			if nonNegDays < length {
				res += length - nonNegDays
			}
		} else {
			res += length
		}
	}

	fmt.Fprintln(out, res)
}
