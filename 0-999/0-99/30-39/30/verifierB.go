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

func daysInMonth(m, y int) int {
	switch m {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if y%4 == 0 {
			return 29
		}
		return 28
	default:
		return 0
	}
}

func canParticipate(fd, fm, fy int, tokens [3]int) bool {
	perms := [6][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
	for _, p := range perms {
		d := tokens[p[0]]
		m := tokens[p[1]]
		y := tokens[p[2]]
		if m < 1 || m > 12 {
			continue
		}
		yearBirth := 2000 + y
		if d < 1 || d > daysInMonth(m, yearBirth) {
			continue
		}
		yearFinal := 2000 + fy
		age := yearFinal - yearBirth
		if fm < m || (fm == m && fd < d) {
			age--
		}
		if age >= 18 {
			return true
		}
	}
	return false
}

func solve(final, birth string) string {
	var fd, fm, fy int
	fmt.Sscanf(final, "%d.%d.%d", &fd, &fm, &fy)
	var t0, t1, t2 int
	fmt.Sscanf(birth, "%d.%d.%d", &t0, &t1, &t2)
	if canParticipate(fd, fm, fy, [3]int{t0, t1, t2}) {
		return "YES"
	}
	return "NO"
}

func randomDate(rng *rand.Rand) string {
	y := rng.Intn(99) + 1
	m := rng.Intn(12) + 1
	d := rng.Intn(daysInMonth(m, 2000+y)) + 1
	return fmt.Sprintf("%02d.%02d.%02d", d, m, y)
}

func generateCase(rng *rand.Rand) (string, string) {
	final := randomDate(rng)
	birth := randomDate(rng)
	input := final + "\n" + birth + "\n"
	expected := solve(final, birth)
	return input, expected
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	if outStr != expected {
		return fmt.Errorf("expected %s got %s", expected, outStr)
	}
	return nil
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
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
