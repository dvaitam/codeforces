package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
)

func genPaths(m int, a []int) [][]int {
    orient := make([]int,m)
    cur := make([]int,m)
    path := [][]int{}
    var rec func(level int)
    rec = func(level int) {
        if level==0 {
            s := make([]int,m)
            copy(s,cur)
            path = append(path,s)
            return
        }
        n := a[level-1]
        for i:=0;i<n;i++ {
            if orient[level-1]==0 { cur[level-1]=i } else { cur[level-1]=n-1-i }
            rec(level-1)
            if i!=n-1 && level>=2 { orient[level-2]^=1 }
        }
    }
    rec(m)
    return path
}

func diffInstr(x,y []int) string {
    m := len(x)
    for i:=0;i<m;i++ {
        if x[i]!=y[i] {
            if y[i]==x[i]+1 { return fmt.Sprintf("inc %d", i+1) }
            return fmt.Sprintf("dec %d", i+1)
        }
    }
    return ""
}

func expectedG(m int, a,b []int) string {
    for i:=0;i<m;i++ { b[i]-- }
    path := genPaths(m,a)
    p := len(path)
    pos :=0
    for i:=0;i<p;i++ {
        ok := true
        for j:=0;j<m;j++ { if path[i][j]!=b[j] { ok=false; break } }
        if ok { pos=i; break }
    }
    rotated := append(path[pos:], path[:pos]...)
    var sb strings.Builder
    sb.WriteString("Path\n")
    for i:=0;i<p-1;i++ {
        sb.WriteString(diffInstr(rotated[i],rotated[i+1]))
        sb.WriteString("\n")
    }
    return strings.TrimSpace(sb.String())
}

func genTestsG() []string {
    rand.Seed(7)
    tests := make([]string,0,100)
    for len(tests)<100 {
        m := rand.Intn(3)+1
        a := make([]int,m)
        b := make([]int,m)
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d\n", m))
        for i:=0;i<m;i++ {
            a[i] = rand.Intn(3)+1
            sb.WriteString(fmt.Sprintf("%d ", a[i]))
        }
        sb.WriteString("\n")
        for i:=0;i<m;i++ {
            b[i] = rand.Intn(a[i])+1
            sb.WriteString(fmt.Sprintf("%d ", b[i]))
        }
        sb.WriteString("\n")
        tests = append(tests,sb.String())
    }
    return tests
}

func runBinary(bin string, input string) (string,error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func parseInputG(t string) (int,[]int,[]int) {
    lines := strings.Split(strings.TrimSpace(t),"\n")
    var m int
    fmt.Sscanf(lines[0], "%d", &m)
    aParts := strings.Fields(lines[1])
    bParts := strings.Fields(lines[2])
    a := make([]int,m)
    b := make([]int,m)
    for i:=0;i<m;i++ {
        fmt.Sscanf(aParts[i], "%d", &a[i])
        fmt.Sscanf(bParts[i], "%d", &b[i])
    }
    return m,a,b
}

func main(){
    if len(os.Args)!=2 { fmt.Fprintf(os.Stderr,"Usage: go run verifierG.go <binary>\n"); os.Exit(1) }
    bin := os.Args[1]
    tests := genTestsG()
    for idx,t := range tests {
        m,a,b := parseInputG(t)
        want := expectedG(m, append([]int(nil), a...), append([]int(nil), b...))
        got, err := runBinary(bin,t)
        if err != nil {
            fmt.Printf("Test %d: runtime error: %v\n", idx+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got)!=want {
            fmt.Printf("Test %d failed.\nInput:\n%s\nExpected:\n%s\nGot:\n%s\n", idx+1, t, want, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}

