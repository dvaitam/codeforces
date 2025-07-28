package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strings"
)

type test struct {
    input    string
    expected string
}

func solve(input string) string {
    r := strings.NewReader(strings.TrimSpace(input))
    var t int
    fmt.Fscan(r, &t)
    var out strings.Builder
    for ; t > 0; t-- {
        var s1, s2, s3, s4 int
        fmt.Fscan(r, &s1, &s2, &s3, &s4)
        arr := []int{s1, s2, s3, s4}
        sort.Ints(arr)
        top1 := s1
        if s2 > top1 {
            top1 = s2
        }
        top2 := s3
        if s4 > top2 {
            top2 = s4
        }
        if (top1 == arr[3] && top2 == arr[2]) || (top1 == arr[2] && top2 == arr[3]) {
            out.WriteString("YES\n")
        } else {
            out.WriteString("NO\n")
        }
    }
    return out.String()
}

func generateTests() []test {
    rand.Seed(1)
    var tests []test
    fixed := [][4]int{
        {7, 3, 9, 1},
        {5, 6, 3, 2},
        {1, 2, 3, 4},
        {10, 1, 8, 9},
        {9, 8, 7, 6},
    }
    for _, arr := range fixed {
        inp := fmt.Sprintf("1\n%d %d %d %d\n", arr[0], arr[1], arr[2], arr[3])
        tests = append(tests, test{inp, solve(inp)})
    }
    for len(tests) < 100 {
        arr := [4]int{rand.Intn(100) + 1, rand.Intn(100) + 1, rand.Intn(100) + 1, rand.Intn(100) + 1}
        inp := fmt.Sprintf("1\n%d %d %d %d\n", arr[0], arr[1], arr[2], arr[3])
        tests = append(tests, test{inp, solve(inp)})
    }
    return tests
}

func runBinary(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
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
        got, err := runBinary(bin, t.input)
        if err != nil {
            fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
            os.Exit(1)
        }
        if got != strings.TrimSpace(t.expected) {
            fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%sGot:%s\n", i+1, t.input, t.expected, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}

