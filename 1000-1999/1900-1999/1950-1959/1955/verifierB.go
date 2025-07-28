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

func buildRef() (string, error) {
	ref := "refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1955B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v: %s", err, string(out))
	}
	return "./" + ref, nil
}

func runProgram(bin, input string) (string, error) {
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

func genValid(rng *rand.Rand) (string, []int64, int, int64, int64) {
	n := rng.Intn(4) + 2
	c := rng.Int63n(10) + 1
	d := rng.Int63n(10) + 1
	base := rng.Int63n(50) + 1
	arr := make([]int64, n*n)
	idx := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			arr[idx] = base + int64(i)*c + int64(j)*d
			idx++
		}
	}
	// shuffle arr
	for i := len(arr) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d %d\n", n, c, d)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String(), arr, n, c, d
}

func genInvalid(rng *rand.Rand) string {
	_, arr, n, c, d := genValid(rng)
	// modify one value
	arr[0]++
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d %d\n", n, c, d)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
	cases := make([]string, 0, 100)
	for len(cases) < 50 {
		tc, _, _, _, _ := genValid(rng)
		cases = append(cases, tc)
	}
	for len(cases) < 100 {
		cases = append(cases, genInvalid(rng))
	}

	for i, tc := range cases {
		exp, err := runProgram(ref, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProgram(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %s got %s\n", i+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
