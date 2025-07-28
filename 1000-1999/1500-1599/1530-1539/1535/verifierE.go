package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
)

type test struct {
    input    string
    expected string
}

const maxQ = 300000 + 5

func solve(input string) string {
    r := strings.NewReader(strings.TrimSpace(input))
    var q int
    var amount [maxQ]int64
    var cost [maxQ]int64
    var parent [maxQ]int
    var dsu [maxQ]int
    fmt.Fscan(r, &q, &amount[0], &cost[0])
    parent[0] = 0
    dsu[0] = 0
    var out strings.Builder
    idx := 1
    for ; idx <= q; idx++ {
        var t int
        fmt.Fscan(r, &t)
        if t == 1 {
            var p int
            var a, c int64
            fmt.Fscan(r, &p, &a, &c)
            parent[idx] = p
            amount[idx] = a
            cost[idx] = c
            dsu[idx] = idx
        } else {
            var v int
            var w int64
            fmt.Fscan(r, &v, &w)
            bought := int64(0)
            spent := int64(0)
            for w > 0 {
                x := v
                for dsu[x] != x {
                    dsu[x] = dsu[dsu[x]]
                    x = dsu[x]
                }
                if amount[x] == 0 {
                    if x == 0 {
                        break
                    }
                    dsu[x] = parent[x]
                    v = dsu[x]
                    continue
                }
                take := w
                if take > amount[x] {
                    take = amount[x]
                }
                amount[x] -= take
                w -= take
                bought += take
                spent += take * cost[x]
                if amount[x] == 0 && x != 0 {
                    dsu[x] = parent[x]
                    v = dsu[x]
                }
            }
            out.WriteString(fmt.Sprintf("%d %d\n", bought, spent))
        }
    }
    return out.String()
}

func generateTests() []test {
    rand.Seed(5)
    var tests []test
    // simple fixed test
    inp := "1 5 2\n2 0 2\n"
    tests = append(tests, test{inp, solve(inp)})
    for len(tests) < 100 {
        q := rand.Intn(10) + 1
        a0 := rand.Intn(5) + 1
        c0 := rand.Intn(5) + 1
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d %d %d\n", q, a0, c0))
        idx := 1
        for i := 0; i < q; i++ {
            if rand.Intn(2) == 0 {
                // add vertex
                p := rand.Intn(idx)
                a := rand.Intn(5) + 1
                c := rand.Intn(5) + c0 + 1
                sb.WriteString(fmt.Sprintf("1 %d %d %d\n", p, a, c))
                idx++
            } else {
                v := rand.Intn(idx)
                w := rand.Intn(5) + 1
                sb.WriteString(fmt.Sprintf("2 %d %d\n", v, w))
            }
        }
        inp := sb.String()
        tests = append(tests, test{inp, solve(inp)})
    }
    return tests
}

func runBinary(bin, input string) (string, error) {
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
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := generateTests()
    for i, t := range tests {
        got, err := runBinary(bin, t.input)
        if err != nil {
            fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
            os.Exit(1)
        }
        if got != strings.TrimSpace(t.expected) {
            fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%sGot:%s\n", i+1, t.input, t.expected, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}

