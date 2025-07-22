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

func runBinary(bin, input string) (string, error) {
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

func expectedAnswer(lines []string) string {
	labels := []string{"A", "B", "C", "D"}
	lens := make([]int, 4)
	for i := 0; i < 4; i++ {
		line := strings.TrimSpace(lines[i])
		if len(line) >= 2 {
			lens[i] = len(line[2:])
		}
	}
	greatIdx := -1
	count := 0
	for i := 0; i < 4; i++ {
		shorter := true
		longer := true
		for j := 0; j < 4; j++ {
			if i == j {
				continue
			}
			if lens[i]*2 > lens[j] {
				shorter = false
			}
			if lens[i] < 2*lens[j] {
				longer = false
			}
		}
		if shorter || longer {
			count++
			greatIdx = i
		}
	}
	if count == 1 {
		return labels[greatIdx]
	}
	return "C"
}

func randString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, string) {
	labels := []string{"A", "B", "C", "D"}
	lines := make([]string, 4)
	for i := 0; i < 4; i++ {
		l := rng.Intn(100) + 1
		lines[i] = fmt.Sprintf("%s.%s", labels[i], randString(rng, l))
	}
	input := strings.Join(lines, "\n") + "\n"
	exp := expectedAnswer(lines)
	return input, exp
}

func manualCase(descs [4]string) (string, string) {
	labels := []string{"A", "B", "C", "D"}
	lines := make([]string, 4)
	for i := 0; i < 4; i++ {
		lines[i] = fmt.Sprintf("%s.%s", labels[i], descs[i])
	}
	input := strings.Join(lines, "\n") + "\n"
	exp := expectedAnswer(lines)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases [][2]string
	in1, ex1 := manualCase([4]string{"short", "mediumlength", "another", "tiny"})
	cases = append(cases, [2]string{in1, ex1})
	in2, ex2 := manualCase([4]string{"aaaaaaaaaa", "bbbbbbbbbbbbbbbbbbbbbbbb", "cccccccccc", "dddddddddd"})
	cases = append(cases, [2]string{in2, ex2})
	in3, ex3 := manualCase([4]string{"same", "same", "same", "same"})
	cases = append(cases, [2]string{in3, ex3})
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		cases = append(cases, [2]string{in, exp})
	}
	for idx, tc := range cases {
		out, err := runBinary(bin, tc[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc[0])
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc[1] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc[1], out, tc[0])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
