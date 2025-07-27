package main

import (
    "bufio"
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "sort"
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

type pair struct{ x,y int64 }

func solveF(input string) string {
    in := bufio.NewReader(strings.NewReader(input))
    var t int
    fmt.Fscan(in, &t)
    var out strings.Builder
    for ; t>0; t-- {
        var n int
        fmt.Fscan(in, &n)
        a := make([]int64, n)
        b := make([]int64, n)
        for i:=0;i<n;i++ { fmt.Fscan(in,&a[i]) }
        for i:=0;i<n;i++ { fmt.Fscan(in,&b[i]) }
        if n%2==1 && a[n/2]!=b[n/2] {
            out.WriteString("No\n")
            continue
        }
        m := n/2
        pa := make([]pair,m)
        pb := make([]pair,m)
        for i:=0;i<m;i++ {
            x1,x2 := a[i], a[n-1-i]
            if x1 > x2 { x1,x2 = x2,x1 }
            pa[i]=pair{x1,x2}
            y1,y2 := b[i], b[n-1-i]
            if y1>y2 { y1,y2 = y2,y1 }
            pb[i]=pair{y1,y2}
        }
        sort.Slice(pa, func(i,j int) bool {
            if pa[i].x==pa[j].x { return pa[i].y < pa[j].y }
            return pa[i].x < pa[j].x
        })
        sort.Slice(pb, func(i,j int) bool {
            if pb[i].x==pb[j].x { return pb[i].y < pb[j].y }
            return pb[i].x < pb[j].x
        })
        ok := true
        for i:=0;i<m;i++ { if pa[i]!=pb[i] { ok=false; break } }
        if ok { out.WriteString("Yes\n") } else { out.WriteString("No\n") }
    }
    return strings.TrimSpace(out.String())
}

func genTests() []test {
    r := rand.New(rand.NewSource(1365))
    tests := make([]test,0,100)
    for len(tests)<100 {
        n := r.Intn(6)+1
        var sb strings.Builder
        sb.WriteString("1\n")
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for i:=0;i<n;i++ { if i>0 {sb.WriteByte(' ')}; sb.WriteString(fmt.Sprintf("%d", r.Int63n(20))) }
        sb.WriteByte('\n')
        for i:=0;i<n;i++ { if i>0 {sb.WriteByte(' ')}; sb.WriteString(fmt.Sprintf("%d", r.Int63n(20))) }
        sb.WriteByte('\n')
        input := sb.String()
        expected := solveF(input)
        tests = append(tests, test{input, expected})
    }
    return tests
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTests()
    for i,t := range tests {
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

