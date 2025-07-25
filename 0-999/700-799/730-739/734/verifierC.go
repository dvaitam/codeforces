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

func expected(n, m, k int, x, s int64, a, b, c, d []int64) string {
	ans := int64(n) * x
	for i := 0; i <= m; i++ {
		cost1 := b[i]
		if cost1 > s {
			continue
		}
		rem := s - cost1
		idx := sort.Search(len(d), func(j int) bool { return d[j] > rem })
		j := idx - 1
		if j < 0 {
			j = 0
		}
		need := int64(n) - c[j]
		if need < 0 {
			need = 0
		}
		t := need * a[i]
		if t < ans {
			ans = t
		}
	}
	return fmt.Sprint(ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	m := rng.Intn(5) + 1
	k := rng.Intn(5) + 1
	x := int64(rng.Intn(20) + 2)
	s := int64(rng.Intn(100) + 1)

	a := make([]int64, m+1)
	b := make([]int64, m+1)
	a[0] = x
	b[0] = 0
	for i := 1; i <= m; i++ {
		a[i] = int64(rng.Intn(int(x-1)) + 1)
	}
	for i := 1; i <= m; i++ {
		b[i] = int64(rng.Intn(100) + 1)
	}

	c := make([]int64, k+1)
	d := make([]int64, k+1)
	c[0] = 0
	d[0] = 0
	curC := int64(0)
	curD := int64(0)
	for i := 1; i <= k; i++ {
		curC += int64(rng.Intn(n) + 1)
		if curC > int64(n) {
			curC = int64(n)
		}
		c[i] = curC
	}
	for i := 1; i <= k; i++ {
		curD += int64(rng.Intn(100) + 1)
		d[i] = curD
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	sb.WriteString(fmt.Sprintf("%d %d\n", x, s))
	for i := 1; i <= m; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(a[i]))
	}
	sb.WriteByte('\n')
	for i := 1; i <= m; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(b[i]))
	}
	sb.WriteByte('\n')
	for i := 1; i <= k; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(c[i]))
	}
	sb.WriteByte('\n')
	for i := 1; i <= k; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(d[i]))
	}
	sb.WriteByte('\n')

	input := sb.String()
	exp := expected(n, m, k, x, s, a, b, c, d)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
