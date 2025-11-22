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

type Test struct{ n int64 }

func (t Test) Input() string {
	return fmt.Sprintf("1\n%d\n", t.n)
}

func runProg(bin, input string) (string, error) {
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

func genTest(rng *rand.Rand) Test {
	n := rng.Int63n(1_000_000_000-2) + 2
	return Test{n}
}

// checkOutput validates a contestant's output for a single test case.
// The solution prints k followed by k non-negative integers d[i].
// Let v[0] = 1 and v[i] = v[i-1] + d[i]. The construction is valid when
// each v[i] <= 2 * v[i-1] and the total sum of all v values equals n.
func checkOutput(n int64, out string) error {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}

	k, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil || k < 0 {
		return fmt.Errorf("invalid k value: %s", fields[0])
	}

	if int(k)+1 != len(fields) {
		return fmt.Errorf("expected %d numbers, got %d", k+1, len(fields))
	}

	last := int64(1)
	total := last
	for i := 0; i < int(k); i++ {
		d, err := strconv.ParseInt(fields[i+1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid number at position %d", i+1)
		}
		if d < 0 {
			return fmt.Errorf("negative increment at position %d", i+1)
		}

		cur := last + d
		if cur > last*2 {
			return fmt.Errorf("value %d exceeds twice previous value %d", cur, last)
		}
		total += cur
		last = cur
	}

	if total != n {
		return fmt.Errorf("sum of values %d does not equal %d", total, n)
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		tc := genTest(rng)
		got, err := runProg(bin, tc.Input())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := checkOutput(tc.n, strings.TrimSpace(got)); err != nil {
			fmt.Printf("case %d failed\ninput:\n%serror: %v\noutput: %s\n", i+1, tc.Input(), err, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
