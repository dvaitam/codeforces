package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCaseC struct {
	k float64
	d float64
	t float64
}

func generateC(rng *rand.Rand) testCaseC {
	k := float64(rng.Intn(1000) + 1)
	d := float64(rng.Intn(1000) + 1)
	t := float64(rng.Intn(1000) + 1)
	return testCaseC{k: k, d: d, t: t}
}

func expectedC(k, d, t float64) float64 {
	cycle := math.Ceil(k/d) * d
	totalNeed := 2 * t
	perCycle := k + cycle
	full := math.Floor(totalNeed / perCycle)
	rem := totalNeed - perCycle*full
	var extra float64
	if rem <= 2*k {
		extra = rem / 2
	} else {
		extra = k + (rem - 2*k)
	}
	return full*cycle + extra
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCaseC) error {
	input := fmt.Sprintf("%g %g %g\n", tc.k, tc.d, tc.t)
	expected := expectedC(tc.k, tc.d, tc.t)
	gotStr, err := run(bin, input)
	if err != nil {
		return err
	}
	got, err := strconv.ParseFloat(strings.TrimSpace(gotStr), 64)
	if err != nil {
		return fmt.Errorf("invalid float output: %v", err)
	}
	if math.Abs(got-expected) > 1e-6 {
		return fmt.Errorf("expected %.6f got %s", expected, gotStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateC(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
