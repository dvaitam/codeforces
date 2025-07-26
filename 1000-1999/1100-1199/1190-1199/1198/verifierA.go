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

func run(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveA(n int, I int64, a []int) int {
	sort.Ints(a)
	vals := []int{a[0]}
	cnt := []int{1}
	for i := 1; i < n; i++ {
		if a[i] == vals[len(vals)-1] {
			cnt[len(cnt)-1]++
		} else {
			vals = append(vals, a[i])
			cnt = append(cnt, 1)
		}
	}
	m := len(vals)
	B := (I * 8) / int64(n)
	if B >= 31 {
		return 0
	}
	W := 1 << B
	if W >= m {
		return 0
	}
	pref := make([]int, m+1)
	for i := 0; i < m; i++ {
		pref[i+1] = pref[i] + cnt[i]
	}
	best := 0
	for i := 0; i+W <= m; i++ {
		sum := pref[i+W] - pref[i]
		if sum > best {
			best = sum
		}
	}
	return n - best
}

func generateCase(rng *rand.Rand) ([]byte, int) {
	n := rng.Intn(20) + 1
	I := int64(rng.Intn(20) + 1)
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(40)
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, "%d %d\n", n, I)
	for i, v := range a {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprint(&b, v)
	}
	b.WriteByte('\n')
	return b.Bytes(), solveA(n, I, a)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, expect := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, string(input))
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		if got != fmt.Sprint(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", i, expect, got, string(input))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
