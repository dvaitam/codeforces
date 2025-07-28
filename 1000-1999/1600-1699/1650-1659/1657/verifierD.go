package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type unit struct {
	c int
	d int64
	h int64
}
type monster struct {
	D int64
	H int64
}

type testCaseD struct {
	n        int
	C        int
	units    []unit
	monsters []monster
}

func genTestsD() []testCaseD {
	rng := rand.New(rand.NewSource(45))
	tests := make([]testCaseD, 100)
	for i := range tests {
		n := rng.Intn(4) + 1
		C := rng.Intn(20) + 5
		units := make([]unit, n)
		for j := 0; j < n; j++ {
			c := rng.Intn(C) + 1
			d := int64(rng.Intn(10) + 1)
			h := int64(rng.Intn(10) + 1)
			units[j] = unit{c, d, h}
		}
		m := rng.Intn(3) + 1
		mons := make([]monster, m)
		for j := 0; j < m; j++ {
			D := int64(rng.Intn(15) + 1)
			H := int64(rng.Intn(30) + 1)
			mons[j] = monster{D, H}
		}
		tests[i] = testCaseD{n, C, units, mons}
	}
	return tests
}

func compileRef() (string, error) {
	exe := filepath.Join(os.TempDir(), "1657D_ref")
	cmd := exec.Command("go", "build", "-o", exe, "1657D.go")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("compile reference failed: %v\n%s", err, out)
	}
	return exe, nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rand.Seed(time.Now().UnixNano())
	tests := genTestsD()
	for i, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.C)
		for _, u := range tc.units {
			fmt.Fprintf(&sb, "%d %d %d\n", u.c, u.d, u.h)
		}
		fmt.Fprintf(&sb, "%d\n", len(tc.monsters))
		for _, m := range tc.monsters {
			fmt.Fprintf(&sb, "%d %d\n", m.D, m.H)
		}
		input := sb.String()
		exp, err := runBinary(ref, input)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("test %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), strings.TrimSpace(got), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
