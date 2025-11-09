package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
    "strconv"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func solveCase(n int, l int, arr []int) string {
	sort.Ints(arr)
	radius := math.Max(float64(arr[0]), float64(l-arr[n-1]))
	for i := 0; i+1 < n; i++ {
		gap := float64(arr[i+1]-arr[i]) / 2.0
		if gap > radius {
			radius = gap
		}
	}
	return fmt.Sprintf("%.12f", radius)
}

func parseFirstFloat(s string) (float64, error) {
    fields := strings.Fields(s)
    if len(fields) == 0 {
        return 0, fmt.Errorf("empty output")
    }
    return strconv.ParseFloat(fields[0], 64)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	l := rng.Intn(1000) + 1
	arr := make([]int, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, l)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(l + 1)
		if i > 0 {
			fmt.Fprintf(&sb, " ")
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteString("\n")
	expect := solveCase(n, l, arr)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		outTrim := strings.TrimSpace(out)
        expTrim := strings.TrimSpace(exp)
        if outTrim != expTrim {
            of, oerr := parseFirstFloat(outTrim)
            ef, eerr := parseFirstFloat(expTrim)
            if oerr == nil && eerr == nil {
                diff := math.Abs(of - ef)
                if diff <= 1e-6 {
                    // acceptable numeric tolerance
                    goto CONTINUE_CASE
                }
            }
            fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expTrim, outTrim, in)
            os.Exit(1)
        }
        CONTINUE_CASE:
	}
	fmt.Println("All tests passed")
}
