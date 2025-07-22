package main

import (
    "bufio"
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strings"
)

type team struct {
    name   string
    rating int
}

type testE struct {
    n int
    x int
    a int
    b int
    c int
    teams []team
}

func genTestsE() []testE {
    rng := rand.New(rand.NewSource(46))
    tests := make([]testE, 100)
    for i := range tests {
        m := (rng.Intn(3)+1)*4 // 4,8,12,16
        x := rng.Intn(100) + 1
        a := rng.Intn(100) + 1
        b := rng.Intn(100) + 1
        c := rng.Intn(100) + 1
        teams := make([]team, m)
        ratings := rng.Perm(m*5 + 50)
        for j := 0; j < m; j++ {
            teams[j] = team{fmt.Sprintf("T%d_%d", i, j), ratings[j]}
        }
        tests[i] = testE{n: m, x: x, a: a, b: b, c: c, teams: teams}
    }
    return tests
}

func solveE(tc testE) [][]string {
    x := tc.x
    a := tc.a
    b := tc.b
    c := tc.c
    n := tc.n
    teams := append([]team(nil), tc.teams...)
    sort.Slice(teams, func(i, j int) bool { return teams[i].rating > teams[j].rating })
    m := n / 4
    baskets := make([][]team, 4)
    for i := 0; i < 4; i++ {
        baskets[i] = append([]team(nil), teams[i*m:(i+1)*m]...)
    }
    rng := func() int { x = (x*a + b) % c; return x }
    groups := make([][]team, m)
    for gi := 0; gi < m-1; gi++ {
        var group []team
        for bi := 0; bi < 4; bi++ {
            k := rng()
            idx := k % len(baskets[bi])
            group = append(group, baskets[bi][idx])
            baskets[bi] = append(baskets[bi][:idx], baskets[bi][idx+1:]...)
        }
        sort.Slice(group, func(i, j int) bool { return group[i].rating > group[j].rating })
        groups[gi] = group
    }
    var last []team
    for bi := 0; bi < 4; bi++ {
        last = append(last, baskets[bi][0])
    }
    sort.Slice(last, func(i, j int) bool { return last[i].rating > last[j].rating })
    groups[m-1] = last
    res := make([][]string, m)
    for i, g := range groups {
        names := make([]string, len(g))
        for j, t := range g {
            names[j] = t.name
        }
        res[i] = names
    }
    return res
}

func run(bin string, input string) (string, error) {
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
    return out.String(), nil
}

func runCase(bin string, tc testE) error {
    var sb strings.Builder
    fmt.Fprintf(&sb, "%d\n%d %d %d %d\n", tc.n, tc.x, tc.a, tc.b, tc.c)
    for _, t := range tc.teams {
        fmt.Fprintf(&sb, "%s %d\n", t.name, t.rating)
    }
    out, err := run(bin, sb.String())
    if err != nil { return err }
    expected := solveE(tc)
    scanner := bufio.NewScanner(strings.NewReader(out))
    for i, group := range expected {
        if !scanner.Scan() { return fmt.Errorf("missing group letter") }
        letter := strings.TrimSpace(scanner.Text())
        if letter != string('A'+i) { return fmt.Errorf("wrong group letter") }
        for range group {
            if !scanner.Scan() { return fmt.Errorf("missing team name") }
            name := strings.TrimSpace(scanner.Text())
            valid := false
            for _, nm := range group {
                if nm == name { valid = true; break }
            }
            if !valid { return fmt.Errorf("unexpected team name %s", name) }
        }
    }
    if scanner.Scan() { return fmt.Errorf("extra output") }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTestsE()
    for i, tc := range tests {
        if err := runCase(bin, tc); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(tests))
}

