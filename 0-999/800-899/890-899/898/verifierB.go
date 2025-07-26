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

type testCaseB struct {
	n, a, b  int64
	expectOk bool
	expectX  int64
	expectY  int64
}

func solveB(n, a, b int64) (bool, int64, int64) {
	for x := int64(0); x*a <= n; x++ {
		rem := n - x*a
		if rem%b == 0 {
			return true, x, rem / b
		}
	}
	return false, 0, 0
}

func genCaseB(rng *rand.Rand) (string, testCaseB) {
	n := int64(rng.Intn(10000) + 1)
	a := int64(rng.Intn(1000) + 1)
	b := int64(rng.Intn(1000) + 1)
	ok, x, y := solveB(n, a, b)
	tc := testCaseB{n: n, a: a, b: b, expectOk: ok, expectX: x, expectY: y}
	in := fmt.Sprintf("%d\n%d\n%d\n", n, a, b)
	return in, tc
}

func runCaseB(bin string, in string, tc testCaseB) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	out := strings.Fields(strings.TrimSpace(buf.String()))
	if tc.expectOk {
		if len(out) != 3 || strings.ToUpper(out[0]) != "YES" {
			return fmt.Errorf("expected YES and two numbers, got %v", out)
		}
		x, err1 := strconv.ParseInt(out[1], 10, 64)
		y, err2 := strconv.ParseInt(out[2], 10, 64)
		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid numbers in output: %v %v", err1, err2)
		}
		if x < 0 || y < 0 || x*tc.a+y*tc.b != tc.n {
			return fmt.Errorf("output %d %d does not satisfy equation", x, y)
		}
	} else {
		if len(out) != 1 || strings.ToUpper(out[0]) != "NO" {
			return fmt.Errorf("expected NO, got %v", out)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, tc := genCaseB(rng)
		if err := runCaseB(bin, in, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
