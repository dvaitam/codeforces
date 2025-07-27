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

type test struct{
    input string
    expected string
}

func runBinary(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func solveB(input string) string {
    in := bufio.NewReader(strings.NewReader(input))
    var t int
    fmt.Fscan(in, &t)
    var out strings.Builder
    for ; t>0; t-- {
        var n int
        fmt.Fscan(in, &n)
        a := make([]int, n)
        for i:=0; i<n; i++ { fmt.Fscan(in, &a[i]) }
        has0, has1 := false, false
        b := make([]int, n)
        for i:=0; i<n; i++ {
            fmt.Fscan(in, &b[i])
            if b[i]==0 { has0 = true } else { has1 = true }
        }
        sorted := true
        for i:=1; i<n; i++ {
            if a[i] < a[i-1] { sorted=false; break }
        }
        if sorted || (has0 && has1) {
            out.WriteString("Yes\n")
        } else {
            out.WriteString("No\n")
        }
    }
    return strings.TrimSpace(out.String())
}

func genTests() []test {
    r := rand.New(rand.NewSource(1365))
    tests := make([]test, 0, 100)
    for len(tests) < 100 {
        n := r.Intn(5)+1
        var sb strings.Builder
        sb.WriteString("1\n")
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for i:=0; i<n; i++ {
            if i>0 { sb.WriteByte(' ') }
            sb.WriteString(fmt.Sprintf("%d", r.Intn(10)))
        }
        sb.WriteByte('\n')
        for i:=0; i<n; i++ {
            if i>0 { sb.WriteByte(' ') }
            if r.Intn(2)==0 { sb.WriteByte('0') } else { sb.WriteByte('1') }
        }
        sb.WriteByte('\n')
        input := sb.String()
        expected := solveB(input)
        tests = append(tests, test{input, expected})
    }
    return tests
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTests()
    for i,t := range tests {
        out, err := runBinary(bin, t.input)
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

