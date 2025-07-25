package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strconv"
)

func solve(n int64) int64 {
    return n*(n+1)/2 + 1
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    data, err := os.ReadFile("testcasesA.txt")
    if err != nil {
        fmt.Println("could not read testcasesA.txt:", err)
        os.Exit(1)
    }

    // parse expected results
    scan := bufio.NewScanner(bytes.NewReader(data))
    scan.Split(bufio.ScanWords)
    if !scan.Scan() {
        fmt.Println("invalid test file")
        os.Exit(1)
    }
    t, _ := strconv.Atoi(scan.Text())
    expected := make([]int64, t)
    for i := 0; i < t; i++ {
        if !scan.Scan() {
            fmt.Println("bad test file")
            os.Exit(1)
        }
        n, _ := strconv.ParseInt(scan.Text(), 10, 64)
        expected[i] = solve(n)
    }

    cmd := exec.Command(bin)
    cmd.Stdin = bytes.NewReader(data)
    var out bytes.Buffer
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        fmt.Println("execution failed:", err)
        os.Exit(1)
    }

    outScan := bufio.NewScanner(bytes.NewReader(out.Bytes()))
    outScan.Split(bufio.ScanWords)
    for i := 0; i < t; i++ {
        if !outScan.Scan() {
            fmt.Printf("missing output for test %d\n", i+1)
            os.Exit(1)
        }
        got, _ := strconv.ParseInt(outScan.Text(), 10, 64)
        if got != expected[i] {
            fmt.Printf("test %d failed: expected %d got %d\n", i+1, expected[i], got)
            os.Exit(1)
        }
    }
    if outScan.Scan() {
        fmt.Println("extra output detected")
        os.Exit(1)
    }
    fmt.Println("All tests passed")
}

