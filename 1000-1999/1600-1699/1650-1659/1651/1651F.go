package main

import (
    "bufio"
    "fmt"
    "os"
)

// Naive simulation for educational purposes. This solution iterates over
// every monster and every tower. It runs in O(n*q) time and is not optimized
// for the original constraints.
func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var n int
    fmt.Fscan(reader, &n)
    c := make([]int64, n)
    r := make([]int64, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &c[i], &r[i])
    }
    var q int
    fmt.Fscan(reader, &q)
    t := make([]int64, q)
    h := make([]int64, q)
    for i := 0; i < q; i++ {
        fmt.Fscan(reader, &t[i], &h[i])
    }

    mana := make([]int64, n)
    last := make([]int64, n)
    for i := 0; i < n; i++ {
        mana[i] = c[i]
        last[i] = 0
    }

    var total int64
    for j := 0; j < q; j++ {
        health := h[j]
        for i := 0; i < n && health > 0; i++ {
            arrival := t[j] + int64(i)
            // regenerate
            delta := arrival - last[i]
            if delta > 0 {
                mana[i] += r[i] * delta
                if mana[i] > c[i] {
                    mana[i] = c[i]
                }
                last[i] = arrival
            }
            // attack
            if mana[i] > 0 {
                if mana[i] >= health {
                    mana[i] -= health
                    health = 0
                } else {
                    health -= mana[i]
                    mana[i] = 0
                }
            }
        }
        total += health
    }

    fmt.Fprintln(writer, total)
}

