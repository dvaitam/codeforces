package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strings"
)

type TestCase struct {
    input string
    ans   string
}

type Interval struct {
    l, r int
}

type Clique struct {
    unionL, unionR int
    comL, comR     int
}

type Color struct {
    cliques []Clique
}

func (c *Color) canUse(x Interval) bool {
    count := 0
    for _, cl := range c.cliques {
        if x.r < cl.unionL || x.l > cl.unionR {
            continue
        }
        if x.r >= cl.comL && x.l <= cl.comR {
            count++
            if count > 1 {
                return false
            }
        } else {
            return false
        }
    }
    return true
}

func (c *Color) add(x Interval) {
    for i := range c.cliques {
        cl := &c.cliques[i]
        if x.r >= cl.comL && x.l <= cl.comR {
            if x.l < cl.unionL {
                cl.unionL = x.l
            }
            if x.r > cl.unionR {
                cl.unionR = x.r
            }
            if x.l > cl.comL {
                cl.comL = x.l
            }
            if x.r < cl.comR {
                cl.comR = x.r
            }
            return
        }
    }
    c.cliques = append(c.cliques, Clique{unionL: x.l, unionR: x.r, comL: x.l, comR: x.r})
}

func colorCount(intervals []Interval) int {
    sort.Slice(intervals, func(i, j int) bool {
        if intervals[i].l != intervals[j].l {
            return intervals[i].l < intervals[j].l
        }
        return intervals[i].r < intervals[j].r
    })
    var colors []Color
    for _, iv := range intervals {
        placed := false
        for ci := range colors {
            if colors[ci].canUse(iv) {
                colors[ci].add(iv)
                placed = true
                break
            }
        }
        if !placed {
            var c Color
            c.cliques = []Clique{{unionL: iv.l, unionR: iv.r, comL: iv.l, comR: iv.r}}
            colors = append(colors, c)
        }
    }
    return len(colors)
}

func genTests() []TestCase {
    r := rand.New(rand.NewSource(10))
    cases := make([]TestCase, 100)
    for i := 0; i < 100; i++ {
        n := r.Intn(4) + 1
        ints := make([]Interval, n)
        for j := 0; j < n; j++ {
            a := r.Intn(10)
            b := a + r.Intn(5) + 1
            ints[j] = Interval{a, b}
        }
        ans := colorCount(ints)
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for j, iv := range ints {
            sb.WriteString(fmt.Sprintf("%d %d\n", iv.l, iv.r))
            if j == n-1 {
            }
        }
        cases[i] = TestCase{input: sb.String(), ans: fmt.Sprintf("%d", ans)}
    }
    return cases
}

func run(bin, in string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(in)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierJ.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTests()
    for i, tc := range tests {
        got, err := run(bin, tc.input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != tc.ans {
            fmt.Fprintf(os.Stderr, "test %d failed: input %q expected %q got %q\n", i+1, tc.input, tc.ans, got)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

