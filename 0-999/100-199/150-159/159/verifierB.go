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

func solveB(n, m int, markers [][2]int, caps [][2]int) string {
	const maxD = 1000
	sumM := make([]int, maxD+1)
	sumC := make([]int, maxD+1)
	md := make([]map[int]int, maxD+1)
	cd := make([]map[int]int, maxD+1)
	for d := 0; d <= maxD; d++ {
		md[d] = make(map[int]int)
		cd[d] = make(map[int]int)
	}
	for _, p := range markers {
		x, y := p[0], p[1]
		sumM[y]++
		md[y][x]++
	}
	for _, c := range caps {
		a, b := c[0], c[1]
		sumC[b]++
		cd[b][a]++
	}
	total := 0
	beautiful := 0
	for d := 0; d <= maxD; d++ {
		if sumM[d] == 0 || sumC[d] == 0 {
			continue
		}
		t := sumM[d]
		if sumC[d] < t {
			t = sumC[d]
		}
		total += t
		b := 0
		for color, cntM := range md[d] {
			cntC := cd[d][color]
			if cntC < cntM {
				b += cntC
			} else {
				b += cntM
			}
		}
		if b > t {
			b = t
		}
		beautiful += b
	}
	return fmt.Sprintf("%d %d", total, beautiful)
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
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		markers := make([][2]int, n)
		caps := make([][2]int, m)
		for j := 0; j < n; j++ {
			markers[j][0] = rng.Intn(5) + 1
			markers[j][1] = rng.Intn(5) + 1
		}
		for j := 0; j < m; j++ {
			caps[j][0] = rng.Intn(5) + 1
			caps[j][1] = rng.Intn(5) + 1
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, p := range markers {
			sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
		for _, c := range caps {
			sb.WriteString(fmt.Sprintf("%d %d\n", c[0], c[1]))
		}
		input := sb.String()
		expected := solveB(n, m, markers, caps)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
