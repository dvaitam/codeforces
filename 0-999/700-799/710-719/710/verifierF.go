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
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "710F.go")
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

func randString(rng *rand.Rand) string {
	l := rng.Intn(5) + 1
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		b[i] = byte('a' + rng.Intn(3))
	}
	return string(b)
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(715))
	cases := make([]Case, 100)
	for i := range cases {
		m := rng.Intn(20) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", m))
		set := map[string]bool{}
		for j := 0; j < m; j++ {
			t := rng.Intn(3) + 1
			s := randString(rng)
			if t == 1 {
				set[s] = true
			} else if t == 2 {
				if len(set) == 0 {
					t = 1
					set[s] = true
				} else {
					for str := range set {
						s = str
						break
					}
					delete(set, s)
				}
			}
			sb.WriteString(fmt.Sprintf("%d %s\n", t, s))
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
	if expected != strings.TrimSpace(got) {
		return fmt.Errorf("expected %q got %q", expected, got)
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
