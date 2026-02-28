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
	ref := "./refF1.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1379F1.go")
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Case struct{ input string }

func genCases() []Case {
	rng := rand.New(rand.NewSource(6))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(4) + 1
		m := rng.Intn(4) + 1
		// collect all white cells (i+j even)
		var whites [][2]int
		for r := 1; r <= 2*n; r++ {
			for c := 1; c <= 2*m; c++ {
				if (r+c)%2 == 0 {
					whites = append(whites, [2]int{r, c})
				}
			}
		}
		// shuffle and pick up to q unique white cells
		rng.Shuffle(len(whites), func(a, b int) {
			whites[a], whites[b] = whites[b], whites[a]
		})
		q := rng.Intn(4) + 1
		if q > len(whites) {
			q = len(whites)
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, m, q)
		for j := 0; j < q; j++ {
			fmt.Fprintf(&sb, "%d %d\n", whites[j][0], whites[j][1])
		}
		cases[i] = Case{sb.String()}
	}
	return cases
}

func runCase(bin, ref string, c Case) error {
	expect, err := runBinary(ref, c.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runBinary(bin, c.input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
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
