package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const mod int64 = 998244353
const testcasesB64 = "MTAwCjQgNgoyIDEgMQozIDEgNAo0IDEgMgo0IDEwCjIgMSA4CjMgMiAxCjQgMSA1CjQgNgoyIDEgMQozIDIgNQo0IDEgMwo1IDMKMiAxIDMKMyAyIDIKNCAzIDIKNSAyIDMKMyAxMAoyIDEgMgozIDEgOAo0IDcKMiAxIDMKMyAyIDIKNCAzIDUKNSA5CjIgMSA4CjMgMSA0CjQgMyA1CjUgMiAxCjIgMQoyIDEgMQozIDEwCjIgMSA3CjMgMSA2CjQgOQoyIDEgMQozIDEgNQo0IDIgMgo0IDUKMiAxIDEKMyAxIDMKNCAyIDQKMiAzCjIgMSAyCjQgOAoyIDEgMQozIDEgOAo0IDEgOAo0IDgKMiAxIDEKMyAyIDYKNCAzIDEKMyAxMAoyIDEgNgozIDEgMwo0IDEKMiAxIDEKMyAxIDEKNCAxIDEKNCA2CjIgMSAzCjMgMiAxCjQgMiAyCjQgOAoyIDEgNAozIDIgMwo0IDEgNgoyIDQKMiAxIDQKMiA4CjIgMSA0CjQgOAoyIDEgNgozIDEgMgo0IDMgMwozIDEKMiAxIDEKMyAxIDEKNCAxCjIgMSAxCjMgMiAxCjQgMSAxCjIgMTAKMiAxIDQKNCA0CjIgMSAzCjMgMiAzCjQgMiAxCjQgMwoyIDEgMgozIDEgMQo0IDMgMwoyIDQKMiAxIDMKNCAzCjIgMSAzCjMgMSAxCjQgMSAyCjQgOQoyIDEgMQozIDEgMQo0IDIgMQoyIDcKMiAxIDEKNSA3CjIgMSA2CjMgMiA1CjQgMiAzCjUgMyA3CjMgOQoyIDEgOAozIDEgMgo1IDIKMiAxIDIKMyAyIDIKNCAxIDEKNSAyIDIKNCA3CjIgMSA2CjMgMiA3CjQgMiA0CjUgMQoyIDEgMQozIDEgMQo0IDIgMQo1IDIgMQo0IDEwCjIgMSA5CjMgMSA3CjQgMyAyCjQgOQoyIDEgNQozIDIgNwo0IDMgNwozIDEwCjIgMSA5CjMgMSA4CjUgMgoyIDEgMQozIDIgMgo0IDMgMgo1IDIgMQoyIDEwCjIgMSAyCjUgNgoyIDEgNAozIDEgMwo0IDIgMQo1IDEgNgo1IDYKMiAxIDYKMyAxIDEKNCAxIDIKNSAxIDIKMyA4CjIgMSA3CjMgMiA0CjMgNgoyIDEgNQozIDIgNAo0IDIKMiAxIDIKMyAxIDIKNCAxIDIKMiAxCjIgMSAxCjQgMQoyIDEgMQozIDIgMQo0IDMgMQo1IDcKMiAxIDUKMyAxIDYKNCAxIDQKNSAyIDIKNSA2CjIgMSAyCjMgMiA2CjQgMSAyCjUgMSAxCjIgMwoyIDEgMwo1IDMKMiAxIDEKMyAyIDIKNCAzIDMKNSA0IDMKMyA4CjIgMSA4CjMgMiA2CjMgMwoyIDEgMwozIDEgMwo1IDQKMiAxIDQKMyAyIDEKNCAxIDQKNSAzIDEKNSAyCjIgMSAyCjMgMiAxCjQgMyAxCjUgMiAyCjIgOQoyIDEgMQo0IDkKMiAxIDgKMyAyIDkKNCAyIDYKNSA2CjIgMSAzCjMgMSAxCjQgMyA1CjUgMyA2CjQgMQoyIDEgMQozIDEgMQo0IDMgMQo0IDgKMiAxIDMKMyAxIDUKNCAyIDMKNCA1CjIgMSAxCjMgMSA1CjQgMSA1CjMgOAoyIDEgMwozIDIgNAo1IDMKMiAxIDIKMyAxIDEKNCAzIDEKNSA0IDIKNSA0CjIgMSAyCjMgMSAzCjQgMSAzCjUgMSAxCjMgOAoyIDEgOAozIDIgOAozIDcKMiAxIDQKMyAyIDIKNSAxCjIgMSAxCjMgMiAxCjQgMyAxCjUgMSAxCjUgMTAKMiAxIDYKMyAyIDUKNCAxIDUKNSA0IDgKMyAxMAoyIDEgOQozIDIgMwo0IDkKMiAxIDEKMyAxIDEKNCAyIDIKMiA4CjIgMSA0CjQgNQoyIDEgMQozIDEgNQo0IDEgMQozIDEwCjIgMSAyCjMgMSAxMAozIDcKMiAxIDEKMyAxIDMKNCA4CjIgMSA4CjMgMiA3CjQgMyA2CjQgMwoyIDEgMgozIDEgMgo0IDIgMgoyIDUKMiAxIDIKNSAyCjIgMSAxCjMgMiAxCjQgMyAyCjUgMSAxCjUgNAoyIDEgNAozIDIgMQo0IDEgMQo1IDEgMwo1IDEKMiAxIDEKMyAxIDEKNCAzIDEKNSAzIDEKNCA3CjIgMSA1CjMgMiA0CjQgMiA0CjIgOQoyIDEgMQo1IDgKMiAxIDYKMyAyIDUKNCAyIDEKNSAxIDYKNCAyCjIgMSAyCjMgMSAxCjQgMyAxCjUgNAoyIDEgNAozIDEgMwo0IDIgMQo1IDQgMgozIDEKMiAxIDEKMyAxIDEKNSAyCjIgMSAyCjMgMiAxCjQgMSAyCjUgMSAyCjIgNwoyIDEgMgo1IDkKMiAxIDQKMyAxIDEKNCAzIDMKNSA0IDEKNSAxMAoyIDEgNAozIDEgMQo0IDEgMQo1IDMgMgo1IDEwCjIgMSAxMAozIDEgMgo0IDEgOAo1IDIgMgozIDcKMiAxIDcKMyAyIDMKMiAxMAoyIDEgNgozIDEKMiAxIDEKMyAxIDEKNCAzCjIgMSAxCjMgMiAzCjQgMSAyCjMgMTAKMiAxIDMKMyAxIDEwCjUgMwoyIDEgMgozIDEgMQo0IDEgMQo1IDEgMwo0IDcKMiAxIDYKMyAxIDIKNCAxIDYKNCA5CjIgMSA3CjMgMSAzCjQgMyAxCjMgOAoyIDEgNQozIDIgNwo="

