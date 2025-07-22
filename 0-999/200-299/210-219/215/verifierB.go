package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func buildRef() (string, error) {
	ref := "./refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "215B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func randDistinct(rng *rand.Rand, count, max int) []int {
	set := make(map[int]bool)
	res := make([]int, 0, count)
	for len(res) < count {
		v := rng.Intn(max) + 1
		if !set[v] {
			set[v] = true
			res = append(res, v)
		}
	}
	sort.Ints(res)
	return res
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	k := rng.Intn(5) + 1
	xs := randDistinct(rng, n, 50)
	ys := randDistinct(rng, m, 50)
	zs := randDistinct(rng, k, 50)
	A := rng.Intn(10) + 1
	B := rng.Intn(10) + 1

	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	for _, v := range xs {
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(strconv.Itoa(m))
	for _, v := range ys {
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(strconv.Itoa(k))
	for _, v := range zs {
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d %d\n", A, B))
	return sb.String()
}

func runCase(exe, ref, input string) error {
	cmdRef := exec.Command(ref)
	cmdRef.Stdin = strings.NewReader(input)
	var refOut bytes.Buffer
	cmdRef.Stdout = &refOut
	cmdRef.Stderr = &refOut
	if err := cmdRef.Run(); err != nil {
		return fmt.Errorf("reference runtime error: %v\n%s", err, refOut.String())
	}
	var expected float64
	if _, err := fmt.Sscan(strings.TrimSpace(refOut.String()), &expected); err != nil {
		return fmt.Errorf("failed to parse reference output: %v", err)
	}

	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got float64
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	if math.Abs(got-expected) > 1e-6 {
		return fmt.Errorf("expected %.10f got %.10f", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		if err := runCase(exe, ref, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
