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

func solveCase(a []int) string {
	n := len(a)
	sum := 0
	for _, v := range a {
		sum ^= v
	}
	var sb strings.Builder
	if n%2 == 1 {
		sb.WriteString("YES\n")
		sb.WriteString(fmt.Sprintf("%d\n", n-2))
		for i := 0; i+2 < n; i += 2 {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", i+1, i+2, i+3))
		}
		for i := n - 3; i-2 >= 0; i -= 2 {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", i+1, i, i-1))
		}
	} else {
		if sum == 0 {
			sb.WriteString("YES\n")
			sb.WriteString(fmt.Sprintf("%d\n", n-3))
			n--
			for i := 0; i+2 < n; i += 2 {
				sb.WriteString(fmt.Sprintf("%d %d %d\n", i+1, i+2, i+3))
			}
			for i := n - 3; i-2 >= 0; i -= 2 {
				sb.WriteString(fmt.Sprintf("%d %d %d\n", i+1, i, i-1))
			}
		} else {
			sb.WriteString("NO")
			return sb.String()
		}
	}
	res := sb.String()
	res = strings.TrimRight(res, "\n")
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 3
	var in strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", n))
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(100) + 1
		in.WriteString(fmt.Sprintf("%d", arr[i]))
		if i+1 < n {
			in.WriteByte(' ')
		}
	}
	in.WriteByte('\n')
	exp := solveCase(append([]int(nil), arr...))
	return in.String(), exp
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\n--- got:\n%s", exp, got)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
