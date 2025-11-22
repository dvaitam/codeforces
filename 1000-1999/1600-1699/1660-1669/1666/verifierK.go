package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type edge struct {
	u int
	v int
	w int64
}

type testCase struct {
	input string
	n     int
	a     int
	b     int
	edges []edge
}

func buildRef() (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "refK.bin")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "1666K.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func makeCase(n, a, b int, edges []edge) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	fmt.Fprintf(&sb, "%d %d\n", a, b)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.w)
	}
	// copy edges so later mutations won't leak
	cpy := make([]edge, len(edges))
	copy(cpy, edges)
	return testCase{
		input: sb.String(),
		n:     n,
		a:     a,
		b:     b,
		edges: cpy,
	}
}

func randCase(r *rand.Rand, nMax int) testCase {
	if nMax < 2 {
		nMax = 2
	}
	n := r.Intn(nMax-1) + 2
	a := r.Intn(n) + 1
	b := r.Intn(n) + 1
	for b == a {
		b = r.Intn(n) + 1
	}
	maxEdges := n * (n - 1) / 2
	if maxEdges > 2000 {
		maxEdges = 2000
	}
	m := r.Intn(maxEdges + 1)
	edges := make([]edge, 0, m)
	used := make(map[[2]int]struct{})
	for len(edges) < m {
		u := r.Intn(n) + 1
		v := r.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		key := [2]int{u, v}
		if _, ok := used[key]; ok {
			continue
		}
		used[key] = struct{}{}
		w := r.Int63n(1_000_000_000) + 1
		edges = append(edges, edge{u, v, w})
	}
	return makeCase(n, a, b, edges)
}

func genCases() []testCase {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		makeCase(2, 1, 2, nil),
		makeCase(3, 1, 2, []edge{{1, 2, 5}, {2, 3, 5}, {1, 3, 1}}),
		makeCase(4, 2, 3, []edge{{1, 2, 10}, {2, 3, 1}, {3, 4, 10}, {1, 4, 1}}),
		makeCase(5, 4, 5, []edge{{1, 2, 1_000_000_000}, {2, 3, 1}, {3, 4, 2}, {4, 5, 3}}),
	}
	for i := 0; i < 40; i++ {
		cases = append(cases, randCase(r, 15))
	}
	for i := 0; i < 20; i++ {
		cases = append(cases, randCase(r, 80))
	}
	for i := 0; i < 10; i++ {
		cases = append(cases, randCase(r, 1000))
	}
	return cases
}

func edgeCost(c1, c2 byte, w int64) int64 {
	if c1 == 'A' && c2 == 'A' {
		return 2 * w
	}
	if c1 == 'B' && c2 == 'B' {
		return 2 * w
	}
	if (c1 == 'A' && c2 == 'C') || (c1 == 'C' && c2 == 'A') {
		return w
	}
	if (c1 == 'B' && c2 == 'C') || (c1 == 'C' && c2 == 'B') {
		return w
	}
	return 0
}

func computeCost(assign string, edges []edge) int64 {
	asg := []byte(assign)
	var total int64
	for _, e := range edges {
		total += edgeCost(asg[e.u-1], asg[e.v-1], e.w)
	}
	return total
}

func parseOutput(out string, tc testCase) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) < 2 {
		return 0, fmt.Errorf("expected cost and partition, got: %q", out)
	}
	cost, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse cost %q: %v", fields[0], err)
	}
	assign := strings.Join(fields[1:], "")
	if len(assign) != tc.n {
		return 0, fmt.Errorf("partition length %d does not match n=%d", len(assign), tc.n)
	}
	for i := 0; i < len(assign); i++ {
		if assign[i] != 'A' && assign[i] != 'B' && assign[i] != 'C' {
			return 0, fmt.Errorf("invalid character %q in partition", assign[i])
		}
	}
	if assign[tc.a-1] != 'A' {
		return 0, fmt.Errorf("town %d must be in A", tc.a)
	}
	if assign[tc.b-1] != 'B' {
		return 0, fmt.Errorf("town %d must be in B", tc.b)
	}
	actual := computeCost(assign, tc.edges)
	if actual != cost {
		return 0, fmt.Errorf("reported cost %d but computed cost is %d", cost, actual)
	}
	return cost, nil
}

func checkCase(bin, ref string, tc testCase) error {
	refOut, err := runBinary(ref, tc.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	expected, err := parseOutput(refOut, tc)
	if err != nil {
		return fmt.Errorf("reference output invalid: %v", err)
	}

	out, err := runBinary(bin, tc.input)
	if err != nil {
		return err
	}
	got, err := parseOutput(out, tc)
	if err != nil {
		return err
	}
	if got != expected {
		return fmt.Errorf("expected minimum cost %d but got %d", expected, got)
	}
	return nil
}

func main() {
	exitCode := 0
	cleanup := func() {}
	defer func() {
		cleanup()
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()

	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierK.go /path/to/binary")
		exitCode = 1
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	cleanup = func() { _ = os.Remove(ref) }

	cases := genCases()
	for i, tc := range cases {
		if err := checkCase(bin, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			exitCode = 1
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
