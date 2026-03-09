package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleD1")
	cmd := exec.Command("go", "build", "-o", oracle, "1249D1.go")
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
	n := r.Intn(20) + 1
	k := r.Intn(n) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		l := r.Intn(200) + 1
		rgt := l + r.Intn(200-l+1)
		sb.WriteString(fmt.Sprintf("%d %d\n", l, rgt))
	}
	return sb.String()
}

// check validates that the candidate output is a correct answer given the input
// and that its count matches the oracle's count (minimality).
func check(caseNum int, input, oracleOut, candOut string) error {
	inLines := strings.Fields(input)
	idx := 0
	n, _ := strconv.Atoi(inLines[idx])
	idx++
	k, _ := strconv.Atoi(inLines[idx])
	idx++
	ls := make([]int, n+1)
	rs := make([]int, n+1)
	for i := 1; i <= n; i++ {
		ls[i], _ = strconv.Atoi(inLines[idx])
		idx++
		rs[i], _ = strconv.Atoi(inLines[idx])
		idx++
	}

	// Parse oracle count (first line)
	oLines := strings.Split(strings.TrimSpace(oracleOut), "\n")
	oracleM, err := strconv.Atoi(strings.TrimSpace(oLines[0]))
	if err != nil {
		return fmt.Errorf("case %d: bad oracle output: %v", caseNum, err)
	}

	// Parse candidate output
	cLines := strings.Split(strings.TrimSpace(candOut), "\n")
	if len(cLines) < 1 {
		return fmt.Errorf("case %d: empty candidate output", caseNum)
	}
	candM, err := strconv.Atoi(strings.TrimSpace(cLines[0]))
	if err != nil {
		return fmt.Errorf("case %d: bad candidate count: %v", caseNum, err)
	}
	if candM != oracleM {
		return fmt.Errorf("case %d: candidate removed %d segments, oracle removed %d (not minimum)", caseNum, candM, oracleM)
	}

	removed := make(map[int]bool)
	if candM > 0 {
		if len(cLines) < 2 {
			return fmt.Errorf("case %d: missing indices line", caseNum)
		}
		parts := strings.Fields(cLines[1])
		if len(parts) != candM {
			return fmt.Errorf("case %d: expected %d indices, got %d", caseNum, candM, len(parts))
		}
		for _, p := range parts {
			v, err := strconv.Atoi(p)
			if err != nil || v < 1 || v > n {
				return fmt.Errorf("case %d: invalid index %q", caseNum, p)
			}
			if removed[v] {
				return fmt.Errorf("case %d: duplicate index %d", caseNum, v)
			}
			removed[v] = true
		}
	}

	// Verify no bad points remain
	for x := 1; x <= 200; x++ {
		cnt := 0
		for i := 1; i <= n; i++ {
			if !removed[i] && ls[i] <= x && x <= rs[i] {
				cnt++
			}
		}
		if cnt > k {
			return fmt.Errorf("case %d: point %d still covered by %d > %d segments after removal", caseNum, x, cnt, k)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
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
		if err := check(i, input, expect, got); err != nil {
			fmt.Printf("%v\ninput:\n%sexpected:\n%s\ngot:\n%s\n", err, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
