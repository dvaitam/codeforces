package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleA")
	cmd := exec.Command("go", "build", "-o", oracle, "407A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	a := rng.Intn(1000) + 1
	b := rng.Intn(1000) + 1
	return fmt.Sprintf("%d %d\n", a, b)
}

// check validates that the solution output is correct for leg lengths a, b.
// oracleYes indicates whether the oracle found a solution.
func check(output string, a, b int, oracleYes bool) error {
	lines := strings.Fields(output)
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	verdict := lines[0]
	if verdict == "NO" {
		if oracleYes {
			return fmt.Errorf("solution says NO but answer is YES")
		}
		return nil
	}
	if verdict != "YES" {
		return fmt.Errorf("first token must be YES or NO, got %q", verdict)
	}
	if !oracleYes {
		return fmt.Errorf("solution says YES but answer is NO")
	}
	// Parse 3 vertices: 6 integers after YES
	nums := lines[1:]
	if len(nums) < 6 {
		return fmt.Errorf("expected 6 integers for 3 vertices, got %d tokens", len(nums))
	}
	coords := make([]int, 6)
	for i := 0; i < 6; i++ {
		v, err := strconv.Atoi(nums[i])
		if err != nil {
			return fmt.Errorf("invalid integer %q", nums[i])
		}
		if v > 1_000_000_000 || v < -1_000_000_000 {
			return fmt.Errorf("coordinate %d out of range", v)
		}
		coords[i] = v
	}
	px, py := coords[0], coords[1]
	qx, qy := coords[2], coords[3]
	rx, ry := coords[4], coords[5]

	// Check no side is axis-parallel
	sides := [][4]int{
		{px, py, qx, qy},
		{qx, qy, rx, ry},
		{px, py, rx, ry},
	}
	for _, s := range sides {
		dx := s[2] - s[0]
		dy := s[3] - s[1]
		if dx == 0 || dy == 0 {
			return fmt.Errorf("side (%d,%d)-(%d,%d) is axis-parallel", s[0], s[1], s[2], s[3])
		}
	}

	// Check right triangle with legs a, b: try right angle at each vertex
	a2, b2 := a*a, b*b
	type pt struct{ x, y int }
	pts := [3]pt{{px, py}, {qx, qy}, {rx, ry}}
	valid := false
	for i := 0; i < 3; i++ {
		v := pts[i]
		u := pts[(i+1)%3]
		w := pts[(i+2)%3]
		ux, uy := u.x-v.x, u.y-v.y
		wx, wy := w.x-v.x, w.y-v.y
		dot := ux*wx + uy*wy
		if dot != 0 {
			continue
		}
		len1 := ux*ux + uy*uy
		len2 := wx*wx + wy*wy
		if (len1 == a2 && len2 == b2) || (len1 == b2 && len2 == a2) {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("vertices do not form a right triangle with legs %d and %d", a, b)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		var a, b int
		fmt.Sscan(input, &a, &b)

		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		oracleYes := strings.HasPrefix(exp, "YES")

		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := check(got, a, b, oracleYes); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\nexpected %s\ngot %s\ninput:\n%s", i+1, err, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
