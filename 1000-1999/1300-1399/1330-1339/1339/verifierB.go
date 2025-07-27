package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func genCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(20) + 3
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = int(rng.Int63n(2_000_000_001) - 1_000_000_000)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	return sb.String(), arr
}

func checkPermutation(orig, out []int) bool {
	if len(orig) != len(out) {
		return false
	}
	m := make(map[int]int)
	for _, v := range orig {
		m[v]++
	}
	for _, v := range out {
		m[v]--
		if m[v] < 0 {
			return false
		}
	}
	return true
}

func isValid(arr []int) bool {
	if len(arr) < 2 {
		return true
	}
	prev := abs(arr[0] - arr[1])
	for i := 1; i < len(arr)-1; i++ {
		cur := abs(arr[i] - arr[i+1])
		if prev > cur {
			return false
		}
		prev = cur
	}
	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func runCase(exe string, input string, orig []int) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(buf.String()))
	scanner.Split(bufio.ScanWords)
	var outArr []int
	for scanner.Scan() {
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("invalid integer in output: %v", err)
		}
		outArr = append(outArr, v)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %v", err)
	}
	if len(outArr) != len(orig) {
		return fmt.Errorf("expected %d numbers got %d", len(orig), len(outArr))
	}
	if !checkPermutation(orig, outArr) {
		return fmt.Errorf("output is not a permutation of input")
	}
	if !isValid(outArr) {
		return fmt.Errorf("adjacent differences not nondecreasing")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, arr := genCase(rng)
		if err := runCase(exe, in, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
