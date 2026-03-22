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

func run(bin string, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func uniqueRandInts(rng *rand.Rand, n int, low, high int64) []int64 {
	m := make(map[int64]struct{})
	res := make([]int64, 0, n)
	for len(res) < n {
		v := rng.Int63n(high-low+1) + low
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		res = append(res, v)
	}
	return res
}

func buildRef() (string, error) {
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1415F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runCase(bin, ref string, input string) error {
	expect, err := run(ref, input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := run(bin, input)
	if err != nil {
		return err
	}
	if expect != got {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	total := 0
	// edge case
	{
		var sb strings.Builder
		sb.WriteString("1\n1 0\n")
		if err := runCase(bin, ref, sb.String()); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total+1, err)
			os.Exit(1)
		}
		total++
	}
	for total < 100 {
		n := rng.Intn(6) + 1
		tVals := make([]int64, n+1)
		xVals := make([]int64, n+1)
		curT := int64(0)
		for i := 1; i <= n; i++ {
			curT += int64(rng.Intn(5) + 1)
			tVals[i] = curT
		}
		coords := uniqueRandInts(rng, n, -10, 10)
		for i := 1; i <= n; i++ {
			xVals[i] = coords[i-1]
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 1; i <= n; i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", tVals[i], xVals[i]))
		}
		if err := runCase(bin, ref, sb.String()); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total+1, err)
			os.Exit(1)
		}
		total++
	}
	fmt.Printf("All %d tests passed\n", total)
}
