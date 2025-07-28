package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", oracle, "1630B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
	}
	return oracle, nil
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func generateCase(rng *rand.Rand) (string, []int, int) {
	n := rng.Intn(20) + 1
	k := rng.Intn(n) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, k)
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String(), a, k
}

func parsePair(line string) (int, int, error) {
	parts := strings.Fields(line)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid line %q", line)
	}
	var x, y int
	if _, err := fmt.Sscan(parts[0], &x); err != nil {
		return 0, 0, err
	}
	if _, err := fmt.Sscan(parts[1], &y); err != nil {
		return 0, 0, err
	}
	return x, y, nil
}

func checkSegs(a []int, segs [][2]int, x, y int, k int) error {
	if len(segs) != k {
		return fmt.Errorf("expected %d segments got %d", k, len(segs))
	}
	pos := 1
	n := len(a)
	for i, s := range segs {
		l, r := s[0], s[1]
		if l != pos || l > r || l < 1 || r > n {
			return fmt.Errorf("invalid segment %d: %d %d", i+1, l, r)
		}
		inside := 0
		for j := l - 1; j < r; j++ {
			if a[j] >= x && a[j] <= y {
				inside++
			}
		}
		outside := (r - l + 1) - inside
		if inside <= outside {
			return fmt.Errorf("segment %d does not satisfy condition", i+1)
		}
		pos = r + 1
	}
	if pos != n+1 {
		return fmt.Errorf("segments do not cover array")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
	for i := 0; i < 100; i++ {
		input, arr, k := generateCase(rng)
		oracleOut, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		candOut, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		oracleLines := strings.Split(strings.TrimSpace(oracleOut), "\n")
		if len(oracleLines) < 1 {
			fmt.Fprintf(os.Stderr, "oracle output malformed on case %d\n", i+1)
			os.Exit(1)
		}
		xOpt, yOpt, err := parsePair(oracleLines[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output malformed on case %d\n", i+1)
			os.Exit(1)
		}
		bestDiff := yOpt - xOpt

		candLines := strings.Split(strings.TrimSpace(candOut), "\n")
		if len(candLines) < 1 {
			fmt.Fprintf(os.Stderr, "case %d: empty output\n", i+1)
			os.Exit(1)
		}
		x, y, err := parsePair(candLines[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid first line\n", i+1)
			os.Exit(1)
		}
		diff := y - x
		if diff != bestDiff {
			fmt.Fprintf(os.Stderr, "case %d: range difference %d not optimal (expected %d)\n", i+1, diff, bestDiff)
			os.Exit(1)
		}
		segs := make([][2]int, 0, k)
		for _, line := range candLines[1:] {
			if strings.TrimSpace(line) == "" {
				continue
			}
			l, r, err := parsePair(line)
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d: invalid segment line %q\n", i+1, line)
				os.Exit(1)
			}
			segs = append(segs, [2]int{l, r})
		}
		if err := checkSegs(arr, segs, x, y, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, input, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
