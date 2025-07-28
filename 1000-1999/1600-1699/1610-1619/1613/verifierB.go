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

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(9) + 2 // 2..10
	perm := rng.Perm(1000)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = perm[i] + 1
	}
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	input := fmt.Sprintf("1\n%d\n%s\n", n, sb.String())
	return input, a
}

func checkCase(a []int, out string) error {
	pairs := len(a) / 2
	scan := bufio.NewScanner(strings.NewReader(out))
	scan.Split(bufio.ScanWords)
	arrSet := make(map[int]bool)
	for _, v := range a {
		arrSet[v] = true
	}
	used := make(map[string]bool)
	for i := 0; i < pairs; i++ {
		if !scan.Scan() {
			return fmt.Errorf("not enough numbers")
		}
		x, err := strconv.Atoi(scan.Text())
		if err != nil {
			return fmt.Errorf("invalid int")
		}
		if !scan.Scan() {
			return fmt.Errorf("not enough numbers")
		}
		y, err := strconv.Atoi(scan.Text())
		if err != nil {
			return fmt.Errorf("invalid int")
		}
		if x == y {
			return fmt.Errorf("x == y")
		}
		if !arrSet[x] || !arrSet[y] {
			return fmt.Errorf("numbers not in array")
		}
		mod := x % y
		if arrSet[mod] {
			return fmt.Errorf("x mod y present")
		}
		key := fmt.Sprintf("%d,%d", x, y)
		if used[key] {
			return fmt.Errorf("duplicate pair")
		}
		used[key] = true
	}
	if scan.Scan() {
		return fmt.Errorf("extra output")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, arr := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if err := checkCase(arr, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
