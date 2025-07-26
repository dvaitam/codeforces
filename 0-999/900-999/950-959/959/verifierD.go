package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runCandidate(bin string, input string) (string, error) {
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

const limit = 2000005

func linearSieve(n int) ([]int, []int) {
	spf := make([]int, n+1)
	primes := make([]int, 0)
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			spf[i] = i
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p > spf[i] || i*p > n {
				break
			}
			spf[i*p] = p
		}
	}
	return primes, spf
}

func solveD(data string) string {
	reader := bufio.NewReader(strings.NewReader(data))
	var n int
	fmt.Fscan(reader, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	primes, spf := linearSieve(limit)
	used := make([]bool, limit+1)
	ans := make([]int, n)
	modified := false
	primeIdx := 0
	for i := 0; i < n; i++ {
		if !modified {
			x := a[i]
			for ; x <= limit; x++ {
				ok := true
				t := x
				for t > 1 {
					p := spf[t]
					if used[p] {
						ok = false
						break
					}
					for t%p == 0 {
						t /= p
					}
				}
				if ok {
					break
				}
			}
			ans[i] = x
			t := x
			for t > 1 {
				p := spf[t]
				used[p] = true
				for t%p == 0 {
					t /= p
				}
			}
			if x > a[i] {
				modified = true
			}
		} else {
			for primeIdx < len(primes) && used[primes[primeIdx]] {
				primeIdx++
			}
			if primeIdx < len(primes) {
				ans[i] = primes[primeIdx]
				used[primes[primeIdx]] = true
				primeIdx++
			} else {
				ans[i] = 2
			}
		}
	}
	var sb strings.Builder
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return sb.String()
}

func genCaseD(rng *rand.Rand) string {
	n := rng.Intn(8) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(rng.Intn(50) + 1))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCaseD(rng)
		expect := solveD(in)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