const refSource = `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	u, v, w int
}

func power(base, exp int64) int64 {
	res := int64(1)
	base %= 998244353
	for exp > 0 {
		if exp%2 == 1 {
			res = (res * base) % 998244353
		}
		base = (base * base) % 998244353
		exp /= 2
	}
	return res
}

func find(parent []int, i int) int {
	root := i
	for root != parent[root] {
		root = parent[root]
	}
	curr := i
	for curr != root {
		nxt := parent[curr]
		parent[curr] = root
		curr = nxt
	}
	return root
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	buf := make([]byte, 1024*1024*10)
	scanner.Buffer(buf, len(buf))

	nextInt := func() int {
		scanner.Scan()
		res := 0
		for _, b := range scanner.Bytes() {
			res = res*10 + int(b-'0')
		}
		return res
	}

	if !scanner.Scan() {
		return
	}
	t := 0
	for _, b := range scanner.Bytes() {
		t = t*10 + int(b-'0')
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for i := 0; i < t; i++ {
		n := nextInt()
		S := nextInt()

		edges := make([]Edge, n-1)
		for j := 0; j < n-1; j++ {
			edges[j].u = nextInt()
			edges[j].v = nextInt()
			edges[j].w = nextInt()
		}

		sort.Slice(edges, func(a, b int) bool {
			return edges[a].w < edges[b].w
		})

		parent := make([]int, n+1)
		size := make([]int64, n+1)
		for j := 1; j <= n; j++ {
			parent[j] = j
			size[j] = 1
		}

		ans := int64(1)
		for j := 0; j < n-1; j++ {
			u := edges[j].u
			v := edges[j].v
			w := edges[j].w

			rootU := find(parent, u)
			rootV := find(parent, v)

			if rootU != rootV {
				pairs := size[rootU]*size[rootV] - 1
				ways := int64(S - w + 1)
				ans = (ans * power(ways, pairs)) % 998244353
				parent[rootV] = rootU
				size[rootU] += size[rootV]
			}
		}
		fmt.Fprintln(out, ans)
	}
}
`

