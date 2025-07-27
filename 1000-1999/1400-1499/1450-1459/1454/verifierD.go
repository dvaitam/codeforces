package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func isPrime(n uint64) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	for i := uint64(5); i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

func primeFac(n uint64) (uint64, uint64) {
	m := n
	var c, p uint64
	var cnt uint64
	for m%2 == 0 {
		cnt++
		m /= 2
	}
	if cnt > c {
		c = cnt
		p = 2
	}
	for i := uint64(3); i*i <= m; i += 2 {
		cnt = 0
		for m%i == 0 {
			cnt++
			m /= i
		}
		if cnt > c {
			c = cnt
			p = i
		}
	}
	if m > 1 && c == 0 {
		c = 1
		p = m
	}
	return c, p
}

func solveCaseD(n uint64) (uint64, []uint64) {
	if isPrime(n) {
		return 1, []uint64{n}
	}
	c, p := primeFac(n)
	seq := make([]uint64, c)
	rem := n
	for i := uint64(0); i < c-1; i++ {
		seq[i] = p
		rem /= p
	}
	seq[c-1] = rem
	return c, seq
}

func generateTests() ([]uint64, string) {
	const t = 100
	r := rand.New(rand.NewSource(4))
	ns := make([]uint64, t)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := uint64(r.Intn(1_000_000_000-2) + 2)
		ns[i] = n
		fmt.Fprintf(&sb, "%d\n", n)
	}
	return ns, sb.String()
}

func verify(ns []uint64, output string) error {
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanWords)
	for idx, n := range ns {
		if !scanner.Scan() {
			return fmt.Errorf("case %d: missing k", idx+1)
		}
		var k uint64
		fmt.Sscan(scanner.Text(), &k)
		seq := make([]uint64, k)
		for i := uint64(0); i < k; i++ {
			if !scanner.Scan() {
				return fmt.Errorf("case %d: incomplete sequence", idx+1)
			}
			fmt.Sscan(scanner.Text(), &seq[i])
		}
		expectedK, expectedSeq := solveCaseD(n)
		if k != expectedK {
			return fmt.Errorf("case %d: expected k=%d got %d", idx+1, expectedK, k)
		}
		for i := uint64(0); i < k; i++ {
			if seq[i] != expectedSeq[i] {
				return fmt.Errorf("case %d: expected %v got %v", idx+1, expectedSeq, seq)
			}
		}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output: %s", scanner.Text())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go <binary>")
		os.Exit(1)
	}
	ns, input := generateTests()
	out, err := runBinary(os.Args[1], input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		os.Exit(1)
	}
	if err := verify(ns, out); err != nil {
		fmt.Fprintln(os.Stderr, "verification failed:", err)
		os.Exit(1)
	}
	fmt.Println("All tests passed for problem D")
}
