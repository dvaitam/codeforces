package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin string, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(8) + 2
	maxEdges := n * (n - 1)
	m := rng.Intn(maxEdges + 1)
	edges := make(map[[2]int]struct{})
	var b strings.Builder
	b.WriteString("1\n")
	fmt.Fprintf(&b, "%d %d\n", n, m)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		key := [2]int{u, v}
		if _, ok := edges[key]; ok {
			continue
		}
		edges[key] = struct{}{}
		fmt.Fprintf(&b, "%d %d\n", u, v)
	}
	return b.String()
}

// checkAnswer validates the candidate output for a single test case.
// n = number of residents/cats, friendships is a set of (resident, cat) pairs.
// output is the candidate's trimmed output for this test case.
// Returns "" on success or an error description.
func checkAnswer(n int, friendships map[[2]int]bool, output string) string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return "empty output"
	}
	first := strings.TrimSpace(lines[0])
	if strings.EqualFold(first, "No") {
		// Verify that "No" is actually correct: check if there's any valid partition.
		// For small n, we can try all possible non-empty subsets of residents as jury.
		// If j+p=n, jury = subset of residents, contestants = complement set of cats.
		// Jury member r must not know any contestant cat c (where c is in P).
		// P = set of cats = {1..n} \ jury set.
		for mask := 1; mask < (1 << n) - 1; mask++ {
			jury := make([]int, 0)
			cats := make([]int, 0)
			for i := 1; i <= n; i++ {
				if mask&(1<<(i-1)) != 0 {
					jury = append(jury, i)
				} else {
					cats = append(cats, i)
				}
			}
			valid := true
			for _, r := range jury {
				for _, c := range cats {
					if friendships[[2]int{r, c}] {
						valid = false
						break
					}
				}
				if !valid {
					break
				}
			}
			if valid {
				return fmt.Sprintf("answered No but valid partition exists: jury=%v cats=%v", jury, cats)
			}
		}
		return ""
	}
	if !strings.EqualFold(first, "Yes") {
		return fmt.Sprintf("first line should be Yes or No, got %q", first)
	}
	if len(lines) < 4 {
		return "Yes answer needs 4 lines"
	}
	parts := strings.Fields(strings.TrimSpace(lines[1]))
	if len(parts) != 2 {
		return fmt.Sprintf("expected 2 numbers on sizes line, got %d", len(parts))
	}
	j, _ := strconv.Atoi(parts[0])
	p, _ := strconv.Atoi(parts[1])
	if j+p != n || j < 1 || p < 1 {
		return fmt.Sprintf("invalid sizes j=%d p=%d n=%d", j, p, n)
	}
	juryParts := strings.Fields(strings.TrimSpace(lines[2]))
	if len(juryParts) != j {
		return fmt.Sprintf("expected %d jury members, got %d", j, len(juryParts))
	}
	catParts := strings.Fields(strings.TrimSpace(lines[3]))
	if len(catParts) != p {
		return fmt.Sprintf("expected %d contestants, got %d", p, len(catParts))
	}
	jurySet := make(map[int]bool)
	for _, s := range juryParts {
		v, _ := strconv.Atoi(s)
		if v < 1 || v > n {
			return fmt.Sprintf("jury member %d out of range", v)
		}
		if jurySet[v] {
			return fmt.Sprintf("duplicate jury member %d", v)
		}
		jurySet[v] = true
	}
	catSet := make(map[int]bool)
	for _, s := range catParts {
		v, _ := strconv.Atoi(s)
		if v < 1 || v > n {
			return fmt.Sprintf("cat %d out of range", v)
		}
		if catSet[v] {
			return fmt.Sprintf("duplicate cat %d", v)
		}
		catSet[v] = true
	}
	// Check: jury members and cat contestants should together cover all n indices
	if len(jurySet)+len(catSet) != n {
		return "jury + cats don't cover all n"
	}
	// Check no jury member knows any contestant cat
	for r := range jurySet {
		for c := range catSet {
			if friendships[[2]int{r, c}] {
				return fmt.Sprintf("jury member %d knows cat %d", r, c)
			}
		}
	}
	return ""
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}

		// Parse the input to extract n, m, and friendships
		lines := strings.Split(strings.TrimSpace(input), "\n")
		// lines[0] = "1" (t=1)
		parts := strings.Fields(lines[1])
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		friendships := make(map[[2]int]bool)
		for k := 0; k < m; k++ {
			ep := strings.Fields(lines[2+k])
			u, _ := strconv.Atoi(ep[0])
			v, _ := strconv.Atoi(ep[1])
			friendships[[2]int{u, v}] = true
		}

		errMsg := checkAnswer(n, friendships, got)
		if errMsg != "" {
			fmt.Fprintf(os.Stderr, "case %d failed: %s\ninput:\n%sgot:\n%s\n", i+1, errMsg, input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
