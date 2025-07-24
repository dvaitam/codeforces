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

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		sum += a[i]
	}
	avg := sum / n

	type meeting struct {
		s int
		b []int
	}
	var ops []meeting

	record := func(start int) {
		b := make([]int, k)
		for j := 0; j < k; j++ {
			b[j] = a[(start+j)%n]
		}
		ops = append(ops, meeting{start, b})
	}

	for i := 0; i < n; i++ {
		for a[i] < avg {
			j := i + 1
			for ; j < i+n && a[j%n] <= avg; j++ {
			}
			if j >= i+n {
				break
			}
			idx := j % n
			x := avg - a[i]
			if diff := a[idx] - avg; diff < x {
				x = diff
			}
			cur := idx
			for cur != i {
				prev := (cur - 1 + n) % n
				a[prev] += x
				a[cur] -= x
				record(prev)
				cur = prev
			}
		}
		for a[i] > avg {
			j := i + 1
			for ; j < i+n && a[j%n] >= avg; j++ {
			}
			if j >= i+n {
				break
			}
			idx := j % n
			x := a[i] - avg
			if diff := avg - a[idx]; diff < x {
				x = diff
			}
			cur := i
			for cur != idx {
				next := (cur + 1) % n
				a[cur] -= x
				a[next] += x
				record(cur)
				cur = next
			}
		}
	}

	fmt.Fprintln(writer, len(ops))
	for _, op := range ops {
		fmt.Fprint(writer, op.s)
		for _, v := range op.b {
			fmt.Fprint(writer, " ", v)
		}
		fmt.Fprintln(writer)
	}
}
