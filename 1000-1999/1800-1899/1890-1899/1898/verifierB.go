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

func expectedB(a []int64) int64 {
	ops := int64(0)
	cur := a[len(a)-1]
	for i := len(a) - 2; i >= 0; i-- {
		if a[i] <= cur {
			cur = a[i]
		} else {
			k := (a[i] + cur - 1) / cur
			ops += k - 1
			cur = a[i] / k
		}
	}
	return ops
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(30) + 1
	a := make([]int64, n)
	for i := range a {
		a[i] = rng.Int63n(1000) + 1
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	exp := expectedB(a)
	return sb.String(), exp
}

func runCase(bin, input string, exp int64) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != 1 {
		return fmt.Errorf("expected single integer output")
	}
	got, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid integer: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
