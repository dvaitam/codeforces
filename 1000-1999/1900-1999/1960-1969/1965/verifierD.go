package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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

func genCase(rng *rand.Rand) (string, int, []int64) {
	n := rng.Intn(3) + 3 // 3..5
	arr := make([]int64, n)
	for i := 0; i < (n+1)/2; i++ {
		v := int64(rng.Intn(5) + 1)
		arr[i] = v
		arr[n-1-i] = v
	}
	var sums []int64
	for i := 0; i < n; i++ {
		s := int64(0)
		for j := i; j < n; j++ {
			s += arr[j]
			sums = append(sums, s)
		}
	}
	idx := rng.Intn(len(sums))
	given := append([]int64(nil), sums[:idx]...)
	given = append(given, sums[idx+1:]...)

	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range given {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), n, given
}

func parseOutput(out string, n int) ([]int64, error) {
	rdr := strings.NewReader(out)
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(rdr, &arr[i]); err != nil {
			return nil, fmt.Errorf("failed to read value %d: %v", i+1, err)
		}
		if arr[i] <= 0 {
			return nil, fmt.Errorf("values must be positive")
		}
	}
	return arr, nil
}

func verify(n int, given []int64, arr []int64) error {
	for i := 0; i < n; i++ {
		if arr[i] != arr[n-1-i] {
			return fmt.Errorf("array not palindrome")
		}
	}
	var sums []int64
	for i := 0; i < n; i++ {
		s := int64(0)
		for j := i; j < n; j++ {
			s += arr[j]
			sums = append(sums, s)
		}
	}
	sort.Slice(sums, func(i, j int) bool { return sums[i] < sums[j] })
	sort.Slice(given, func(i, j int) bool { return given[i] < given[j] })
	i, j := 0, 0
	diff := 0
	for i < len(sums) && j < len(given) {
		if sums[i] == given[j] {
			i++
			j++
		} else {
			diff++
			i++
		}
	}
	diff += len(sums) - i
	diff += len(given) - j
	if diff != 1 {
		return fmt.Errorf("subarray sums mismatch")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, n, given := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		arr, err := parseOutput(out, n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\noutput:\n%s\ninput:\n%s", i+1, err, out, in)
			os.Exit(1)
		}
		if err := verify(n, given, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\noutput:\n%s\ninput:\n%s", i+1, err, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
