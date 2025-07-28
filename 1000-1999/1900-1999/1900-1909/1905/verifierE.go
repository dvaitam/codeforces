package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod int64 = 998244353

func powMod(a, e int64) int64 {
	a %= mod
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

var memo = map[int64][2]int64{}

func xy(n int64) (int64, int64) {
	if n == 1 {
		return 1, 0
	}
	if val, ok := memo[n]; ok {
		return val[0], val[1]
	}
	right := n / 2
	left := n - right
	XL, YL := xy(left)
	XR, YR := xy(right)
	f := powMod(2, n) - powMod(2, left) - powMod(2, right) + 1
	f %= mod
	if f < 0 {
		f += mod
	}
	X := (f + 2*XL + 2*XR) % mod
	Y := (YL + XR + YR) % mod
	memo[n] = [2]int64{X, Y}
	return X, Y
}

func solve(n int64) int64 {
	X, Y := xy(n)
	return (X + Y) % mod
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
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		n, _ := strconv.ParseInt(line, 10, 64)
		want := fmt.Sprintf("%d", solve(n))
		input := fmt.Sprintf("1\n%d\n", n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
