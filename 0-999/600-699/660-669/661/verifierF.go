package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strconv"
)

func sieve(n int) []bool {
    prime := make([]bool, n+1)
    for i := 2; i <= n; i++ { prime[i] = true }
    for p := 2; p*p <= n; p++ {
        if prime[p] {
            for j := p * p; j <= n; j += p {
                prime[j] = false
            }
        }
    }
    return prime
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierF.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    data, err := os.ReadFile("testcasesF.txt")
    if err != nil {
        fmt.Println("could not read testcasesF.txt:", err)
        os.Exit(1)
    }
    scan := bufio.NewScanner(bytes.NewReader(data))
    scan.Split(bufio.ScanWords)
    if !scan.Scan() { fmt.Println("invalid test file"); os.Exit(1) }
    t, _ := strconv.Atoi(scan.Text())
    expected := make([]int, t)
    primes := sieve(500)
    for i := 0; i < t; i++ {
        if !scan.Scan() { fmt.Println("bad test file"); os.Exit(1) }
        l, _ := strconv.Atoi(scan.Text())
        scan.Scan(); r, _ := strconv.Atoi(scan.Text())
        cnt := 0
        for x := l; x <= r; x++ { if primes[x] { cnt++ } }
        expected[i] = cnt
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
        if !outScan.Scan() { fmt.Printf("missing output for test %d\n", i+1); os.Exit(1) }
        got, _ := strconv.Atoi(outScan.Text())
        if got != expected[i] {
            fmt.Printf("test %d failed: expected %d got %d\n", i+1, expected[i], got)
            os.Exit(1)
        }
    }
    if outScan.Scan() { fmt.Println("extra output detected"); os.Exit(1) }
    fmt.Println("All tests passed")
}
