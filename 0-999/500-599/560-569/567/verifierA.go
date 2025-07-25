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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String() + stderr.String(), err
	}
	return out.String(), nil
}

type Test struct {
	n     int
	xs    []int64
	input string
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(20) + 2 // 2..21
	xs := make([]int64, n)
	cur := int64(rng.Intn(10) - 5)
	for i := 0; i < n; i++ {
		cur += int64(rng.Intn(10) + 1) // ensure strictly increasing
		xs[i] = cur
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(xs[i], 10))
	}
	sb.WriteByte('\n')
	return Test{n: n, xs: xs, input: sb.String()}
}

func solve(t Test) string {
	n := t.n
	xs := t.xs
	var sb strings.Builder
	for i := 0; i < n; i++ {
		var minDist, maxDist int64
		if i == 0 {
			minDist = xs[1] - xs[0]
			maxDist = xs[n-1] - xs[0]
		} else if i == n-1 {
			minDist = xs[n-1] - xs[n-2]
			maxDist = xs[n-1] - xs[0]
		} else {
			left := xs[i] - xs[i-1]
			right := xs[i+1] - xs[i]
			if left < right {
				minDist = left
			} else {
				minDist = right
			}
			distFirst := xs[i] - xs[0]
			distLast := xs[n-1] - xs[i]
			if distFirst > distLast {
				maxDist = distFirst
			} else {
				maxDist = distLast
			}
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", minDist, maxDist))
	}
	return strings.TrimSpace(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		t := genTest(rng)
		expected := solve(t)
		out, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\noutput:\n%s", i+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, t.input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("ok 100 tests")
}
