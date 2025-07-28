package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func solveD(n int, m int64, a []int64, b []int) string {
	var total int64
	ones := []int64{}
	twos := []int64{}
	for i := 0; i < n; i++ {
		total += a[i]
		if b[i] == 1 {
			ones = append(ones, a[i])
		} else {
			twos = append(twos, a[i])
		}
	}
	if total < m {
		return "-1"
	}
	sort.Slice(ones, func(i, j int) bool { return ones[i] > ones[j] })
	sort.Slice(twos, func(i, j int) bool { return twos[i] > twos[j] })
	pref1 := make([]int64, len(ones)+1)
	for i := 0; i < len(ones); i++ {
		pref1[i+1] = pref1[i] + ones[i]
	}
	pref2 := make([]int64, len(twos)+1)
	for i := 0; i < len(twos); i++ {
		pref2[i+1] = pref2[i] + twos[i]
	}
	ans := math.MaxInt32
	i := len(ones)
	for j := 0; j <= len(twos); j++ {
		mem := pref2[j]
		for i > 0 && mem+pref1[i-1] >= m {
			i--
		}
		if mem+pref1[i] >= m {
			cost := 2*j + i
			if cost < ans {
				ans = cost
			}
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateD(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	a := make([]int64, n)
	b := make([]int, n)
	var total int64
	for i := 0; i < n; i++ {
		a[i] = rng.Int63n(20) + 1
		total += a[i]
		if rng.Intn(2) == 0 {
			b[i] = 1
		} else {
			b[i] = 2
		}
	}
	m := rng.Int63n(total + 1)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", b[i]))
	}
	sb.WriteByte('\n')
	out := solveD(n, m, a, b)
	return sb.String(), out
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateD(rng)
		got, err := run(bin, in)
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
