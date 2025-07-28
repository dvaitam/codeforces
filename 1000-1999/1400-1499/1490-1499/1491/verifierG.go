package main

import (
    "bufio"
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "time"
)

func run(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("%v\n%s", err, stderr.String())
    }
    return out.String(), nil
}

func genCase(rng *rand.Rand) (string, []int) {
    n := rng.Intn(7) + 3
    perm := rng.Perm(n)
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d\n", n))
    for i, v := range perm {
        if i > 0 { sb.WriteByte(' ') }
        sb.WriteString(fmt.Sprint(v + 1))
    }
    sb.WriteByte('\n')
    return sb.String(), perm
}

func check(n int, perm []int, output string) error {
    scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(output)))
    if !scanner.Scan() {
        return fmt.Errorf("missing q")
    }
    q, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
    if err != nil || q < 0 || q > n+1 {
        return fmt.Errorf("invalid q")
    }
    coins := make([]int, n)
    for i, v := range perm { coins[i] = v + 1 }
    up := make([]bool, n)
    for i := range up { up[i] = true }
    for i := 0; i < q; i++ {
        if !scanner.Scan() {
            return fmt.Errorf("expected %d operations", q)
        }
        fields := strings.Fields(scanner.Text())
        if len(fields) != 2 {
            return fmt.Errorf("invalid operation format")
        }
        a, err1 := strconv.Atoi(fields[0])
        b, err2 := strconv.Atoi(fields[1])
        if err1 != nil || err2 != nil || a < 1 || a > n || b < 1 || b > n || a == b {
            return fmt.Errorf("invalid indices")
        }
        a--; b--
        coins[a], coins[b] = coins[b], coins[a]
        up[a] = !up[a]
        up[b] = !up[b]
    }
    // ensure no extra tokens with non-space content
    for scanner.Scan() {
        if strings.TrimSpace(scanner.Text()) != "" {
            return fmt.Errorf("extra output")
        }
    }
    for i := 0; i < n; i++ {
        if coins[i] != i+1 || !up[i] {
            return fmt.Errorf("final state incorrect")
        }
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
        os.Exit(1)
    }
    candidate := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for tcase := 0; tcase < 100; tcase++ {
        input, perm := genCase(rng)
        out, err := run(candidate, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "runtime error on case %d: %v\ninput:\n%s", tcase+1, err, input)
            os.Exit(1)
        }
        if err := check(len(perm), perm, out); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", tcase+1, err, input, out)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}
