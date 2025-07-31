package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runProg(bin, input string) (string, error) {
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
	rng := rand.New(rand.NewSource(1725))
	cases := make([]Case, 0, 100)
	preset := []struct {
		n    int
		vals []int
	}{{2, []int{1, 2}}, {4, []int{1, 2, 3, 4}}}
	for _, p := range preset {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", p.n))
		for i, v := range p.vals {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		cases = append(cases, Case{sb.String()})
	}
	for len(cases) < 100 {
		n := rng.Intn(3)*2 + 2
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", rng.Intn(10)+1))
		}
		sb.WriteByte('\n')
		cases = append(cases, Case{sb.String()})
	}
	return cases
}

func validate(input, output string) error {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return fmt.Errorf("empty input")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid n")
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i], err = strconv.Atoi(fields[i+1])
		if err != nil {
			return fmt.Errorf("invalid number")
		}
	}
	out := strings.TrimSpace(output)
	if out == "-1" {
		return fmt.Errorf("should be solvable but got -1")
	}
	parts := strings.Fields(out)
	if len(parts) != 2 {
		return fmt.Errorf("expected two tokens")
	}
	z, err := strconv.Atoi(parts[0])
	if err != nil || z < 0 || z > 2 {
		return fmt.Errorf("invalid Z")
	}
	s := parts[1]
	if len(s) != n {
		return fmt.Errorf("invalid string length")
	}
	zeros := 0
	for i := 0; i < n; i++ {
		if s[i] == '0' {
			zeros++
		} else if s[i] != '1' {
			return fmt.Errorf("invalid character in string")
		}
	}
	if zeros*2 != n {
		return fmt.Errorf("need %d zeros got %d", n/2, zeros)
	}
	mods := make([]int, n)
	for i := 0; i < n; i++ {
		mods[i] = a[i] % 3
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if s[i] == s[j] {
				continue
			}
			v := (mods[i]*mods[i] + mods[j]*mods[j]) % 3
			if v == z {
				return fmt.Errorf("pair %d,%d triggers reaction", i+1, j+1)
			}
		}
	}
	return nil
}

func runCase(bin string, c Case) error {
	got, err := runProg(bin, c.input)
	if err != nil {
		return err
	}
	if err := validate(c.input, got); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, c := range cases {
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
