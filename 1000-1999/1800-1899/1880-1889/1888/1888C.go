package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}

		mp1 := make(map[int]int) // last occurrence (1-based)
		mp2 := make(map[int]int) // first occurrence (1-based)

		for i := 1; i <= n; i++ {
			mp1[a[i]] = i
		}

		var ans int64
		num := 0
		for i := 1; i <= n; i++ {
			if mp2[a[i]] == 0 {
				mp2[a[i]] = i
				num++
			}
			if mp1[a[i]] == i {
				ans += int64(num)
			}
		}
		fmt.Fprintln(out, ans)
	}
}
