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
		var n int
		fmt.Fscan(in, &n)
        var s1, s2 string
        fmt.Fscan(in, &s1)
        fmt.Fscan(in, &s2)

        // Build grid as ints and locate span [l, r]
        g := make([][2]int, n+2) // 1-based indexing for convenience
        l, r := 1, n
        for j := 1; j <= n; j++ {
            if s1[j-1] == '*' {
                g[j][0] = 1
            }
            if s2[j-1] == '*' {
                g[j][1] = 1
            }
        }
        // Check if there is any chip
        any := false
        for j := 1; j <= n; j++ {
            if g[j][0] == 1 || g[j][1] == 1 {
                any = true
                break
            }
        }
        if !any {
            fmt.Fprintln(out, 0)
            continue
        }
        for l <= n && g[l][0] == 0 && g[l][1] == 0 {
            l++
        }
        for r >= 1 && g[r][0] == 0 && g[r][1] == 0 {
            r--
        }
        // DP over columns l..r
        f0 := make([]int, n+2)
        f1 := make([]int, n+2)
        // Defaults are 0, matching the common CF solution with final -1
        for i := l; i <= r; i++ {
            // end at top row
            a := f0[i-1] + g[i][1] + 1
            b := f1[i-1] + 2
            if a < b {
                f0[i] = a
            } else {
                f0[i] = b
            }
            // end at bottom row
            c := f1[i-1] + g[i][0] + 1
            d := f0[i-1] + 2
            if c < d {
                f1[i] = c
            } else {
                f1[i] = d
            }
        }
        // subtract 1 to remove the initial extra move counted at l
        ans := f0[r]
        if f1[r] < ans {
            ans = f1[r]
        }
        fmt.Fprintln(out, ans-1)
    }
}
