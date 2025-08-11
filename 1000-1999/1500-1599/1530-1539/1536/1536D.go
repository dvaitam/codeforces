package main

import (
    "bufio"
    "fmt"
    "math/rand"
    "os"
)

const INF = int(2_000_000_000)

type node struct {
    key, pr       int
    left, right   *node
}

func rotateRight(p *node) *node {
    l := p.left
    p.left = l.right
    l.right = p
    return l
}

func rotateLeft(p *node) *node {
    r := p.right
    p.right = r.left
    r.left = p
    return r
}

func insert(p *node, key int) *node {
    if p == nil {
        return &node{key: key, pr: rand.Int()}
    }
    if key < p.key {
        p.left = insert(p.left, key)
        if p.left.pr > p.pr {
            p = rotateRight(p)
        }
    } else if key > p.key {
        p.right = insert(p.right, key)
        if p.right.pr > p.pr {
            p = rotateLeft(p)
        }
    }
    return p
}

func predecessor(p *node, key int) int {
    res := -INF
    for p != nil {
        if key <= p.key {
            p = p.left
        } else {
            res = p.key
            p = p.right
        }
    }
    return res
}

func successor(p *node, key int) int {
    res := INF
    for p != nil {
        if key >= p.key {
            p = p.right
        } else {
            res = p.key
            p = p.left
        }
    }
    return res
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()
    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var n int
        fmt.Fscan(in, &n)
        b := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &b[i])
        }
        l, r := -INF, INF
        var root *node
        ok := true
        for _, x := range b {
            if x < l || x > r {
                ok = false
                break
            }
            root = insert(root, x)
            l = predecessor(root, x)
            r = successor(root, x)
        }
        if ok {
            fmt.Fprintln(out, "YES")
        } else {
            fmt.Fprintln(out, "NO")
        }
    }
}

