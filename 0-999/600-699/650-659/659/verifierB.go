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

// validate checks whether output is a correct answer for the given input.
func validate(input, output string) error {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var n, m int
	fmt.Sscan(lines[0], &n, &m)

	// Group students by region, sorted by score descending.
	regions := make([][]student, m+1)
	for _, line := range lines[1:] {
		var s student
		fmt.Sscan(line, &s.name, &s.region, &s.points)
		regions[s.region] = append(regions[s.region], s)
	}
	for i := 1; i <= m; i++ {
		sort.Slice(regions[i], func(a, b int) bool {
			return regions[i][a].points > regions[i][b].points
		})
	}

	outLines := strings.Split(strings.TrimSpace(output), "\n")
	if len(outLines) != m {
		return fmt.Errorf("expected %d output lines, got %d", m, len(outLines))
	}

	for i := 1; i <= m; i++ {
		r := regions[i]
		line := strings.TrimSpace(outLines[i-1])

		// Ambiguous when 3rd highest ties with 2nd highest.
		ambiguous := len(r) >= 3 && r[1].points == r[2].points

		if ambiguous {
			if line != "?" {
				return fmt.Errorf("region %d: expected '?' (ambiguous), got %q", i, line)
			}
		} else {
			// Unique team: top 2 students, in any order.
			parts := strings.Fields(line)
			if len(parts) != 2 {
				return fmt.Errorf("region %d: expected two names, got %q", i, line)
			}
			want := map[string]bool{r[0].name: true, r[1].name: true}
			got := map[string]bool{parts[0]: true, parts[1]: true}
			for k := range want {
				if !got[k] {
					return fmt.Errorf("region %d: expected team {%s, %s}, got %q",
						i, r[0].name, r[1].name, line)
				}
			}
		}
	}
	return nil
}

func generateTests(rng *rand.Rand) []string {
	tests := make([]string, 200)
	for t := 0; t < 200; t++ {
		m := rng.Intn(5) + 1
		n := 2*m + rng.Intn(6)
		students := make([]student, 0, n)
		nameID := 1
		for r := 1; r <= m; r++ {
			for i := 0; i < 2; i++ {
				students = append(students, student{
					name:   fmt.Sprintf("name%d", nameID),
					region: r,
					points: rng.Intn(801),
				})
				nameID++
			}
		}
		for len(students) < n {
			students = append(students, student{
				name:   fmt.Sprintf("name%d", nameID),
				region: rng.Intn(m) + 1,
				points: rng.Intn(801),
			})
			nameID++
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
	rng := rand.New(rand.NewSource(43))
	tests := generateTests(rng)
	for i, tc := range tests {
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := validate(tc, got); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%sgot:\n%s\n", i+1, err, tc, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
