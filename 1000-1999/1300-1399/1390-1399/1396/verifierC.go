package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func solveC(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	var r1, r2, r3, d int64
	if _, err := fmt.Fscan(reader, &n, &r1, &r2, &r3, &d); err != nil {
		return ""
	}
	const INF int64 = math.MaxInt64 / 4
	dp0 := -d
	dp1 := INF
	for i := 0; i < n; i++ {
		var a int64
		fmt.Fscan(reader, &a)
		t1 := a*r1 + r3
		t2 := min(r2, (a+1)*r1)
		ndp0 := min(dp0+d+t1, dp1+r3+d+t1)
		ndp1 := min(dp0+d+t2, dp1+d+t2)
		dp0, dp1 = ndp0, ndp1
	}
	ans := min(dp0, dp1+r3)
	return fmt.Sprint(ans)
}

func genTest(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	r1 := rng.Intn(5) + 1
	r2 := rng.Intn(5) + 1
	r3 := rng.Intn(5) + 1
	d := rng.Intn(5) + 1
	var buf strings.Builder
	fmt.Fprintf(&buf, "%d %d %d %d %d\n", n, r1, r2, r3, d)
	for i := 0; i < n; i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%d", rng.Intn(5)+1)
	}
	buf.WriteByte('\n')
	return buf.String()
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
	for i := 0; i < 100; i++ {
		input := genTest(rng)
		expect := solveC(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
