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
	n := rng.Intn(8) + 1
	arr := make([]int, n)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(50)
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(arr[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), arr
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else if strings.HasSuffix(bin, ".py") {
		cmd = exec.Command("python3", bin)
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
	return out.String(), nil
}

func checkOutput(arr []int, out string) error {
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(out)))
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	mStr := strings.TrimSpace(scanner.Text())
	m, err := strconv.Atoi(mStr)
	if err != nil {
		return fmt.Errorf("invalid operation count")
	}
	if m < 0 || m > len(arr) {
		return fmt.Errorf("invalid operation count %d", m)
	}
	ops := make([][2]int, m)
	for i := 0; i < m; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("expected %d operations", m)
		}
		var l, r int
		if _, err := fmt.Sscan(scanner.Text(), &l, &r); err != nil {
			return fmt.Errorf("bad operation format")
		}
		if l < 1 || l > len(arr) || r < 1 || r > len(arr) || l >= r {
			return fmt.Errorf("invalid indices")
		}
		ops[i] = [2]int{l - 1, r - 1}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output")
	}
	a := make([]int, len(arr))
	copy(a, arr)
	for _, op := range ops {
		l := op[0]
		r := op[1]
		if (a[l]+a[r])%2 == 1 {
			a[r] = a[l]
		} else {
			a[l] = a[r]
		}
	}
	for i := 1; i < len(a); i++ {
		if a[i] < a[i-1] {
			return fmt.Errorf("array not sorted")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, arr := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := checkOutput(arr, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
