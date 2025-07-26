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
	err := cmd.Run()
	if err != nil {
		return out.String() + errb.String(), err
	}
	return out.String(), nil
}

type Test struct {
	n     int
	U     int
	E     []int
	input string
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(7) + 3
	E := make([]int, n)
	cur := rng.Intn(5)
	for i := 0; i < n; i++ {
		cur += rng.Intn(5) + 1
		E[i] = cur
	}
	U := rng.Intn(E[n-1]-E[0]+1) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, U))
	for i, v := range E {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return Test{n: n, U: U, E: E, input: sb.String()}
}

func solve(t Test) string {
	best := -1.0
	k := 0
	for i := 0; i < t.n; i++ {
		if k < i {
			k = i
		}
		for k+1 < t.n && t.E[k+1]-t.E[i] <= t.U {
			k++
		}
		if k >= i+2 {
			val := float64(t.E[k]-t.E[i+1]) / float64(t.E[k]-t.E[i])
			if val > best {
				best = val
			}
		}
	}
	if best < 0 {
		return "-1"
	}
	return fmt.Sprintf("%.12f", best)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		t := genTest(rng)
		expected := solve(t)
		out, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%s got:%s\n", i+1, t.input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("ok 100 tests")
}
