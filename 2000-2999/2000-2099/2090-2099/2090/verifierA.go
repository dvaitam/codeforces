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

type testCase struct {
	x, y, a int64
}

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func parseAnswers(out string, t int) ([]string, error) {
	reader := strings.NewReader(out)
	ans := make([]string, 0, t)
	for i := 0; i < t; i++ {
		var tok string
		if _, err := fmt.Fscan(reader, &tok); err != nil {
			return nil, fmt.Errorf("output ended early on case %d: %v", i+1, err)
		}
		up := strings.ToUpper(tok)
		if up != "YES" && up != "NO" {
			return nil, fmt.Errorf("invalid verdict '%s' on case %d", tok, i+1)
		}
		ans = append(ans, up)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output after %d answers", t)
	}
	return ans, nil
}

func r64(lo, hi int64) int64 {
	return lo + rand.Int63n(hi-lo+1)
}

func generateCase() testCase {
	mode := rand.Intn(6)
	switch mode {
	case 0:
		return testCase{r64(1, 1_000_000_000), r64(1, 1_000_000_000), r64(1, 1_000_000_000)}
	case 1:
		// Small numbers to mirror samples.
		return testCase{r64(1, 5), r64(1, 5), r64(1, 6)}
	case 2:
		// Ensure B finishes immediately.
		a := r64(1, 30)
		x := a + 1
		return testCase{x, r64(1, 50), a}
	case 3:
		// K much faster.
		a := r64(1, 1_000_000_000)
		x := r64(1, a)
		y := r64(a, 1_000_000_000)
		return testCase{x, y, a}
	case 4:
		// Near threshold in several cycles.
		cycles := r64(1, 10)
		x := r64(1, 1_000_000_000)
		y := r64(1, 1_000_000_000)
		target := cycles*(x+y) + r64(0, x+y-1)
		if target > 1_000_000_000 {
			target = 1_000_000_000
		}
		return testCase{x, y, target}
	default:
		// One digs zero extra beyond first move.
		a := r64(1, 1_000_000_000)
		x := r64(1, 10)
		y := 1
		return testCase{x, y, a}
	}
}

func buildInput() []byte {
	t := rand.Intn(20) + 1
	if rand.Intn(6) == 0 {
		t = 1000
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		tc := generateCase()
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.x, tc.y, tc.a))
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refA.bin"
	if err := exec.Command("go", "build", "-o", ref, "2090A.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())

	for iter := 0; iter < 200; iter++ {
		input := buildInput()
		refOut, err := run(ref, input)
		if err != nil {
			fmt.Println("reference failed on iteration", iter+1, ":", err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		candOut, err := run(cand, input)
		if err != nil {
			fmt.Printf("candidate runtime error on iteration %d: %v\n", iter+1, err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}

		var t int
		if _, err := fmt.Fscan(strings.NewReader(string(input)), &t); err != nil {
			fmt.Println("failed to parse generated input:", err)
			os.Exit(1)
		}

		refAns, err := parseAnswers(refOut, t)
		if err != nil {
			fmt.Println("failed to parse reference output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", refOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		candAns, err := parseAnswers(candOut, t)
		if err != nil {
			fmt.Println("failed to parse candidate output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", candOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}

		if len(refAns) != len(candAns) {
			fmt.Printf("answer count mismatch on iteration %d\n", iter+1)
			fmt.Println("input:\n", string(input))
			fmt.Println("reference:\n", refOut)
			fmt.Println("candidate:\n", candOut)
			os.Exit(1)
		}
		for i := range refAns {
			if refAns[i] != candAns[i] {
				fmt.Printf("wrong answer on iteration %d, case %d\n", iter+1, i+1)
				fmt.Println("input case starts at above input; reference:", refAns[i], "candidate:", candAns[i])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed.")
}
