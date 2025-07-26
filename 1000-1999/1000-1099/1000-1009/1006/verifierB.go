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

func solve(n, k int, arr []int) (int, []int) {
	tmp := make([]int, n)
	copy(tmp, arr)
	sort.Slice(tmp, func(i, j int) bool { return tmp[i] > tmp[j] })
	sum := 0
	maxv := make([]int, k)
	for i := 0; i < k; i++ {
		sum += tmp[i]
		maxv[i] = tmp[i]
	}
	segments := make([]int, 0, k)
	prev := -1
	idx := 0
	for i := 0; i < n && len(maxv) > 0; i++ {
		for j := 0; j < len(maxv); j++ {
			if arr[i] == maxv[j] {
				segments = append(segments, i-prev)
				prev = i
				maxv = append(maxv[:j], maxv[j+1:]...)
				idx = i + 1
				break
			}
		}
	}
	if len(segments) > 0 {
		segments[len(segments)-1] += n - idx
	}
	return sum, segments
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		k := rng.Intn(n) + 1
		arr := make([]int, n)
		input := fmt.Sprintf("%d %d\n", n, k)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(2000) + 1
			if j > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", arr[j])
		}
		input += "\n"
		sum, seg := solve(n, k, arr)
		expected := fmt.Sprintf("%d\n", sum)
		for j, v := range seg {
			if j > 0 {
				expected += " "
			}
			expected += fmt.Sprintf("%d", v)
		}
		expected = strings.TrimSpace(expected)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
