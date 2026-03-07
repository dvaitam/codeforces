package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
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

func buildRef() (string, error) {
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1739F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v: %s", err, out)
	}
	return ref, nil
}

var letters = []rune("abcdefghijkl")

func genString(rng *rand.Rand, length int) string {
	var sb strings.Builder
	prev := -1
	for i := 0; i < length; i++ {
		idx := rng.Intn(len(letters))
		if idx == prev {
			idx = (idx + 1) % len(letters)
		}
		sb.WriteRune(letters[idx])
		prev = idx
	}
	return sb.String()
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		w := rng.Intn(100) + 1
		l := rng.Intn(8) + 2
		s := genString(rng, l)
		sb.WriteString(fmt.Sprintf("%d %s\n", w, s))
	}
	return sb.String()
}

// score computes the optimality of a keyboard for the given words.
// keyboard is a 12-letter permutation of a-l.
// words is a slice of (cost, word) pairs.
func score(keyboard string, words [][2]string) (int, error) {
	if len(keyboard) != 12 {
		return 0, fmt.Errorf("keyboard length %d != 12", len(keyboard))
	}
	pos := make(map[rune]int)
	for i, ch := range keyboard {
		if ch < 'a' || ch > 'l' {
			return 0, fmt.Errorf("invalid character %c in keyboard", ch)
		}
		if _, dup := pos[ch]; dup {
			return 0, fmt.Errorf("duplicate character %c in keyboard", ch)
		}
		pos[ch] = i
	}
	if len(pos) != 12 {
		return 0, fmt.Errorf("keyboard does not contain all 12 letters")
	}
	total := 0
	for _, w := range words {
		var cost int
		fmt.Sscan(w[0], &cost)
		word := w[1]
		easy := true
		for j := 0; j+1 < len(word); j++ {
			p1, p2 := pos[rune(word[j])], pos[rune(word[j+1])]
			diff := p1 - p2
			if diff < 0 {
				diff = -diff
			}
			if diff != 1 {
				easy = false
				break
			}
		}
		if easy {
			total += cost
		}
	}
	return total, nil
}

// parseWords extracts (cost, word) pairs from the test input.
func parseWords(input string) [][2]string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var words [][2]string
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			words = append(words, [2]string{fields[0], fields[1]})
		}
	}
	return words
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(5))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		words := parseWords(input)

		want, err := runProg(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}

		wantScore, err := score(strings.TrimSpace(want), words)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid reference output on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotScore, err := score(strings.TrimSpace(got), words)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid candidate output on test %d: %v\ninput:\n%s\noutput: %s\n", i+1, err, input, got)
			os.Exit(1)
		}
		if gotScore != wantScore {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s (score %d)\nactual:%s (score %d)\n",
				i+1, input, want, wantScore, got, gotScore)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
