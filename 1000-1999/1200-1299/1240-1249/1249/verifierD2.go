package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleD2")
	cmd := exec.Command("go", "build", "-o", oracle, "1249D2.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runProg(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(50) + 1
	k := r.Intn(n) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		l := r.Intn(1000) + 1
		rgt := l + r.Intn(1000-l+1)
		sb.WriteString(fmt.Sprintf("%d %d\n", l, rgt))
	}
	return sb.String()
}

// normalize parses an output string and returns the sorted list of removed
// segment indices. It ensures that the count on the first line matches the
// number of indices that follow so that order in the second line is irrelevant.
func normalize(out string) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	m, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, err
	}
	if m != len(fields)-1 {
		return nil, fmt.Errorf("mismatch count: expect %d numbers, got %d", m, len(fields)-1)
	}
	res := make([]int, m)
	for i := 0; i < m; i++ {
		v, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	sort.Ints(res)
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		expAns, err := normalize(expect)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output parse error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		gotAns, err := normalize(got)
		if err != nil {
			fmt.Printf("case %d failed\ninput:\n%sparse error: %v\n", i, input, err)
			os.Exit(1)
		}
		if len(expAns) != len(gotAns) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, input, expect, got)
			os.Exit(1)
		}
		for j := range expAns {
			if expAns[j] != gotAns[j] {
				fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, input, expect, got)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All 100 tests passed")
}
