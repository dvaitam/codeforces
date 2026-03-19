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

func buildRef() (string, error) {
	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	content, err := os.ReadFile(refSrc)
	if err != nil {
		return "", fmt.Errorf("read reference source: %v", err)
	}
	bin := os.TempDir() + "/1819B_ref.bin"
	if strings.Contains(string(content), "#include") {
		cppSrc := os.TempDir() + "/1819B_ref.cpp"
		if err := os.WriteFile(cppSrc, content, 0644); err != nil {
			return "", fmt.Errorf("write cpp source: %v", err)
		}
		cmd := exec.Command("g++", "-O2", "-o", bin, cppSrc)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build c++ reference failed: %v\n%s", err, string(out))
		}
	} else {
		cmd := exec.Command("go", "build", "-o", bin, refSrc)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
		}
	}
	return bin, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	H := rng.Intn(5) + 1
	W := rng.Intn(5) + 1
	rects := [][2]int{{H, W}}
	for len(rects) < n {
		idx := rng.Intn(len(rects))
		h, w := rects[idx][0], rects[idx][1]
		if h == 1 && w == 1 {
			continue
		}
		if (rng.Intn(2) == 0 && h > 1) || w == 1 {
			x := rng.Intn(h-1) + 1
			rects[idx][0] = h - x
			rects = append(rects, [2]int{x, w})
		} else {
			y := rng.Intn(w-1) + 1
			rects[idx][1] = w - y
			rects = append(rects, [2]int{h, y})
		}
	}
	rng.Shuffle(len(rects), func(i, j int) { rects[i], rects[j] = rects[j], rects[i] })
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(rects)))
	for _, r := range rects {
		sb.WriteString(fmt.Sprintf("%d %d\n", r[0], r[1]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	candPath := os.Args[1]
	refPath, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refPath)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		exp, err := runBinary(refPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(candPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		// Compare as sets of pairs (order may differ)
		expLines := strings.Split(strings.TrimSpace(exp), "\n")
		gotLines := strings.Split(strings.TrimSpace(got), "\n")
		if len(expLines) == 0 || len(gotLines) == 0 || strings.TrimSpace(expLines[0]) != strings.TrimSpace(gotLines[0]) {
			fmt.Printf("case %d failed (count mismatch)\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
		if len(expLines) != len(gotLines) {
			fmt.Printf("case %d failed (line count mismatch)\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
		expSet := make(map[string]bool)
		for _, l := range expLines[1:] {
			expSet[strings.TrimSpace(l)] = true
		}
		match := true
		for _, l := range gotLines[1:] {
			if !expSet[strings.TrimSpace(l)] {
				match = false
				break
			}
		}
		if !match {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
