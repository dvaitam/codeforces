package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func buildRef() (string, error) {
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "601D.go")
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

type Case struct{ input string }

func genCases() []Case {
	rng := rand.New(rand.NewSource(63))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(6) + 1 // 1..6
		vals := make([]int64, n)
		for j := 0; j < n; j++ {
			vals[j] = int64(rng.Intn(10))
		}
		letters := make([]byte, n)
		for j := 0; j < n; j++ {
			letters[j] = byte('a' + rng.Intn(3))
		}
		edges := make([][2]int, n-1)
		for v := 2; v <= n; v++ {
			p := rng.Intn(v-1) + 1
			edges[v-2] = [2]int{p, v}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j, v := range vals {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		sb.WriteString(string(letters))
		sb.WriteByte('\n')
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		cases[i] = Case{sb.String()}
	}
	return cases
}

func runCase(bin, ref string, c Case) error {
	expected, err := runBinary(ref, c.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runBinary(bin, c.input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(expected) != strings.TrimSpace(got) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
