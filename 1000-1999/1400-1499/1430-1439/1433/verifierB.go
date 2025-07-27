package main

import (
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
	return strings.TrimSpace(out.String()), nil
}

type Case struct {
	n   int
	arr []int
}

func genCase(rng *rand.Rand) Case {
	n := rng.Intn(50) + 1
	arr := make([]int, n)
	ones := 0
	for i := range arr {
		if rng.Intn(2) == 1 {
			arr[i] = 1
			ones++
		}
	}
	if ones == 0 {
		idx := rng.Intn(n)
		arr[idx] = 1
	}
	return Case{n: n, arr: arr}
}

func expected(c Case) int {
	first := -1
	last := -1
	for i, v := range c.arr {
		if v == 1 {
			if first == -1 {
				first = i
			}
			last = i
		}
	}
	zeros := 0
	for i := first; i <= last; i++ {
		if c.arr[i] == 0 {
			zeros++
		}
	}
	return zeros
}

func runCase(bin string, c Case) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", c.n))
	for i, v := range c.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	out, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	got, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	want := expected(c)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		c := genCase(rng)
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
