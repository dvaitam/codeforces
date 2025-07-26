package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(9) + 2 // 2..10
	seq := make([]int64, n)
	seq[0] = int64(rng.Intn(20) + 1)
	for i := 1; i < n; i++ {
		if rng.Intn(2) == 0 {
			seq[i] = seq[i-1] * 2
		} else {
			val := seq[i-1]
			if val%3 == 0 {
				seq[i] = val / 3
			} else {
				seq[i] = val * 2
			}
		}
	}
	// shuffle order for input
	perm := rng.Perm(n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, idx := range perm {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(seq[idx]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseInts(fields []string) ([]int64, error) {
	res := make([]int64, len(fields))
	for i, s := range fields {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func check(input, output string) error {
	ins := strings.Fields(input)
	n, _ := strconv.Atoi(ins[0])
	nums, err := parseInts(ins[1:])
	if err != nil {
		return fmt.Errorf("bad input")
	}
	if len(nums) != n {
		return fmt.Errorf("input mismatch")
	}
	outNums, err := parseInts(strings.Fields(output))
	if err != nil {
		return fmt.Errorf("invalid output")
	}
	if len(outNums) != n {
		return fmt.Errorf("expected %d numbers", n)
	}
	used := make(map[int64]bool)
	for _, v := range nums {
		used[v] = false
	}
	for _, v := range outNums {
		if _, ok := used[v]; !ok {
			return fmt.Errorf("number %d not in input", v)
		}
		if used[v] {
			return fmt.Errorf("number %d repeated", v)
		}
		used[v] = true
	}
	for i := 0; i < n-1; i++ {
		x := outNums[i]
		y := outNums[i+1]
		if x*2 != y && !(x%3 == 0 && x/3 == y) {
			return fmt.Errorf("invalid step from %d to %d", x, y)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(4))
	for i := 0; i < 100; i++ {
		in := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(in, out); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%soutput:%s\n", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
