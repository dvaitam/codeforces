package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strings"
)

func expectedE(n int, l,r []int) string {
    p := make([]int,n)
    for i:=0;i<n;i++ { p[i] = i }
    sort.Slice(p, func(i,j int) bool {
        if l[p[i]] == l[p[j]] {
            return r[p[i]] < r[p[j]]
        }
        return l[p[i]] < l[p[j]]
    })
    for i:=1;i<n;i++ {
        if l[p[i-1]] == l[p[i]] { return fmt.Sprintf("%d", p[i-1]+1) }
        if r[p[i-1]] >= r[p[i]] { return fmt.Sprintf("%d", p[i]+1) }
    }
    for i:=1;i+1<n;i++ {
        if r[p[i-1]]+1 >= l[p[i+1]] { return fmt.Sprintf("%d", p[i]+1) }
    }
    return "-1"
}

func genTestsE() []string {
    rand.Seed(5)
    tests := make([]string,0,100)
    for len(tests) < 100 {
        n := rand.Intn(8)+2
        l := make([]int,n)
        r := make([]int,n)
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for i:=0;i<n;i++ {
            l[i] = rand.Intn(100)+1
            r[i] = l[i] + rand.Intn(20)
            sb.WriteString(fmt.Sprintf("%d %d\n", l[i], r[i]))
        }
        tests = append(tests, sb.String())
    }
    return tests
}

func runBinary(bin string, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main(){
    if len(os.Args)!=2 {
        fmt.Fprintf(os.Stderr,"Usage: go run verifierE.go <binary>\n")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTestsE()
    for idx,t := range tests {
        lines := strings.Split(strings.TrimSpace(t),"\n")
        var n int
        fmt.Sscanf(lines[0], "%d", &n)
        l := make([]int,n)
        r := make([]int,n)
        for i:=0;i<n;i++ {
            fmt.Sscanf(lines[1+i], "%d %d", &l[i], &r[i])
        }
        want := expectedE(n,l,r)
        got, err := runBinary(bin,t)
        if err != nil {
            fmt.Printf("Test %d: runtime error: %v\n", idx+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got)!=want {
            fmt.Printf("Test %d failed.\nInput:\n%s\nExpected: %s\nGot: %s\n", idx+1, t, want, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}

