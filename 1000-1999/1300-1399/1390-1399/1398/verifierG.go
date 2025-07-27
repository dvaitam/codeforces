package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solveCase(input string) string {
	sc := bufio.NewScanner(strings.NewReader(input))
	sc.Split(bufio.ScanWords)
	sc.Scan()
	n, _ := strconv.Atoi(sc.Text())
	sc.Scan()
	x, _ := strconv.Atoi(sc.Text())
	_ = x
	sc.Scan()
	y, _ := strconv.Atoi(sc.Text())
	a := make([]int, n+1)
	for i := 0; i <= n; i++ {
		sc.Scan()
		a[i], _ = strconv.Atoi(sc.Text())
	}
	diffs := make(map[int]bool)
	for i := 0; i <= n; i++ {
		for j := 0; j <= n; j++ {
			d := a[i] - a[j]
			if d > 0 {
				diffs[d] = true
			}
		}
	}
	const maxL = 1000000
	best := make([]int, maxL+1)
	for d := range diffs {
		L := 2 * (d + y)
		for m := L; m <= maxL; m += L {
			if best[m] < L {
				best[m] = L
			}
		}
	}
	sc.Scan()
	q, _ := strconv.Atoi(sc.Text())
	out := make([]string, q)
	for i := 0; i < q; i++ {
		sc.Scan()
		li, _ := strconv.Atoi(sc.Text())
		if li <= maxL && best[li] > 0 {
			out[i] = fmt.Sprint(best[li])
		} else {
			out[i] = "-1"
		}
	}
	return strings.Join(out, " ")
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	x := rng.Intn(10) + 1
	y := rng.Intn(5)
	a := make([]int, n+1)
	for i := 0; i <= n; i++ {
		a[i] = rng.Intn(x + 1)
	}
	q := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, x, y))
	for i := 0; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d ", a[i]))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintln(q))
	for i := 0; i < q; i++ {
		sb.WriteString(fmt.Sprintf("%d ", rng.Intn(20)+1))
	}
	sb.WriteByte('\n')
	return sb.String()
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected = strings.TrimSpace(expected)
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		exp := solveCase(in)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
