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

func expectedA(a, b []int) string {
	n := len(a)
	diff := make([]int, n)
	for i := 0; i < n; i++ {
		diff[i] = b[i] - a[i]
		if diff[i] < 0 {
			return "NO"
		}
	}
	i := 0
	for i < n && diff[i] == 0 {
		i++
	}
	if i == n {
		return "YES"
	}
	k := diff[i]
	if k <= 0 {
		return "NO"
	}
	for i < n && diff[i] == k {
		i++
	}
	for i < n {
		if diff[i] != 0 {
			return "NO"
		}
		i++
	}
	return "YES"
}

func generateCaseA(rng *rand.Rand) ([]int, []int) {
	n := rng.Intn(6) + 1
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(11)
		b[i] = rng.Intn(11)
	}
	return a, b
}

func runCaseA(bin string, a, b []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(a)))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteString("\n")
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteString("\n")

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	got = strings.ToUpper(got)
	expect := expectedA(a, b)
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		a, b := generateCaseA(rng)
		if err := runCaseA(bin, a, b); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
