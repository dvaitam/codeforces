package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

func season(m int) string {
    switch m {
    case 12,1,2:
        return "Winter"
    case 3,4,5:
        return "Spring"
    case 6,7,8:
        return "Summer"
    default:
        return "Autumn"
    }
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
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
    t := 0
    fmt.Sscan(scan.Text(), &t)
    expected := make([]string, t)
    for i := 0; i < t; i++ {
        if !scan.Scan() {
            fmt.Println("bad test file")
            os.Exit(1)
        }
        var m int
        fmt.Sscan(scan.Text(), &m)
        expected[i] = season(m)
    }

    cmd := exec.Command(bin)
    cmd.Stdin = bytes.NewReader(data)
    var out bytes.Buffer
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        fmt.Println("execution failed:", err)
        os.Exit(1)
    }

    outLines := strings.Fields(strings.TrimSpace(out.String()))
    if len(outLines) != t {
        fmt.Printf("expected %d lines of output got %d\n", t, len(outLines))
        os.Exit(1)
    }
    for i := 0; i < t; i++ {
        if outLines[i] != expected[i] {
            fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], outLines[i])
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

