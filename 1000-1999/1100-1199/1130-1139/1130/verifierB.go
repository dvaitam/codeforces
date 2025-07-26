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

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

type test struct {
    input    string
    expected string
}

func solve(input string) string {
    fields := strings.Fields(strings.TrimSpace(input))
    if len(fields) == 0 {
        return ""
    }
    idx := 0
    n, _ := strconv.Atoi(fields[idx])
    idx++
    m := 2 * n
    arr := make([]int, m)
    for i := 0; i < m; i++ {
        arr[i], _ = strconv.Atoi(fields[idx+i])
    }
    pos := make([][2]int, n+1)
    pos[0][0], pos[0][1] = 1, 1
    for i := 0; i < m; i++ {
        x := arr[i]
        if pos[x][0] == 0 {
            pos[x][0] = i + 1
        } else {
            pos[x][1] = i + 1
        }
    }
    var ans int64
    for i := 1; i <= n; i++ {
        p0, p1 := pos[i][0], pos[i][1]
        q0, q1 := pos[i-1][0], pos[i-1][1]
        cost1 := abs(p1-q1) + abs(p0-q0)
        cost2 := abs(p0-q1) + abs(p1-q0)
        if cost1 < cost2 {
            ans += int64(cost1)
        } else {
            ans += int64(cost2)
        }
    }
    return fmt.Sprint(ans)
}

func generateTests() []test {
    rand.Seed(1130)
    var tests []test
    fixed := []int{1, 2, 3, 4, 5}
    for _, n := range fixed {
        arr := make([]int, 2*n)
        for i := 0; i < n; i++ {
            arr[i] = i + 1
            arr[i+n] = i + 1
        }
        var sb strings.Builder
        sb.WriteString(fmt.Sprint(n))
        sb.WriteByte('\n')
        for i, v := range arr {
            if i > 0 {
                sb.WriteByte(' ')
            }
            sb.WriteString(fmt.Sprint(v))
        }
        sb.WriteByte('\n')
        inp := sb.String()
        tests = append(tests, test{inp, solve(inp)})
    }
    for len(tests) < 100 {
        n := rand.Intn(10) + 1
        arr := make([]int, 2*n)
        for i := 0; i < n; i++ {
            arr[i] = i + 1
            arr[i+n] = i + 1
        }
        rand.Shuffle(2*n, func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
        var sb strings.Builder
        sb.WriteString(fmt.Sprint(n))
        sb.WriteByte('\n')
        for i, v := range arr {
            if i > 0 {
                sb.WriteByte(' ')
            }
            sb.WriteString(fmt.Sprint(v))
        }
        sb.WriteByte('\n')
        inp := sb.String()
        tests = append(tests, test{inp, solve(inp)})
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
        fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
            fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected: %s\nGot: %s\n", i+1, t.input, strings.TrimSpace(t.expected), got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}

