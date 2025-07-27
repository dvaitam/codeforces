package main

import (
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

func solve(n int) int {
    cards := make([]int,0)
    for h:=1;;h++{
        val := (3*h*h + h)/2
        if val>1e9{break}
        cards = append(cards,val)
    }
    cnt := 0
    for n >= 2 {
        best := -1
        for _,v := range cards{
            if v<=n{best=v}else{break}
        }
        if best==-1 {break}
        n-=best
        cnt++
    }
    return cnt
}

func main(){
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := make([]int,0)
    for i:=1;i<=100;i++{
        tests = append(tests,i)
    }
    extra := []int{150,200,300,500,1000,5000,10000,50000,100000,1000000}
    tests = append(tests, extra...)
    input := fmt.Sprintf("%d\n", len(tests))
    expected := make([]int,0,len(tests))
    for _,n := range tests{
        input += fmt.Sprintf("%d\n", n)
        expected = append(expected, solve(n))
    }
    out,err := run(bin,input)
    if err != nil {
        fmt.Printf("error executing %s: %v\nOutput:\n%s", bin, err, out)
        os.Exit(1)
    }
    lines := strings.Fields(strings.TrimSpace(out))
    if len(lines)!=len(expected){
        fmt.Printf("expected %d lines, got %d\n", len(expected), len(lines))
        os.Exit(1)
    }
    for i,exp := range expected{
        if lines[i] != fmt.Sprint(exp) {
            fmt.Printf("mismatch on test %d: expected %d got %s\n", i+1, exp, lines[i])
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed!")
}

func run(bin, input string)(string,error){
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    return out.String(), err
}
