package main

import (
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

type test struct{n,m int}

func solve(n,m int) string {
    if n==1 || m==1 || (n==2 && m==2) {
        return "YES"
    }
    return "NO"
}

func main(){
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := make([]test,0)
    // generate 120 test cases
    for i:=1;i<=20;i++{
        for j:=1;j<=6;j++{
            tests = append(tests,test{i,j})
        }
    }
    // add some larger cases
    large := []test{{100,100},{1,1000},{1000,1},{20000,20000}}
    tests = append(tests, large...)
    input := fmt.Sprintf("%d\n", len(tests))
    expected := make([]string,0,len(tests))
    for _,t := range tests{
        input += fmt.Sprintf("%d %d\n", t.n,t.m)
        expected = append(expected, solve(t.n,t.m))
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
        if strings.ToUpper(lines[i]) != exp {
            fmt.Printf("mismatch on test %d: expected %s got %s\n", i+1, exp, lines[i])
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed!")
}

func run(bin,stringInput string)(string,error){
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(stringInput)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    return out.String(), err
}
