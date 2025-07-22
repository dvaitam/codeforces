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

func divisors(x int64) []int64 {
	var d []int64
	for i := int64(1); i*i <= x; i++ {
		if x%i == 0 {
			d = append(d, i)
			if i*i != x {
				d = append(d, x/i)
			}
		}
	}
	return d
}

func expected(arr []int64) int64 {
	n := len(arr)
	half := (n + 1) / 2
	cnt := make(map[int64]int)
	for _, v := range arr {
		for _, d := range divisors(v) {
			cnt[d]++
		}
	}
	best := int64(1)
	for d, c := range cnt {
		if c >= half && d > best {
			best = d
		}
	}
	return best
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(100) + 1)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	exp := fmt.Sprintf("%d", expected(arr))
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
