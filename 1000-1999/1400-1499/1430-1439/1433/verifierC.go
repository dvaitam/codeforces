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
	n := rng.Intn(20) + 2
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(10) + 1
	}
	return Case{n: n, arr: arr}
}

func dominantIndices(c Case) []int {
	mx := c.arr[0]
	for _, v := range c.arr {
		if v > mx {
			mx = v
		}
	}
	res := []int{}
	for i, v := range c.arr {
		if v == mx {
			if (i > 0 && c.arr[i-1] < v) || (i+1 < len(c.arr) && c.arr[i+1] < v) {
				res = append(res, i+1)
			}
		}
	}
	return res
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
	choices := dominantIndices(c)
	if len(choices) == 0 {
		if got != -1 {
			return fmt.Errorf("expected -1 got %d", got)
		}
		return nil
	}
	for _, idx := range choices {
		if got == idx {
			return nil
		}
	}
	return fmt.Errorf("output %d not one of valid indices %v", got, choices)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
