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

func solveCase(arr []int64) int64 {
	var sum, max int64
	for _, v := range arr {
		sum += v
		if v > max {
			max = v
		}
	}
	if sum == 0 {
		return 0
	}
	diff := max - (sum - max)
	if diff <= 0 {
		return 1
	}
	return diff
}

func buildCase(arr []int64) (string, int64) {
	n := len(arr)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), solveCase(arr)
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(50) + 1
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = rng.Int63n(20)
	}
	return buildCase(arr)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []struct {
		input string
		want  int64
	}

	// fixed simple cases
	cases = append(cases, func() struct {
		input string
		want  int64
	} {
		arr := []int64{0, 0}
		in, w := buildCase(arr)
		return struct {
			input string
			want  int64
		}{in, w}
	}())
	cases = append(cases, func() struct {
		input string
		want  int64
	} {
		arr := []int64{1, 1}
		in, w := buildCase(arr)
		return struct {
			input string
			want  int64
		}{in, w}
	}())
	cases = append(cases, func() struct {
		input string
		want  int64
	} {
		arr := []int64{4, 0, 0}
		in, w := buildCase(arr)
		return struct {
			input string
			want  int64
		}{in, w}
	}())

	for i := 0; i < 100; i++ {
		in, w := generateCase(rng)
		cases = append(cases, struct {
			input string
			want  int64
		}{in, w})
	}

	for idx, tc := range cases {
		got, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		var val int64
		if _, err := fmt.Sscan(got, &val); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: bad output %q\n", idx+1, got)
			os.Exit(1)
		}
		if val != tc.want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", idx+1, tc.want, val, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
