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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		var s string
		fmt.Fscan(in, &n, &k)
		fmt.Fscan(in, &s)
		freq := make([]int, 26)
		for i := 0; i < n; i++ {
			freq[s[i]-'a']++
		}
		pairs := 0
		singles := 0
		for _, v := range freq {
			pairs += v / 2
			singles += v % 2
		}
		m := pairs / k
		ans := 2 * m
		if singles+2*(pairs%k) >= k {
			ans++
		}
		fmt.Fprintln(out, ans)
	}
}
