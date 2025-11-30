package main

import (
    "bytes"
    "compress/gzip"
    "encoding/base64"
    "fmt"
    "io"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

type Edge struct {
    id    int
    group int
    u     int
    v     int
    dist  int
    cap   int
}

type Flow struct {
    id   int
    s    int
    t    int
    rate int
}

type testCase struct {
    nodeCount   int
    edgeCount   int
    constrCount int
    flowCount   int
    edges       []Edge
    flows       []Flow
}

// Gzipped+base64 testcases data from testcasesA.txt.
const encodedTestcases = `
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
`

// solveCase mirrors 1576A.go and returns the output string for a single test.
func solveCase(tc testCase) string {
    edges := tc.edges
    n := tc.nodeCount
    adj := make(map[[2]int][]int)
    for i, e := range edges {
        adj[[2]int{e.u, e.v}] = append(adj[[2]int{e.u, e.v}], i)
        adj[[2]int{e.v, e.u}] = append(adj[[2]int{e.v, e.u}], i)
    }

    const SFL = 200
    const GFL = 100

    nodeFlowCount := make([]int, n)
    groupFlowCount := make(map[int]int)
    edgeCapLeft := make([]int, len(edges))
    for i, e := range edges {
        edgeCapLeft[i] = e.cap
    }

    type Path struct {
        flowID int
        edgeID int
    }
    var results []Path

    for _, f := range tc.flows {
        if nodeFlowCount[f.s] >= SFL || nodeFlowCount[f.t] >= SFL {
            continue
        }
        candidates := adj[[2]int{f.s, f.t}]
        for _, ei := range candidates {
            e := edges[ei]
            if edgeCapLeft[ei] < f.rate {
                continue
            }
            if groupFlowCount[e.group] >= GFL {
                continue
            }
            edgeCapLeft[ei] -= f.rate
            nodeFlowCount[f.s]++
            nodeFlowCount[f.t]++
            groupFlowCount[e.group]++
            results = append(results, Path{flowID: f.id, edgeID: e.id})
            break
        }
    }

    var sb strings.Builder
    fmt.Fprintln(&sb, len(results))
    for _, p := range results {
    fmt.Fprintf(&sb, "%d %d\n", p.flowID, p.edgeID)
    }
    return strings.TrimSpace(sb.String())
}

func decodeTestcases() (string, error) {
    data, err := base64.StdEncoding.DecodeString(encodedTestcases)
    if err != nil {
        return "", err
    }
    r, err := gzip.NewReader(bytes.NewReader(data))
    if err != nil {
        return "", err
    }
    defer r.Close()
    var out bytes.Buffer
    if _, err := io.Copy(&out, r); err != nil {
        return "", err
    }
    return out.String(), nil
}

func parseTestcases() ([]testCase, error) {
    raw, err := decodeTestcases()
    if err != nil {
        return nil, err
    }
    fields := strings.Fields(raw)
    if len(fields) == 0 {
        return nil, fmt.Errorf("no data")
    }
    pos := 0
    t, err := strconv.Atoi(fields[pos])
    if err != nil {
        return nil, fmt.Errorf("bad t: %v", err)
    }
    pos++
    cases := make([]testCase, 0, t)
    for caseIdx := 0; caseIdx < t; caseIdx++ {
        if pos+3 >= len(fields) {
            return nil, fmt.Errorf("case %d header truncated", caseIdx+1)
        }
        nodeCount, _ := strconv.Atoi(fields[pos])
        edgeCount, _ := strconv.Atoi(fields[pos+1])
        constrCount, _ := strconv.Atoi(fields[pos+2])
        flowCount, _ := strconv.Atoi(fields[pos+3])
        pos += 4
        edges := make([]Edge, edgeCount)
        for i := 0; i < edgeCount; i++ {
            if pos+5 >= len(fields) {
                return nil, fmt.Errorf("case %d edges truncated", caseIdx+1)
            }
            id, _ := strconv.Atoi(fields[pos])
            group, _ := strconv.Atoi(fields[pos+1])
            u, _ := strconv.Atoi(fields[pos+2])
            v, _ := strconv.Atoi(fields[pos+3])
            dist, _ := strconv.Atoi(fields[pos+4])
            cap, _ := strconv.Atoi(fields[pos+5])
            edges[i] = Edge{id: id, group: group, u: u, v: v, dist: dist, cap: cap}
            pos += 6
        }
        pos += 3 * constrCount // skip constraints
        flows := make([]Flow, flowCount)
        for i := 0; i < flowCount; i++ {
            if pos+3 >= len(fields) {
                return nil, fmt.Errorf("case %d flows truncated", caseIdx+1)
            }
            id, _ := strconv.Atoi(fields[pos])
            s, _ := strconv.Atoi(fields[pos+1])
            tgt, _ := strconv.Atoi(fields[pos+2])
            rate, _ := strconv.Atoi(fields[pos+3])
            flows[i] = Flow{id: id, s: s, t: tgt, rate: rate}
            pos += 4
        }
        cases = append(cases, testCase{
            nodeCount:   nodeCount,
            edgeCount:   edgeCount,
            constrCount: constrCount,
            flowCount:   flowCount,
            edges:       edges,
            flows:       flows,
        })
    }
    return cases, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
    var sb strings.Builder
    fmt.Fprintf(&sb, "%d %d %d %d\n", tc.nodeCount, tc.edgeCount, tc.constrCount, tc.flowCount)
    for _, e := range tc.edges {
        fmt.Fprintf(&sb, "%d %d %d %d %d %d\n", e.id, e.group, e.u, e.v, e.dist, e.cap)
    }
    for i := 0; i < tc.constrCount; i++ {
        fmt.Fprintln(&sb, "0 0 0")
    }
    for _, f := range tc.flows {
        fmt.Fprintf(&sb, "%d %d %d %d\n", f.id, f.s, f.t, f.rate)
    }
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(sb.String())
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
        fmt.Println("usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]

    tests, err := parseTestcases()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    for idx, tc := range tests {
        expect := solveCase(tc)
        got, err := runCandidate(bin, tc)
        if err != nil {
            fmt.Printf("case %d failed: %v\n", idx+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != strings.TrimSpace(expect) {
            fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(tests))
}
