package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	mod       = 1000000007
	refSource = "./2022E2.go"
)

type testCase struct {
	name  string
	input string
	t     int
}

type dsu struct {
	parent      []int
	xorToParent []int
	size        []int
	comps       int
}

func newDSU(n int) *dsu {
	parent := make([]int, n)
	xorToParent := make([]int, n)
	size := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &dsu{parent: parent, xorToParent: xorToParent, size: size, comps: n}
}

func (d *dsu) find(x int) (int, int) {
	if d.parent[x] == x {
		return x, 0
	}
	root, xr := d.find(d.parent[x])
	total := xr ^ d.xorToParent[x]
	d.parent[x] = root
	d.xorToParent[x] = total
	return root, total
}

func (d *dsu) union(a, b, val int) bool {
	ra, xa := d.find(a)
	rb, xb := d.find(b)
	if ra == rb {
		return (xa ^ xb) == val
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
		xa, xb = xb, xa
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	d.xorToParent[rb] = xa ^ xb ^ val
	d.comps--
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		expected, err := simulate(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to simulate test %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}

		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refAns, err := parseOutput(len(expected), refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		if !equalAnswers(refAns, expected) {
			fmt.Fprintf(os.Stderr, "reference output mismatch simulation on test %d (%s)\ninput:\n%soutput:\n%s", idx+1, tc.name, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseOutput(len(expected), candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if !equalAnswers(candAns, expected) {
			fmt.Fprintf(os.Stderr, "candidate answer mismatch on test %d (%s)\ninput:\n%soutput:\n%s", idx+1, tc.name, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2022E2-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2022E2.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
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
	return out.String(), err
}

func parseOutput(expected int, output string) ([]int, error) {
	fields := strings.Fields(output)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d values, got %d", expected, len(fields))
	}
	ans := make([]int, expected)
	for i, tok := range fields {
		val, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		val %= mod
		if val < 0 {
			val += mod
		}
		ans[i] = val
	}
	return ans, nil
}

func equalAnswers(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if ((a[i]-b[i])%mod+mod)%mod != 0 {
			return false
		}
	}
	return true
}

func simulate(input string) ([]int, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	results := make([]int, 0)
	for ; t > 0; t-- {
		var n, m, k, q int
		if _, err := fmt.Fscan(reader, &n, &m, &k, &q); err != nil {
			return nil, err
		}
		d := newDSU(n + m)
		inconsistent := false
		process := func() {
			if inconsistent {
				results = append(results, 0)
				return
			}
			exp := (d.comps - 1) * 30
			results = append(results, modPow2(exp))
		}
		for i := 0; i < k; i++ {
			var r, c, v int
			fmt.Fscan(reader, &r, &c, &v)
			r--
			c--
			if !inconsistent {
				if !d.union(r, n+c, v) {
					inconsistent = true
				}
			}
		}
		process()
		for i := 0; i < q; i++ {
			var r, c, v int
			fmt.Fscan(reader, &r, &c, &v)
			r--
			c--
			if !inconsistent {
				if !d.union(r, n+c, v) {
					inconsistent = true
				}
			}
			process()
		}
	}
	return results, nil
}

func modPow2(exp int) int {
	exp %= mod - 1
	if exp < 0 {
		exp += mod - 1
	}
	res := 1
	base := 2
	for exp > 0 {
		if exp&1 == 1 {
			res = (res * base) % mod
		}
		base = (base * base) % mod
		exp >>= 1
	}
	return res
}

func buildTests() []testCase {
	tests := []testCase{
		{name: "simple", input: buildInputCase(1, []testInstance{
			{
				n: 2, m: 2, k: 1, q: 1,
				assignments: []assignment{
					{1, 1, 0},
					{2, 2, 0},
				},
			},
		})},
	}

	rng := rand.New(rand.NewSource(123456789))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

type assignment struct {
	r, c, v int
}

type testInstance struct {
	n, m, k, q  int
	assignments []assignment
}

func buildInputCase(t int, instances []testInstance) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for _, inst := range instances {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", inst.n, inst.m, inst.k, inst.q))
		for _, a := range inst.assignments {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", a.r, a.c, a.v))
		}
	}
	return sb.String()
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(3) + 1
	instances := make([]testInstance, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(5) + 2
		m := rng.Intn(5) + 2
		maxCells := n * m
		k := rng.Intn(maxCells/2 + 1)
		q := rng.Intn(maxCells - k + 1)
		used := make(map[[2]int]struct{})
		assignCount := k + q
		assigns := make([]assignment, 0, assignCount)
		for len(assigns) < assignCount {
			r := rng.Intn(n) + 1
			c := rng.Intn(m) + 1
			key := [2]int{r, c}
			if _, ok := used[key]; ok {
				continue
			}
			used[key] = struct{}{}
			v := rng.Intn(1 << 10)
			assigns = append(assigns, assignment{r: r, c: c, v: v})
		}
		instances[i] = testInstance{
			n: n, m: m, k: k, q: q,
			assignments: assigns,
		}
	}
	return testCase{
		name:  fmt.Sprintf("random_%d", idx+1),
		input: buildInputCase(t, instances),
		t:     t,
	}
}
