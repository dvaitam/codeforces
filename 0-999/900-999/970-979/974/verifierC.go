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

type pair struct{ x, y int }

type triple struct{ a, b, c int }

func countPairs(arr []int) int64 {
	n := len(arr)
	if n < 3 {
		return 0
	}
	triples := make([]triple, n-2)
	for i := 0; i < n-2; i++ {
		triples[i] = triple{arr[i], arr[i+1], arr[i+2]}
	}
	var ans int64
	ab := make(map[pair]int)
	bc := make(map[pair]int)
	ac := make(map[pair]int)
	abc := make(map[triple]int)
	for _, t := range triples {
		k1 := pair{t.a, t.b}
		k2 := pair{t.b, t.c}
		k3 := pair{t.a, t.c}
		ans += int64(ab[k1] - abc[t])
		ans += int64(bc[k2] - abc[t])
		ans += int64(ac[k3] - abc[t])
		ab[k1]++
		bc[k2]++
		ac[k3]++
		abc[t]++
	}
	return ans
}

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

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append(os.Args[:1], os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 3
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(5)
		}
		input := fmt.Sprintf("1\n%d\n", n)
		for j := 0; j < n; j++ {
			input += fmt.Sprintf("%d ", arr[j])
		}
		input += "\n"
		expected := fmt.Sprintf("%d", countPairs(arr))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
