package main

import (
    "bufio"
    "fmt"
    "math"
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

        // Collect chip positions
        type pt struct{ r, c int }
        chips := make([]pt, 0)
        for i := 0; i < n; i++ {
            if s1[i] == '*' {
                chips = append(chips, pt{0, i})
            }
            if s2[i] == '*' {
                chips = append(chips, pt{1, i})
            }
        }

        if len(chips) <= 1 {
            fmt.Fprintln(out, 0)
            continue
        }

        // Exact oracle: try all target cells (2 * n) and take minimal sum of L1 distances
        best := math.MaxInt32
        for r := 0; r < 2; r++ {
            for c := 0; c < n; c++ {
                cur := 0
                for _, p := range chips {
                    dr := p.r - r
                    if dr < 0 {
                        dr = -dr
                    }
                    dc := p.c - c
                    if dc < 0 {
                        dc = -dc
                    }
                    cur += dr + dc
                }
                if cur < best {
                    best = cur
                }
            }
        }
        fmt.Fprintln(out, best)
    }
}
