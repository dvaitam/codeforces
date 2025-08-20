package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strconv"
    "strings"
    "time"
)

func run(bin, input string) (string, error) {
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
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	p := rng.Intn(n) + 1
	k := rng.Intn(p) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, p, k)
	for i := 0; i < n; i++ {
		a := rng.Intn(50) + 1
		b := rng.Intn(50) + 1
		fmt.Fprintf(&sb, "%d %d\n", a, b)
	}
	return sb.String()
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    bin := os.Args[1]
    ref := "332C.go"
    for i := 0; i < 100; i++ {
        in := generateCase(rng)
        exp, err := run(ref, in)
        if err != nil {
            fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
            os.Exit(1)
        }
        out, err := run(bin, in)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
            os.Exit(1)
        }
        // Validate by objective equivalence, not exact order
        n, p, k, a, b := parseInput(in)
        expIdx, err := parseIndices(exp, n, p)
        if err != nil {
            fmt.Fprintf(os.Stderr, "reference parse error on case %d: %v\ninput:\n%sref:\n%s", i+1, err, in, exp)
            os.Exit(1)
        }
        gotIdx, err := parseIndices(out, n, p)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d parse error: %v\ninput:\n%soutput:\n%s", i+1, err, in, out)
            os.Exit(1)
        }
        expHair, expDisp := evalSet(expIdx, k, a, b)
        gotHair, gotDisp := evalSet(gotIdx, k, a, b)
        if gotHair != expHair || gotDisp != expDisp {
            fmt.Fprintf(os.Stderr, "case %d wrong answer\nexpected hair=%d displeasure=%d\n   got hair=%d displeasure=%d\nexpected indices:\n%s\n     got indices:\n%s\ninput:\n%s", i+1, expHair, expDisp, gotHair, gotDisp, strings.TrimSpace(exp), strings.TrimSpace(out), in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

func parseInput(in string) (n, p, k int, a, b []int) {
    lines := strings.Split(strings.TrimSpace(in), "\n")
    header := strings.Fields(lines[0])
    n, _ = strconv.Atoi(header[0])
    p, _ = strconv.Atoi(header[1])
    k, _ = strconv.Atoi(header[2])
    a = make([]int, n+1)
    b = make([]int, n+1)
    for i := 0; i < n; i++ {
        parts := strings.Fields(lines[1+i])
        ai, _ := strconv.Atoi(parts[0])
        bi, _ := strconv.Atoi(parts[1])
        a[i+1] = ai
        b[i+1] = bi
    }
    return
}

func parseIndices(out string, n, p int) ([]int, error) {
    fields := strings.Fields(out)
    if len(fields) < p {
        return nil, fmt.Errorf("expected %d indices, got %d", p, len(fields))
    }
    idx := make([]int, 0, p)
    used := make(map[int]bool)
    for i := 0; i < p; i++ {
        v, err := strconv.Atoi(fields[i])
        if err != nil { return nil, fmt.Errorf("invalid integer: %v", err) }
        if v < 1 || v > n { return nil, fmt.Errorf("index out of range: %d", v) }
        if used[v] { return nil, fmt.Errorf("duplicate index: %d", v) }
        used[v] = true
        idx = append(idx, v)
    }
    return idx, nil
}

func evalSet(idx []int, k int, a, b []int) (hair, displeasure int) {
    // Determine which k orders will be obeyed: maximize b, tie-break minimal a
    type item struct{ id, a, b int }
    items := make([]item, len(idx))
    totalB := 0
    for i, id := range idx {
        items[i] = item{id: id, a: a[id], b: b[id]}
        totalB += b[id]
    }
    sort.Slice(items, func(i, j int) bool {
        if items[i].b != items[j].b { return items[i].b > items[j].b }
        if items[i].a != items[j].a { return items[i].a < items[j].a }
        return items[i].id < items[j].id
    })
    hair = 0
    obeyB := 0
    for i := 0; i < k; i++ {
        hair += items[i].a
        obeyB += items[i].b
    }
    displeasure = totalB - obeyB
    return
}
