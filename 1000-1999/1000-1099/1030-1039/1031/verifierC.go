package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func generateCase(rng *rand.Rand) (string, int, int, int) {
	a := rng.Intn(1000)
	b := rng.Intn(1000)
	total := a + b
	k := int(math.Sqrt(2 * float64(total)))
	for k*(k+1)/2 > total {
		k--
	}
	for (k+1)*(k+2)/2 <= total {
		k++
	}
	input := fmt.Sprintf("%d %d\n", a, b)
	return input, a, b, k
}

func parseOutput(out string) ([]int, []int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 3 {
		return nil, nil, fmt.Errorf("expected at least 3 lines")
	}
	n1, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, nil, fmt.Errorf("invalid n1")
	}
	day1 := []int{}
	if n1 > 0 {
		fields := strings.Fields(lines[1])
		if len(fields) != n1 {
			return nil, nil, fmt.Errorf("expected %d numbers in day1", n1)
		}
		day1 = make([]int, n1)
		for i := 0; i < n1; i++ {
			v, err := strconv.Atoi(fields[i])
			if err != nil {
				return nil, nil, fmt.Errorf("bad int")
			}
			day1[i] = v
		}
	}
	idx := 1
	if n1 > 0 {
		idx = 2
	}
	if len(lines) <= idx {
		return nil, nil, fmt.Errorf("missing third line")
	}
	n2, err := strconv.Atoi(strings.TrimSpace(lines[idx]))
	if err != nil {
		return nil, nil, fmt.Errorf("invalid n2")
	}
	day2 := []int{}
	if n2 > 0 {
		if len(lines) <= idx+1 {
			return nil, nil, fmt.Errorf("missing day2 numbers")
		}
		fields := strings.Fields(lines[idx+1])
		if len(fields) != n2 {
			return nil, nil, fmt.Errorf("expected %d numbers in day2", n2)
		}
		day2 = make([]int, n2)
		for i := 0; i < n2; i++ {
			v, err := strconv.Atoi(fields[i])
			if err != nil {
				return nil, nil, fmt.Errorf("bad int")
			}
			day2[i] = v
		}
	}
	return day1, day2, nil
}

func checkCase(a, b, k int, day1, day2 []int) error {
	used := make(map[int]bool)
	sum1 := 0
	for _, v := range day1 {
		if v < 1 || v > k {
			return fmt.Errorf("invalid number %d", v)
		}
		if used[v] {
			return fmt.Errorf("duplicate %d", v)
		}
		used[v] = true
		sum1 += v
	}
	sum2 := 0
	for _, v := range day2 {
		if v < 1 || v > k {
			return fmt.Errorf("invalid number %d", v)
		}
		if used[v] {
			return fmt.Errorf("duplicate %d", v)
		}
		used[v] = true
		sum2 += v
	}
	if len(used) != k {
		return fmt.Errorf("expected %d total numbers got %d", k, len(used))
	}
	if sum1 > a || sum2 > b {
		return fmt.Errorf("sums exceed limits")
	}
	return nil
}

func runCase(bin, input string, a, b, k int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, errBuf.String())
	}
	day1, day2, err := parseOutput(out.String())
	if err != nil {
		return err
	}
	if err := checkCase(a, b, k, day1, day2); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		input, a, b, k := generateCase(rng)
		if err := runCase(bin, input, a, b, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
