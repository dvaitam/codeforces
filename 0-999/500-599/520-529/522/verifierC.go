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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type testCaseC struct {
	input    string
	expected string
}

func computeC(m, k int, a []int, events [][2]int) string {
	cnt := make([]int, k+1)
	unknown := 0
	for _, ev := range events {
		ti, _ := ev[0], ev[1]
		if ti > 0 && ti <= k {
			cnt[ti]++
		} else {
			unknown++
		}
	}
	var sb strings.Builder
	for i := 1; i <= k; i++ {
		if cnt[i]+unknown >= a[i] {
			sb.WriteByte('Y')
		} else {
			sb.WriteByte('N')
		}
	}
	return sb.String()
}

func generateCaseC() testCaseC {
	m := rand.Intn(4) + 2 // 2..5
	k := rand.Intn(4) + 1 // 1..5
	a := make([]int, k+1)
	sum := 0
	for i := 1; i <= k; i++ {
		a[i] = rand.Intn(3) + 1
		sum += a[i]
	}
	if sum < m {
		a[1] += m - sum
		sum = m
	}
	events := make([][2]int, m-1)
	for i := 0; i < m-1; i++ {
		ti := rand.Intn(k + 1)
		ri := rand.Intn(2)
		events[i] = [2]int{ti, ri}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", m, k))
	for i := 1; i <= k; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	for _, ev := range events {
		sb.WriteString(fmt.Sprintf("%d %d\n", ev[0], ev[1]))
	}
	return testCaseC{input: sb.String(), expected: computeC(m, k, a, events)}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 100; i++ {
		tc := generateCaseC()
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, tc.input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d: expected %s got %s\ninput:\n%s", i, tc.expected, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
