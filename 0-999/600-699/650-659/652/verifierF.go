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

type pair struct { pos int64; idx int }

func floorDiv(a, b int64) int64 {
    if a >= 0 { return a / b }
    return -((-a + b - 1) / b)
}

func mod(a, b int64) int64 {
    r := a % b
    if r < 0 { r += b }
    return r
}

func solveF(in io.Reader) string {
    reader := bufio.NewReader(in)
    var n int
    var m, t int64
    if _, err := fmt.Fscan(reader, &n, &m, &t); err != nil {
        return ""
    }
    pos0 := make([]int64, n)
    dir := make([]int, n)
    initOrder := make([]pair, n)
    for i := 0; i < n; i++ {
        var s int64
        var d string
        fmt.Fscan(reader, &s, &d)
        pos0[i] = s - 1
        if d == "R" { dir[i] = 1 } else { dir[i] = -1 }
        initOrder[i] = pair{pos: pos0[i], idx: i}
    }
    sort.Slice(initOrder, func(i,j int) bool { return initOrder[i].pos < initOrder[j].pos })
    final := make([]int64, n)
    var shiftSum int64
    for i := 0; i < n; i++ {
        move := pos0[i] + int64(dir[i])*t
        final[i] = mod(move, m)
        shiftSum += floorDiv(move, m)
    }
    sort.Slice(final, func(i,j int) bool { return final[i] < final[j] })
    shift := mod(shiftSum, int64(n))
    res := make([]int64, n)
    for k := 0; k < n; k++ {
        idx := initOrder[k].idx
        res[idx] = final[(int64(k)+shift)%int64(n)] + 1
    }
    out := make([]string, n)
    for i, v := range res { out[i] = fmt.Sprint(v) }
    return strings.Join(out, " ")
}

func genTests() []string {
    rng := rand.New(rand.NewSource(6))
    tests := make([]string, 100)
    for i := 0; i < 100; i++ {
        n := rng.Intn(5) + 2
        m := int64(rng.Intn(90) + 10)
        t := int64(rng.Intn(1000))
        used := map[int]bool{}
        positions := make([]int64, n)
        for j := 0; j < n; j++ {
            p := rng.Intn(int(m)) + 1
            for used[p] { p = rng.Intn(int(m)) + 1 }
            used[p] = true
            positions[j] = int64(p)
        }
        dirs := make([]string, n)
        for j := 0; j < n; j++ {
            if rng.Intn(2) == 0 { dirs[j] = "L" } else { dirs[j] = "R" }
        }
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, t))
        for j := 0; j < n; j++ {
            sb.WriteString(fmt.Sprintf("%d %s\n", positions[j], dirs[j]))
        }
        tests[i] = sb.String()
    }
    return tests
}

func runBinary(bin, input string) (string, error) {
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
        fmt.Println("Usage: go run verifierF.go /path/to/binary")
        return
    }
    bin := os.Args[1]
    tests := genTests()
    for i, tc := range tests {
        expected := solveF(strings.NewReader(tc))
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

