package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(n int, a, b []int, d []int) string {
	cnt := map[int]int{}
	for _, x := range d {
		cnt[x]++
	}
	need := map[int]int{}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			need[b[i]]++
		}
	}
	for val, c := range need {
		if cnt[val] < c {
			return "NO"
		}
	}
	return "YES"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(5) + 1
		if rng.Intn(2) == 0 {
			b[i] = a[i]
		} else {
			b[i] = rng.Intn(5) + 1
		}
	}
	m := rng.Intn(8) + 1
	d := make([]int, m)
	for i := 0; i < m; i++ {
		d[i] = rng.Intn(5) + 1
	}
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	lineA := sb.String()
	sb.Reset()
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	lineB := sb.String()
	sb.Reset()
	for i, v := range d {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	lineD := sb.String()
	input := fmt.Sprintf("1\n%d\n%s\n%s\n%d\n%s\n", n, lineA, lineB, m, lineD)
	expect := solveCase(n, a, b, d)
	return input, expect
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
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
