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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveB(arr []int, x int) int {
	const maxV = 100000
	input := make([]bool, maxV+1)
	andArr := make([]bool, maxV+1)
	ans := 3
	for _, y := range arr {
		if y <= maxV && input[y] {
			ans = min(ans, 0)
		}
		if y <= maxV && andArr[y] {
			ans = min(ans, 1)
		}
		t := y & x
		if t <= maxV && input[t] {
			ans = min(ans, 1)
		}
		if t <= maxV && andArr[t] {
			ans = min(ans, 2)
		}
		if y <= maxV {
			input[y] = true
		}
		if t <= maxV {
			andArr[t] = true
		}
	}
	if ans == 3 {
		return -1
	}
	return ans
}

func generateCaseB(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 2
	x := rng.Intn(100000) + 1
	arr := make([]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, x))
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(100000) + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	expect := solveB(arr, x)
	return sb.String(), fmt.Sprint(expect)
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

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCaseB(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
