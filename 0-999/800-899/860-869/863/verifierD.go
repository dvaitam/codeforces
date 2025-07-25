package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
)

func expectedD(n,q,m int, a []int, t,l,r []int, b []int) string {
    res := make([]int, m)
    for idx:=0; idx<m; idx++ {
        pos := b[idx]
        for i:=q-1; i>=0; i-- {
            if pos < l[i] || pos > r[i] {
                continue
            }
            if t[i]==1 {
                if pos == l[i] {
                    pos = r[i]
                } else {
                    pos--
                }
            } else {
                pos = l[i] + r[i] - pos
            }
        }
        res[idx] = a[pos]
    }
    var sb strings.Builder
    for i,v := range res {
        if i>0 { sb.WriteByte(' ') }
        sb.WriteString(fmt.Sprintf("%d", v))
    }
    return sb.String()
}

func genTestsD() []string {
    rand.Seed(4)
    tests := make([]string,0,100)
    for len(tests) < 100 {
        n := rand.Intn(5)+1
        q := rand.Intn(5)+1
        m := rand.Intn(5)+1
        a := make([]int,n+1)
        for i:=1;i<=n;i++ { a[i] = rand.Intn(100)+1 }
        tArr := make([]int,q)
        lArr := make([]int,q)
        rArr := make([]int,q)
        for i:=0;i<q;i++ {
            tArr[i] = rand.Intn(2)+1
            lArr[i] = rand.Intn(n)+1
            rArr[i] = rand.Intn(n-lArr[i]+1)+lArr[i]
        }
        b := make([]int,m)
        for i:=0;i<m;i++ { b[i] = rand.Intn(n)+1 }
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d %d %d\n", n,q,m))
        for i:=1;i<=n;i++ {
            if i>1 { sb.WriteByte(' ') }
            sb.WriteString(fmt.Sprintf("%d", a[i]))
        }
        sb.WriteString("\n")
        for i:=0;i<q;i++ {
            sb.WriteString(fmt.Sprintf("%d %d %d\n", tArr[i], lArr[i], rArr[i]))
        }
        for i:=0;i<m;i++ {
            if i>0 { sb.WriteByte(' ') }
            sb.WriteString(fmt.Sprintf("%d", b[i]))
        }
        sb.WriteString("\n")
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

func parseInts(line string) []int {
    parts := strings.Fields(line)
    res := make([]int,len(parts))
    for i,p := range parts { fmt.Sscanf(p, "%d", &res[i]) }
    return res
}

func main(){
    if len(os.Args)!=2 {
        fmt.Fprintf(os.Stderr,"Usage: go run verifierD.go <binary>\n")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTestsD()
    for idx,t := range tests {
        lines := strings.Split(strings.TrimSpace(t),"\n")
        var n,q,m int
        fmt.Sscanf(lines[0],"%d %d %d", &n,&q,&m)
        a := make([]int,n+1)
        vals := parseInts(lines[1])
        copy(a[1:],vals)
        tArr := make([]int,q)
        lArr := make([]int,q)
        rArr := make([]int,q)
        for i:=0;i<q;i++ {
            fmt.Sscanf(lines[2+i],"%d %d %d", &tArr[i], &lArr[i], &rArr[i])
        }
        b := parseInts(lines[2+q])
        want := expectedD(n,q,m,a,tArr,lArr,rArr,b)
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

