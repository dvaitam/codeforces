package main

import (
    "bufio"
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
)

type test struct{ input, expected string }

func runBin(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else { cmd = exec.Command(bin) }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var errb bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errb
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func solveC(input string) string {
    in := bufio.NewReader(strings.NewReader(input))
    var n int
    fmt.Fscan(in, &n)
    a := make([]int, n)
    b := make([]int, n)
    for i:=0; i<n; i++ { fmt.Fscan(in, &a[i]) }
    for i:=0; i<n; i++ { fmt.Fscan(in, &b[i]) }
    pos := make([]int, n+1)
    for i:=0; i<n; i++ { pos[a[i]] = i }
    cnt := make(map[int]int)
    best := 0
    for i:=0; i<n; i++ {
        shift := (i - pos[b[i]] + n) % n
        cnt[shift]++
        if cnt[shift] > best { best = cnt[shift] }
    }
    return fmt.Sprintf("%d", best)
}

func genTests() []test {
    r := rand.New(rand.NewSource(1365))
    tests := make([]test,0,100)
    for len(tests) < 100 {
        n := r.Intn(6)+1
        perm := rand.Perm(n)
        a := make([]int, n)
        for i:=0;i<n;i++ { a[i] = perm[i]+1 }
        perm2 := rand.Perm(n)
        b := make([]int, n)
        for i:=0;i<n;i++ { b[i] = perm2[i]+1 }
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for i, v := range a { if i>0 { sb.WriteByte(' ') }; sb.WriteString(fmt.Sprintf("%d", v)) }
        sb.WriteByte('\n')
        for i, v := range b { if i>0 { sb.WriteByte(' ') }; sb.WriteString(fmt.Sprintf("%d", v)) }
        sb.WriteByte('\n')
        input := sb.String()
        expected := solveC(input)
        tests = append(tests, test{input, expected})
    }
    return tests
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTests()
    for i, t := range tests {
        out, err := runBin(bin, t.input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, t.input)
            os.Exit(1)
        }
        if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
            fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, out)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}

