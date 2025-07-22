package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "time"
)

func solveB(n, m int, dances [][3]int) []int {
    a := make([]int, n+1)
    for i := 0; i < m; i++ {
        t := dances[i]
        used := [4]bool{}
        for j := 0; j < 3; j++ {
            used[a[t[j]]] = true
        }
        for j := 0; j < 3; j++ {
            if a[t[j]] == 0 {
                for k := 1; k <= 3; k++ {
                    if !used[k] {
                        a[t[j]] = k
                        used[k] = true
                        break
                    }
                }
            }
        }
    }
    return a[1:]
}

func genCase(rng *rand.Rand) (string, string) {
    n := rng.Intn(6) + 3  // 3..8
    m := rng.Intn(8) + 1  // 1..8
    dances := make([][3]int, m)
    var sb strings.Builder
    fmt.Fprintf(&sb, "%d %d\n", n, m)
    for i := 0; i < m; i++ {
        perm := rng.Perm(n)
        dances[i] = [3]int{perm[0]+1, perm[1]+1, perm[2]+1}
        fmt.Fprintf(&sb, "%d %d %d\n", dances[i][0], dances[i][1], dances[i][2])
    }
    sb.WriteByte('\n')
    colors := solveB(n, m, dances)
    var out strings.Builder
    for i, c := range colors {
        if i > 0 {
            out.WriteByte(' ')
        }
        fmt.Fprintf(&out, "%d", c)
    }
    out.WriteByte('\n')
    return sb.String(), out.String()
}

func runCase(bin, input, expected string) error {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    got := strings.TrimSpace(out.String())
    if got != strings.TrimSpace(expected) {
        return fmt.Errorf("expected %s got %s", strings.TrimSpace(expected), got)
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < 100; i++ {
        in, exp := genCase(rng)
        if err := runCase(bin, in, exp); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

