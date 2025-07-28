package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func expectedC(n, k int, xs []int) int {
	pos := []int{}
	neg := []int{}
	maxAbs := 0
	for _, x := range xs {
		if x > 0 {
			pos = append(pos, x)
			if x > maxAbs {
				maxAbs = x
			}
		} else if x < 0 {
			v := -x
			neg = append(neg, v)
			if v > maxAbs {
				maxAbs = v
			}
		}
	}
	sort.Slice(pos, func(i, j int) bool { return pos[i] > pos[j] })
	sort.Slice(neg, func(i, j int) bool { return neg[i] > neg[j] })
	total := 0
	for i := 0; i < len(pos); i += k {
		total += pos[i] * 2
	}
	for i := 0; i < len(neg); i += k {
		total += neg[i] * 2
	}
	total -= maxAbs
	return total
}

func runBinary(bin, input string) (string, error) {
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

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(10) + 1
	k := rng.Intn(n) + 1
	xs := make([]int, n)
	for i := 0; i < n; i++ {
		xs[i] = rng.Intn(41) - 20
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, k)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(xs[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), expectedC(n, k, xs)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output %q\ninput:\n%s", i+1, out, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
