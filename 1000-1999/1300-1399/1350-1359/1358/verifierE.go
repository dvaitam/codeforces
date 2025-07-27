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

func expected(n int, a []int64, x int64) string {
	m := (n + 1) / 2
	s := int64(0)
	for i := 0; i < m; i++ {
		s += a[i]
	}
	if s+int64(n/2)*x <= 0 && x >= 0 {
		return "-1"
	}
	le := n
	half := m
	for i := 0; le*2 > n && i+le <= n; i++ {
		for le*2 > n && s+x*int64(le+i-half) <= 0 {
			le--
		}
		s -= a[i]
	}
	if le*2 > n {
		return fmt.Sprintf("%d", le)
	}
	return "-1"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 1; t <= 100; t++ {
		n := rng.Intn(20) + 2 // ensure at least 2
		m := (n + 1) / 2
		a := make([]int64, m)
		for i := 0; i < m; i++ {
			a[i] = int64(rng.Intn(20) - 10)
		}
		x := int64(rng.Intn(20) - 10)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i*2 < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[i]))
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d\n", x))
		want := expected(n, append([]int64(nil), a...), x)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", t, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
