package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Case struct {
	n     int
	T     int64
	times []int64
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(1194))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(10) + 1
		times := make([]int64, n)
		var sum int64
		for j := 0; j < n; j++ {
			t := rng.Int63n(20) + 1
			times[j] = t
			sum += t
		}
		T := sum/2 + rng.Int63n(sum/2+1)
		cases[i] = Case{n, T, times}
	}
	return cases
}

func buildRef() (string, error) {
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1194F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin, ref string, c Case) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", c.n, c.T)
	for i, t := range c.times {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", t)
	}
	sb.WriteByte('\n')
	input := sb.String()
	expected, err := runBinary(ref, input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	if expected != strings.TrimSpace(got) {
		return fmt.Errorf("expected %s got %s", expected, got)
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

	cases := genCases()
	for i, c := range cases {
		if err := runCase(bin, ref, c); err != nil {
			var sb strings.Builder
			fmt.Fprintf(&sb, "%d %d\n", c.n, c.T)
			for j, t := range c.times {
				if j > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", t)
			}
			sb.WriteByte('\n')
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
