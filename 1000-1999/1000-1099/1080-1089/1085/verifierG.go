package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func buildRef() (string, error) {
	src := os.Getenv("REFERENCE_SOURCE_PATH")
	if src == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	bin := fmt.Sprintf("%s/ref1085G.bin", os.TempDir())
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return bin, nil
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

type Case struct {
	n    int
	grid [][]int
}

// genDerangementRow generates a random permutation of 1..n that differs from prev at every position.
func genDerangementRow(r *rand.Rand, n int, prev []int) []int {
	for {
		perm := r.Perm(n)
		row := make([]int, n)
		ok := true
		for j := 0; j < n; j++ {
			row[j] = perm[j] + 1
			if row[j] == prev[j] {
				ok = false
				break
			}
		}
		if ok {
			return row
		}
	}
}

func genCases() []Case {
	r := rand.New(rand.NewSource(1085))
	cases := make([]Case, 100)
	for i := range cases {
		n := r.Intn(4) + 2 // 2..5
		grid := make([][]int, n)
		// First row: random permutation of 1..n
		perm := r.Perm(n)
		grid[0] = make([]int, n)
		for j := 0; j < n; j++ {
			grid[0][j] = perm[j] + 1
		}
		// Subsequent rows: permutations that differ from previous row at every position
		for row := 1; row < n; row++ {
			grid[row] = genDerangementRow(r, n, grid[row-1])
		}
		cases[i] = Case{n, grid}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rand.Seed(time.Now().UnixNano())
	cases := genCases()
	for i, c := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", c.n)
		for _, row := range c.grid {
			for j, x := range row {
				if j > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", x)
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		want, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
