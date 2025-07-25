package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runProgram(bin string, input string) (string, error) {
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

func generatePolygon(rng *rand.Rand, n int) [][2]int64 {
	pts := make([][2]int64, n)
	radius := float64(1000)
	for i := 0; i < n; i++ {
		angle := 2 * math.Pi * float64(i) / float64(n)
		x := int64(math.Round(radius * math.Cos(angle)))
		y := int64(math.Round(radius * math.Sin(angle)))
		pts[i] = [2]int64{x, y}
	}
	return pts
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	refBin := "ref611G.bin"
	if err := exec.Command("go", "build", "-o", refBin, "611G.go").Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(7))
	for t := 0; t < 100; t++ {
		n := rng.Intn(6) + 4
		pts := generatePolygon(rng, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%d %d\n", pts[i][0]+int64(rng.Intn(3)), pts[i][1]+int64(rng.Intn(3)))
		}
		input := sb.String()
		want, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
			os.Exit(1)
		}
		out, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(want) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", t+1, want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
