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

func runCandidate(bin, input string) (string, error) {
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

func solve(n, k int, arr []int) int64 {
	sort.Ints(arr)
	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + int64(arr[i])
	}
	var best int64
	for i := 0; i <= k; i++ {
		left := 2 * i
		if left > n {
			break
		}
		right := n - (k - i)
		if right < left {
			continue
		}
		sum := prefix[right] - prefix[left]
		if sum > best {
			best = sum
		}
	}
	return best
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	k := rng.Intn(n/2 + 1)
	arr := make([]int, n)
	used := map[int]bool{}
	for i := 0; i < n; i++ {
		v := rng.Intn(100)
		for used[v] {
			v = rng.Intn(100)
		}
		used[v] = true
		arr[i] = v
	}
	input := fmt.Sprintf("1\n%d %d\n", n, k)
	for i := 0; i < n; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", arr[i])
	}
	input += "\n"
	return input, fmt.Sprintf("%d", solve(n, k, arr))
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
