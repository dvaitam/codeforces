package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func isValid(a []int, l, r int) bool {
	n := len(a)
	freq := make(map[int]int)
	for i := l - 1; i < r; i++ {
		freq[a[i]]++
	}
	for i := 0; i < n; i++ {
		j := n - 1 - i
		if i > j {
			break
		}
		inI := i >= l-1 && i < r
		inJ := j >= l-1 && j < r
		if !inI && !inJ {
			if a[i] != a[j] {
				return false
			}
			continue
		}
		if inI && !inJ {
			need := a[j]
			if freq[need] == 0 {
				return false
			}
			freq[need]--
		} else if !inI && inJ {
			need := a[i]
			if freq[need] == 0 {
				return false
			}
			freq[need]--
		}
	}
	odd := 0
	for _, v := range freq {
		if v%2 == 1 {
			odd++
		}
	}
	return odd <= 1
}

func solveE(a []int) int64 {
	n := len(a)
	var ans int64
	for l := 1; l <= n; l++ {
		for r := l; r <= n; r++ {
			if isValid(a, l, r) {
				ans++
			}
		}
	}
	return ans
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	expect := fmt.Sprintf("%d\n", solveE(arr))
	return sb.String(), expect
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
