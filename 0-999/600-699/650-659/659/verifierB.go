package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type student struct {
	name   string
	region int
	points int
}

func solve(input string) string {
	var n, m int
	fmt.Sscan(strings.SplitN(input, "\n", 2)[0], &n, &m)
	lines := strings.Split(strings.TrimSpace(input), "\n")[1:]
	top := make([][3]student, m+1)
	for i := range top {
		for j := range top[i] {
			top[i][j].points = -1
		}
	}
	for _, line := range lines {
		var s student
		fmt.Sscan(line, &s.name, &s.region, &s.points)
		r := s.region
		if s.points >= top[r][0].points {
			top[r][2] = top[r][1]
			top[r][1] = top[r][0]
			top[r][0] = s
		} else if s.points >= top[r][1].points {
			top[r][2] = top[r][1]
			top[r][1] = s
		} else if s.points >= top[r][2].points {
			top[r][2] = s
		}
	}
	var b strings.Builder
	for i := 1; i <= m; i++ {
		if top[i][1].points == top[i][2].points {
			b.WriteString("?\n")
		} else {
			fmt.Fprintf(&b, "%s %s\n", top[i][0].name, top[i][1].name)
		}
	}
	return b.String()
}

func generateTests() []string {
	rand.Seed(43)
	tests := make([]string, 100)
	for t := 0; t < 100; t++ {
		m := rand.Intn(5) + 1
		n := 2*m + rand.Intn(6)
		students := make([]student, 0, n)
		nameID := 1
		for r := 1; r <= m; r++ {
			for i := 0; i < 2; i++ {
				s := student{
					name:   fmt.Sprintf("name%d", nameID),
					region: r,
					points: rand.Intn(801),
				}
				nameID++
				students = append(students, s)
			}
		}
		for len(students) < n {
			s := student{
				name:   fmt.Sprintf("name%d", nameID),
				region: rand.Intn(m) + 1,
				points: rand.Intn(801),
			}
			nameID++
			students = append(students, s)
		}
		sort.Slice(students, func(i, j int) bool { return students[i].name < students[j].name })
		var b strings.Builder
		fmt.Fprintf(&b, "%d %d\n", n, m)
		for _, s := range students {
			fmt.Fprintf(&b, "%s %d %d\n", s.name, s.region, s.points)
		}
		tests[t] = b.String()
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		expect := strings.TrimSpace(solve(t))
		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if expect != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, t, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
