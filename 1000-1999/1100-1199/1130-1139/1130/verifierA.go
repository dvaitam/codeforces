package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

type test struct {
    input string
}

func parseInput(input string) (n int, arr []int) {
    fields := strings.Fields(strings.TrimSpace(input))
    if len(fields) == 0 {
        return 0, nil
    }
    n, _ = strconv.Atoi(fields[0])
    arr = make([]int, n)
    for i := 0; i < n; i++ {
        arr[i], _ = strconv.Atoi(fields[i+1])
    }
    return
}

func isValidAnswer(arr []int, out string) bool {
    d, err := strconv.Atoi(strings.TrimSpace(out))
    if err != nil {
        return false
    }
    if d < -1000 || d > 1000 {
        return false
    }
    n := len(arr)
    posCount, negCount := 0, 0
    for _, v := range arr {
        if v > 0 {
            posCount++
        } else if v < 0 {
            negCount++
        }
    }
    half := (n + 1) / 2
    if posCount < half && negCount < half {
        return d == 0
    }
    if d == 0 {
        return false
    }
    if d > 0 {
        return posCount >= half
    }
    return negCount >= half
}

func generateTests() []test {
    rand.Seed(1130)
    var tests []test
    fixed := [][]int{
        {1, 1},
        {1, -1},
        {3, 1, -2, 3},
        {5, -1, -2, -3, 0, 0},
        {4, 0, 0, 0, 0},
    }
    for _, arr := range fixed {
        n := arr[0]
        data := arr[1:]
        var sb strings.Builder
        sb.WriteString(fmt.Sprint(n))
        sb.WriteByte('\n')
        for i, v := range data {
            if i > 0 {
                sb.WriteByte(' ')
            }
            sb.WriteString(fmt.Sprint(v))
        }
        sb.WriteByte('\n')
        inp := sb.String()
        tests = append(tests, test{inp})
    }
    for len(tests) < 100 {
        n := rand.Intn(20) + 1
        var sb strings.Builder
        sb.WriteString(fmt.Sprint(n))
        sb.WriteByte('\n')
        for i := 0; i < n; i++ {
            if i > 0 {
                sb.WriteByte(' ')
            }
            val := rand.Intn(2001) - 1000
            sb.WriteString(fmt.Sprint(val))
        }
        sb.WriteByte('\n')
        tests = append(tests, test{sb.String()})
    }
    return tests
}

func runBinary(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := generateTests()
    for i, t := range tests {
        n, arr := parseInput(t.input)
        _ = n
        got, err := runBinary(bin, t.input)
        if err != nil {
            fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
            os.Exit(1)
        }
        if !isValidAnswer(arr, got) {
            fmt.Printf("Wrong answer on test %d\nInput:\n%sOutput: %s\n", i+1, t.input, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}

