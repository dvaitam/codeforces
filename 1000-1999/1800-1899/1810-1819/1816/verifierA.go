package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runExe(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	a := rng.Int63n(1_000_000_000) + 1
	b := rng.Int63n(1_000_000_000) + 1
	return fmt.Sprintf("1\n%d %d\n", a, b)
}

func gcd64(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func verifyOutput(input, output string) error {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return fmt.Errorf("failed to read t from input: %v", err)
	}
	scan := bufio.NewScanner(strings.NewReader(output))
	scan.Split(bufio.ScanWords)
	for caseIdx := 1; caseIdx <= t; caseIdx++ {
		var a64, b64 int64
		if _, err := fmt.Fscan(in, &a64, &b64); err != nil {
			return fmt.Errorf("case %d: failed to read a b from input: %v", caseIdx, err)
		}
		if !scan.Scan() {
			return fmt.Errorf("case %d: missing n", caseIdx)
		}
		var n int
		if _, err := fmt.Sscan(scan.Text(), &n); err != nil {
			return fmt.Errorf("case %d: bad n", caseIdx)
		}
		if n < 1 || n > 2 {
			return fmt.Errorf("case %d: n must be 1 or 2, got %d", caseIdx, n)
		}
		prevX, prevY := int64(0), int64(0)
		var x, y int64
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				return fmt.Errorf("case %d: missing x at jump %d", caseIdx, i+1)
			}
			if _, err := fmt.Sscan(scan.Text(), &x); err != nil {
				return fmt.Errorf("case %d: bad x at jump %d", caseIdx, i+1)
			}
			if !scan.Scan() {
				return fmt.Errorf("case %d: missing y at jump %d", caseIdx, i+1)
			}
			if _, err := fmt.Sscan(scan.Text(), &y); err != nil {
				return fmt.Errorf("case %d: bad y at jump %d", caseIdx, i+1)
			}
			if x < 0 || y < 0 || x > 1_000_000_000 || y > 1_000_000_000 {
				return fmt.Errorf("case %d: coordinates out of range at jump %d: %d %d", caseIdx, i+1, x, y)
			}
			dx := x - prevX
			dy := y - prevY
			g := gcd64(dx, dy)
			if g != 1 {
				return fmt.Errorf("case %d: jump %d not primitive: gcd(|%d|,|%d|)=%d", caseIdx, i+1, dx, dy, g)
			}
			prevX, prevY = x, y
		}
		if x != a64 || y != b64 {
			return fmt.Errorf("case %d: last point must be a b, expected %d %d got %d %d", caseIdx, a64, b64, x, y)
		}
	}
	if scan.Scan() {
		return fmt.Errorf("extra output: %s", scan.Text())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		out, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := verifyOutput(input, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
