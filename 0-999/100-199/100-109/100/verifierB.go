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

type TestCase struct {
    input string
    ans   string
}

func solve(nums []int) string {
    n := len(nums)
    for i := 0; i < n; i++ {
        for j := i + 1; j < n; j++ {
            a, b := nums[i], nums[j]
            ok := false
            if a != 0 && b%a == 0 {
                ok = true
            } else if b != 0 && a%b == 0 {
                ok = true
            }
            if !ok {
                return "NOT FRIENDS"
            }
        }
    }
    return "FRIENDS"
}

func genTests() []TestCase {
    r := rand.New(rand.NewSource(2))
    cases := make([]TestCase, 100)
    for i := 0; i < 100; i++ {
        n := r.Intn(5) + 1
        arr := make([]int, n)
        for j := 0; j < n; j++ {
            v := r.Intn(21) - 10
            if v == 0 {
                v = 1
            }
            arr[j] = v
        }
        sort.Ints(arr)
        parts := make([]string, n)
        for j, v := range arr {
            parts[j] = fmt.Sprintf("%d", v)
        }
        input := fmt.Sprintf("%d\n%s\n", n, strings.Join(parts, ","))
        ans := solve(arr)
        cases[i] = TestCase{input: input, ans: ans}
    }
    return cases
}

func run(bin, in string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(in)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTests()
    for i, tc := range tests {
        got, err := run(bin, tc.input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != tc.ans {
            fmt.Fprintf(os.Stderr, "test %d failed: input %q expected %q got %q\n", i+1, tc.input, tc.ans, got)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

