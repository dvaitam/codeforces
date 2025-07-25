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

func checkSolution(n, m int, w [][]int64) bool {
	b := make([]int64, m)
	for j := 0; j < m; j++ {
		b[j] = w[0][j]
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = w[i][0] - b[0]
	}
	gcd := func(a, b int64) int64 {
		for b != 0 {
			a, b = b, a%b
		}
		return a
	}
	abs := func(x int64) int64 {
		if x < 0 {
			return -x
		}
		return x
	}
	var g int64
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			g = gcd(g, abs(w[i][j]-a[i]-b[j]))
		}
	}
	if g == 0 {
		g = 1000000001
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if g <= w[i][j] {
				return false
			}
		}
	}
	return true
}

func runCase(bin string, w [][]int64) error {
	n := len(w)
	m := len(w[0])
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(w[i][j], 10))
		}
		sb.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	first := strings.TrimSpace(scanner.Text())
	possible := checkSolution(n, m, w)
	if first == "NO" {
		if possible {
			return fmt.Errorf("should be YES")
		}
		if scanner.Scan() {
			return fmt.Errorf("extra output after NO")
		}
		return nil
	}
	if first != "YES" {
		return fmt.Errorf("expected YES or NO")
	}
	if !possible {
		return fmt.Errorf("should be NO")
	}
	if !scanner.Scan() {
		return fmt.Errorf("missing k")
	}
	k, err := strconv.ParseInt(strings.TrimSpace(scanner.Text()), 10, 64)
	if err != nil {
		return fmt.Errorf("bad k: %v", err)
	}
	if k <= 0 {
		return fmt.Errorf("k must be positive")
	}
	if !scanner.Scan() {
		return fmt.Errorf("missing a line")
	}
	fields := strings.Fields(scanner.Text())
	if len(fields) != n {
		return fmt.Errorf("expected %d numbers for a", n)
	}
	a := make([]int64, n)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return fmt.Errorf("bad a: %v", err)
		}
		a[i] = val
	}
	if !scanner.Scan() {
		return fmt.Errorf("missing b line")
	}
	fields = strings.Fields(scanner.Text())
	if len(fields) != m {
		return fmt.Errorf("expected %d numbers for b", m)
	}
	b := make([]int64, m)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return fmt.Errorf("bad b: %v", err)
		}
		b[i] = val
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output")
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if ((a[i]+b[j])%k+k)%k != w[i][j] {
				return fmt.Errorf("wrong matrix value at %d,%d", i+1, j+1)
			}
		}
	}
	return nil
}

func generateMatrix(rng *rand.Rand) [][]int64 {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	k := rng.Int63n(1000) + 2
	a := make([]int64, n)
	b := make([]int64, m)
	for i := 0; i < n; i++ {
		a[i] = rng.Int63n(k)
	}
	for j := 0; j < m; j++ {
		b[j] = rng.Int63n(k)
	}
	w := make([][]int64, n)
	for i := 0; i < n; i++ {
		w[i] = make([]int64, m)
		for j := 0; j < m; j++ {
			w[i][j] = (a[i] + b[j]) % k
		}
	}
	if rng.Intn(2) == 0 {
		// invalidate
		i := rng.Intn(n)
		j := rng.Intn(m)
		w[i][j] = (w[i][j] + 1) % k
		if w[i][j] == (a[i]+b[j])%k {
			w[i][j] = (w[i][j] + 1) % k
		}
	}
	return w
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		w := generateMatrix(rng)
		if err := runCase(bin, w); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
