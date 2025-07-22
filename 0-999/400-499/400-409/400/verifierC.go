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

func transform(n, m int64, x, y, z int64, points [][2]int64) [][2]int64 {
	r1 := x % 4
	f := y % 2
	r2 := z % 4
	res := make([][2]int64, len(points))
	for idx, p := range points {
		i, j := p[0], p[1]
		h, w := n, m
		for t := int64(0); t < r1; t++ {
			i, j = j, h-i+1
			h, w = w, h
		}
		if f == 1 {
			j = w - j + 1
		}
		for t := int64(0); t < r2; t++ {
			i, j = w-j+1, i
			h, w = w, h
		}
		res[idx] = [2]int64{i, j}
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tcase := 0; tcase < 100; tcase++ {
		n := int64(rng.Intn(10) + 1)
		m := int64(rng.Intn(10) + 1)
		x := int64(rng.Intn(20))
		y := int64(rng.Intn(20))
		z := int64(rng.Intn(20))
		p := rng.Intn(5) + 1
		points := make([][2]int64, p)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d %d %d %d\n", n, m, x, y, z, p))
		for i := 0; i < p; i++ {
			r := int64(rng.Intn(int(n)) + 1)
			c := int64(rng.Intn(int(m)) + 1)
			points[i] = [2]int64{r, c}
			sb.WriteString(fmt.Sprintf("%d %d\n", r, c))
		}
		input := sb.String()
		transformed := transform(n, m, x, y, z, points)
		var exp strings.Builder
		for i := 0; i < p; i++ {
			exp.WriteString(fmt.Sprintf("%d %d\n", transformed[i][0], transformed[i][1]))
		}
		expectedOut := strings.TrimSpace(exp.String())
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tcase+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expectedOut {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", tcase+1, expectedOut, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
