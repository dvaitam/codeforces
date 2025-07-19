package main

import (
	"bufio"
	"fmt"
	"os"
)

var prime []bool

func init() {
	const N = 10000
	prime = make([]bool, N+1)
	for i := 2; i <= N; i++ {
		prime[i] = true
	}
	for i := 2; i*i <= N; i++ {
		if prime[i] {
			for j := i * i; j <= N; j += i {
				prime[j] = false
			}
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var tt int
	fmt.Fscan(reader, &tt)
	for tt > 0 {
		tt--
		var n int
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)
		cnt := make([]int, 10)
		a := make([]int, n)
		for i, ch := range s {
			d := int(ch - '0')
			a[i] = d
			cnt[d]++
		}
		// single-digit non-prime
		if cnt[1] > 0 || cnt[4] > 0 || cnt[6] > 0 || cnt[8] > 0 || cnt[9] > 0 {
			fmt.Fprintln(writer, 1)
			switch {
			case cnt[1] > 0:
				fmt.Fprintln(writer, 1)
			case cnt[4] > 0:
				fmt.Fprintln(writer, 4)
			case cnt[6] > 0:
				fmt.Fprintln(writer, 6)
			case cnt[8] > 0:
				fmt.Fprintln(writer, 8)
			default:
				fmt.Fprintln(writer, 9)
			}
			continue
		}
		// two-digit
		found := false
		for i := 0; i < n && !found; i++ {
			for j := i + 1; j < n; j++ {
				num := a[i]*10 + a[j]
				if !prime[num] {
					fmt.Fprintln(writer, 2)
					fmt.Fprintf(writer, "%d%d\n", a[i], a[j])
					found = true
					break
				}
			}
		}
		if found {
			continue
		}
		// three-digit
		for i := 0; i < n && !found; i++ {
			for j := i + 1; j < n && !found; j++ {
				for k := j + 1; k < n; k++ {
					num := a[i]*100 + a[j]*10 + a[k]
					if !prime[num] {
						fmt.Fprintln(writer, 3)
						fmt.Fprintf(writer, "%d%d%d\n", a[i], a[j], a[k])
						found = true
						break
					}
				}
			}
		}
	}
}
