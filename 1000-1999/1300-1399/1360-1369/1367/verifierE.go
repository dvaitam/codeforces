package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strconv"
)

func gcd(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}

func solveCase(n, k int, s string) int {
    var freq [26]int
    for _, ch := range s {
        freq[ch-'a']++
    }
    best := 0
    for length := 1; length <= n; length++ {
        g := gcd(length, k)
        cycleLen := length / g
        cycles := 0
        for _, cnt := range freq {
            cycles += cnt / cycleLen
        }
        if cycles >= g && length > best {
            best = length
        }
    }
    return best
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierE.go /path/to/binary")
        os.Exit(1)
    }
    data, err := os.ReadFile("testcasesE.txt")
    if err != nil {
        fmt.Println("could not read testcasesE.txt:", err)
        os.Exit(1)
    }
    scan := bufio.NewScanner(bytes.NewReader(data))
    scan.Split(bufio.ScanWords)
    if !scan.Scan() {
        fmt.Println("invalid test file")
        os.Exit(1)
    }
    t, _ := strconv.Atoi(scan.Text())
    cases := make([]struct{ n, k int; s string }, t)
    for i := 0; i < t; i++ {
        scan.Scan()
        n, _ := strconv.Atoi(scan.Text())
        scan.Scan()
        k, _ := strconv.Atoi(scan.Text())
        scan.Scan()
        s := scan.Text()
        cases[i] = struct{ n, k int; s string }{n, k, s}
    }
    expected := make([]int, t)
    for i, c := range cases {
        expected[i] = solveCase(c.n, c.k, c.s)
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

