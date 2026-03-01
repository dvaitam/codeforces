package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

func buildRef() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1085E.go")
	bin := filepath.Join(os.TempDir(), "ref1085E.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return bin, nil
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Case struct {
	k int
	s string
	a string
	b string
}

func parseAnswer(out string) (ok bool, perm string, err error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return false, "", fmt.Errorf("empty output")
	}
	if fields[0] == "NO" {
		return false, "", nil
	}
	if fields[0] != "YES" {
		return false, "", fmt.Errorf("first token must be YES/NO, got %q", fields[0])
	}
	if len(fields) < 2 {
		return true, "", fmt.Errorf("YES without permutation")
	}
	return true, fields[1], nil
}

func validatePermutation(c Case, perm string) error {
	if len(perm) != c.k {
		return fmt.Errorf("permutation length mismatch: got %d want %d", len(perm), c.k)
	}
	b := []byte(perm)
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	for i := 0; i < c.k; i++ {
		want := byte('a' + i)
		if b[i] != want {
			return fmt.Errorf("permutation is not a bijection over first %d letters", c.k)
		}
	}
	t := make([]byte, len(c.s))
	for i := range c.s {
		idx := c.s[i] - 'a'
		if idx < 0 || int(idx) >= c.k {
			return fmt.Errorf("source string has char out of range: %q", c.s[i])
		}
		t[i] = perm[idx]
	}
	ts := string(t)
	if ts < c.a || ts > c.b {
		return fmt.Errorf("mapped string %q is outside [%q, %q]", ts, c.a, c.b)
	}
	return nil
}

func genCases() []Case {
	r := rand.New(rand.NewSource(1085))
	cases := make([]Case, 100)
	letters := "abcdefghijklmnopqrstuvwxyz"
	for i := range cases {
		k := r.Intn(5) + 1
		n := r.Intn(5) + 1
		var sb strings.Builder
		for j := 0; j < n; j++ {
			sb.WriteByte(letters[r.Intn(k)])
		}
		s := sb.String()
		sb.Reset()
		for j := 0; j < n; j++ {
			sb.WriteByte(letters[r.Intn(k)])
		}
		a := sb.String()
		sb.Reset()
		for j := 0; j < n; j++ {
			sb.WriteByte(letters[r.Intn(k)])
		}
		b := sb.String()
		if a > b {
			a, b = b, a
		}
		cases[i] = Case{k, s, a, b}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rand.Seed(time.Now().UnixNano())
	cases := genCases()
	for i, c := range cases {
		input := fmt.Sprintf("1\n%d\n%s\n%s\n%s\n", c.k, c.s, c.a, c.b)
		want, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		wantOK, _, err := parseAnswer(want)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on case %d: %v\noutput: %s\n", i+1, err, want)
			os.Exit(1)
		}
		gotOK, gotPerm, err := parseAnswer(got)
		if err != nil {
			fmt.Printf("case %d failed\ninput:\n%sgot invalid output: %v\nraw output: %s\n", i+1, input, err, got)
			os.Exit(1)
		}
		if wantOK != gotOK {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
		if gotOK {
			if err := validatePermutation(c, gotPerm); err != nil {
				fmt.Printf("case %d failed\ninput:\n%sinvalid YES answer: %v\ngot: %s\n", i+1, input, err, got)
				os.Exit(1)
			}
		}
		if !wantOK && gotOK {
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
