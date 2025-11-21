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
)

const testCount = 200

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "477E.go")
	tmp, err := os.CreateTemp("", "oracle477E")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()
	os.Remove(path)
	cmd := exec.Command("go", "build", "-o", path, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	return path, nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) (string, int) {
	n := r.Intn(40) + 1
	if r.Intn(4) == 0 {
		n = r.Intn(400) + 1
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		val := r.Intn(200) + 1
		if r.Intn(5) == 0 {
			val = r.Intn(100000000) + 1
		}
		a[i] = val
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteByte('\n')
	q := r.Intn(60) + 1
	if r.Intn(5) == 0 {
		q = r.Intn(300) + 1
	}
	fmt.Fprintf(&sb, "%d\n", q)
	for i := 0; i < q; i++ {
		r1 := r.Intn(n) + 1
		r2 := r.Intn(n) + 1
		c1 := r.Intn(a[r1-1] + 1)
		c2 := r.Intn(a[r2-1] + 1)
		if r.Intn(6) == 0 {
			c1 = 0
		} else if r.Intn(6) == 0 {
			c1 = a[r1-1]
		}
		if r.Intn(6) == 0 {
			c2 = 0
		} else if r.Intn(6) == 0 {
			c2 = a[r2-1]
		}
		fmt.Fprintf(&sb, "%d %d %d %d\n", r1, c1, r2, c2)
	}
	return sb.String(), q
}

func parseOutputs(out string, q int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != q {
		return nil, fmt.Errorf("expected %d numbers, got %d", q, len(fields))
	}
	res := make([]int64, q)
	for i := 0; i < q; i++ {
		val, err := strconv.ParseInt(fields[i], 10, 64)
		if err != nil {
			return nil, err
		}
		res[i] = val
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	r := rand.New(rand.NewSource(1))
	for i := 0; i < testCount; i++ {
		input, q := genCase(r)
		expectStr, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotStr, err := runBinary(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expectVals, err := parseOutputs(expectStr, q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotVals, err := parseOutputs(gotStr, q)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", i+1, input, err)
			os.Exit(1)
		}
		for j := 0; j < q; j++ {
			if expectVals[j] != gotVals[j] {
				fmt.Printf("test %d failed\ninput:\n%s\nexpected: %d\ngot: %d\n", i+1, input, expectVals[j], gotVals[j])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
