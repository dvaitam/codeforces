package main

import (
    "fmt"
    "math/rand"
    "bytes"
    "os"
    "os/exec"
    "strings"
    "time"
)

type event struct {
    name  string
    score int
}

func winner(events []event) string {
    totals := make(map[string]int)
    for _, e := range events {
        totals[e.name] += e.score
    }
    maxScore := -1 << 60
    for _, v := range totals {
        if v > maxScore {
            maxScore = v
        }
    }
    cumulative := make(map[string]int)
    for _, e := range events {
        cumulative[e.name] += e.score
        if cumulative[e.name] >= maxScore && totals[e.name] == maxScore {
            return e.name
        }
    }
    return ""
}

func generateCase(rng *rand.Rand) (string, string) {
    players := []string{"alice", "bob", "carol", "dave", "eve"}
    n := rng.Intn(10) + 1
    events := make([]event, n)
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d\n", n))
    for i := 0; i < n; i++ {
        name := players[rng.Intn(len(players))]
        score := rng.Intn(21) - 10
        events[i] = event{name, score}
        sb.WriteString(fmt.Sprintf("%s %d\n", name, score))
    }
    // ensure at least one positive final score
    totals := make(map[string]int)
    for _, e := range events {
        totals[e.name] += e.score
    }
    positive := false
    for _, v := range totals {
        if v > 0 {
            positive = true
            break
        }
    }
    if !positive {
        events[0].score = rng.Intn(10) + 1
        sb.Reset()
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for _, e := range events {
            sb.WriteString(fmt.Sprintf("%s %d\n", e.name, e.score))
        }
    }
    return sb.String(), winner(events)
}

func runCase(exe string, input, expected string) error {
    cmd := exec.Command(exe)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    got := strings.TrimSpace(out.String())
    if got != expected {
        return fmt.Errorf("expected %s got %s", expected, got)
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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

