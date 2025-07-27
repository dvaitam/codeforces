package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		next := make([]int, n)
		prev := make([]int, n)
		for i := 0; i < n; i++ {
			next[i] = (i + 1) % n
			prev[i] = (i - 1 + n) % n
		}

		queue := make([]int, 0)
		for i := 0; i < n; i++ {
			j := next[i]
			if gcd(a[i], a[j]) == 1 {
				queue = append(queue, j)
			}
		}

		removed := make([]bool, n)
		ans := make([]int, 0)
		head := 0
		for head < len(queue) {
			j := queue[head]
			head++
			if removed[j] {
				continue
			}
			i := prev[j]
			if removed[i] {
				continue
			}
			if gcd(a[i], a[j]) != 1 {
				continue
			}
			removed[j] = true
			ans = append(ans, j+1)
			ni := next[j]
			next[i] = ni
			prev[ni] = i
			if n > 1 && gcd(a[i], a[ni]) == 1 {
				queue = append(queue, ni)
			}
			n--
			if n == 0 {
				break
			}
		}

		fmt.Fprint(out, len(ans))
		for _, v := range ans {
			fmt.Fprintf(out, " %d", v)
		}
		fmt.Fprintln(out)
	}
}
