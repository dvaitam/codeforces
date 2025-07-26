package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "time"
)

type Node struct {
    left, right *Node
    sum  int64
    lazy int8 // -1 none, 0 or 1 set
}

func newNode() *Node { return &Node{lazy: -1} }

func push(node *Node, l, r int64) {
    if node.lazy == -1 || l == r {
        return
    }
    mid := (l + r) >> 1
    if node.left == nil { node.left = newNode() }
    if node.right == nil { node.right = newNode() }
    val := int64(node.lazy)
    node.left.sum = (mid - l + 1) * val
    node.left.lazy = node.lazy
    node.right.sum = (r - mid) * val
    node.right.lazy = node.lazy
    node.lazy = -1
}

func update(node *Node, l, r, ql, qr int64, val int8) {
    if node == nil || ql > r || qr < l { return }
    if ql <= l && r <= qr {
        node.sum = (r - l + 1) * int64(val)
        node.lazy = val
        node.left, node.right = nil, nil
        return
    }
    push(node, l, r)
    mid := (l + r) >> 1
    if ql <= mid {
        if node.left == nil { node.left = newNode() }
        update(node.left, l, mid, ql, qr, val)
    }
    if qr > mid {
        if node.right == nil { node.right = newNode() }
        update(node.right, mid+1, r, ql, qr, val)
    }
    node.sum = 0
    if node.left != nil { node.sum += node.left.sum }
    if node.right != nil { node.sum += node.right.sum }
}

func solveE(n int64, ops [][3]int64) string {
    root := &Node{sum: n, lazy: 1}
    var sb strings.Builder
    for _, op := range ops {
        l, r, k := op[0], op[1], op[2]
        val := int8(0)
        if k == 2 {
            val = 1
        }
        update(root, 1, n, l, r, val)
        sb.WriteString(fmt.Sprintf("%d\n", root.sum))
    }
    return strings.TrimRight(sb.String(), "\n")
}

func generateE(rng *rand.Rand) (string, string) {
    n := int64(rng.Intn(1000) + 1)
    q := int64(rng.Intn(20) + 1)
    ops := make([][3]int64, q)
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
    for i := int64(0); i < q; i++ {
        l := int64(rng.Intn(int(n)) + 1)
        r := int64(rng.Intn(int(n-l+1)) + int(l))
        k := int64(rng.Intn(2) + 1)
        ops[i] = [3]int64{l, r, k}
        sb.WriteString(fmt.Sprintf("%d %d %d\n", l, r, k))
    }
    return sb.String(), solveE(n, ops)
}

func runCase(bin, input, exp string) error {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    got := strings.TrimSpace(out.String())
    if got != exp {
        return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
        os.Exit(1)
    }
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    bin := os.Args[1]
    for i := 0; i < 100; i++ {
        in, exp := generateE(rng)
        if err := runCase(bin, in, exp); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

