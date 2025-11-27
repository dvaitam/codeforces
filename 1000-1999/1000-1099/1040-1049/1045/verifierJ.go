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

const refSource = "1000-1999/1000-1099/1040-1049/1045/1045J.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierJ.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for i, input := range tests {
		refOut, err := runExecutable(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if normalize(refOut) != normalize(candOut) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1045J-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	src := os.Getenv("REFERENCE_SOURCE_PATH")
	if strings.TrimSpace(src) == "" {
		src = refSource
	}
	absSrc, err := filepath.Abs(src)
	if err != nil {
		return "", err
	}
	srcDir := filepath.Dir(absSrc)
	srcFile := filepath.Base(absSrc)

	cmd := exec.Command("go", "build", "-o", tmp.Name(), srcFile)
	cmd.Dir = srcDir
	cmd.Env = append(os.Environ(),
		"GO111MODULE=off",
		"GOWORK=off",
	)
	var combined bytes.Buffer
	cmd.Stdout = &combined
	cmd.Stderr = &combined
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, combined.String())
	}
	return tmp.Name(), nil
}

func runExecutable(path, input string) (string, error) {
	cmd := exec.Command(path)
	return execute(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return execute(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func execute(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func normalize(out string) string {
	fields := strings.Fields(out)
	return strings.Join(fields, "\n")
}

func buildTests() []string {
	tests := []string{
		simpleTest(),
		starTest(),
		chainTest(8),
	}

	randomConfigs := []struct {
		n, q int
		seed int64
	}{
		{5, 6, 1},
		{10, 20, 2},
		{30, 40, 3},
		{60, 80, 4},
		{120, 150, 5},
		{200, 200, 6},
		{400, 400, time.Now().UnixNano()},
	}
	for _, cfg := range randomConfigs {
		tests = append(tests, randomTest(cfg.n, cfg.q, cfg.seed))
	}
	return tests
}

func simpleTest() string {
	var sb strings.Builder
	sb.WriteString("5\n")
	sb.WriteString("1 2 a\n")
	sb.WriteString("2 3 b\n")
	sb.WriteString("3 4 c\n")
	sb.WriteString("4 5 d\n")
	sb.WriteString("6\n")
	sb.WriteString("1 5 abcd\n")
	sb.WriteString("5 1 dcba\n")
	sb.WriteString("2 4 bc\n")
	sb.WriteString("2 4 cb\n")
	sb.WriteString("3 3 a\n")
	sb.WriteString("1 5 zz\n")
	return sb.String()
}

func starTest() string {
	var sb strings.Builder
	sb.WriteString("6\n")
	for i := 2; i <= 6; i++ {
		sb.WriteString(fmt.Sprintf("1 %d %c\n", i, 'a'+(i-2)))
	}
	sb.WriteString("5\n")
	sb.WriteString("2 3 a\n")
	sb.WriteString("2 3 b\n")
	sb.WriteString("2 6 ab\n")
	sb.WriteString("4 5 cd\n")
	sb.WriteString("6 6 c\n")
	return sb.String()
}

func chainTest(n int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 2; i <= n; i++ {
		ch := 'a'
		if i%2 == 0 {
			ch = 'b'
		}
		sb.WriteString(fmt.Sprintf("%d %d %c\n", i-1, i, ch))
	}
	sb.WriteString("5\n")
	sb.WriteString(fmt.Sprintf("1 %d ab\n", n))
	sb.WriteString(fmt.Sprintf("%d 1 ba\n", n))
	sb.WriteString(fmt.Sprintf("2 %d b\n", n))
	sb.WriteString(fmt.Sprintf("3 %d a\n", n))
	sb.WriteString("1 1 aba\n")
	return sb.String()
}

func randomTest(n, q int, seed int64) string {
	r := rand.New(rand.NewSource(seed))
	var sb strings.Builder
	if n < 2 {
		n = 2
	}
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for v := 2; v <= n; v++ {
		parent := r.Intn(v-1) + 1
		ch := 'a' + rune(r.Intn(26))
		sb.WriteString(fmt.Sprintf("%d %d %c\n", parent, v, ch))
	}
	if q < 1 {
		q = 1
	}
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		u := r.Intn(n) + 1
		v := r.Intn(n) + 1
		length := r.Intn(100) + 1
		sb.WriteString(fmt.Sprintf("%d %d %s\n", u, v, randomString(r, length)))
	}
	return sb.String()
}

func randomString(r *rand.Rand, length int) string {
	data := make([]byte, length)
	for i := range data {
		data[i] = byte('a' + r.Intn(26))
	}
	return string(data)
}
