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
	src := filepath.Join(dir, "420C.go")
	tmp, err := os.CreateTemp("", "oracle420C")
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

func run(bin, input string) (string, error) {
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

func randExcluding(r *rand.Rand, n int, forbid ...int) int {
	for {
		x := r.Intn(n) + 1
		ok := true
		for _, f := range forbid {
			if x == f {
				ok = false
				break
			}
		}
		if ok {
			return x
		}
	}
}

func genCase(r *rand.Rand) string {
	n := 3 + r.Intn(80)
	if r.Intn(5) == 0 {
		n = 3 + r.Intn(1000)
	}
	p := r.Intn(n + 1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, p)
	target := r.Intn(n) + 1
	for i := 1; i <= n; i++ {
		var x, y int
		if r.Intn(4) == 0 {
			x = randExcluding(r, n, i)
			y = randExcluding(r, n, i, x)
		} else if r.Intn(2) == 0 {
			x = target
			if x == i {
				x = randExcluding(r, n, i)
			}
			y = randExcluding(r, n, i, x)
		} else {
			x = randExcluding(r, n, i)
			if r.Intn(3) == 0 {
				y = (x % n) + 1
				if y == i || y == x {
					y = randExcluding(r, n, i, x)
				}
			} else {
				y = randExcluding(r, n, i, x)
			}
		}
		fmt.Fprintf(&sb, "%d %d\n", x, y)
	}
	return sb.String()
}

func parseOutput(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	if len(fields) > 1 {
		return 0, fmt.Errorf("too many numbers in output")
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
		input := genCase(r)
		expectStr, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotStr, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expectVal, err := parseOutput(expectStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotVal, err := parseOutput(gotStr)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", i+1, input, err)
			os.Exit(1)
		}
		if expectVal != gotVal {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected: %d\ngot: %d\n", i+1, input, expectVal, gotVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
