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

func countDistinct(a []int) int {
	seen := make(map[int]struct{})
	for _, v := range a {
		seen[v] = struct{}{}
	}
	return len(seen)
}

func isMinimalSegment(a []int, k, l, r int) bool {
	n := len(a)
	if l < 1 || r < 1 || l > r || r > n {
		return false
	}
	l--
	r--
	if countDistinct(a[l:r+1]) != k {
		return false
	}
	if l < r && countDistinct(a[l+1:r+1]) == k {
		return false
	}
	if l < r && countDistinct(a[l:r]) == k {
		return false
	}
	return true
}

func existsSegment(a []int, k int) bool {
	n := len(a)
	for l := 1; l <= n; l++ {
		for r := l; r <= n; r++ {
			if isMinimalSegment(a, k, l, r) {
				return true
			}
		}
	}
	return false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(6) + 1
		}
		k := rng.Intn(n) + 1
		has := existsSegment(a, k)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for j, v := range a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		fields := strings.Fields(got)
		if len(fields) != 2 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected two numbers got %q\ninput:\n%s", i+1, got, input)
			os.Exit(1)
		}
		lOut, err1 := strconv.Atoi(fields[0])
		rOut, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: output not integers %q\ninput:\n%s", i+1, got, input)
			os.Exit(1)
		}
		if !has {
			if !(lOut == -1 && rOut == -1) {
				fmt.Fprintf(os.Stderr, "case %d failed: expected -1 -1 got %s\ninput:\n%s", i+1, got, input)
				os.Exit(1)
			}
			continue
		}
		if !isMinimalSegment(a, k, lOut, rOut) {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid segment %s\ninput:\n%s", i+1, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
