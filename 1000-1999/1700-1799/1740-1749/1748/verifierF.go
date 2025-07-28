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

func runProg(bin, input string) (string, error) {
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
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(20) + 2
	return fmt.Sprintf("%d\n", n)
}

func verify(n int, out string) error {
	scan := bufio.NewScanner(strings.NewReader(out))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return fmt.Errorf("missing k")
	}
	k, err := strconv.Atoi(scan.Text())
	if err != nil || k < 0 || k > 250000 {
		return fmt.Errorf("invalid k")
	}
	ops := make([]int, k)
	for i := 0; i < k; i++ {
		if !scan.Scan() {
			return fmt.Errorf("missing op %d", i)
		}
		v, err := strconv.Atoi(scan.Text())
		if err != nil || v < 0 || v >= n {
			return fmt.Errorf("bad op value")
		}
		ops[i] = v
	}
	if scan.Scan() {
		return fmt.Errorf("extra output")
	}
	arr := make([]uint64, n)
	for i := 0; i < n; i++ {
		arr[i] = 1 << i
	}
	for _, v := range ops {
		arr[v] ^= arr[(v+1)%n]
	}
	for i := 0; i < n; i++ {
		if arr[i] != 1<<(n-1-i) {
			return fmt.Errorf("result incorrect")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		out, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n%s", i, err, input)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(strings.TrimSpace(input))
		if err := verify(n, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
