package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "0-999/200-299/200-209/207/207D3.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD3.go /path/to/candidate")
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
	tmp, err := os.CreateTemp("", "207D3-ref-*")
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
	return []string{
		formatDoc(101, "Cabinet Briefing",
			"President and prime minister met parliament leaders at the kremlin embassy to discuss new foreign policy and defense strategy during the summit.\n"),
		formatDoc(202, "Derby Recap",
			"The championship match ended 3:2 as the striker scored the winning goal, the fans filled the stadium, and the coach praised the lineup and tactics after the game.\n"),
		formatDoc(303, "Market Pulse",
			"Global markets rallied as investors tracked GDP forecasts, oil prices, billion dollar mergers, and currency percent swings across the trading sessions.\n"),
		formatDoc(404, "Logistics Outlook",
			"Export contracts, shipping logistics, and warehouse financing dominated the board meeting where analysts cited percent growth and rising capital demand.\n"),
		formatDoc(505, "Cup Preview",
			"Supporters debated the club roster before the tournament, noting the goalkeeper injury, the captain's return, and a possible 1-0 scoreline in the semifinals.\n"),
	}
}

func formatDoc(id int, title, body string) string {
	if !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	return fmt.Sprintf("%d\n%s\n%s", id, title, body)
}
