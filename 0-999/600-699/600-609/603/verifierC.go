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

func grundyOdd(x int64) int {
	switch x {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 0
	case 3:
		return 1
	case 4:
		return 2
	}
	if x%2 == 1 {
		return 0
	}
	g := grundyOdd(x / 2)
	if g == 1 {
		return 2
	}
	return 1
}

func grundyEven(x int64) int {
	if x == 1 {
		return 1
	}
	if x == 2 {
		return 2
	}
	if x%2 == 1 {
		return 0
	}
	g := grundyEven(x / 2)
	if g == 1 {
		return 2
	}
	return 1
}

func solveCase(n int, k int64, arr []int64) string {
	res := 0
	for _, a := range arr {
		var g int
		if k%2 == 0 {
			g = grundyEven(a)
		} else {
			g = grundyOdd(a)
		}
		res ^= g
	}
	if res != 0 {
		return "Kevin"
	}
	return "Nicky"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	k := int64(rng.Intn(10) + 1)
	arr := make([]int64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(20) + 1)
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(arr[i]))
	}
	sb.WriteByte('\n')
	expected := solveCase(n, k, arr)
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
