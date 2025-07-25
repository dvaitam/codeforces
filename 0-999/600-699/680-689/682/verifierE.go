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

type point struct{ x, y int64 }

type testCase struct {
	n   int
	S   int64
	pts []point
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return out.String(), nil
}

func area2(a, b, c point) int64 {
	v := (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
	if v < 0 {
		return -v
	}
	return v
}

func pointInTriangle(p, a, b, c point) bool {
	s1 := (b.x-a.x)*(p.y-a.y) - (b.y-a.y)*(p.x-a.x)
	s2 := (c.x-b.x)*(p.y-b.y) - (c.y-b.y)*(p.x-b.x)
	s3 := (a.x-c.x)*(p.y-c.y) - (a.y-c.y)*(p.x-c.x)
	hasNeg := (s1 < 0) || (s2 < 0) || (s3 < 0)
	hasPos := (s1 > 0) || (s2 > 0) || (s3 > 0)
	return !(hasNeg && hasPos)
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.S))
	for i, p := range tc.pts {
		sb.WriteString(fmt.Sprintf("%d %d", p.x, p.y))
		if i+1 < len(tc.pts) {
			sb.WriteByte('\n')
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	lines := strings.Fields(out)
	if len(lines) < 6 {
		return fmt.Errorf("output should contain six integers")
	}
	var tri [3]point
	for i := 0; i < 3; i++ {
		if _, err := fmt.Sscan(lines[2*i], &tri[i].x); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		if _, err := fmt.Sscan(lines[2*i+1], &tri[i].y); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
	}
	a, b, c := tri[0], tri[1], tri[2]
	area := area2(a, b, c)
	if area == 0 || area > 8*tc.S {
		return fmt.Errorf("triangle area invalid")
	}
	for _, p := range tc.pts {
		if !pointInTriangle(p, a, b, c) {
			return fmt.Errorf("point (%d,%d) not inside triangle", p.x, p.y)
		}
	}
	return nil
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 3
	pts := make([]point, n)
	for i := range pts {
		pts[i] = point{int64(rng.Intn(200) - 100), int64(rng.Intn(200) - 100)}
	}
	S := int64(1000000000)
	return testCase{n: n, S: S, pts: pts}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, generateRandomCase(rng))
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, buildInput(tc))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
