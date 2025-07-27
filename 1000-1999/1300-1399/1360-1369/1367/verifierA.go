package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

func solve(b string) string {
    n := len(b)
    if n == 0 {
        return ""
    }
    a := make([]byte, 0, n/2+1)
    a = append(a, b[0])
    for i := 1; i < n-1; i += 2 {
        a = append(a, b[i])
    }
    a = append(a, b[n-1])
    return string(a)
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    data, err := os.ReadFile("testcasesA.txt")
    if err != nil {
        fmt.Println("could not read testcasesA.txt:", err)
        os.Exit(1)
    }
    scan := bufio.NewScanner(bytes.NewReader(data))
    scan.Split(bufio.ScanLines)
    if !scan.Scan() {
        fmt.Println("empty test file")
        os.Exit(1)
    }
    var t int
    fmt.Sscan(scan.Text(), &t)
    inputs := make([]string, 0, t)
    for scan.Scan() {
        inputs = append(inputs, strings.TrimSpace(scan.Text()))
    }
    if len(inputs) != t {
        fmt.Println("invalid test count")
        os.Exit(1)
    }
    expected := make([]string, t)
    for i := 0; i < t; i++ {
        expected[i] = solve(inputs[i])
    }
    cmd := exec.Command(os.Args[1])
    cmd.Stdin = bytes.NewReader(data)
    out, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Println("execution failed:", err)
        os.Exit(1)
    }
    outScan := bufio.NewScanner(bytes.NewReader(out))
    outScan.Split(bufio.ScanWords)
    for i := 0; i < t; i++ {
        if !outScan.Scan() {
            fmt.Printf("missing output for test %d\n", i+1)
            os.Exit(1)
        }
        got := outScan.Text()
        if got != expected[i] {
            fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
            os.Exit(1)
        }
    }
    if outScan.Scan() {
        fmt.Println("extra output detected")
        os.Exit(1)
    }
    fmt.Println("All tests passed!")
}

