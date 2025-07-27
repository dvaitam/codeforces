package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "sort"
    "strconv"
    "strings"
)

type segment struct{ l, r int }

type testCaseF struct {
    segs []segment
}

func parseTestcasesF(path string) ([]testCaseF, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    var cases []testCaseF
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }
        fields := strings.Fields(line)
        if len(fields) < 1 {
            return nil, fmt.Errorf("bad line")
        }
        n, _ := strconv.Atoi(fields[0])
        if len(fields)-1 != 2*n {
            return nil, fmt.Errorf("expected %d numbers", 2*n)
        }
        segs := make([]segment, n)
        for i := 0; i < n; i++ {
            l, _ := strconv.Atoi(fields[1+2*i])
            r, _ := strconv.Atoi(fields[2+2*i])
            segs[i] = segment{l, r}
        }
        cases = append(cases, testCaseF{segs})
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    return cases, nil
}

func solveF(segs []segment) int {
    n := len(segs)
    xs := make([]int, 0, 2*n)
    for _, s := range segs {
        xs = append(xs, s.l, s.r)
    }
    sort.Ints(xs)
    m := 0
    last := -1
    for _, v := range xs {
        if v != last {
            xs[m] = v
            m++
            last = v
        }
    }
    xs = xs[:m]
    idx := make(map[int]int, m)
    for i, v := range xs {
        idx[v] = i + 1
    }
    endsAt := make([][]int, m+2)
    for _, s := range segs {
        l := idx[s.l]
        r := idx[s.r]
        endsAt[r] = append(endsAt[r], l)
    }
    f := make([][]int16, m+2)
    for i := 0; i <= m+1; i++ {
        f[i] = make([]int16, m+2)
    }
    for width := 0; width < m; width++ {
        for i := 1; i+width <= m; i++ {
            j := i + width
            var best int16
            if j > i {
                best = f[i][j-1]
            }
            for _, l0 := range endsAt[j] {
                if l0 < i {
                    continue
                }
                cur := int16(1)
                if l0 > i {
                    cur += f[i][l0-1]
                }
                if l0 < j {
                    cur += f[l0+1][j-1]
                }
                if cur > best {
                    best = cur
                }
            }
            f[i][j] = best
        }
    }
    if m > 0 {
        return int(f[1][m])
    }
    return 0
}

func run(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var errBuf bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errBuf
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierF.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    cases, err := parseTestcasesF("testcasesF.txt")
    if err != nil {
        fmt.Println("failed to parse testcases:", err)
        os.Exit(1)
    }
    for idx, tc := range cases {
        var sb strings.Builder
        n := len(tc.segs)
        sb.WriteString("1\n")
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for _, s := range tc.segs {
            sb.WriteString(fmt.Sprintf("%d %d\n", s.l, s.r))
        }
        expected := strconv.Itoa(solveF(tc.segs))
        got, err := run(bin, sb.String())
        if err != nil {
            fmt.Printf("case %d failed: %v\n", idx+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != expected {
            fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(cases))
}

