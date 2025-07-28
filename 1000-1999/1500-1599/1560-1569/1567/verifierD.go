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

func solveCase(s int64, n int) []int64 {
	base := int64(1)
	for base <= s {
		base *= 10
	}
	base /= 10
	res := make([]int64, 0, n)
	rem := s
	curBase := base
	left := n
	for left > 0 {
		if left == 1 {
			res = append(res, rem)
			break
		}
		if rem-curBase >= int64(left-1) {
			res = append(res, curBase)
			rem -= curBase
			left--
		} else {
			curBase /= 10
		}
	}
	return res
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genTest(rng *rand.Rand) (int64, int) {
	s := int64(rng.Intn(1_000_000_000) + 1)
	n := rng.Intn(min(100, int(s))) + 1
	return s, n
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func format(res []int64) string {
	var sb strings.Builder
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s, n := genTest(rng)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", s, n))
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input.String())
			os.Exit(1)
		}
		expected := format(solveCase(s, n))
		if out != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expected, out, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
