package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierG.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    data, err := os.ReadFile("testcasesG.txt")
    if err != nil {
        fmt.Println("could not read testcasesG.txt:", err)
        os.Exit(1)
    }
    scan := bufio.NewScanner(bytes.NewReader(data))
    scan.Split(bufio.ScanWords)
    if !scan.Scan() { fmt.Println("invalid test file"); os.Exit(1) }
    var t int
    fmt.Sscan(scan.Text(), &t)
    expected := make([]string, t)
    for i:=0;i<t;i++{
        if !scan.Scan(){ fmt.Println("bad test file"); os.Exit(1) }
        expected[i] = "g_" + scan.Text()
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
    if len(outs)!=t { fmt.Printf("expected %d lines got %d\n", t,len(outs)); os.Exit(1) }
    for i:=0;i<t;i++{
        if outs[i]!=expected[i]{
            fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], outs[i])
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}
