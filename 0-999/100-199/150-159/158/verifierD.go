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

func expectedD(vals []int64) int64 {
	n := len(vals)
	ans := int64(-1 << 60)
	for k := 3; k <= n; k++ {
		if n%k != 0 {
			continue
		}
		step := n / k
		for off := 0; off < step; off++ {
			var sum int64
			for j := 0; j < k; j++ {
				sum += vals[off+j*step]
			}
			if sum > ans {
				ans = sum
			}
		}
	}
	return ans
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(30) + 3
	vals := make([]int64, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		vals[i] = int64(rng.Intn(2001) - 1000)
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", vals[i])
	}
	sb.WriteByte('\n')
	return sb.String(), expectedD(vals)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input, exp := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", t+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != fmt.Sprint(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", t+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
