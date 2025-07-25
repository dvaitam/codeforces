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

func solveD(arr []int) string {
	n := len(arr)
	counts := make(map[int]int)
	for _, v := range arr {
		counts[v]++
	}
	ans := 0
	if counts[0] > ans {
		ans = counts[0]
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			a, b := arr[i], arr[j]
			if a == 0 && b == 0 {
				continue
			}
			used := map[int]int{a: 1, b: 1}
			if used[a] > counts[a] || used[b] > counts[b] {
				continue
			}
			length := 2
			x, y := a, b
			for {
				next := x + y
				if counts[next]-used[next] <= 0 {
					break
				}
				used[next]++
				length++
				x, y = y, next
			}
			if length > ans {
				ans = length
			}
		}
	}
	return fmt.Sprintf("%d\n", ans)
}

func generateCaseD(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 2
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(21) - 10
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	expect := solveD(arr)
	return sb.String(), expect
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCaseD(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
