package main

import (
    "bytes"
    "fmt"
    "math"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strings"
)

func expectedB(n int, w []int) int {
    sort.Ints(w)
    ans := math.MaxInt32
    for i := 0; i < 2*n; i++ {
        for j := i + 1; j < 2*n; j++ {
            arr := make([]int, 0, 2*n-2)
            for k := 0; k < 2*n; k++ {
                if k == i || k == j {
                    continue
                }
                arr = append(arr, w[k])
            }
            sum := 0
            for k := 0; k < len(arr); k += 2 {
                sum += arr[k+1] - arr[k]
            }
            if sum < ans {
                ans = sum
            }
        }
    }
    return ans
}

func genTestsB() []string {
    rand.Seed(2)
    tests := make([]string, 0, 100)
    for len(tests) < 100 {
        n := rand.Intn(49) + 2
        w := make([]int, 2*n)
        for i := range w {
            w[i] = rand.Intn(1000) + 1
        }
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for i, val := range w {
            if i > 0 {
                sb.WriteString(" ")
            }
            sb.WriteString(fmt.Sprintf("%d", val))
        }
        sb.WriteString("\n")
        tests = append(tests, sb.String())
    }
    return tests
}

func runBinary(bin string, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func parseWeights(line string) []int {
    parts := strings.Fields(line)
    res := make([]int, len(parts))
    for i, p := range parts {
        fmt.Sscanf(p, "%d", &res[i])
    }
    return res
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: go run verifierB.go <binary>\n")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTestsB()
    for idx, t := range tests {
        lines := strings.Split(strings.TrimSpace(t), "\n")
        var n int
        fmt.Sscanf(lines[0], "%d", &n)
        weights := parseWeights(lines[1])
        want := fmt.Sprintf("%d", expectedB(n, append([]int(nil), weights...)))
        got, err := runBinary(bin, t)
        if err != nil {
            fmt.Printf("Test %d: runtime error: %v\n", idx+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != want {
            fmt.Printf("Test %d failed.\nInput:\n%s\nExpected: %s\nGot: %s\n", idx+1, t, want, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}

