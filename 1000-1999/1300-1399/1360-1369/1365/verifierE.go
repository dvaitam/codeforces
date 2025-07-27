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

func solveE(input string) string {
    in := bufio.NewReader(strings.NewReader(input))
    var n int
    fmt.Fscan(in, &n)
    arr := make([]int64, n)
    for i:=0;i<n;i++ { fmt.Fscan(in,&arr[i]) }
    var ans int64
    for i:=0;i<n;i++ {
        if arr[i] > ans { ans = arr[i] }
        for j:=i+1;j<n;j++ {
            v := arr[i] | arr[j]
            if v > ans { ans = v }
            for k:=j+1;k<n;k++ {
                v3 := v | arr[k]
                if v3 > ans { ans = v3 }
            }
        }
    }
    return fmt.Sprintf("%d", ans)
}

func genTests() []test {
    r := rand.New(rand.NewSource(1365))
    tests := make([]test,0,100)
    for len(tests) < 100 {
        n := r.Intn(6)+1
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for i:=0;i<n;i++ {
            if i>0 { sb.WriteByte(' ') }
            sb.WriteString(fmt.Sprintf("%d", r.Int63n(512)))
        }
        sb.WriteByte('\n')
        input := sb.String()
        expected := solveE(input)
        tests = append(tests, test{input, expected})
    }
    return tests
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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

