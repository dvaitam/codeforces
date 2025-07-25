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

func compileRef() (string, error) {
    ref := "refE"
    cmd := exec.Command("go", "build", "-o", ref, "707E.go")
    if err := cmd.Run(); err != nil {
        return "", err
    }
    return ref, nil
}

func run(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func generateTest() string {
    n := rand.Intn(3) + 1
    m := rand.Intn(3) + 1
    k := rand.Intn(3) + 1
    used := make(map[[2]int]bool)
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
    for gi := 0; gi < k; gi++ {
        sb.WriteString("1\n")
        var x, y int
        for {
            x = rand.Intn(n) + 1
            y = rand.Intn(m) + 1
            if !used[[2]int{x,y}] {
                used[[2]int{x,y}] = true
                break
            }
        }
        w := rand.Intn(10) + 1
        sb.WriteString(fmt.Sprintf("%d %d %d\n", x, y, w))
    }
    q := rand.Intn(20) + 1
    sb.WriteString(fmt.Sprintf("%d\n", q))
    for i := 0; i < q; i++ {
        if rand.Intn(2) == 0 {
            sb.WriteString(fmt.Sprintf("SWITCH %d\n", rand.Intn(k)+1))
        } else {
            x1 := rand.Intn(n) + 1
            y1 := rand.Intn(m) + 1
            x2 := rand.Intn(n-x1+1) + x1
            y2 := rand.Intn(m-y1+1) + y1
            sb.WriteString(fmt.Sprintf("ASK %d %d %d %d\n", x1, y1, x2, y2))
        }
    }
    return sb.String()
}

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
        os.Exit(1)
    }
    rand.Seed(time.Now().UnixNano())
    ref, err := compileRef()
    if err != nil {
        fmt.Fprintln(os.Stderr, "could not build reference solution:", err)
        os.Exit(1)
    }
    defer os.Remove(ref)

    bin := os.Args[1]
    for t := 0; t < 100; t++ {
        input := generateTest()
        want, err := run(ref, input)
        if err != nil {
            fmt.Fprintln(os.Stderr, "reference failed:", err)
            os.Exit(1)
        }
        got, err := run(bin, input)
        if err != nil {
            fmt.Fprintln(os.Stderr, "test", t+1, "error running binary:", err)
            os.Exit(1)
        }
        if got != want {
            fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected: %s\nactual: %s\n", t+1, input, want, got)
            os.Exit(1)
        }
    }
    fmt.Println("ok")
}

