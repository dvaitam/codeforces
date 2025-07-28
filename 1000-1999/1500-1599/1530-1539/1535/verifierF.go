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

type test struct {
    input    string
    expected string
}

func solve(input string) string {
    r := strings.NewReader(strings.TrimSpace(input))
    var n int
    fmt.Fscan(r, &n)
    strs := make([]string, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(r, &strs[i])
    }
    groups := make(map[string][]string)
    for _, s := range strs {
        bs := []byte(s)
        sort.Slice(bs, func(i, j int) bool { return bs[i] < bs[j] })
        key := string(bs)
        groups[key] = append(groups[key], s)
    }
    totalPairs := int64(n) * (int64(n) - 1) / 2
    var intraPairs int64
    var sumSame int64
    for _, grp := range groups {
        g := len(grp)
        if g <= 1 {
            continue
        }
        intraPairs += int64(g) * (int64(g) - 1) / 2
        L := len(grp[0])
        brute := func() {
            arr := make([][]byte, g)
            for i, s := range grp {
                arr[i] = []byte(s)
            }
            var cnt1 int64
            for i := 0; i < g; i++ {
                a := arr[i]
                for j := i + 1; j < g; j++ {
                    b := arr[j]
                    l := 0
                    for l < L && a[l] == b[l] {
                        l++
                    }
                    if l == L {
                        continue
                    }
                    r := L - 1
                    for r > l && a[r] == b[r] {
                        r--
                    }
                    var ca [26]int
                    for k := l; k <= r; k++ {
                        ca[a[k]-'a']++
                        ca[b[k]-'a']--
                    }
                    ok := true
                    for _, v := range ca {
                        if v != 0 {
                            ok = false
                            break
                        }
                    }
                    if !ok {
                        continue
                    }
                    bs := b
                    sortedB := true
                    for k := l; k < r; k++ {
                        if bs[k] > bs[k+1] {
                            sortedB = false
                            break
                        }
                    }
                    if sortedB {
                        cnt1++
                        continue
                    }
                    as := a
                    sortedA := true
                    for k := l; k < r; k++ {
                        if as[k] > as[k+1] {
                            sortedA = false
                            break
                        }
                    }
                    if sortedA {
                        cnt1++
                    }
                }
            }
            sumSame += cnt1 + (int64(g)*(int64(g)-1)/2-cnt1)*2
        }
        enumeration := func() {
            idx := make(map[string]int, g)
            for i, s := range grp {
                idx[s] = i
            }
            L := len(grp[0])
            var cnt1 int64
            neighborSeen := make([]bool, g)
            sbuf := make([]byte, L)
            for i, s := range grp {
                bs := []byte(s)
                for j := range neighborSeen {
                    neighborSeen[j] = false
                }
                for l := 0; l < L; l++ {
                    for r := l + 1; r < L; r++ {
                        sortedSeg := true
                        for k := l; k < r; k++ {
                            if bs[k] > bs[k+1] {
                                sortedSeg = false
                                break
                            }
                        }
                        if sortedSeg {
                            continue
                        }
                        var cnt [26]int
                        for k := l; k <= r; k++ {
                            cnt[bs[k]-'a']++
                        }
                        copy(sbuf, bs)
                        p := l
                        for c := 0; c < 26; c++ {
                            for t := 0; t < cnt[c]; t++ {
                                sbuf[p] = byte('a' + c)
                                p++
                            }
                        }
                        tkey := string(sbuf)
                        if j, ok := idx[tkey]; ok && j > i && !neighborSeen[j] {
                            neighborSeen[j] = true
                            cnt1++
                        }
                    }
                }
            }
            sumSame += cnt1 + (int64(g)*(int64(g)-1)/2-cnt1)*2
        }
        if len(grp[0]) <= 20 {
            enumeration()
        } else {
            brute()
        }
    }
    sum := sumSame + (totalPairs-intraPairs)*1337
    return fmt.Sprintf("%d\n", sum)
}

func generateTests() []test {
    rand.Seed(6)
    var tests []test
    fixed := [][]string{{"ab", "ab"}, {"ba", "ab"}, {"a", "b"}}
    for _, arr := range fixed {
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
        for _, s := range arr {
            sb.WriteString(fmt.Sprintf("%s\n", s))
        }
        inp := sb.String()
        tests = append(tests, test{inp, solve(inp)})
    }
    for len(tests) < 100 {
        n := rand.Intn(5) + 1
        L := rand.Intn(5) + 1
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for i := 0; i < n; i++ {
            for j := 0; j < L; j++ {
                sb.WriteByte(byte('a' + rand.Intn(3)))
            }
            sb.WriteByte('\n')
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
        fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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

