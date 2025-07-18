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

func dist(s, t string) int {
    dx := int(s[0]) - int(t[0])
    if dx < 0 {
        dx = -dx
    }
    dy := int(s[1]) - int(t[1])
    if dy < 0 {
        dy = -dy
    }
    if dx > dy {
        return dx
    }
    return dy
}

func generateCase(rng *rand.Rand) (string, int) {
    letters := "abcdefgh"
    nums := "12345678"
    s := []byte{letters[rng.Intn(8)], nums[rng.Intn(8)]}
    t := []byte{letters[rng.Intn(8)], nums[rng.Intn(8)]}
    input := fmt.Sprintf("%s\n%s\n", s, t)
    return input, dist(string(s), string(t))
}

func runCase(exe string, input string, expected int) error {
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

