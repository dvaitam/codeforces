package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strings"
    "time"
)

type boat struct {
    typ int
    cap int
}

func solve(n, v int, bs []boat) int {
    var one, two []int
    for _, b := range bs {
        if b.typ == 1 {
            one = append(one, b.cap)
        } else {
            two = append(two, b.cap)
        }
    }
    sort.Slice(one, func(i, j int) bool { return one[i] > one[j] })
    sort.Slice(two, func(i, j int) bool { return two[i] > two[j] })
    pref1 := make([]int, len(one)+1)
    for i, v := range one {
        pref1[i+1] = pref1[i] + v
    }
    pref2 := make([]int, len(two)+1)
    for i, v := range two {
        pref2[i+1] = pref2[i] + v
    }
    best := 0
    for k := 0; k <= len(two) && 2*k <= v; k++ {
        rem := v - 2*k
        if rem > len(one) {
            rem = len(one)
        }
        val := pref2[k] + pref1[rem]
        if val > best {
            best = val
        }
    }
    return best
}

func generateCase(rng *rand.Rand) (string, int) {
    n := rng.Intn(6) + 1
    v := rng.Intn(10) + 1
    boats := make([]boat, n)
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d %d\n", n, v))
    for i := 0; i < n; i++ {
        typ := rng.Intn(2) + 1
        cap := rng.Intn(20) + 1
        boats[i] = boat{typ, cap}
        sb.WriteString(fmt.Sprintf("%d %d\n", typ, cap))
    }
    return sb.String(), solve(n, v, boats)
}

func runCase(exe, input string, expected int) error {
    cmd := exec.Command(exe)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    var got int
    if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
        return fmt.Errorf("bad output: %v", err)
    }
    if got != expected {
        return fmt.Errorf("expected %d got %d", expected, got)
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    exe := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < 100; i++ {
        in, exp := generateCase(rng)
        if err := runCase(exe, in, exp); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

