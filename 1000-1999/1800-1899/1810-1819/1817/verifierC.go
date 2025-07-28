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

const numTestsC = 100
const MOD int64 = 1000000007

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "binC*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), path)
		if out, err := cmd.CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, nil, nil
}

func prepareOracle() (string, func(), error) {
	tmp, err := os.CreateTemp("", "oracleC*")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), "1817C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", nil, fmt.Errorf("go build oracle failed: %v: %s", err, out)
	}
	return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), err
}

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func evalPoly(coef []int64, x int64) int64 {
	res := int64(0)
	pow := int64(1)
	for _, c := range coef {
		res = (res + c*pow) % MOD
		pow = pow * x % MOD
	}
	return res
}

func genCase(rng *rand.Rand) string {
	d := rng.Intn(3) + 1
	coef := make([]int64, d+1)
	for i := range coef {
		coef[i] = rng.Int63n(10) + 1
	}
	s := rng.Int63n(10)
	A := make([]int64, d+1)
	B := make([]int64, d+1)
	for i := 0; i <= d; i++ {
		A[i] = evalPoly(coef, int64(i))
		B[i] = evalPoly(coef, int64(i)+s)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", d)
	for i, v := range A {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range B {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin, oracle, input string) error {
	exp, err := run(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	got, err := run(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	if strings.TrimSpace(got) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin, cleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if cleanup != nil {
		defer cleanup()
	}
	oracle, cleanOracle, err := prepareOracle()
	if err != nil {
		fmt.Println("oracle compile error:", err)
		return
	}
	defer cleanOracle()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < numTestsC; i++ {
		in := genCase(rng)
		if err := runCase(bin, oracle, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			return
		}
	}
	fmt.Println("All tests passed")
}
