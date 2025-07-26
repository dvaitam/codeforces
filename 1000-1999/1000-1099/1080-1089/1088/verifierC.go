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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func genCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(10) + 1
	a := make([]int, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(100000)
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteByte('\n')
	return sb.String(), a
}

func applyOps(a []int, ops []string) error {
	for idx, line := range ops {
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return fmt.Errorf("operation %d: expected 3 numbers", idx+1)
		}
		t, err1 := strconv.Atoi(parts[0])
		i, err2 := strconv.Atoi(parts[1])
		x, err3 := strconv.Atoi(parts[2])
		if err1 != nil || err2 != nil || err3 != nil {
			return fmt.Errorf("operation %d: invalid integer", idx+1)
		}
		if i < 1 || i > len(a) {
			return fmt.Errorf("operation %d: index out of range", idx+1)
		}
		if t == 1 {
			if x < 0 || x > 1000000 {
				return fmt.Errorf("operation %d: x out of range", idx+1)
			}
			for j := 0; j < i; j++ {
				a[j] += x
			}
		} else if t == 2 {
			if x <= 0 || x > 1000000 {
				return fmt.Errorf("operation %d: x out of range", idx+1)
			}
			for j := 0; j < i; j++ {
				a[j] %= x
			}
		} else {
			return fmt.Errorf("operation %d: invalid type", idx+1)
		}
	}
	for j := 1; j < len(a); j++ {
		if a[j-1] >= a[j] {
			return fmt.Errorf("final array not strictly increasing")
		}
	}
	return nil
}

func check(input, output string) error {
	scan := bufio.NewScanner(strings.NewReader(input))
	scan.Split(bufio.ScanWords)
	scan.Scan()
	n, _ := strconv.Atoi(scan.Text())
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		scan.Scan()
		arr[i], _ = strconv.Atoi(scan.Text())
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	m, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("invalid operations count")
	}
	if m > n+1 || m < 0 {
		return fmt.Errorf("invalid number of operations")
	}
	if len(lines)-1 != m {
		return fmt.Errorf("expected %d operation lines", m)
	}
	ops := lines[1:]
	if err := applyOps(arr, ops); err != nil {
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
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		in, _ := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, in)
			os.Exit(1)
		}
		if err := check(in, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", t+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
