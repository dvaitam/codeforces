package main

import (
    "bufio"
    "bytes"
    "context"
    "fmt"
    "io"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "time"
)

func solveC(in io.Reader) string {
    reader := bufio.NewReader(in)
    var n, m int
    if _, err := fmt.Fscan(reader, &n, &m); err != nil {
        return ""
    }
    perm := make([]int, n)
    pos := make([]int, n+1)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &perm[i])
        pos[perm[i]] = i + 1
    }
    limit := make([]int, n+2)
    for i := range limit {
        limit[i] = n
    }
    for i := 0; i < m; i++ {
        var a, b int
        fmt.Fscan(reader, &a, &b)
        x := pos[a]
        y := pos[b]
        if x > y { x, y = y, x }
        if y-1 < limit[x] { limit[x] = y - 1 }
    }
    for i := n - 1; i >= 1; i-- {
        if limit[i+1] < limit[i] { limit[i] = limit[i+1] }
    }
    var ans int64
    for i := 1; i <= n; i++ {
        ans += int64(limit[i]-i+1)
    }
    return fmt.Sprint(ans)
}

func genTests() []string {
    r := rand.New(rand.NewSource(3))
    tests := make([]string, 100)
    for i := 0; i < 100; i++ {
        n := r.Intn(8) + 1
        m := r.Intn(n)
        perm := r.Perm(n)
        for j := 0; j < n; j++ { perm[j]++ }
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
        for j, v := range perm {
            if j > 0 { sb.WriteByte(' ') }
            sb.WriteString(fmt.Sprint(v))
        }
        sb.WriteByte('\n')
        pairs := make([][2]int, 0, m)
        for j := 0; j < m; j++ {
            a := r.Intn(n) + 1
            b := r.Intn(n) + 1
            for a == b {
                b = r.Intn(n) + 1
            }
            pairs = append(pairs, [2]int{a, b})
        }
        for _, p := range pairs {
            sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
        }
        tests[i] = sb.String()
    }
    return tests
}

func runBinary(bin string, input string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    cmd := exec.CommandContext(ctx, bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run verifierC.go /path/to/binary")
        return
    }
    bin := os.Args[1]
    tests := genTests()
    for i, tc := range tests {
        expected := solveC(strings.NewReader(tc))
        actual, err := runBinary(bin, tc)
        if err != nil {
            fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
            return
        }
        if actual != strings.TrimSpace(expected) {
            fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", i+1, tc, expected, actual)
            return
        }
    }
    fmt.Println("All tests passed!")
}

