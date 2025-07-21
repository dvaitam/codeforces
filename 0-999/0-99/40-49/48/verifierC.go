package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func floor(x float64) float64 {
	fx := math.Floor(x)
	if x < 0 && fx != x {
		return fx - 1
	}
	return fx
}

func ceil(x float64) float64 {
	fx := math.Floor(x)
	if x > fx {
		return fx + 1
	}
	return fx
}

func solve(stops []int) string {
	n := len(stops)
	L := 10.0 * float64(stops[0])
	R := (10.0*float64(stops[0]) + 10.0)
	for i := 2; i <= n; i++ {
		si := float64(stops[i-1])
		li := 10.0 * si / float64(i)
		ri := (10.0*si + 10.0) / float64(i)
		if li > L {
			L = li
		}
		if ri < R {
			R = ri
		}
	}
	if L < 10 {
		L = 10
	}
	last := stops[n-1]
	eps := 1e-12
	uMin := ((float64(n+1) * L) - 10.0) / 10.0
	uMax := ((float64(n+1) * R) - 10.0) / 10.0
	mMin := int(floor(uMin + eps))
	cand1 := mMin + 1
	if cand1 < last+1 {
		cand1 = last + 1
	}
	mMax1 := int(ceil(uMax - eps))
	cand2 := mMax1
	if cand2 < last+1 {
		cand2 = last + 1
	}
	if cand1 == cand2 {
		return fmt.Sprintf("unique\n%d", cand1)
	}
	return "not unique"
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(10) + 1
	stops := make([]int, n)
	last := 0
	for i := 0; i < n; i++ {
		last += r.Intn(5) + 1
		stops[i] = last
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(stops[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), solve(stops)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(r)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
