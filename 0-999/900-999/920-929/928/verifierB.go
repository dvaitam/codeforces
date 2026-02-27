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

// bruteForce computes the answer for all starting messages by simulation.
// a is 1-indexed (a[0] unused); returns answers for messages 1..n.
func bruteForce(n, k int, a []int) []int {
	result := make([]int, n+1)
	for start := 1; start <= n; start++ {
		seen := make([]bool, n+1)
		cur := start
		for cur != 0 {
			lo := cur - k
			if lo < 1 {
				lo = 1
			}
			hi := cur + k
			if hi > n {
				hi = n
			}
			for j := lo; j <= hi; j++ {
				seen[j] = true
			}
			cur = a[cur]
		}
		count := 0
		for j := 1; j <= n; j++ {
			if seen[j] {
				count++
			}
		}
		result[start] = count
	}
	return result[1:]
}

func runProg(bin, input string) (string, error) {
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

func genCase(rng *rand.Rand) (string, int, int, []int) {
	n := rng.Intn(20) + 1
	k := rng.Intn(n + 1)
	a := make([]int, n+1) // 1-indexed
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		v := rng.Intn(i) // 0..i-1, satisfies a[i] < i
		a[i] = v
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String(), n, k, a
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, n, k, a := genCase(rng)
		want := bruteForce(n, k, a)

		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}

		fields := strings.Fields(got)
		if len(fields) != n {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d values, got %d\ninput:\n%s", i+1, n, len(fields), input)
			os.Exit(1)
		}
		wantStrs := make([]string, n)
		for j, v := range want {
			wantStrs[j] = strconv.Itoa(v)
		}
		if strings.Join(fields, " ") != strings.Join(wantStrs, " ") {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s",
				i+1, strings.Join(wantStrs, " "), strings.Join(fields, " "), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
