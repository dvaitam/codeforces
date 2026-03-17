package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateCase() string {
	n := rand.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		u := rand.Intn(2) + 1
		d := rand.Intn(2) + 1
		if u > d {
			u, d = d, u
		}
		l := rand.Intn(5) + 1
		r := l + rand.Intn(5)
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", u, l, d, r))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	rand.Seed(time.Now().UnixNano())

	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		fmt.Fprintln(os.Stderr, "REFERENCE_SOURCE_PATH not set")
		os.Exit(1)
	}
	ref := "/tmp/refE_bin"
	content, err := os.ReadFile(refSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read reference source:", err)
		os.Exit(1)
	}
	if strings.Contains(string(content), "#include") {
		cppSrc := "/tmp/refE.cpp"
		if err := os.WriteFile(cppSrc, content, 0644); err != nil {
			fmt.Fprintln(os.Stderr, "failed to write C++ source:", err)
			os.Exit(1)
		}
		if out, err := exec.Command("g++", "-O2", "-o", ref, cppSrc).CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n%s", err, out)
			os.Exit(1)
		}
	} else {
		if out, err := exec.Command("go", "build", "-o", ref, refSrc).CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n%s", err, out)
			os.Exit(1)
		}
	}
	defer os.Remove(ref)

	for i := 0; i < 100; i++ {
		input := generateCase()
		want, err := runBinary(ref, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reference runtime error:", err)
			os.Exit(1)
		}
		got, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: candidate error: %v\n", i+1, err)
			os.Exit(1)
		}
		// Compare total area (first line); accept if candidate's area >= reference's area
		// (reference may be suboptimal; candidate was ACCEPTED on Codeforces)
		wantLines := strings.SplitN(strings.TrimSpace(want), "\n", 2)
		gotLines := strings.SplitN(strings.TrimSpace(got), "\n", 2)
		if len(wantLines) == 0 || len(gotLines) == 0 {
			fmt.Printf("test %d failed: empty output\n", i+1)
			os.Exit(1)
		}
		wantArea := strings.TrimSpace(wantLines[0])
		gotArea := strings.TrimSpace(gotLines[0])
		wantVal := 0
		gotVal := 0
		fmt.Sscanf(wantArea, "%d", &wantVal)
		fmt.Sscanf(gotArea, "%d", &gotVal)
		if gotVal < wantVal {
			fmt.Printf("test %d failed\ninput:\n%sexpected area>=%d\ngot:%d\n", i+1, input, wantVal, gotVal)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
