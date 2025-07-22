package main

import (
    "bufio"
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strconv"
    "strings"
)

type testB struct {
    n int
    k int
    a []int
}

func genTestsB() []testB {
    rng := rand.New(rand.NewSource(43))
    tests := make([]testB, 100)
    for i := range tests {
        n := rng.Intn(50) + 1
        k := rng.Intn(n) + 1
        arr := make([]int, n)
        for j := range arr {
            arr[j] = rng.Intn(101)
        }
        tests[i] = testB{n: n, k: k, a: arr}
    }
    return tests
}

func kthLargest(a []int, k int) int {
    b := append([]int(nil), a...)
    sort.Ints(b)
    return b[len(b)-k]
}

func run(bin string, input string) (string, error) {
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
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    return out.String(), nil
}

func runCase(bin string, tc testB) error {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
    for i, v := range tc.a {
        if i > 0 { sb.WriteByte(' ') }
        sb.WriteString(strconv.Itoa(v))
    }
    sb.WriteByte('\n')
    out, err := run(bin, sb.String())
    if err != nil { return err }
    scanner := bufio.NewScanner(strings.NewReader(out))
    scanner.Split(bufio.ScanWords)
    if !scanner.Scan() { return fmt.Errorf("missing threshold") }
    threshStr := scanner.Text()
    t, err := strconv.Atoi(threshStr)
    if err != nil { return fmt.Errorf("bad threshold") }
    indices := make([]int, 0, tc.k)
    for i := 0; i < tc.k; i++ {
        if !scanner.Scan() { return fmt.Errorf("missing index") }
        v, err := strconv.Atoi(scanner.Text())
        if err != nil { return fmt.Errorf("bad index") }
        indices = append(indices, v)
    }
    if scanner.Scan() { return fmt.Errorf("extra output") }
    seen := make(map[int]bool)
    minVal := 101
    for _, idx := range indices {
        if idx < 1 || idx > tc.n { return fmt.Errorf("index out of range") }
        if seen[idx] { return fmt.Errorf("duplicate index") }
        seen[idx] = true
        val := tc.a[idx-1]
        if val < minVal { minVal = val }
    }
    if minVal != t { return fmt.Errorf("reported threshold mismatch") }
    if t != kthLargest(tc.a, tc.k) { return fmt.Errorf("threshold not optimal") }
    if len(indices) != tc.k { return fmt.Errorf("incorrect number of indices") }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTestsB()
    for i, tc := range tests {
        if err := runCase(bin, tc); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(tests))
}

