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

func compileRef() (string, error) {
    ref := "refB"
    cmd := exec.Command("go", "build", "-o", ref, "707B.go")
    if err := cmd.Run(); err != nil {
        return "", err
    }
    return ref, nil
}

func run(bin string, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func generateTest() string {
    n := rand.Intn(10) + 2
    m := rand.Intn(n*(n-1)/2) + 1
    k := rand.Intn(n + 1) // can be 0
    edges := make(map[[2]int]bool)
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
    for len(edges) < m {
        u := rand.Intn(n) + 1
        v := rand.Intn(n) + 1
        if u == v { continue }
        if u > v { u, v = v, u }
        key := [2]int{u,v}
        if edges[key] { continue }
        edges[key] = true
        l := rand.Intn(100) + 1
        sb.WriteString(fmt.Sprintf("%d %d %d\n", u, v, l))
    }
    if k > 0 {
        used := make(map[int]bool)
        for i := 0; i < k; i++ {
            var v int
            for {
                v = rand.Intn(n) + 1
                if !used[v] { break }
            }
            used[v] = true
            if i > 0 { sb.WriteByte(' ') }
            sb.WriteString(fmt.Sprintf("%d", v))
        }
        sb.WriteByte('\n')
    }
    return sb.String()
}

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    rand.Seed(time.Now().UnixNano())
    ref, err := compileRef()
    if err != nil {
        fmt.Fprintln(os.Stderr, "could not build reference solution:", err)
        os.Exit(1)
    }
    defer os.Remove(ref)

    bin := os.Args[1]
    for t := 0; t < 100; t++ {
        input := generateTest()
        want, err := run(ref, input)
        if err != nil {
            fmt.Fprintln(os.Stderr, "reference failed:", err)
            os.Exit(1)
        }
        got, err := run(bin, input)
        if err != nil {
            fmt.Fprintln(os.Stderr, "test", t+1, "error running binary:", err)
            os.Exit(1)
        }
        if got != want {
            fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected: %s\nactual: %s\n", t+1, input, want, got)
            os.Exit(1)
        }
    }
    fmt.Println("ok")
}

