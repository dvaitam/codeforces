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

func solveA(input string) string {
    in := bufio.NewReader(strings.NewReader(input))
    var t int
    fmt.Fscan(in, &t)
    var out strings.Builder
    for ; t>0; t-- {
        var n, m int
        fmt.Fscan(in, &n, &m)
        rows := make([]int, n)
        cols := make([]int, m)
        for i:=0; i<n; i++ {
            for j:=0; j<m; j++ {
                var x int
                fmt.Fscan(in, &x)
                if x==1 {
                    rows[i]=1
                    cols[j]=1
                }
            }
        }
        er, ec := 0,0
        for _,v := range rows { if v==0 { er++ } }
        for _,v := range cols { if v==0 { ec++ } }
        moves := er
        if ec < moves { moves = ec }
        if moves%2==1 { out.WriteString("Ashish\n") } else { out.WriteString("Vivek\n") }
    }
    return strings.TrimSpace(out.String())
}

func genTests() []test {
    r := rand.New(rand.NewSource(1365))
    tests := make([]test, 0, 100)
    for len(tests) < 100 {
        n := r.Intn(4)+1
        m := r.Intn(4)+1
        var sb strings.Builder
        sb.WriteString("1\n")
        sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
        for i:=0; i<n; i++ {
            for j:=0; j<m; j++ {
                if j>0 { sb.WriteByte(' ') }
                sb.WriteString(fmt.Sprintf("%d", r.Intn(2)))
            }
            sb.WriteByte('\n')
        }
        input := sb.String()
        expected := solveA(input)
        tests = append(tests, test{input, expected})
    }
    return tests
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTests()
    for i, t := range tests {
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

