package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

var primes = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47}

func genSeq(n int) []int {
	t := 2 * n * n
	A := make([]int, n)
	AA := make([]int, n)
	A[0] = 1
	m := 1
	for k := 0; m < n && k < len(primes); k++ {
		p := primes[k]
		mm := 0
		d := p
		l := m
		cnt := 0
		for mm < n && d <= t {
			var i int
			for i = 0; i < l && mm < n; i++ {
				if d*A[i] <= t {
					AA[mm] = A[i] * d
					mm++
					if d > 1 {
						cnt++
					}
				} else {
					break
				}
			}
			if d == p {
				l = i
			}
			if d == p {
				d = 1
			} else if d == 1 {
				d = p * p
			} else {
				d *= p
			}
		}
		for cnt >= (mm+2)/2 && l < m && mm < n {
			AA[mm] = A[l]
			mm++
			l++
		}
		if mm == m {
			break
		}
		m = mm
		for i := 0; i < m; i++ {
			A[i] = AA[i]
		}
		sort.Ints(A[:m])
	}
	return A[:n]
}

func runCandidate(bin, input string) (string, error) {
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

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	seq := genSeq(n)
	var exp strings.Builder
	for i, v := range seq {
		if i > 0 {
			exp.WriteByte(' ')
		}
		exp.WriteString(fmt.Sprintf("%d", v))
	}
	input := fmt.Sprintf("%d\n", n)
	return input, exp.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
