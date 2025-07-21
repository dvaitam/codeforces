package main

import (
    "bytes"
    "fmt"
    "math"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "time"
)

func run(bin, input string) (string, error) {
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
    return strings.TrimSpace(out.String()), nil
}

func compute(n, m, a, b int, h [][]int64) [][3]int64 {
    occ := make([][]bool, n)
    for i := range occ {
        occ[i] = make([]bool, m)
    }
    res := make([][3]int64, 0)
    area := int64(a * b)
    for {
        bestCost := int64(math.MaxInt64)
        bi, bj := -1, -1
        for i := 0; i+a <= n; i++ {
            for j := 0; j+b <= m; j++ {
                overl := false
                for x := i; x < i+a && !overl; x++ {
                    for y := j; y < j+b; y++ {
                        if occ[x][y] {
                            overl = true
                            break
                        }
                    }
                }
                if overl {
                    continue
                }
                minv := h[i][j]
                sum := int64(0)
                for x := i; x < i+a; x++ {
                    for y := j; y < j+b; y++ {
                        v := h[x][y]
                        if v < minv {
                            minv = v
                        }
                        sum += v
                    }
                }
                cost := sum - int64(minv)*area
                if cost < bestCost || (cost == bestCost && (i < bi || (i == bi && j < bj))) {
                    bestCost = cost
                    bi, bj = i, j
                }
            }
        }
        if bi == -1 {
            break
        }
        res = append(res, [3]int64{int64(bi + 1), int64(bj + 1), bestCost})
        for x := bi; x < bi+a; x++ {
            for y := bj; y < bj+b; y++ {
                occ[x][y] = true
            }
        }
    }
    return res
}

func genCase(rng *rand.Rand) (string, string) {
    n := rng.Intn(4) + 2  // 2..5
    m := rng.Intn(4) + 2
    a := rng.Intn(n) + 1
    b := rng.Intn(m) + 1
    h := make([][]int64, n)
    for i := range h {
        h[i] = make([]int64, m)
        for j := range h[i] {
            h[i][j] = int64(rng.Intn(20))
        }
    }
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, m, a, b))
    for i := 0; i < n; i++ {
        for j := 0; j < m; j++ {
            if j > 0 {
                sb.WriteByte(' ')
            }
            sb.WriteString(fmt.Sprintf("%d", h[i][j]))
        }
        sb.WriteByte('\n')
    }
    res := compute(n, m, a, b, h)
    var exp strings.Builder
    exp.WriteString(fmt.Sprintf("%d\n", len(res)))
    for _, v := range res {
        exp.WriteString(fmt.Sprintf("%d %d %d\n", v[0], v[1], v[2]))
    }
    return sb.String(), strings.TrimSpace(exp.String())
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < 100; i++ {
        in, exp := genCase(rng)
        out, err := run(bin, in)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
            os.Exit(1)
        }
        if strings.TrimSpace(out) != exp {
            fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

