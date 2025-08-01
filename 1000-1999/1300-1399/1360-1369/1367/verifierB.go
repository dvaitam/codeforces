package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strconv"
)

func solveCase(arr []int) int {
    mismatchEven := 0
    mismatchOdd := 0
    for i, v := range arr {
        if i%2 != v%2 {
            if i%2 == 0 {
                mismatchEven++
            } else {
                mismatchOdd++
            }
        }
    }
    if mismatchEven == mismatchOdd {
        return mismatchEven
    }
    return -1
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    data, err := os.ReadFile("testcasesB.txt")
    if err != nil {
        fmt.Println("could not read testcasesB.txt:", err)
        os.Exit(1)
    }
    scan := bufio.NewScanner(bytes.NewReader(data))
    scan.Split(bufio.ScanWords)
    if !scan.Scan() {
        fmt.Println("invalid test file")
        os.Exit(1)
    }
    t, _ := strconv.Atoi(scan.Text())
    inputs := make([][]int, t)
    for i := 0; i < t; i++ {
        if !scan.Scan() {
            fmt.Println("missing n")
            os.Exit(1)
        }
        n, _ := strconv.Atoi(scan.Text())
        arr := make([]int, n)
        for j := 0; j < n; j++ {
            if !scan.Scan() {
                fmt.Println("missing array element")
                os.Exit(1)
            }
            arr[j], _ = strconv.Atoi(scan.Text())
        }
        inputs[i] = arr
    }
    expected := make([]int, t)
    for i := 0; i < t; i++ {
        expected[i] = solveCase(inputs[i])
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

