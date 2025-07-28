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
	ref := "./refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1845B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Case struct {
	xa, ya, xb, yb, xc, yc int
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(45))
	cases := make([]Case, 120)
	for i := range cases {
		xa := rng.Intn(20) + 1
		ya := rng.Intn(20) + 1
		var xb, yb int
		for {
			xb = rng.Intn(20) + 1
			yb = rng.Intn(20) + 1
			if xb != xa || yb != ya {
				break
			}
		}
		var xc, yc int
		for {
			xc = rng.Intn(20) + 1
			yc = rng.Intn(20) + 1
			if (xc != xa || yc != ya) && (xc != xb || yc != yb) {
				break
			}
		}
		cases[i] = Case{xa, ya, xb, yb, xc, yc}
	}
	return cases
}

func runCase(bin, ref string, c Case) error {
	input := fmt.Sprintf("1\n%d %d\n%d %d\n%d %d\n", c.xa, c.ya, c.xb, c.yb, c.xc, c.yc)
	exp, err := run(ref, input)
	if err != nil {
		return fmt.Errorf("reference: %v", err)
	}
	got, err := run(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(exp) != strings.TrimSpace(got) {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
