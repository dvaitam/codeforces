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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func expected(n int, ink []int64) int {
	L := int64(n) * 7 / gcd(int64(n), 7)
	pos := make([][]int64, n)
	for d := int64(1); d <= L; d++ {
		pen := int((d - 1) % int64(n))
		dow := (d - 1) % 7
		if dow != 6 {
			pos[pen] = append(pos[pen], d)
		}
	}
	bestPen := -1
	bestDay := int64(1<<63 - 1)
	for i := 0; i < n; i++ {
		arr := pos[i]
		cnt := int64(len(arr))
		if cnt == 0 {
			continue
		}
		k := ink[i]
		q := (k - 1) / cnt
		p := arr[(k-1)%cnt]
		day := q*L + p
		if day < bestDay {
			bestDay = day
			bestPen = i + 1
		}
	}
	return bestPen
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
	inks := make([]int64, n)
	for i := 0; i < n; i++ {
		inks[i] = int64(rng.Intn(10) + 1)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", inks[i])
	}
	b.WriteByte('\n')
	return b.String(), expected(n, inks)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
