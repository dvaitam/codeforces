package main

import (
    "bufio"
    "fmt"
    "os"
)

func mex(set map[int]bool) int {
    for i := 0; ; i++ {
        if !set[i] {
            return i
        }
    }
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var n, l, r int
    if _, err := fmt.Fscan(reader, &n, &l, &r); err != nil {
        return
    }

    adj := make([][]int, n)
    for i := 0; i < n; i++ {
        var u, v int
        fmt.Fscan(reader, &u, &v)
        u--
        v--
        adj[u] = append(adj[u], v)
        adj[v] = append(adj[v], u)
    }

    visited := make([]bool, n)
    cycles := []int{}
    for i := 0; i < n; i++ {
        if visited[i] {
            continue
        }
        cur := i
        prev := -1
        length := 0
        for {
            visited[cur] = true
            length++
            next := adj[cur][0]
            if next == prev {
                if len(adj[cur]) > 1 {
                    next = adj[cur][1]
                }
            }
            prev, cur = cur, next
            if cur == i {
                break
            }
        }
        cycles = append(cycles, length)
    }

    maxLen := 0
    for _, v := range cycles {
        if v > maxLen {
            maxLen = v
        }
    }

    if maxLen < l {
        fmt.Fprintln(writer, "Bob")
        return
    }

    grundy := make([]int, maxLen+1)
    for len := 1; len <= maxLen; len++ {
        reachable := make(map[int]bool)
        for k := l; k <= r && k <= len; k++ {
            for start := 0; start <= len-k; start++ {
                left := start
                right := len - k - start
                val := grundy[left] ^ grundy[right]
                reachable[val] = true
            }
        }
        grundy[len] = mex(reachable)
    }

    total := 0
    for _, c := range cycles {
        reachable := make(map[int]bool)
        if c >= l && c <= r {
            reachable[0] = true
        }
        for k := l; k <= r && k < c; k++ {
            reachable[grundy[c-k]] = true
        }
        g := mex(reachable)
        total ^= g
    }

    if total != 0 {
        fmt.Fprintln(writer, "Alice")
    } else {
        fmt.Fprintln(writer, "Bob")
    }
}

