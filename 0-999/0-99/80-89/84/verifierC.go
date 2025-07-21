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

type circle struct {
	x int
	r int
}

type shot struct {
	x int
	y int
}

func expected(n int, circles []circle, shots []shot) (int, []int) {
	ans := make([]int, n)
	for i := range ans {
		ans[i] = -1
	}
	count := 0
	for idx, s := range shots {
		for j, c := range circles {
			dx := s.x - c.x
			dy := s.y
			if dx*dx+dy*dy <= c.r*c.r {
				if ans[j] == -1 {
					ans[j] = idx + 1
					count++
				}
			}
		}
	}
	return count, ans
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

func formatOutput(count int, ans []int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", count))
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return strings.TrimSpace(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(5) + 1
		circles := make([]circle, n)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			x := rng.Intn(21) - 10
			r := rng.Intn(5) + 1
			circles[i] = circle{x: x, r: r}
			input.WriteString(fmt.Sprintf("%d %d\n", x, r))
		}
		m := rng.Intn(5) + 1
		input.WriteString(fmt.Sprintf("%d\n", m))
		shots := make([]shot, m)
		for i := 0; i < m; i++ {
			sx := rng.Intn(21) - 10
			sy := rng.Intn(21) - 10
			shots[i] = shot{x: sx, y: sy}
			input.WriteString(fmt.Sprintf("%d %d\n", sx, sy))
		}
		count, ans := expected(n, circles, shots)
		expectedOut := formatOutput(count, ans)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input.String())
			os.Exit(1)
		}
		if got != expectedOut {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\ngot:\n%s\ninput:\n%s", t+1, expectedOut, got, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
