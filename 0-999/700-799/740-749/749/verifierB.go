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

func buildOracle() (string, error) {
	refPath := os.Getenv("REFERENCE_SOURCE_PATH")
	if refPath == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	oracle := filepath.Join(os.TempDir(), "oracleB")
	// Detect language by reading file content
	content, err := os.ReadFile(refPath)
	if err != nil {
		return "", fmt.Errorf("read reference: %v", err)
	}
	if strings.Contains(string(content), "#include") {
		// C++ reference saved as .go by worker; copy to .cpp and compile
		cppPath := filepath.Join(os.TempDir(), "ref749B.cpp")
		if err := os.WriteFile(cppPath, content, 0644); err != nil {
			return "", err
		}
		cmd := exec.Command("g++", "-O2", "-o", oracle, cppPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build oracle (c++) failed: %v\n%s", err, out)
		}
	} else {
		cmd := exec.Command("go", "build", "-o", oracle, refPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
		}
	}
	return oracle, nil
}

func genCase(rng *rand.Rand) string {
	x1 := rng.Intn(41) - 20
	y1 := rng.Intn(41) - 20
	x2 := rng.Intn(41) - 20
	y2 := rng.Intn(41) - 20
	x3 := rng.Intn(41) - 20
	y3 := rng.Intn(41) - 20
	return fmt.Sprintf("%d %d %d %d %d %d\n", x1, y1, x2, y2, x3, y3)
}

func parsePointSet(s string) (int, map[string]bool) {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	if len(lines) == 0 {
		return 0, nil
	}
	count := 0
	fmt.Sscanf(lines[0], "%d", &count)
	set := make(map[string]bool)
	for _, line := range lines[1:] {
		set[strings.TrimSpace(line)] = true
	}
	return count, set
}

func runCase(bin, oracle, input string) error {
	cmdO := exec.Command(oracle)
	cmdO.Stdin = strings.NewReader(input)
	var outO bytes.Buffer
	cmdO.Stdout = &outO
	if err := cmdO.Run(); err != nil {
		return fmt.Errorf("oracle run error: %v", err)
	}
	expCount, expSet := parsePointSet(outO.String())

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	gotCount, gotSet := parsePointSet(out.String())
	if gotCount != expCount {
		return fmt.Errorf("expected count %d got %d", expCount, gotCount)
	}
	if len(gotSet) != len(expSet) {
		return fmt.Errorf("expected %d unique points got %d", len(expSet), len(gotSet))
	}
	for p := range expSet {
		if !gotSet[p] {
			return fmt.Errorf("missing point %s in output", p)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		if err := runCase(bin, oracle, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
