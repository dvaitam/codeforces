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

func solveE(n, k int, arr, c []int64) string {
	a := make([]int64, n+1)
	copy(a, arr)
	cc := make([]int64, n)
	copy(cc, c)
	for i := n - 1; i >= 0; i-- {
		mul := a[i] * 100
		if mul > cc[i] {
			a[i] = cc[i]
		} else {
			a[i] = mul
			if a[i+1] > a[i] {
				if cc[i] < a[i+1] {
					a[i] = cc[i]
				} else {
					a[i] = a[i+1]
				}
			}
		}
	}
	slice := a[:n]
	sort.Slice(slice, func(i, j int) bool { return slice[i] < slice[j] })
	ans := float64(slice[0])
	tot := 1.0
	limit := n - k
	for i := 0; i < limit; i++ {
		tot *= float64(limit-i) / float64(n-i)
		if tot < 1e-18 {
			break
		}
		ans += tot * float64(slice[i+1]-slice[i])
	}
	return fmt.Sprintf("%f\n", ans)
}

func generateCaseE(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	k := rng.Intn(n) + 1
	arr := make([]int64, n)
	c := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(10) + 1)
	}
	for i := 0; i < n; i++ {
		c[i] = int64(rng.Intn(10) + 1)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for i, v := range c {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	expect := solveE(n, k, arr, c)
	return sb.String(), expect
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCaseE(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
