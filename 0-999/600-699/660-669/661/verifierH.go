package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierH.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    data, err := os.ReadFile("testcasesH.txt")
    if err != nil {
        fmt.Println("could not read testcasesH.txt:", err)
        os.Exit(1)
    }
    scan := bufio.NewScanner(bytes.NewReader(data))
    scan.Split(bufio.ScanWords)
    if !scan.Scan() { fmt.Println("invalid test file"); os.Exit(1) }
    t, _ := strconv.Atoi(scan.Text())
    expected := make([]string, t)
    for i:=0;i<t;i++{
        if !scan.Scan(){ fmt.Println("bad test file"); os.Exit(1) }
        a,_:=strconv.Atoi(scan.Text()); scan.Scan(); b,_:=strconv.Atoi(scan.Text()); scan.Scan(); c,_:=strconv.Atoi(scan.Text()); scan.Scan(); d,_:=strconv.Atoi(scan.Text())
        expected[i] = fmt.Sprintf("%d %d %d %d", c,a,d,b)
    }
    cmd := exec.Command(bin)
    cmd.Stdin = bytes.NewReader(data)
    var out bytes.Buffer
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        fmt.Println("execution failed:", err)
        os.Exit(1)
    }
    outs := strings.Fields(strings.TrimSpace(out.String()))
    if len(outs)!=4*t { fmt.Printf("expected %d numbers got %d\n", 4*t,len(outs)); os.Exit(1) }
    for i:=0;i<t;i++{
        got := strings.Join(outs[i*4:(i+1)*4], " ")
        if got!=expected[i]{
            fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}
