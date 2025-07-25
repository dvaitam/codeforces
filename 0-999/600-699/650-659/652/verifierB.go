package main

import (
    "bufio"
    "bytes"
    "context"
    "fmt"
    "io"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strings"
    "time"
)

type void struct{}

func solveB(in io.Reader) string {
    reader := bufio.NewReader(in)
    var n int
    if _, err := fmt.Fscan(reader, &n); err != nil {
        return ""
    }
    arr := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &arr[i])
    }
    sort.Ints(arr)
    ans := make([]int, n)
    ptr := n - 1
    for i := 1; i < n; i += 2 {
        ans[i] = arr[ptr]
        ptr--
    }
    for i := 0; i < n; i += 2 {
        ans[i] = arr[ptr]
        ptr--
    }
    out := make([]string, n)
    for i, v := range ans {
        out[i] = fmt.Sprint(v)
    }
    return strings.Join(out, " ")
}

func genTests() []string {
    r := rand.New(rand.NewSource(2))
    tests := make([]string, 100)
    for i := 0; i < 100; i++ {
        n := r.Intn(10) + 1
        arr := make([]int, n)
        for j := range arr {
            arr[j] = r.Intn(1000) + 1
        }
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for j, v := range arr {
            if j > 0 {
                sb.WriteByte(' ')
            }
            sb.WriteString(fmt.Sprint(v))
        }
        sb.WriteByte('\n')
        tests[i] = sb.String()
    }
    return tests
}

func runBinary(bin string, input string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    cmd := exec.CommandContext(ctx, bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run verifierB.go /path/to/binary")
        return
    }
    bin := os.Args[1]
    tests := genTests()
    for i, tc := range tests {
        expected := solveB(strings.NewReader(tc))
        actual, err := runBinary(bin, tc)
        if err != nil {
            fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
            return
        }
        if actual != strings.TrimSpace(expected) {
            fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", i+1, tc, expected, actual)
            return
        }
    }
    fmt.Println("All tests passed!")
}

