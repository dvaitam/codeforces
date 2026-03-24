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

const refSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

func nextInt(reader *bufio.Reader) int {
	var res int
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return 0
		}
		if b >= '0' && b <= '9' {
			res = int(b - '0')
			break
		}
	}
	for {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}
		if b >= '0' && b <= '9' {
			res = res*10 + int(b-'0')
		} else {
			break
		}
	}
	return res
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 65536)
	writer := bufio.NewWriterSize(os.Stdout, 65536)
	defer writer.Flush()

	n := nextInt(reader)
	if n == 0 {
		return
	}

	MOD := int64(998244353)

	type Edge struct{ u, v int }
	maxW := n * (n - 1) / 2
	edgeList := make([]Edge, maxW+1)

	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			w := nextInt(reader)
			if i < j {
				edgeList[w] = Edge{i, j}
			}
		}
	}

	parent := make([]int, n+1)
	size := make([]int, n+1)
	edges := make([]int, n+1)
	children := make([][]int, n+1)
	poly := make([][]int, 2*n+1)

	for i := 1; i <= n; i++ {
		parent[i] = i
		size[i] = 1
		edges[i] = 0
		children[i] = []int{i}
		poly[i] = []int{0, 1}
	}

	var find func(i int) int
	find = func(i int) int {
		if parent[i] == i {
			return i
		}
		parent[i] = find(parent[i])
		return parent[i]
	}

	nodeCount := n

	for w := 1; w <= maxW; w++ {
		u := edgeList[w].u
		v := edgeList[w].v

		ru := find(u)
		rv := find(v)

		var r int
		if ru != rv {
			if size[ru] < size[rv] {
				ru, rv = rv, ru
			}
			parent[rv] = ru
			size[ru] += size[rv]
			edges[ru] += edges[rv] + 1
			children[ru] = append(children[ru], children[rv]...)
			r = ru
		} else {
			edges[ru]++
			r = ru
		}

		if edges[r] == size[r]*(size[r]-1)/2 {
			nodeCount++
			G := nodeCount

			res := []int{1}
			for _, c := range children[r] {
				pc := poly[c]
				newRes := make([]int, len(res)+len(pc)-1)
				for i := 0; i < len(res); i++ {
					if res[i] == 0 {
						continue
					}
					for j := 0; j < len(pc); j++ {
						newRes[i+j] = int((int64(newRes[i+j]) + int64(res[i])*int64(pc[j])) % MOD)
					}
				}
				res = newRes
			}

			if len(res) < 2 {
				temp := make([]int, 2)
				copy(temp, res)
				res = temp
			}
			res[1] = int((int64(res[1]) + 1) % MOD)

			poly[G] = res
			children[r] = []int{G}
		}
	}

	ans := poly[nodeCount]
	for i := 1; i <= n; i++ {
		if i < len(ans) {
			fmt.Fprintf(writer, "%d ", ans[i])
		} else {
			fmt.Fprintf(writer, "0 ")
		}
	}
	fmt.Fprintln(writer)
}
`

func runProgram(bin, input string) (string, error) {
	if _, err := os.Stat(bin); err == nil && !strings.Contains(bin, "/") {
		bin = "./" + bin
	}
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildReferenceBinary() (string, error) {
	srcFile, err := os.CreateTemp("", "cf-1408G-src-*.go")
	if err != nil {
		return "", err
	}
	if _, err := srcFile.WriteString(refSource); err != nil {
		srcFile.Close()
		os.Remove(srcFile.Name())
		return "", err
	}
	srcFile.Close()
	defer os.Remove(srcFile.Name())

	tmp, err := os.CreateTemp("", "cf-1408G-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), srcFile.Name())
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func randomCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 2
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if i == j {
				mat[i][j] = 0
			} else {
				mat[i][j] = rng.Intn(5)
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(mat[i][j]))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refBin, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input := randomCase(rng)
		want, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
			os.Exit(1)
		}
		out, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(want) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", t+1, want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
