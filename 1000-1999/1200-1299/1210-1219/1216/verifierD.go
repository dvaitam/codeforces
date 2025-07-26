package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func expected(arr []int64) (int64, int64) {
	maxA := arr[0]
	for _, v := range arr {
		if v > maxA {
			maxA = v
		}
	}
	g := int64(0)
	for _, v := range arr {
		g = gcd(g, maxA-v)
	}
	var y int64
	for _, v := range arr {
		y += (maxA - v) / g
	}
	return y, g
}

func genCase(rng *rand.Rand) (string, []int64) {
	n := rng.Intn(10) + 1
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(100) + 1)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), arr
}

func parseOutput(out string) (int64, int64, error) {
	reader := strings.NewReader(out)
	var y, g int64
	if _, err := fmt.Fscan(reader, &y, &g); err != nil {
		return 0, 0, fmt.Errorf("parse output: %v", err)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err != io.EOF {
		return 0, 0, fmt.Errorf("extra output")
	}
	return y, g, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, arr := genCase(rng)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		y, g, err := parseOutput(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d bad output: %v\ninput:\n%s\noutput:\n%s", i+1, err, input, out)
			os.Exit(1)
		}
		expY, expG := expected(arr)
		if y != expY || g != expG {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d %d got %d %d\ninput:\n%s", i+1, expY, expG, y, g, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
