package main

import (
    "bufio"
    "bytes"
    "context"
    "fmt"
    "io"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strings"
    "time"
)

type Segment struct {
    l   int
    r   int
    idx int
}

type BIT struct {
    n    int
    tree []int
}

func NewBIT(n int) *BIT { return &BIT{n: n, tree: make([]int, n+2)} }
func (b *BIT) Add(i, d int) {
    for i <= b.n {
        b.tree[i] += d
        i += i & -i
    }
}
func (b *BIT) Sum(i int) int {
    s := 0
    for i > 0 {
        s += b.tree[i]
        i -= i & -i
    }
    return s
}

func solveD(in io.Reader) string {
    reader := bufio.NewReader(in)
    var n int
    if _, err := fmt.Fscan(reader, &n); err != nil {
        return ""
    }
    segs := make([]Segment, n)
    rights := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &segs[i].l, &segs[i].r)
        segs[i].idx = i
        rights[i] = segs[i].r
    }
    sortedRights := make([]int, n)
    copy(sortedRights, rights)
    sort.Ints(sortedRights)
    mp := make(map[int]int, n)
    for i, v := range sortedRights {
        mp[v] = i + 1
    }
    for i := range segs {
        segs[i].r = mp[segs[i].r]
    }
    sort.Slice(segs, func(i, j int) bool { return segs[i].l < segs[j].l })
    bit := NewBIT(n)
    ans := make([]int, n)
    for i := n - 1; i >= 0; i-- {
        rp := segs[i].r
        ans[segs[i].idx] = bit.Sum(rp)
        bit.Add(rp, 1)
    }
    var buf strings.Builder
    for i := 0; i < n; i++ {
        buf.WriteString(fmt.Sprintln(ans[i]))
    }
    return strings.TrimSpace(buf.String())
}

func genTests() []string {
    rng := rand.New(rand.NewSource(4))
    tests := make([]string, 100)
    for i := 0; i < 100; i++ {
        n := rng.Intn(8) + 1
        usedL := map[int]bool{}
        usedR := map[int]bool{}
        segs := make([][2]int, n)
        for j := 0; j < n; j++ {
            l := rng.Intn(100) - 50
            for usedL[l] {
                l = rng.Intn(100) - 50
            }
            usedL[l] = true
            rgt := l + rng.Intn(20) + 1
            for usedR[rgt] {
                rgt = l + rng.Intn(20) + 1
            }
            usedR[rgt] = true
            segs[j] = [2]int{l, rgt}
        }
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for j := 0; j < n; j++ {
            sb.WriteString(fmt.Sprintf("%d %d\n", segs[j][0], segs[j][1]))
        }
        tests[i] = sb.String()
    }
    return tests
}

func runBinary(bin, input string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    cmd := exec.CommandContext(ctx, bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run verifierD.go /path/to/binary")
        return
    }
    bin := os.Args[1]
    tests := genTests()
    for i, tc := range tests {
        expected := solveD(strings.NewReader(tc))
        actual, err := runBinary(bin, tc)
        if err != nil {
            fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
            return
        }
        if actual != strings.TrimSpace(expected) {
            fmt.Printf("Test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, tc, expected, actual)
            return
        }
    }
    fmt.Println("All tests passed!")
}

