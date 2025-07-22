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

func expected(input string) string {
	fields := strings.Fields(input)
	idx := 0
	n := toInt64(fields[idx])
	idx++
	m := toInt64(fields[idx])
	idx++
	xc := toInt64(fields[idx])
	idx++
	yc := toInt64(fields[idx])
	idx++
	k := toInt(fields[idx])
	idx++
	var total int64
	const inf int64 = 1 << 60
	for i := 0; i < k; i++ {
		dx := toInt64(fields[idx])
		idx++
		dy := toInt64(fields[idx])
		idx++
		var t1, t2 int64
		if dx > 0 {
			t1 = (n - xc) / dx
		} else if dx < 0 {
			t1 = (1 - xc) / dx
		} else {
			t1 = inf
		}
		if dy > 0 {
			t2 = (m - yc) / dy
		} else if dy < 0 {
			t2 = (1 - yc) / dy
		} else {
			t2 = inf
		}
		t := t1
		if t2 < t {
			t = t2
		}
		if t < 0 {
			t = 0
		}
		total += t
		xc += t * dx
		yc += t * dy
	}
	return fmt.Sprintf("%d", total)
}

func toInt(s string) int {
	var v int
	fmt.Sscan(s, &v)
	return v
}

func toInt64(s string) int64 {
	var v int64
	fmt.Sscan(s, &v)
	return v
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

func generateCase(rng *rand.Rand) string {
	n := rng.Int63n(1000) + 1
	m := rng.Int63n(1000) + 1
	xc := rng.Int63n(n) + 1
	yc := rng.Int63n(m) + 1
	k := rng.Intn(10) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	fmt.Fprintf(&sb, "%d %d\n", xc, yc)
	fmt.Fprintf(&sb, "%d\n", k)
	for i := 0; i < k; i++ {
		dx := rng.Int63n(21) - 10
		dy := rng.Int63n(21) - 10
		if dx == 0 && dy == 0 {
			dx = 1
		}
		fmt.Fprintf(&sb, "%d %d\n", dx, dy)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		exp := expected(strings.ReplaceAll(tc, "\n", " "))
		got, err := run(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
