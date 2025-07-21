package main

import (
   "fmt"
)

func main() {
    var x1, x2, a, b int64
    if _, err := fmt.Scan(&x1, &x2, &a, &b); err != nil {
        return
    }
    // Direct kill by first
    d := x2 - x1
    if d >= a && d <= b {
        fmt.Println("FIRST")
        fmt.Println(x2)
        return
    }
    // If can stay indefinitely
    if a <= 0 && b >= 0 {
        fmt.Println("DRAW")
        return
    }
    // Setup absolute move lengths
    A := a
    B := b
    if A < 0 {
        A = -A
    }
    if B < 0 {
        B = -B
    }
    if A > B {
        A, B = B, A
    }
    // eventual win by modular reduction
    D := d
    if D < 0 {
        D = -D
    }
    c := A + B
    if c != 0 {
        r := D % c
        if r >= A && r <= B {
            // first can force win
            // choose move towards opponent of length r
            var move int64 = r
            if d < 0 {
                move = -r
            }
            fmt.Println("FIRST")
            fmt.Println(x1 + move)
            return
        }
    }
    // check if second can force win: all first moves lead to kill-range for second
    // risky intervals: d1 in [d - B, d - A] or [d + A, d + B]
    // intersect with [a,b]
    // compute union length of risky set
    var intervals [][2]int64
    // interval1
    l1 := d - B
    r1 := d - A
    if l1 <= r1 {
        if l1 < a {
            l1 = a
        }
        if r1 > b {
            r1 = b
        }
        if l1 <= r1 {
            intervals = append(intervals, [2]int64{l1, r1})
        }
    }
    // interval2
    l2 := d + A
    r2 := d + B
    if l2 <= r2 {
        if l2 < a {
            l2 = a
        }
        if r2 > b {
            r2 = b
        }
        if l2 <= r2 {
            intervals = append(intervals, [2]int64{l2, r2})
        }
    }
    // merge intervals and count covered points
    covered := int64(0)
    if len(intervals) > 0 {
        // sort by start
        for i := 0; i < len(intervals)-1; i++ {
            for j := i + 1; j < len(intervals); j++ {
                if intervals[j][0] < intervals[i][0] {
                    intervals[i], intervals[j] = intervals[j], intervals[i]
                }
            }
        }
        cur := intervals[0]
        for i := 1; i < len(intervals); i++ {
            nxt := intervals[i]
            if nxt[0] <= cur[1]+1 {
                if nxt[1] > cur[1] {
                    cur[1] = nxt[1]
                }
            } else {
                covered += cur[1] - cur[0] + 1
                cur = nxt
            }
        }
        covered += cur[1] - cur[0] + 1
    }
    total := b - a + 1
    if covered == total {
        fmt.Println("SECOND")
    } else {
        fmt.Println("DRAW")
    }
}
