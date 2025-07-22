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

func runCandidate(bin, input string) (string, error) {
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

func solveCase(a []int) string {
	n := len(a)
	l := -1
	for i := 0; i < n-1; i++ {
		if a[i] > a[i+1] {
			l = i
			break
		}
	}
	if l == -1 {
		return "yes\n1 1"
	}
	r := -1
	for j := n - 1; j > 0; j-- {
		if a[j-1] > a[j] {
			r = j
			break
		}
	}
	b := append([]int(nil), a...)
	for i, j := l, r; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	ok := true
	for i := 0; i < n-1; i++ {
		if b[i] > b[i+1] {
			ok = false
			break
		}
	}
	if ok {
		return fmt.Sprintf("yes\n%d %d", l+1, r+1)
	}
	return "no"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	vals := rng.Perm(n)
	for i := range vals {
		vals[i]++
	}
	// shuffle more to make unsorted arrays too
	for i := 0; i < n; i++ {
		j := rng.Intn(n)
		vals[i], vals[j] = vals[j], vals[i]
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	expect := solveCase(vals)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
