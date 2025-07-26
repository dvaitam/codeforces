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

func expectedC(n, m int, ops [][2]int) string {
	var sum int64
	for i := 0; i < m; i++ {
		x := int64(ops[i][0])
		d := int64(ops[i][1])
		sum += x * int64(n)
		if d >= 0 {
			sum += d * int64(n-1) * int64(n) / 2
		} else {
			mid := n/2 + 1
			front := int64(mid - 1)
			last := int64(n - mid)
			sum += d * (front*(front+1)/2 + last*(last+1)/2)
		}
	}
	ans := float64(sum) / float64(n)
	return fmt.Sprintf("%.15f", ans)
}

func genCaseC(rng *rand.Rand) (int, int, [][2]int) {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	ops := make([][2]int, m)
	for i := 0; i < m; i++ {
		ops[i][0] = rng.Intn(2001) - 1000
		ops[i][1] = rng.Intn(2001) - 1000
	}
	return n, m, ops
}

func runCaseC(bin string, n, m int, ops [][2]int) error {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte(' ')
	sb.WriteString(strconv.Itoa(m))
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(strconv.Itoa(ops[i][0]))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(ops[i][1]))
	}
	sb.WriteByte('\n')
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	exp := expectedC(n, m, ops)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, ops := genCaseC(rng)
		if err := runCaseC(bin, n, m, ops); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
