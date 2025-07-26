package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func buildRef() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1085F.go")
	bin := filepath.Join(os.TempDir(), "ref1085F.bin")
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
	n   int
	q   int
	s   string
	ops [][2]interface{}
}

func genCases() []Case {
	r := rand.New(rand.NewSource(1085))
	letters := [3]byte{'R', 'P', 'S'}
	cases := make([]Case, 100)
	for i := range cases {
		n := r.Intn(8) + 2
		q := r.Intn(8) + 1
		var sb strings.Builder
		for j := 0; j < n; j++ {
			sb.WriteByte(letters[r.Intn(3)])
		}
		s := sb.String()
		ops := make([][2]interface{}, q)
		for j := 0; j < q; j++ {
			pos := r.Intn(n) + 1
			ch := letters[r.Intn(3)]
			ops[j] = [2]interface{}{pos, ch}
		}
		cases[i] = Case{n, q, s, ops}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n%s\n", c.n, c.q, c.s)
		for _, op := range c.ops {
			fmt.Fprintf(&sb, "%d %c\n", op[0], op[1])
		}
		input := sb.String()
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
		if want != got {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
