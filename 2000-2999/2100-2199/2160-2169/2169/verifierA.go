package main

import (
    "bufio"
    "bytes"
    "fmt"
    "io"
    "math"
    "os"
    "os/exec"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
)

const refSource = "2000-2999/2100-2199/2160-2169/2169/2169A.go"

type testCase struct {
    n int
    a int64
    vals []int64
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    candidate := os.Args[1]

    cases, inputBytes, err := readInput()
    if err != nil {
        fmt.Fprintln(os.Stderr, "failed to read input:", err)
        os.Exit(1)
    }

    refBin, cleanup, err := buildReference()
    if err != nil {
        fmt.Fprintln(os.Stderr, "failed to build reference:", err)
        os.Exit(1)
    }
    defer cleanup()

    refOut, err := runProgram(refBin, inputBytes)
    if err != nil {
        fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
        os.Exit(1)
    }

    candOut, err := runProgram(candidate, inputBytes)
    if err != nil {
        fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
        os.Exit(1)
    }

    refAns := strings.Fields(refOut)
    candAns := strings.Fields(candOut)
    if len(refAns) != len(cases) || len(candAns) != len(cases) {
        fmt.Fprintf(os.Stderr, "output length mismatch: expected %d answers\nreference output:\n%s\ncandidate output:\n%s\n", len(cases), refOut, candOut)
        os.Exit(1)
    }

    for idx, tc := range cases {
        bestScore := computeBestScore(tc)
        candB, err := strconv.ParseInt(candAns[idx], 10, 64)
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: invalid candidate integer %q\n", idx+1, candAns[idx])
            os.Exit(1)
        }
        candScore := evaluateScore(tc, candB)
        if candScore < bestScore {
            fmt.Fprintf(os.Stderr, "test %d failed: candidate score %d < optimal %d (candidate b=%d)\n", idx+1, candScore, bestScore, candB)
            os.Exit(1)
        }
    }

    fmt.Println("Accepted")
}

func readInput() ([]testCase, []byte, error) {
    data, err := io.ReadAll(os.Stdin)
    if err != nil {
        return nil, nil, err
    }
    reader := bufio.NewReader(bytes.NewReader(data))
    var t int
    if _, err := fmt.Fscan(reader, &t); err != nil {
        return nil, nil, err
    }
    cases := make([]testCase, t)
    for i := 0; i < t; i++ {
        var n int
        var a int64
        if _, err := fmt.Fscan(reader, &n, &a); err != nil {
            return nil, nil, err
        }
        vals := make([]int64, n)
        for j := 0; j < n; j++ {
            if _, err := fmt.Fscan(reader, &vals[j]); err != nil {
                return nil, nil, err
            }
        }
        cases[i] = testCase{n: n, a: a, vals: vals}
    }
    return cases, data, nil
}

func runProgram(bin string, input []byte) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = bytes.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    return out.String(), err
}

func buildReference() (string, func(), error) {
    dir, err := os.MkdirTemp("", "ref2169A-")
    if err != nil {
        return "", nil, err
    }
    bin := filepath.Join(dir, "ref2169A.bin")
    cmd := exec.Command("go", "build", "-o", bin, refSource)
    var stderr bytes.Buffer
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        os.RemoveAll(dir)
        return "", nil, fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
    }
    return bin, func() { os.RemoveAll(dir) }, nil
}

func computeBestScore(tc testCase) int {
    if tc.n == 0 {
        return 0
    }
    events := make(map[int64]int)
    for _, v := range tc.vals {
        dist := abs(v - tc.a)
        if dist == 0 {
            continue
        }
        left := v - dist + 1
        right := v + dist - 1
        if left < 0 {
            left = 0
        }
        if right < 0 {
            continue
        }
        events[left]++
        events[right+1]--
    }
    if len(events) == 0 {
        return 0
    }
    keys := make([]int64, 0, len(events))
    for k := range events {
        keys = append(keys, k)
    }
    sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
    cur := 0
    best := 0
    for _, x := range keys {
        cur += events[x]
        if cur > best {
            best = cur
        }
    }
    return best
}

func evaluateScore(tc testCase, b int64) int {
    score := 0
    for _, v := range tc.vals {
        distA := abs(v - tc.a)
        distB := abs(v - b)
        if distB < distA || (distB == distA && v > tc.a) {
            score++
        }
    }
    return score
}

func abs(x int64) int64 {
    if x < 0 {
        return -x
    }
    return x
}
