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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func minimalSegment(a []int, k int) (int, int) {
	n := len(a)
	count := make(map[int]int)
	distinct := 0
	for _, v := range a {
		if count[v] == 0 {
			distinct++
		}
		count[v]++
	}
	if distinct < k {
		return -1, -1
	}
	l := 0
	for ; l < n; l++ {
		v := a[l]
		count[v]--
		if count[v] == 0 {
			distinct--
			if distinct < k {
				count[v]++
				distinct++
				break
			}
		}
	}
	r := n - 1
	for ; r >= 0; r-- {
		v := a[r]
		count[v]--
		if count[v] == 0 {
			distinct--
			if distinct < k {
				count[v]++
				distinct++
				break
			}
		}
	}
	return l + 1, r + 1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(6) + 1
		}
		k := rng.Intn(n) + 1
		l, r := minimalSegment(a, k)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for j, v := range a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := fmt.Sprintf("%d %d", l, r)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
