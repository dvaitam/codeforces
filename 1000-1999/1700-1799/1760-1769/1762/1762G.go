package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(r *bufio.Reader, w *bufio.Writer) {
	var n int
	fmt.Fscan(r, &n)
	a := make([]int, n+2)
	buc := make([]int, n+2)
	cnt := make([]int, n+2)
	p := make([]int, n+2)
	for i := 1; i <= n; i++ {
		fmt.Fscan(r, &a[i])
		buc[a[i]]++
		cnt[buc[a[i]]]++
	}
	mx := n
	for mx > 0 && cnt[mx] == 0 {
		mx--
	}
	if mx > (n+1)/2 {
		fmt.Fprintln(w, "NO")
		return
	}
	l, r := 0, 0
	for i := 1; i <= n; i++ {
		if mx == (n-l+1)/2 {
			l--
			for j := 1; j <= n; j++ {
				if buc[j] == mx {
					for k := i; k <= n; k++ {
						if a[k] == j {
							l += 2
							p[l] = k
						} else {
							r += 2
							p[r] = k
						}
					}
				}
			}
			break
		}
		if a[p[l]] == a[i] {
			r += 2
			p[r] = i
		} else {
			l++
			p[l] = i
			if l < r {
				l++
			} else {
				r = l
			}
		}
		// decrease count of current bucket
		old := buc[a[i]]
		cnt[old]--
		buc[a[i]]--
		for mx > 0 && cnt[mx] == 0 {
			mx--
		}
	}
	fmt.Fprintln(w, "YES")
	for i := 1; i <= n; i++ {
		fmt.Fprint(w, p[i], " ")
	}
	fmt.Fprintln(w)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var T int
	fmt.Fscan(reader, &T)
	for T > 0 {
		solve(reader, writer)
		T--
	}
}
