package main

import (
    "bufio"
    "fmt"
    "os"
)

// This solution performs a straightforward search with aggressive pruning.
// The input length is at most 20000, so we iterate over all substrings,
// quickly discard those that fail obvious necessary conditions, and only
// run a full O(len) verification when the conditions are met.
//
// Necessary conditions we use before the expensive check:
// - even length (the grid must contain an even number of cells);
// - counts of U/D and L/R are equal (total displacement must be zero);
// - at least one dimension of the grid is even (otherwise product is odd);
// - trivial boundary invalidity is ruled out by the final check.
//
// Full check for a candidate substring and a specific grid size n x m:
// - ensure all moves stay inside the grid;
// - accumulate indegree for every cell; valid iff each indegree is exactly 1.
// If any divisor (dimension) works, the substring is counted once.

func main() {
    in := bufio.NewReader(os.Stdin)
    var s string
    fmt.Fscan(in, &s)
    n := len(s)

    // Prefix counts for quick balance checks
    pref := make([][4]int, n+1) // order: U D L R
    idx := func(c byte) int {
        switch c {
        case 'U':
            return 0
        case 'D':
            return 1
        case 'L':
            return 2
        default:
            return 3
        }
    }
    for i := 0; i < n; i++ {
        pref[i+1] = pref[i]
        pref[i+1][idx(s[i])]++
    }

    // Precompute divisors for every length once
    divisors := make([][]int, n+1)
    for l := 1; l <= n; l++ {
        for d := 1; d*d <= l; d++ {
            if l%d == 0 {
                divisors[l] = append(divisors[l], d)
                if d*d != l {
                    divisors[l] = append(divisors[l], l/d)
                }
            }
        }
    }

    ans := 0
    // Helper: balance check using prefix counts
    balanced := func(l, r int) bool {
        u := pref[r][0] - pref[l][0]
        d := pref[r][1] - pref[l][1]
        if u != d {
            return false
        }
        L := pref[r][2] - pref[l][2]
        R := pref[r][3] - pref[l][3]
        return L == R
    }

    // Validation for a given substring and grid size
    check := func(start, length, rows, cols int) bool {
        dirs := []struct{ dr, dc int }{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
        indeg := make([]byte, length)
        for idxAbs := 0; idxAbs < length; idxAbs++ {
            ch := s[start+idxAbs]
            var dir int
            switch ch {
            case 'U':
                dir = 0
            case 'D':
                dir = 1
            case 'L':
                dir = 2
            default:
                dir = 3
            }
            r := idxAbs / cols
            c := idxAbs % cols
            nr := r + dirs[dir].dr
            nc := c + dirs[dir].dc
            if nr < 0 || nr >= rows || nc < 0 || nc >= cols {
                return false
            }
            target := nr*cols + nc
            indeg[target]++
        }
        for _, v := range indeg {
            if v != 1 {
                return false
            }
        }
        return true
    }

    for l := 0; l < n; l++ {
        for r := l + 1; r <= n; r++ {
            length := r - l
            if length%2 == 1 {
                continue
            }
            if !balanced(l, r) {
                continue
            }
            ok := false
            for _, rows := range divisors[length] {
                cols := length / rows
                if rows%2 == 1 && cols%2 == 1 {
                    continue // product odd, impossible
                }
                if check(l, length, rows, cols) {
                    ok = true
                    break
                }
            }
            if ok {
                ans++
            }
        }
    }

    fmt.Println(ans)
}
