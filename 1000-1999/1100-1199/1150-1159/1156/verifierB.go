package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1156B.go")
	bin := filepath.Join(os.TempDir(), "oracle1156B.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return bin, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func randString(r *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + r.Intn(26))
	}
	return string(b)
}

func genCase(r *rand.Rand) string {
	t := r.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := r.Intn(10) + 1
		sb.WriteString(randString(r, n))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		input := genCase(r)
		want, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}

		inLines := strings.Split(strings.TrimSpace(input), "\n")
		wantLines := strings.Split(want, "\n")
		gotLines := strings.Split(got, "\n")
		t, _ := strconv.Atoi(inLines[0])
		if len(wantLines) != t || len(gotLines) != t {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
		for caseIdx := 0; caseIdx < t; caseIdx++ {
			s := inLines[caseIdx+1]
			w := wantLines[caseIdx]
			g := gotLines[caseIdx]
			if w == "No answer" {
				if g != "No answer" {
					fmt.Printf("test %d case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, caseIdx+1, input, w, g)
					os.Exit(1)
				}
				continue
			}
			if g == "No answer" || !validOutput(s, g) {
				fmt.Printf("test %d case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, caseIdx+1, input, w, g)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}

func validOutput(in, out string) bool {
	if len(in) != len(out) {
		return false
	}
	freq := make([]int, 26)
	for i := 0; i < len(in); i++ {
		freq[in[i]-'a']++
	}
	for i := 0; i < len(out); i++ {
		idx := out[i] - 'a'
		if idx < 0 || idx >= 26 {
			return false
		}
		freq[idx]--
		if freq[idx] < 0 {
			return false
		}
		if i > 0 {
			diff := int(out[i]) - int(out[i-1])
			if diff == 1 || diff == -1 {
				return false
			}
		}
	}
	for _, v := range freq {
		if v != 0 {
			return false
		}
	}
	return true
}
