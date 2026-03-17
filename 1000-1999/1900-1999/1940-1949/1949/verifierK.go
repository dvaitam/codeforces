package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, string(out))
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runProg(exe string, input []byte) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func genTest() []byte {
	n := rand.Intn(5) + 3
	na := rand.Intn(n-2) + 1
	nb := rand.Intn(n-na-1) + 1
	nc := n - na - nb
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, na, nb, nc))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d", rand.Intn(20)+1))
		if i+1 == n {
			sb.WriteByte('\n')
		} else {
			sb.WriteByte(' ')
		}
	}
	return []byte(sb.String())
}

// solve implements the correct algorithm (ported from the accepted C++ solution).
// Returns "YES\n..." or "NO".
func solve(n int, cnt [3]int, a []int) string {
	// Sort ascending
	sorted := make([]int, n)
	copy(sorted, a)
	sort.Ints(sorted) // ascending

	// prefix sums
	pre := make([]int64, n+1)
	for i := 0; i < n; i++ {
		pre[i+1] = pre[i] + int64(sorted[i])
	}

	var Sm int64
	for _, v := range sorted {
		Sm += int64(v)
	}
	Sm = (Sm - 1) / 2

	// Check: for each group, the smallest cnt[i] elements must not exceed Sm
	for i := 0; i < 3; i++ {
		if pre[cnt[i]] > Sm {
			return "NO"
		}
	}

	// Check: at least one group can take the largest element along with cnt[i]-1 smallest
	flg := false
	for i := 0; i < 3; i++ {
		if pre[cnt[i]-1]+int64(sorted[n-1]) <= Sm {
			flg = true
		}
	}
	if !flg {
		return "NO"
	}

	// Greedy assignment from largest to smallest
	sm := [3]int64{}
	num := [3]int{}
	ans := [3][]int{{}, {}, {}}

	for i := n - 1; i >= 0; i-- {
		id := -1
		for j := 0; j < 3; j++ {
			if num[j] == cnt[j] {
				continue
			}
			// Check: can we place sorted[i] into group j and still fill the rest?
			// Need: sm[j] + sorted[i] + sum_of_smallest_(cnt[j]-num[j]-1) <= Sm
			remaining := cnt[j] - num[j] - 1
			if remaining < 0 {
				continue
			}
			if sm[j]+int64(sorted[i])+pre[remaining] > Sm {
				continue
			}
			// Also check feasibility: after placing sorted[i] in group j,
			// for each other group k that still needs elements,
			// verify it can still be filled (with the next largest element sorted[i-1]).
			ok := (i == 0) // if last element, trivially ok
			for k := 0; k < 3; k++ {
				need := cnt[k] - num[k]
				if k == j {
					need--
				}
				if need <= 0 {
					continue
				}
				// group k needs 'need' more elements, the largest available would be sorted[i-1]
				extra := int64(0)
				if k == j {
					extra = int64(sorted[i])
				}
				if i >= 1 && sm[k]+extra+int64(sorted[i-1])+pre[need-1] <= Sm {
					ok = true
				}
			}
			if ok {
				id = j
			}
		}
		if id == -1 {
			return "NO"
		}
		ans[id] = append(ans[id], sorted[i])
		sm[id] += int64(sorted[i])
		num[id]++
	}

	var sb strings.Builder
	sb.WriteString("YES")
	for i := 0; i < 3; i++ {
		sb.WriteByte('\n')
		for j, v := range ans[i] {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
	}
	return sb.String()
}

// validateAnswer checks if a candidate's YES answer is actually correct.
func validateAnswer(n int, cnt [3]int, a []int, output string) string {
	lines := strings.Split(output, "\n")
	if len(lines) < 1 {
		return "empty output"
	}
	verdict := strings.TrimSpace(lines[0])
	if verdict == "NO" {
		return "" // we check NO vs YES separately
	}
	if verdict != "YES" {
		return fmt.Sprintf("unknown verdict: %s", verdict)
	}
	if len(lines) < 4 {
		return "YES but not enough lines"
	}

	// Parse 3 groups
	groups := [3][]int{}
	for i := 0; i < 3; i++ {
		parts := strings.Fields(lines[i+1])
		if len(parts) != cnt[i] {
			return fmt.Sprintf("group %d: expected %d elements, got %d", i+1, cnt[i], len(parts))
		}
		for _, p := range parts {
			v, err := strconv.Atoi(p)
			if err != nil {
				return fmt.Sprintf("group %d: bad number %s", i+1, p)
			}
			groups[i] = append(groups[i], v)
		}
	}

	// Check that all elements form a permutation of the input
	got := []int{}
	for i := 0; i < 3; i++ {
		got = append(got, groups[i]...)
	}
	sort.Ints(got)
	expect := make([]int, len(a))
	copy(expect, a)
	sort.Ints(expect)
	if len(got) != len(expect) {
		return "wrong number of elements"
	}
	for i := range got {
		if got[i] != expect[i] {
			return "elements don't match input"
		}
	}

	// Check triangle inequality with positive area
	sums := [3]int64{}
	for i := 0; i < 3; i++ {
		for _, v := range groups[i] {
			sums[i] += int64(v)
		}
	}
	sa, sb, sc := sums[0], sums[1], sums[2]
	if sa+sb <= sc || sa+sc <= sb || sb+sc <= sa {
		return fmt.Sprintf("triangle inequality violated: sums %d %d %d", sa, sb, sc)
	}

	return ""
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierK.go /path/to/binary")
		return
	}
	exe, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer cleanup()

	for i := 1; i <= 100; i++ {
		in := genTest()

		// Parse the test case for our embedded solver
		fields := strings.Fields(string(in))
		// fields: "1" n na nb nc x1 x2 ... xn
		n, _ := strconv.Atoi(fields[1])
		na, _ := strconv.Atoi(fields[2])
		nb, _ := strconv.Atoi(fields[3])
		nc, _ := strconv.Atoi(fields[4])
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j], _ = strconv.Atoi(fields[5+j])
		}
		cnt := [3]int{na, nb, nc}

		// Get reference answer from embedded solver
		refAnswer := solve(n, cnt, a)
		refVerdict := strings.TrimSpace(strings.Split(refAnswer, "\n")[0])

		// Run candidate
		got, err := runProg(exe, in)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n%s", i, err, got)
			os.Exit(1)
		}
		gotVerdict := strings.TrimSpace(strings.Split(got, "\n")[0])

		// Check verdict agreement
		if refVerdict == "NO" && gotVerdict == "NO" {
			continue
		}
		if refVerdict == "YES" && gotVerdict == "NO" {
			fmt.Printf("wrong answer on test %d: reference says YES but candidate says NO\ninput:\n%s\nreference:\n%s\n", i, string(in), refAnswer)
			os.Exit(1)
		}
		if refVerdict == "NO" && gotVerdict == "YES" {
			// Candidate says YES but reference says NO - validate candidate answer
			errMsg := validateAnswer(n, cnt, a, got)
			if errMsg != "" {
				fmt.Printf("wrong answer on test %d: candidate says YES but answer invalid: %s\ninput:\n%s\ncandidate:\n%s\n", i, errMsg, string(in), got)
				os.Exit(1)
			}
			// Candidate found a valid answer that reference missed - that's ok
			continue
		}

		// Both say YES - validate candidate answer semantically
		errMsg := validateAnswer(n, cnt, a, got)
		if errMsg != "" {
			fmt.Printf("wrong answer on test %d: %s\ninput:\n%s\ncandidate:\n%s\n", i, errMsg, string(in), got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