func runCase(bin, input string) (string, error) {
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

func buildRef() (string, func()) {
	tmpDir, err := os.MkdirTemp("", "ref1857G")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	srcPath := filepath.Join(tmpDir, "ref.go")
	os.WriteFile(srcPath, []byte(refSource), 0644)
	binPath := filepath.Join(tmpDir, "ref")
	cmd := exec.Command("go", "build", "-o", binPath, srcPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build ref: %v\n%s\n", err, string(out))
		os.Exit(1)
	}
	return binPath, func() { os.RemoveAll(tmpDir) }
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refBin, cleanup := buildRef()
	defer cleanup()

	inputs := loadInputs()
	for idx, input := range inputs {
		refOut, err := runCase(refBin, input)
		if err != nil {
			fmt.Printf("reference failed on case %d: %v\n%s", idx+1, err, refOut)
			os.Exit(1)
		}
		got, err := runCase(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", idx+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(refOut) {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, refOut, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}

func loadInputs() []string {
	data, err := base64.StdEncoding.DecodeString(testcasesB64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode embedded testcases: %v\n", err)
		os.Exit(1)
	}
	tokens := strings.Fields(string(bytes.TrimSpace(data)))
	if len(tokens) == 0 {
		fmt.Fprintln(os.Stderr, "no testcases found")
		os.Exit(1)
	}
	t, err := strconv.Atoi(tokens[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid test count header\n")
		os.Exit(1)
	}
	pos := 1
	var inputs []string
	for caseNum := 1; caseNum <= t; caseNum++ {
		if pos+1 >= len(tokens) {
			fmt.Fprintf(os.Stderr, "case %d missing header\n", caseNum)
			os.Exit(1)
		}
		n, errN := strconv.Atoi(tokens[pos])
		S := tokens[pos+1]
		if errN != nil {
			fmt.Fprintf(os.Stderr, "invalid header on case %d\n", caseNum)
			os.Exit(1)
		}
		pos += 2
		edges := make([][3]string, n-1)
		if pos+3*(n-1) > len(tokens) {
			fmt.Fprintf(os.Stderr, "case %d missing edges\n", caseNum)
			os.Exit(1)
		}
		for i := 0; i < n-1; i++ {
			edges[i] = [3]string{tokens[pos+3*i], tokens[pos+3*i+1], tokens[pos+3*i+2]}
		}
		pos += 3 * (n - 1)
		lines := []string{"1", fmt.Sprintf("%d %s", n, S)}
		for _, e := range edges {
			lines = append(lines, fmt.Sprintf("%s %s %s", e[0], e[1], e[2]))
		}
		inputs = append(inputs, strings.Join(lines, "\n")+"\n")
	}
	return inputs
}
