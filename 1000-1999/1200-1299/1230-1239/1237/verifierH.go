package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildRef() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	ref := filepath.Join(dir, "oracleH")
	cmd := exec.Command("go", "build", "-o", ref, "1237H.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out)
	}
	return ref, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
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

func genCase(r *rand.Rand) string {
	t := r.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := r.Intn(3) + 1
		var s strings.Builder
		var b strings.Builder
		for j := 0; j < 2*n; j++ {
			if r.Intn(2) == 0 {
				s.WriteByte('0')
				b.WriteByte('0')
			} else {
				s.WriteByte('1')
				b.WriteByte('1')
			}
		}
		sb.WriteString(fmt.Sprintf("%s %s\n", s.String(), b.String()))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
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
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		want, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("wrong answer on test %d\ninput:%sexpected:%s got:%s\n", i, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
