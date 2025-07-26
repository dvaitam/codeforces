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

type vec struct{ x, y int }

type caseC struct {
	vs []vec
}

func genCaseC(rng *rand.Rand) caseC {
	n := rng.Intn(100) + 1
	vs := make([]vec, n)
	for i := 0; i < n; i++ {
		vs[i] = vec{rng.Intn(2000001) - 1000000, rng.Intn(2000001) - 1000000}
	}
	return caseC{vs: vs}
}

func runCaseC(bin string, tc caseC) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tc.vs))
	for _, v := range tc.vs {
		fmt.Fprintf(&sb, "%d %d\n", v.x, v.y)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != len(tc.vs) {
		return fmt.Errorf("expected %d numbers got %d", len(tc.vs), len(fields))
	}
	sx, sy := int64(0), int64(0)
	for i, f := range fields {
		if f != "1" && f != "-1" {
			return fmt.Errorf("invalid sign %q", f)
		}
		sign := int64(1)
		if f == "-1" {
			sign = -1
		}
		sx += sign * int64(tc.vs[i].x)
		sy += sign * int64(tc.vs[i].y)
	}
	lim := int64(1500000)
	if sx*sx+sy*sy > lim*lim {
		return fmt.Errorf("vector length too large: %d", sx*sx+sy*sy)
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
		tc := genCaseC(rng)
		if err := runCaseC(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
