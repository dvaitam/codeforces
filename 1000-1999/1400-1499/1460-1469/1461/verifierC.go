package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1461C.go")
	oracle := filepath.Join(os.TempDir(), "oracle1461C.bin")
	cmd := exec.Command("go", "build", "-o", oracle, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func genCase(r *rand.Rand) string {
	t := r.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := r.Intn(100) + 1
		m := r.Intn(100) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		perm := r.Perm(n)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", perm[j]+1))
		}
		sb.WriteByte('\n')
		for j := 0; j < m; j++ {
			rpos := r.Intn(n) + 1
			p := r.Float64()
			sb.WriteString(fmt.Sprintf("%d %.6f\n", rpos, p))
		}
	}
	return sb.String()
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

func compareOutputs(expect, got string) error {
	expFields := strings.Fields(expect)
	gotFields := strings.Fields(got)

	if len(expFields) != len(gotFields) {
		return fmt.Errorf("expected %d tokens, got %d", len(expFields), len(gotFields))
	}

	for i := range expFields {
		expVal, err1 := strconv.ParseFloat(expFields[i], 64)
		gotVal, err2 := strconv.ParseFloat(gotFields[i], 64)

		if err1 == nil && err2 == nil {
			if math.Abs(expVal-gotVal) > 1e-6 {
				return fmt.Errorf("at token %d: expected %s, got %s (diff > 1e-6)", i+1, expFields[i], gotFields[i])
			}
		} else {
			if expFields[i] != gotFields[i] {
				return fmt.Errorf("at token %d: expected %s, got %s", i+1, expFields[i], gotFields[i])
			}
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
		input := genCase(rng)
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := compareOutputs(expect, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, err, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}