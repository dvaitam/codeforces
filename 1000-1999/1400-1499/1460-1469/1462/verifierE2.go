package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const MOD int64 = 1000000007

var fact [200005]int64
var invfact [200005]int64

func modpow(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func init() {
	fact[0] = 1
	for i := 1; i < len(fact); i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invfact[len(fact)-1] = modpow(fact[len(fact)-1], MOD-2)
	for i := len(fact) - 1; i > 0; i-- {
		invfact[i-1] = invfact[i] * int64(i) % MOD
	}
}

func comb(n, k int) int64 {
	if k < 0 || n < k {
		return 0
	}
	return fact[n] * invfact[k] % MOD * invfact[n-k] % MOD
}

func expectedE2(n, m, k int, arr []int) int64 {
	sort.Ints(arr)
	l := 0
	var ans int64
	for r := 0; r < n; r++ {
		for l < r && arr[r]-arr[l] > k {
			l++
		}
		w := r - l
		if w >= m-1 {
			ans = (ans + comb(w, m-1)) % MOD
		}
	}
	return ans
}

func runCase(bin string, n, m, k int, arr []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(arr[i]))
	}
	sb.WriteByte('\n')

	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(gotStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid output %q", gotStr)
	}
	exp := expectedE2(n, m, k, arr)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesE2.txt")
	if err != nil {
		fmt.Println("could not open testcasesE2.txt:", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scanner.Text())
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scanner.Text())
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		m, _ := strconv.Atoi(scanner.Text())
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		kval, _ := strconv.Atoi(scanner.Text())
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if !scanner.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			v, _ := strconv.Atoi(scanner.Text())
			arr[j] = v
		}
		if err := runCase(bin, n, m, kval, arr); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
