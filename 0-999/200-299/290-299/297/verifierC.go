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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randUnique(rng *rand.Rand, n int) []int {
	m := make(map[int]struct{})
	res := make([]int, 0, n)
	for len(res) < n {
		v := rng.Intn(50)
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		res = append(res, v)
	}
	return res
}

func dupCount(arr []int) int {
	freq := make(map[int]int)
	rem := 0
	for _, v := range arr {
		freq[v]++
	}
	for _, f := range freq {
		if f > 1 {
			rem += f - 1
		}
	}
	return rem
}

func verifyOutput(s []int, out string) error {
	tokens := strings.Fields(out)
	if len(tokens) < 1 {
		return fmt.Errorf("no output")
	}
	if strings.ToUpper(tokens[0]) != "YES" {
		return fmt.Errorf("first token should be YES")
	}
	if len(tokens) < 1+2*len(s) {
		return fmt.Errorf("not enough numbers")
	}
	a := make([]int, len(s))
	b := make([]int, len(s))
	idx := 1
	for i := 0; i < len(s); i++ {
		fmt.Sscan(tokens[idx], &a[i])
		idx++
	}
	for i := 0; i < len(s); i++ {
		fmt.Sscan(tokens[idx], &b[i])
		idx++
	}
	for i := 0; i < len(s); i++ {
		if a[i] < 0 || b[i] < 0 {
			return fmt.Errorf("negative value")
		}
		if a[i]+b[i] != s[i] {
			return fmt.Errorf("sum mismatch at %d", i)
		}
	}
	limit := (len(s) + 2) / 3
	if dupCount(a) > limit || dupCount(b) > limit {
		return fmt.Errorf("not almost unique")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		n := rng.Intn(8) + 1
		s := randUnique(rng, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j, v := range s {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := verifyOutput(s, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
