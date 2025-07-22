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

func solveA(arr []int, x, y int) int {
    sum2 := 0
    sum1 := 0
    for _, v := range arr {
        sum1 += v
    }
    for i, v := range arr {
        sum2 += v
        sum1 -= v
        if sum2 >= x && sum2 <= y && sum1 >= x && sum1 <= y {
            return i + 2
        }
    }
    return 0
}

func genCase(rng *rand.Rand) (string, string) {
    m := rng.Intn(9) + 2 // 2..10
    arr := make([]int, m)
    total := 0
    for i := 0; i < m; i++ {
        arr[i] = rng.Intn(21)
        total += arr[i]
    }
    if total == 0 {
        idx := rng.Intn(m)
        arr[idx] = rng.Intn(5) + 1
        total = arr[idx]
    }
    x := rng.Intn(30) + 1
    y := rng.Intn(11) + x // ensure y>=x
    if y > 40 {
        y = 40
    }
    var sb strings.Builder
    fmt.Fprintf(&sb, "%d\n", m)
    for i, v := range arr {
        if i > 0 {
            sb.WriteByte(' ')
        }
        fmt.Fprintf(&sb, "%d", v)
    }
    fmt.Fprintf(&sb, "\n%d %d\n", x, y)

    exp := fmt.Sprintf("%d\n", solveA(arr, x, y))
    return sb.String(), exp
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
        fmt.Println("usage: go run verifierA.go /path/to/binary")
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
