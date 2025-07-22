package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type constraint struct{ l, r, q int }

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

func solve(n int, cons []constraint) (bool, []int) {
	const B = 30
	diff := make([][]int, B)
	for b := 0; b < B; b++ {
		diff[b] = make([]int, n+1)
	}
	for _, c := range cons {
		for b := 0; b < B; b++ {
			if (c.q>>b)&1 == 1 {
				diff[b][c.l-1]++
				diff[b][c.r]--
			}
		}
	}
	a := make([]int, n)
	for b := 0; b < B; b++ {
		cnt := 0
		for i := 0; i < n; i++ {
			cnt += diff[b][i]
			if cnt > 0 {
				a[i] |= 1 << b
			}
		}
	}
	size := 1
	for size < n {
		size <<= 1
	}
	allOnes := (1 << B) - 1
	seg := make([]int, 2*size)
	for i := 0; i < n; i++ {
		seg[size+i] = a[i]
	}
	for i := n; i < size; i++ {
		seg[size+i] = allOnes
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = seg[2*i] & seg[2*i+1]
	}
	query := func(l, r int) int {
		l += size
		r += size
		res := allOnes
		for l <= r {
			if l&1 == 1 {
				res &= seg[l]
				l++
			}
			if r&1 == 0 {
				res &= seg[r]
				r--
			}
			l >>= 1
			r >>= 1
		}
		return res
	}
	for _, c := range cons {
		if query(c.l-1, c.r-1) != c.q {
			return false, nil
		}
	}
	return true, a
}

func generateCase(rng *rand.Rand) (int, []constraint) {
	n := rng.Intn(20) + 1
	m := rng.Intn(20) + 1
	cons := make([]constraint, m)
	for i := 0; i < m; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n) + 1
		if l > r {
			l, r = r, l
		}
		q := rng.Intn(1 << 30)
		cons[i] = constraint{l, r, q}
	}
	return n, cons
}

func check(n int, cons []constraint, out string) error {
	scan := bufio.NewScanner(strings.NewReader(out))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return fmt.Errorf("no output")
	}
	tok := strings.ToUpper(scan.Text())
	expectOK, _ := solve(n, cons)
	if tok == "NO" {
		if expectOK {
			return fmt.Errorf("solution exists but got NO")
		}
		if scan.Scan() {
			return fmt.Errorf("extra output after NO")
		}
		return nil
	}
	if tok != "YES" {
		return fmt.Errorf("expected YES/NO")
	}
	arr := make([]int, 0, n)
	for scan.Scan() {
		v, err := strconv.Atoi(scan.Text())
		if err != nil {
			return fmt.Errorf("invalid integer")
		}
		arr = append(arr, v)
	}
	if len(arr) != n {
		return fmt.Errorf("expected %d numbers got %d", n, len(arr))
	}
	if !expectOK {
		return fmt.Errorf("expected NO but got YES")
	}
	for _, c := range cons {
		andVal := arr[c.l-1]
		for i := c.l; i <= c.r; i++ {
			andVal &= arr[i-1]
		}
		if andVal != c.q {
			return fmt.Errorf("constraint failed")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, cons := generateCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, len(cons)))
		for _, c := range cons {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", c.l, c.r, c.q))
		}
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(n, cons, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
