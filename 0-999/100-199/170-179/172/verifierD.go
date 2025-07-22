package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(a, n int) uint64 {
	N := a + n - 1
	rem := make([]int32, N+1)
	for i := 1; i <= N; i++ {
		rem[i] = int32(i)
	}
	lim := int(math.Sqrt(float64(N)))
	for p := 2; p <= lim; p++ {
		sq := p * p
		for j := sq; j <= N; j += sq {
			for rem[j]%int32(sq) == 0 {
				rem[j] /= int32(sq)
			}
		}
	}
	var total uint64
	for x := a; x <= N; x++ {
		total += uint64(rem[x])
	}
	return total
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "test %d malformed\n", idx)
			os.Exit(1)
		}
		a, _ := strconv.Atoi(parts[0])
		n, _ := strconv.Atoi(parts[1])
		want := expected(a, n)
		input := fmt.Sprintf("%d %d\n", a, n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		g, _ := strconv.ParseUint(strings.TrimSpace(got), 10, 64)
		if g != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, want, g)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
