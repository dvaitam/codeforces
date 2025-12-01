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

const refSource = "./1090J.go"

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
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n",
				i+1, input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1090J-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
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
	return strings.TrimSpace(out)
}

func buildTests() []string {
	tests := []string{
		"aba\naa\n",
		"aaaaaaaaa\na\n",
		"a\na\n",
		"abcde\nvwxyz\n",
		"zzzzzzzzzz\nzzzzzzzzzz\n",
		"abababab\nbaba\n",
	}

	randomConfigs := []struct {
		lenS int
		lenT int
		seed int64
	}{
		{5, 5, 1},
		{10, 8, 2},
		{20, 18, 3},
		{50, 50, 4},
		{100, 100, 5},
		{200, 180, 6},
		{500, 500, time.Now().UnixNano()},
	}

	for _, cfg := range randomConfigs {
		tests = append(tests, randomTest(cfg.lenS, cfg.lenT, cfg.seed))
	}
	return tests
}

func randomTest(lenS, lenT int, seed int64) string {
	r := rand.New(rand.NewSource(seed))
	if lenS < 1 {
		lenS = 1
	}
	if lenT < 1 {
		lenT = 1
	}
	return fmt.Sprintf("%s\n%s\n", randomString(r, lenS), randomString(r, lenT))
}

func randomString(r *rand.Rand, length int) string {
	data := make([]byte, length)
	for i := range data {
		data[i] = byte('a' + r.Intn(26))
	}
	return string(data)
}
