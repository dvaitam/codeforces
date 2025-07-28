package main

import (
	"bytes"
	"fmt"
	"math"
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

func generateCase(rng *rand.Rand) int {
	return rng.Intn(50) + 3
}

func runCase(bin string, n int) error {
	input := fmt.Sprintf("1\n%d\n", n)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	m, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if m > n+5 {
		return fmt.Errorf("too many operations")
	}
	if len(fields) != 1+2*m {
		return fmt.Errorf("expected %d numbers got %d", 1+2*m, len(fields))
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = i
	}
	idx := 1
	for i := 0; i < m; i++ {
		x, err1 := strconv.Atoi(fields[idx])
		y, err2 := strconv.Atoi(fields[idx+1])
		idx += 2
		if err1 != nil || err2 != nil {
			return fmt.Errorf("bad indices")
		}
		if x < 1 || x > n || y < 1 || y > n || x == y {
			return fmt.Errorf("invalid indices")
		}
		ax := a[x]
		ay := a[y]
		a[x] = int(math.Ceil(float64(ax) / float64(ay)))
	}
	ones := 0
	twos := 0
	for i := 1; i <= n; i++ {
		if a[i] == 1 {
			ones++
		} else if a[i] == 2 {
			twos++
		}
	}
	if ones != n-1 || twos != 1 {
		return fmt.Errorf("final array incorrect")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := generateCase(rng)
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
