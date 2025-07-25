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

func solveA(x, y []int) string {
	present := make(map[int]struct{}, len(x)+len(y))
	for _, v := range x {
		present[v] = struct{}{}
	}
	for _, v := range y {
		present[v] = struct{}{}
	}
	count := 0
	for _, a := range x {
		for _, b := range y {
			if _, ok := present[a^b]; ok {
				count++
			}
		}
	}
	if count%2 == 0 {
		return "Karen"
	}
	return "Koyomi"
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		used := make(map[int]bool)
		x := make([]int, n)
		y := make([]int, n)
		for j := 0; j < n; {
			v := rng.Intn(2_000_000) + 1
			if !used[v] {
				used[v] = true
				x[j] = v
				j++
			}
		}
		for j := 0; j < n; {
			v := rng.Intn(2_000_000) + 1
			if !used[v] {
				used[v] = true
				y[j] = v
				j++
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j, v := range x {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for j, v := range y {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := solveA(x, y)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
