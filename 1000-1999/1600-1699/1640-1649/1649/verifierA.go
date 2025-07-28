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

func solveCase(arr []int) int {
	first := -1
	last := -1
	for i, v := range arr {
		if v == 0 {
			if first == -1 {
				first = i
			}
			last = i
		}
	}
	if first == -1 {
		return 0
	}
	return last - first + 2
}

func buildCase(arr []int) (string, int) {
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

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(98) + 2 // 2..99
	arr := make([]int, n)
	arr[0] = 1
	arr[n-1] = 1
	for i := 1; i < n-1; i++ {
		arr[i] = rng.Intn(2)
	}
	return buildCase(arr)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []struct {
		input string
		want  int
	}

	// some fixed simple cases
	cases = append(cases, func() struct {
		input string
		want  int
	} {
		arr := []int{1, 1}
		in, w := buildCase(arr)
		return struct {
			input string
			want  int
		}{in, w}
	}())
	cases = append(cases, func() struct {
		input string
		want  int
	} {
		arr := []int{1, 0, 1}
		in, w := buildCase(arr)
		return struct {
			input string
			want  int
		}{in, w}
	}())
	cases = append(cases, func() struct {
		input string
		want  int
	} {
		arr := []int{1, 0, 0, 1}
		in, w := buildCase(arr)
		return struct {
			input string
			want  int
		}{in, w}
	}())
	cases = append(cases, func() struct {
		input string
		want  int
	} {
		arr := []int{1, 1, 0, 1}
		in, w := buildCase(arr)
		return struct {
			input string
			want  int
		}{in, w}
	}())

	for i := 0; i < 100; i++ {
		in, w := generateCase(rng)
		cases = append(cases, struct {
			input string
			want  int
		}{in, w})
	}

	for idx, tc := range cases {
		got, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		var val int
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
