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

const MOD int64 = 998244353

var fact, inv []int64

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

func initComb(n int) {
	fact = make([]int64, n+1)
	inv = make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	inv[n] = modPow(fact[n], MOD-2)
	for i := n; i > 0; i-- {
		inv[i-1] = inv[i] * int64(i) % MOD
	}
}

func C(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * inv[k] % MOD * inv[n-k] % MOD
}

type BIT struct {
	n    int
	tree []int
}

func (b *BIT) init(n int) { b.n = n; b.tree = make([]int, n+2) }
func (b *BIT) add(i, delta int) {
	for i <= b.n {
		b.tree[i] += delta
		i += i & -i
	}
}
func (b *BIT) kth(k int) int {
	idx := 0
	bit := 1
	for bit<<1 <= b.n {
		bit <<= 1
	}
	for bit > 0 {
		next := idx + bit
		if next <= b.n && b.tree[next] < k {
			k -= b.tree[next]
			idx = next
		}
		bit >>= 1
	}
	return idx + 1
}

func solve(n int, pairs [][2]int) int64 {
	j := make([]int, n+1)
	for i := 1; i <= n; i++ {
		j[i] = i
	}
	for _, p := range pairs {
		j[p[0]] = p[1]
	}
	bit := BIT{}
	bit.init(n)
	for i := 1; i <= n; i++ {
		bit.add(i, 1)
	}
	p := make([]int, n+1)
	for i := n; i >= 1; i-- {
		pos := bit.kth(j[i])
		p[pos] = i
		bit.add(pos, -1)
	}
	r := 0
	for i := 1; i < n; i++ {
		if p[i] > p[i+1] {
			r++
		}
	}
	return C(2*n-1-r, n)
}

func runCase(bin string, n int, pairs [][2]int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(pairs)))
	for _, pr := range pairs {
		sb.WriteString(fmt.Sprintf("%d %d\n", pr[0], pr[1]))
	}
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
	got := strings.TrimSpace(out.String())
	exp := fmt.Sprintf("%d", solve(n, pairs))
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	initComb(400000)
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Println("could not open testcasesD.txt:", err)
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
	for caseNum := 0; caseNum < t; caseNum++ {
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
		pairs := make([][2]int, m)
		for i := 0; i < m; i++ {
			if !scanner.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			x, _ := strconv.Atoi(scanner.Text())
			if !scanner.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			y, _ := strconv.Atoi(scanner.Text())
			pairs[i] = [2]int{x, y}
		}
		if err := runCase(bin, n, pairs); err != nil {
			fmt.Printf("case %d failed: %v\n", caseNum+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
