package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type item struct {
	a int64
	b int64
}

type Test struct {
	n     int
	items []item
}

func (tc Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for _, it := range tc.items {
		sb.WriteString(fmt.Sprintf("%d %d\n", it.a, it.b))
	}
	return sb.String()
}

func expected(tc Test) string {
	items := make([]item, tc.n)
	copy(items, tc.items)
	sort.Slice(items, func(i, j int) bool { return items[i].b < items[j].b })
	l := 0
	r := tc.n - 1
	var bought, cost int64
	for l <= r {
		for l <= r && items[l].a == 0 {
			l++
		}
		for l <= r && items[r].a == 0 {
			r--
		}
		if l > r {
			break
		}
		if bought >= items[l].b {
			cost += items[l].a
			bought += items[l].a
			items[l].a = 0
			l++
			continue
		}
		need := items[l].b - bought
		if items[r].a <= need {
			cost += items[r].a * 2
			bought += items[r].a
			items[r].a = 0
			r--
		} else {
			cost += need * 2
			bought += need
			items[r].a -= need
		}
	}
	return fmt.Sprintf("%d", cost)
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
	items := make([]item, n)
	for i := 0; i < n; i++ {
		items[i] = item{a: rng.Int63n(20), b: rng.Int63n(50)}
	}
	return Test{n: n, items: items}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
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
