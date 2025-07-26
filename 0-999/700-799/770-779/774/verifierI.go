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

func match(s string, pos int, t string) int {
	j := pos
	for i := 0; i < len(t) && j < len(s); i++ {
		if t[i] == s[j] {
			j++
		}
	}
	return j
}

func expected(arr []string, s string) int {
	m := len(s)
	dist := make([]int, m+1)
	for i := range dist {
		dist[i] = -1
	}
	queue := make([]int, 0, m+1)
	dist[0] = 0
	queue = append(queue, 0)
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]
		if pos == m {
			return dist[pos]
		}
		for _, t := range arr {
			np := match(s, pos, t)
			if np > pos && dist[np] == -1 {
				dist[np] = dist[pos] + 1
				queue = append(queue, np)
			}
		}
	}
	return -1
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

func genCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(5) + 1
	if rng.Float64() < 0.1 {
		n = rng.Intn(10) + 1
	}
	arr := make([]string, n)
	for i := range arr {
		l := rng.Intn(4) + 1
		b := make([]byte, l)
		for j := 0; j < l; j++ {
			b[j] = byte('a' + rng.Intn(3))
		}
		arr[i] = string(b)
	}
	slen := rng.Intn(10) + 1
	if rng.Float64() < 0.2 {
		slen = rng.Intn(25) + 1
	}
	sb := make([]byte, slen)
	for i := 0; i < slen; i++ {
		sb[i] = byte('a' + rng.Intn(3))
	}
	s := string(sb)
	var bld strings.Builder
	fmt.Fprintf(&bld, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&bld, "%s\n", arr[i])
	}
	fmt.Fprintf(&bld, "%s\n", s)
	return bld.String(), expected(arr, s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
