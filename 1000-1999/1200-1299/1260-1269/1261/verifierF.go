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

const testCount = 120
const mod = 998244353

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1261F.go")
	tmp, err := os.CreateTemp("", "oracle1261F")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()
	os.Remove(path)
	cmd := exec.Command("go", "build", "-o", path, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(out))
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

func genSegments(r *rand.Rand, maxSeg int, maxVal uint64) ([][2]uint64, int) {
	count := 1 + r.Intn(maxSeg)
	segs := make([][2]uint64, count)
	for i := 0; i < count; i++ {
		l := r.Uint64()%(maxVal/2+1) + 1
		length := r.Uint64()%uint64(1+int(maxVal/10)) + 1
		rVal := l + length
		if rVal > maxVal {
			rVal = maxVal
		}
		if rVal < l {
			l, rVal = rVal, l
		}
		segs[i] = [2]uint64{l, rVal}
	}
	return segs, count
}

func formatInput(segA [][2]uint64, segB [][2]uint64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(segA))
	for _, seg := range segA {
		fmt.Fprintf(&sb, "%d %d\n", seg[0], seg[1])
	}
	fmt.Fprintf(&sb, "%d\n", len(segB))
	for _, seg := range segB {
		fmt.Fprintf(&sb, "%d %d\n", seg[0], seg[1])
	}
	return sb.String()
}

func parseOutput(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %d", len(fields))
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, err
	}
	val %= mod
	if val < 0 {
		val += mod
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
	for t := 0; t < testCount; t++ {
		segA, _ := genSegments(r, 5, uint64(1e6))
		segB, _ := genSegments(r, 5, uint64(1e6))
		input := formatInput(segA, segB)
		expectStr, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotStr, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		expectVal, err := parseOutput(expectStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotVal, err := parseOutput(gotStr)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", t+1, input, err)
			os.Exit(1)
		}
		if expectVal != gotVal {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected: %d\ngot: %d\n", t+1, input, expectVal, gotVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
