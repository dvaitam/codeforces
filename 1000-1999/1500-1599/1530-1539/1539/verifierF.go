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

type Test struct {
	n int
	a []int
}

func (tc Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func expected(tc Test) string {
	n := tc.n
	a := tc.a // 0-indexed
	res := make([]int, n)
	for i := 0; i < n; i++ {
		best := 0
		for l := 0; l <= i; l++ {
			for r := i; r < n; r++ {
				// Sort subsegment [l..r], find position of a[i], compute distance to center
				seg := make([]int, r-l+1)
				copy(seg, a[l:r+1])
				length := len(seg)
				// Count elements < a[i] and elements > a[i]
				less := 0
				greater := 0
				equal := 0
				for _, v := range seg {
					if v < a[i] {
						less++
					} else if v > a[i] {
						greater++
					} else {
						equal++
					}
				}
				// Element a[i] can be placed at any position in [less+1, less+equal] (1-indexed)
				// Center position: if length is odd, (length+1)/2; if even, length/2 + 1
				var center int
				if length%2 == 1 {
					center = (length + 1) / 2
				} else {
					center = length/2 + 1
				}
				// Try placing a[i] as far from center as possible
				// Possible positions: less+1 to less+equal
				lo := less + 1
				hi := less + equal
				d1 := center - lo
				if d1 < 0 {
					d1 = -d1
				}
				d2 := center - hi
				if d2 < 0 {
					d2 = -d2
				}
				d := d1
				if d2 > d {
					d = d2
				}
				if d > best {
					best = d
				}
			}
		}
		res[i] = best
	}
	var out strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			out.WriteByte(' ')
		}
		out.WriteString(fmt.Sprintf("%d", res[i]))
	}
	return out.String()
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(20) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(n) + 1
	}
	return Test{n: n, a: arr}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		tc := genTest(rng)
		expect := expected(tc)
		got, err := runProg(bin, tc.Input())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc.Input(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
