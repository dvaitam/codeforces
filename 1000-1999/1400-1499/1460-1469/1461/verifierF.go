package main

import (
	"bytes"
	"fmt"
	"math/big"
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
	oracle := filepath.Join(dir, "oracleF")
	cmd := exec.Command("go", "build", "-o", oracle, "1461F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func genCase(r *rand.Rand) (string, []int, map[rune]bool) {
	n := r.Intn(15) + 1
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = r.Intn(10)
	}
	opsAll := []rune{"+", "-", "*"}
	opsCount := r.Intn(3) + 1
	ops := make([]rune, 0, opsCount)
	used := make(map[rune]bool)
	for len(ops) < opsCount {
		c := opsAll[r.Intn(3)]
		if !used[c] {
			used[c] = true
			ops = append(ops, c)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(nums[i]))
	}
	sb.WriteByte('\n')
	for _, c := range ops {
		sb.WriteRune(c)
	}
	sb.WriteByte('\n')
	allowed := make(map[rune]bool)
	for _, c := range ops {
		allowed[c] = true
	}
	return sb.String(), nums, allowed
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

func evalExpr(nums []int, allowed map[rune]bool, expr string) (*big.Int, error) {
	expr = strings.ReplaceAll(strings.TrimSpace(expr), " ", "")
	if len(expr) != len(nums)+(len(nums)-1) {
		return nil, fmt.Errorf("invalid expression length")
	}
	val := big.NewInt(int64(nums[0]))
	pos := 1
	for i := 1; i < len(nums); i++ {
		op := rune(expr[pos-1])
		if !allowed[op] {
			return nil, fmt.Errorf("invalid operator")
		}
		digit := expr[pos]
		if digit < '0' || digit > '9' {
			return nil, fmt.Errorf("invalid digit")
		}
		if int(digit-'0') != nums[i] {
			return nil, fmt.Errorf("unexpected digit")
		}
		pos += 2
		tmp := big.NewInt(int64(nums[i]))
		switch op {
		case '+':
			val.Add(val, tmp)
		case '-':
			val.Sub(val, tmp)
		case '*':
			val.Mul(val, tmp)
		default:
			return nil, fmt.Errorf("bad op")
		}
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
		input, nums, allowed := genCase(rng)
		expectExpr, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		expectVal, err := evalExpr(nums, allowed, expectExpr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle produced invalid output on case %d", i+1)
			os.Exit(1)
		}
		gotExpr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		gotVal, err := evalExpr(nums, allowed, gotExpr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d invalid output: %v\ninput:%s\noutput:%s", i+1, err, input, gotExpr)
			os.Exit(1)
		}
		if gotVal.Cmp(expectVal) != 0 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected value %s got %s\ninput:%s", i+1, expectVal.String(), gotVal.String(), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
