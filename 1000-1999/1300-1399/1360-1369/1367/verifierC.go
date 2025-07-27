package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strconv"
)

func solveCase(n, k int, s string) int {
    right := make([]int, n)
    next := n * 2
    for i := n - 1; i >= 0; i-- {
        if s[i] == '1' {
            next = i
        }
        right[i] = next
    }
    last := -n * 2
    ans := 0
    for i := 0; i < n; i++ {
        if s[i] == '1' {
            last = i
            continue
        }
        if i-last > k && right[i]-i > k {
            ans++
            last = i
        }
    }
    return ans
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    data, err := os.ReadFile("testcasesC.txt")
    if err != nil {
        fmt.Println("could not read testcasesC.txt:", err)
        os.Exit(1)
    }
    scan := bufio.NewScanner(bytes.NewReader(data))
    scan.Split(bufio.ScanWords)
    if !scan.Scan() {
        fmt.Println("invalid test file")
        os.Exit(1)
    }
    t, _ := strconv.Atoi(scan.Text())
    types := make([]struct{
        n int
        k int
        s string
    }, t)
    for i := 0; i < t; i++ {
        if !scan.Scan() {
            fmt.Println("missing n")
            os.Exit(1)
        }
        n, _ := strconv.Atoi(scan.Text())
        scan.Scan()
        k, _ := strconv.Atoi(scan.Text())
        scan.Scan()
        s := scan.Text()
        types[i] = struct{n int;k int;s string}{n,k,s}
    }
    expected := make([]int, t)
    for i := range types {
        expected[i] = solveCase(types[i].n, types[i].k, types[i].s)
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
        got, _ := strconv.Atoi(outScan.Text())
        if got != expected[i] {
            fmt.Printf("test %d failed: expected %d got %d\n", i+1, expected[i], got)
            os.Exit(1)
        }
    }
    if outScan.Scan() {
        fmt.Println("extra output detected")
        os.Exit(1)
    }
    fmt.Println("All tests passed!")
}

