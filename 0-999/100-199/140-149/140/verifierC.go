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

func tryC(A []int, q int) bool {
	n := len(A)
	if n < 3*q {
		return false
	}
	k := 0
	if A[q-1] == A[n-q] {
		return false
	}
	if A[q-1] == A[q] {
		j := 0
		for ; j < q; j++ {
			if A[j] == A[q] {
				break
			}
		}
		i := q
		for ; i < n; i++ {
			if A[i] != A[q] {
				break
			}
		}
		i -= q
		if i > j {
			k += i - j
		}
	}
	if n-k < 3*q {
		return false
	}
	if A[n-q-1] == A[n-q] {
		j := q
		for ; j < n-q; j++ {
			if A[j] == A[n-q] {
				break
			}
		}
		j -= k + q
		i := n - q
		for ; i < n; i++ {
			if A[i] != A[n-q] {
				break
			}
		}
		i -= (n - q)
		if i > j {
			return false
		}
	}
	return true
}

func outC(A []int, q int) string {
	n := len(A)
	i, j, k := 0, q, n-q
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for cnt := q; cnt > 0; cnt-- {
		for A[i] == A[j] {
			j++
		}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", A[k], A[j], A[i]))
		i++
		j++
		k++
	}
	return strings.TrimRight(sb.String(), "\n")
}

func solveC(nums []int) string {
	A := append([]int(nil), nums...)
	sort.Ints(A)
	l, r := 0, len(A)
	for r > l+1 {
		m := (l + r) / 2
		if tryC(A, m) {
			l = m
		} else {
			r = m
		}
	}
	return outC(A, l)
}

func generateCaseC(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 3
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(20) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(arr[i]))
	}
	sb.WriteByte('\n')
	input := sb.String()
	expect := solveC(arr)
	return input, expect
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
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCaseC(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
