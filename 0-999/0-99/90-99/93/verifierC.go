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

var (
    mul = [4]int{1, 2, 4, 8}
    node = [100]struct{ a, b, c, d int }{}
    target int
    steps int
    builder strings.Builder
    found bool
)

func dfs(x int) {
    if found {
        return
    }
    if x == steps {
        if node[x].a == target {
            builder.WriteString(fmt.Sprintf("%d\n", steps-1))
            for i := 2; i <= steps; i++ {
                reg := byte('a' + i - 1)
                builder.WriteString("lea e")
                builder.WriteByte(reg)
                builder.WriteString("x, [")
                if node[i].b != 0 {
                    base := byte('a' + node[i].b - 1)
                    builder.WriteString(fmt.Sprintf("e%cx + ", base))
                }
                if node[i].d != 0 {
                    builder.WriteString(fmt.Sprintf("%d*", mul[node[i].d]))
                }
                src := byte('a' + node[i].c - 1)
                builder.WriteString(fmt.Sprintf("e%cx]\n", src))
            }
            found = true
        }
        return
    }
    for i := 0; i <= x && !found; i++ {
        for j := 1; j <= x && !found; j++ {
            for k := 0; k < 4 && !found; k++ {
                A := node[i].a + mul[k]*node[j].a
                if A <= node[x].a || A > target {
                    continue
                }
                node[x+1] = struct{ a, b, c, d int }{A, i, j, k}
                dfs(x + 1)
            }
        }
    }
}

func compute(n int) string {
    target = n
    for steps = 1; ; steps++ {
        node[1] = struct{ a, b, c, d int }{1, 0, 0, 0}
        builder.Reset()
        found = false
        dfs(1)
        if found {
            return builder.String()
        }
    }
}

type testCase struct{
    n int
}

func runCase(bin string, tc testCase) error {
    input := fmt.Sprintf("%d\n", tc.n)
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    expected := strings.TrimSpace(compute(tc.n))
    got := strings.TrimSpace(out.String())
    if expected != got {
        return fmt.Errorf("expected:\n%s\n-- got:\n%s", expected, got)
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))

    var cases []testCase
    cases = append(cases, testCase{n: 1})
    cases = append(cases, testCase{n: 2})
    cases = append(cases, testCase{n: 255})
    for i:=0;i<100;i++ {
        n := rng.Intn(255) + 1
        cases = append(cases, testCase{n: n})
    }

    for i, tc := range cases {
        if err := runCase(bin, tc); err != nil {
            fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput: %d\n", i+1, err, tc.n)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

