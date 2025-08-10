package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func runProg(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		v := rng.Intn(100) + 1
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseInput(s string) (int, []int, error) {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return 0, nil, fmt.Errorf("empty input")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, nil, err
	}
	if len(fields) != n+1 {
		return 0, nil, fmt.Errorf("expected %d numbers, got %d", n, len(fields)-1)
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return 0, nil, err
		}
		arr[i] = v
	}
	return n, arr, nil
}

func parseOutput(s string, n int) (int, []int, error) {
	fields := strings.Fields(s)
	if len(fields) != n+1 {
		return 0, nil, fmt.Errorf("expected %d numbers, got %d", n+1, len(fields))
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid k: %v", err)
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return 0, nil, fmt.Errorf("invalid number: %v", err)
		}
		arr[i] = v
	}
	return k, arr, nil
}

func checkCase(input, output string) error {
	n, orig, err := parseInput(input)
	if err != nil {
		return fmt.Errorf("parse input: %v", err)
	}
	k, arr, err := parseOutput(output, n)
	if err != nil {
		return fmt.Errorf("parse output: %v", err)
	}
	if k != (n-1)/2 {
		return fmt.Errorf("reported %d beautiful numbers, expected %d", k, (n-1)/2)
	}
	sortedOrig := append([]int(nil), orig...)
	sortedArr := append([]int(nil), arr...)
	sort.Ints(sortedOrig)
	sort.Ints(sortedArr)
	for i := 0; i < n; i++ {
		if sortedOrig[i] != sortedArr[i] {
			return fmt.Errorf("output numbers are not a permutation of input")
		}
	}
	cnt := 0
	for i := 1; i+1 < n; i++ {
		if arr[i] < arr[i-1] && arr[i] < arr[i+1] {
			cnt++
		}
	}
	if cnt != k {
		return fmt.Errorf("reported %d beautiful numbers, found %d", k, cnt)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := checkCase(input, got); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on case %d\ninput:\n%soutput:\n%s\n%v", i+1, input, got, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
