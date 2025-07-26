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

func solveD(a []int) (string, []int) {
	n := len(a)
	freq := make([]int, n+1)
	for _, v := range a {
		if v >= 1 && v <= n {
			freq[v]++
		}
	}
	missing := make([]int, 0)
	for v := 1; v <= n; v++ {
		if freq[v] == 0 {
			missing = append(missing, v)
		}
	}
	missIdx := 0
	used := make([]bool, n+1)
	changes := 0
	for i := 0; i < n; i++ {
		v := a[i]
		if v < 1 || v > n {
			a[i] = missing[missIdx]
			missIdx++
			changes++
			continue
		}
		if !used[v] {
			if freq[v] == 1 {
				used[v] = true
				continue
			}
			if missIdx < len(missing) && missing[missIdx] < v {
				freq[v]--
				a[i] = missing[missIdx]
				missIdx++
				changes++
			} else {
				used[v] = true
				freq[v]--
			}
		} else {
			freq[v]--
			a[i] = missing[missIdx]
			missIdx++
			changes++
		}
	}
	return fmt.Sprintf("%d", changes), a
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 2
	a := make([]int, n)
	for i := range a {
		if rng.Intn(2) == 0 {
			a[i] = rng.Intn(n*2) + 1
		} else {
			a[i] = rng.Intn(n) + 1
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	expChanges, arr := solveD(append([]int(nil), a...))
	var out strings.Builder
	out.WriteString(expChanges)
	out.WriteByte('\n')
	for i, v := range arr {
		if i > 0 {
			out.WriteByte(' ')
		}
		out.WriteString(fmt.Sprintf("%d", v))
	}
	return sb.String(), out.String()
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
