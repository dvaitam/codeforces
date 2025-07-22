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

const mod = 1000000007

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

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func modInv(x int64) int64 { return modPow(x, mod-2) }

func solveCase(n int, arr []int64) string {
	exp := make(map[int64]int)
	for _, v := range arr {
		x := v
		for p := int64(2); p*p <= x; p++ {
			for x%p == 0 {
				exp[p]++
				x /= p
			}
		}
		if x > 1 {
			exp[x]++
		}
	}
	maxE := 0
	for _, e := range exp {
		if e > maxE {
			maxE = e
		}
	}
	lim := maxE + n - 1
	fact := make([]int64, lim+1)
	inv := make([]int64, lim+1)
	fact[0] = 1
	for i := 1; i <= lim; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	inv[lim] = modInv(fact[lim])
	for i := lim; i >= 1; i-- {
		inv[i-1] = inv[i] * int64(i) % mod
	}
	res := int64(1)
	for _, e := range exp {
		top := e + n - 1
		ways := fact[top] * inv[n-1] % mod * inv[e] % mod
		res = res * ways % mod
	}
	if res < 0 {
		res += mod
	}
	return fmt.Sprintf("%d", res)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			v, _ := strconv.ParseInt(scan.Text(), 10, 64)
			arr[j] = v
		}
		var input bytes.Buffer
		fmt.Fprintln(&input, n)
		for j, v := range arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		expected := solveCase(n, arr)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
