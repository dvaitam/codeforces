package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func floorDiv(u, m int64) int64 {
	if u >= 0 {
		return u / m
	}
	return (u+1)/m - 1
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func solveB(a, b, x1, y1, x2, y2 int64) string {
	u1 := x1 + y1
	u2 := x2 + y2
	v1 := x1 - y1
	v2 := x2 - y2
	iu1 := floorDiv(u1, 2*a)
	iu2 := floorDiv(u2, 2*a)
	iv1 := floorDiv(v1, 2*b)
	iv2 := floorDiv(v2, 2*b)
	pu := abs(iu1 - iu2)
	pv := abs(iv1 - iv2)
	ans := max(pu, pv)
	return fmt.Sprintf("%d", ans)
}

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

func genCase(rng *rand.Rand) (string, [6]int64) {
	var data [6]int64
	a := int64(rng.Intn(100) + 1)
	b := int64(rng.Intn(100) + 1)
	x1 := int64(rng.Intn(2000) - 1000)
	y1 := int64(rng.Intn(2000) - 1000)
	x2 := int64(rng.Intn(2000) - 1000)
	y2 := int64(rng.Intn(2000) - 1000)
	data[0] = a
	data[1] = b
	data[2] = x1
	data[3] = y1
	data[4] = x2
	data[5] = y2
	input := fmt.Sprintf("%d %d %d %d %d %d\n", a, b, x1, y1, x2, y2)
	return input, data
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, data := genCase(rng)
		expect := solveB(data[0], data[1], data[2], data[3], data[4], data[5])
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
