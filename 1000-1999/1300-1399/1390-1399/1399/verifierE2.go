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

type edgeE2 struct {
    to   int
    w    int64
    cost int
}

type testCaseE2 struct {
    n   int
    S   int64
    edges [][4]int
}

func parseTestcasesE2(path string) ([]testCaseE2, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    scanner := bufio.NewScanner(f)
    var cases []testCaseE2
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }
        fields := strings.Fields(line)
        if len(fields) != 2 {
            return nil, fmt.Errorf("bad header")
        }
        n, _ := strconv.Atoi(fields[0])
        S, _ := strconv.ParseInt(fields[1], 10, 64)
        edges := make([][4]int, 0, n-1)
        for i := 0; i < n-1; i++ {
            if !scanner.Scan() {
                return nil, fmt.Errorf("unexpected EOF")
            }
            ep := strings.Fields(scanner.Text())
            if len(ep) != 4 {
                return nil, fmt.Errorf("bad edge line")
            }
            u, _ := strconv.Atoi(ep[0])
            v, _ := strconv.Atoi(ep[1])
            w, _ := strconv.Atoi(ep[2])
            c, _ := strconv.Atoi(ep[3])
            edges = append(edges, [4]int{u, v, w, c})
        }
        cases = append(cases, testCaseE2{n, S, edges})
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    return cases, nil
}

func solveE2(n int, S int64, edges [][4]int) int64 {
    adj := make([][]edgeE2, n+1)
    for _, e := range edges {
        u, v := e[0], e[1]
        w := int64(e[2])
        c := e[3]
        adj[u] = append(adj[u], edgeE2{v, w, c})
        adj[v] = append(adj[v], edgeE2{u, w, c})
    }
    var edgesInfo []struct{ leaves int64; w int64; cost int }
    var total int64
    var dfs func(u, p int) int64
    dfs = func(u, p int) int64 {
        if len(adj[u]) == 1 && p != 0 {
            return 1
        }
        var sum int64
        for _, e := range adj[u] {
            if e.to == p {
                continue
            }
            cnt := dfs(e.to, u)
            total += e.w * cnt
            edgesInfo = append(edgesInfo, struct{ leaves int64; w int64; cost int }{cnt, e.w, e.cost})
            sum += cnt
        }
        return sum
    }
    dfs(1, 0)
    if total <= S {
        return 0
    }
    var d1, d2 []int64
    for _, ei := range edgesInfo {
        w := ei.w
        leaves := ei.leaves
        for w > 0 {
            delta := (w - w/2) * leaves
            if ei.cost == 1 {
                d1 = append(d1, delta)
            } else {
                d2 = append(d2, delta)
            }
            w /= 2
        }
    }
    sort.Slice(d1, func(i, j int) bool { return d1[i] > d1[j] })
    sort.Slice(d2, func(i, j int) bool { return d2[i] > d2[j] })
    p1 := make([]int64, len(d1)+1)
    for i := 0; i < len(d1); i++ {
        p1[i+1] = p1[i] + d1[i]
    }
    p2 := make([]int64, len(d2)+1)
    for i := 0; i < len(d2); i++ {
        p2[i+1] = p2[i] + d2[i]
    }
    need := total - S
    ans := int64(1<<62)
    for k2 := 0; k2 < len(p2); k2++ {
        sum2 := p2[k2]
        if sum2 >= need {
            cost := int64(k2) * 2
            if cost < ans {
                ans = cost
            }
            break
        }
        rem := need - sum2
        idx := sort.Search(len(p1), func(i int) bool { return p1[i] >= rem })
        if idx < len(p1) {
            cost := int64(k2)*2 + int64(idx)
            if cost < ans {
                ans = cost
            }
        }
    }
    return ans
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
        fmt.Println("usage: go run verifierE2.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    cases, err := parseTestcasesE2("testcasesE2.txt")
    if err != nil {
        fmt.Println("failed to parse testcases:", err)
        os.Exit(1)
    }
    for idx, tc := range cases {
        var sb strings.Builder
        sb.WriteString("1\n")
        sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.S))
        for _, e := range tc.edges {
            sb.WriteString(fmt.Sprintf("%d %d %d %d\n", e[0], e[1], e[2], e[3]))
        }
        expected := strconv.FormatInt(solveE2(tc.n, tc.S, tc.edges), 10)
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

