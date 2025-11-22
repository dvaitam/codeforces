package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Segment struct{ l, r, c int }

type Case struct{ segs []Segment }

func genCases() []Case {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Add a deterministic case that previously failed the verifier.
	fixed := Case{segs: []Segment{{12, 20, 1}, {1, 5, 2}, {10, 26, 3}, {17, 29, 1}, {2, 14, 2}}}

	cases := []Case{fixed}
	for i := 0; i < 99; i++ {
		n := rng.Intn(5) + 1
		segs := make([]Segment, n)
		for j := range segs {
			l := rng.Intn(20) + 1
			r := l + rng.Intn(20)
			c := rng.Intn(3) + 1
			segs[j] = Segment{l, r, c}
		}
		cases = append(cases, Case{segs: segs})
	}
	return cases
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func distance(a, b Segment) int {
	if a.r < b.l {
		return b.l - a.r
	}
	if b.r < a.l {
		return a.l - b.r
	}
	return 0
}

func expectedOutputs(c Case) []int {
	res := make([]int, len(c.segs))
	for i, s := range c.segs {
		best := math.MaxInt32
		for j, t := range c.segs {
			if i == j || s.c == t.c {
				continue
			}
			d := distance(s, t)
			if d < best {
				best = d
			}
		}
		res[i] = best
	}
	return res
}

func formatOutputs(ans []int) string {
	var sb strings.Builder
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return sb.String()
}

func parseOutputs(got string, n int) ([]int, error) {
	var res []int
	reader := strings.NewReader(got)
	for len(res) < n {
		var x int
		if _, err := fmt.Fscan(reader, &x); err != nil {
			return nil, fmt.Errorf("unable to parse output: %v", err)
		}
		res = append(res, x)
	}
	return res, nil
}

func runCase(bin string, c Case) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(c.segs)))
	for _, s := range c.segs {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", s.l, s.r, s.c))
	}
	input := sb.String()
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	expected := expectedOutputs(c)
	parsed, err := parseOutputs(got, len(c.segs))
	if err != nil {
		return err
	}
	for i := range expected {
		if expected[i] != parsed[i] {
			return fmt.Errorf("expected %s got %s", formatOutputs(expected), got)
		}
	}
	return nil
}

func buildBinary(bin string) (string, func(), error) {
	if !strings.HasSuffix(bin, ".go") {
		return bin, func() {}, nil
	}
	tmp, err := os.CreateTemp("", "verifierF_*")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()
	if err := exec.Command("go", "build", "-o", tmp.Name(), bin).Run(); err != nil {
		os.Remove(tmp.Name())
		return "", nil, fmt.Errorf("failed to build solution: %v", err)
	}
	return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	binPath, cleanup, err := buildBinary(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	cases := genCases()
	for i, c := range cases {
		if err := runCase(binPath, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %+v\n", i+1, err, c)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
