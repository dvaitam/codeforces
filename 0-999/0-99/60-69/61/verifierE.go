package main

import (
    "bufio"
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strconv"
    "strings"
)

type BIT struct { n int; tree []int }

func NewBIT(n int) *BIT { return &BIT{n:n, tree: make([]int, n+1)} }
func (b *BIT) Add(i, v int) { for ; i<=b.n; i+= i&-i { b.tree[i]+=v } }
func (b *BIT) Sum(i int) int { s:=0; for ; i>0; i-= i&-i { s += b.tree[i] }; return s }

type test struct{ input, expected string }

func solve(input string) string {
    in := bufio.NewReader(strings.NewReader(input))
    var n int
    fmt.Fscan(in, &n)
    vals := make([]int, n)
    for i:=0;i<n;i++ { fmt.Fscan(in, &vals[i]) }
    type pair struct{ val, idx int }
    arr := make([]pair, n)
    for i,v := range vals { arr[i] = pair{v,i} }
    sort.Slice(arr, func(i,j int) bool { return arr[i].val < arr[j].val })
    ranks := make([]int,n)
    for i,p := range arr { ranks[p.idx] = i+1 }
    rightLess := make([]int,n)
    bit := NewBIT(n)
    for j:=n-1; j>=0; j-- { r := ranks[j]; if r>1 { rightLess[j] = bit.Sum(r-1) }; bit.Add(r,1) }
    bit = NewBIT(n)
    var ans int64
    for j:=0; j<n; j++ { r := ranks[j]; leftGt := j - bit.Sum(r); ans += int64(leftGt) * int64(rightLess[j]); bit.Add(r,1) }
    return fmt.Sprintf("%d", ans)
}

func generateTests() []test {
    rand.Seed(46)
    var tests []test
    fixed := []string{
        "3\n3 2 1\n",
        "4\n1 2 3 4\n",
        "5\n5 4 3 2 1\n",
    }
    for _, f := range fixed { tests = append(tests, test{f, solve(f)}) }
    for len(tests) < 100 {
        n := rand.Intn(20)+3
        vals := make([]int,n)
        for i := range vals { vals[i] = rand.Intn(1000*n) }
        // ensure unique
        used := map[int]bool{}
        for i := range vals {
            for {
                v := vals[i]
                if !used[v] {
                    used[v] = true
                    break
                }
                vals[i] = rand.Intn(1000*n)
            }
        }
        var sb strings.Builder
        sb.WriteString(strconv.Itoa(n)+"\n")
        for i:=0;i<n;i++ { sb.WriteString(fmt.Sprintf("%d ", vals[i])) }
        sb.WriteByte('\n')
        inp := sb.String()
        tests = append(tests, test{inp, solve(inp)})
    }
    return tests
}

func runBinary(bin, input string) (string, error) {
    cmd := exec.Command(bin)
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
            fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}

