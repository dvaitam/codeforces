package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// Embedded solver for 1868D
func solveD(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	var t int
	fmt.Fscan(reader, &t)

	type Node struct {
		id  int
		deg int
	}

	for tc := 0; tc < t; tc++ {
		var n int
		fmt.Fscan(reader, &n)

		d := make([]int, n)
		sumD := 0
		nodes := make([]Node, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &d[i])
			sumD += d[i]
			nodes[i] = Node{id: i + 1, deg: d[i]}
		}

		if sumD != 2*n {
			fmt.Fprintln(writer, "No")
			continue
		}

		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].deg > nodes[j].deg
		})

		M := 0
		M3 := 0
		for i := 0; i < n; i++ {
			if nodes[i].deg >= 2 {
				M++
			}
			if nodes[i].deg >= 3 {
				M3++
			}
		}

		if M == n {
			fmt.Fprintln(writer, "Yes")
			for i := 0; i < n; i++ {
				fmt.Fprintln(writer, nodes[i].id, nodes[(i+1)%n].id)
			}
			continue
		}

		prefixD := make([]int, M+1)
		for i := 0; i < M; i++ {
			prefixD[i+1] = prefixD[i] + nodes[i].deg
		}

		foundC := -1
		foundK := -1

		for C := 2; C <= M3; C++ {
			k := M / C
			if k == 0 {
				continue
			}
			if k == 1 {
				if M == C {
					foundC = C
					foundK = k
					break
				}
				continue
			}

			idx := C
			OPrev := prefixD[C] - 2*C
			possible := true
			rem := M - k*C

			for step := 1; step < k; step++ {
				if OPrev < C {
					possible = false
					break
				}
				canAdd := OPrev - C
				add := rem
				if canAdd < add {
					add = canAdd
				}
				LSize := C + add
				rem -= add

				OCurr := prefixD[idx+LSize] - prefixD[idx] - LSize
				OPrev = OCurr
				idx += LSize
			}

			if possible && rem == 0 {
				foundC = C
				foundK = k
				break
			}
		}

		if foundC == -1 {
			fmt.Fprintln(writer, "No")
		} else {
			fmt.Fprintln(writer, "Yes")
			LSizes := make([]int, foundK)
			for i := range LSizes {
				LSizes[i] = foundC
			}
			rem := M - foundK*foundC
			idx := foundC
			OPrev := prefixD[foundC] - 2*foundC

			for step := 1; step < foundK; step++ {
				canAdd := OPrev - foundC
				add := rem
				if canAdd < add {
					add = canAdd
				}
				LSizes[step] += add
				rem -= add
				OCurr := prefixD[idx+LSizes[step]] - prefixD[idx] - LSizes[step]
				OPrev = OCurr
				idx += LSizes[step]
			}

			outDeg := make([]int, M)
			for i := 0; i < foundC; i++ {
				outDeg[i] = nodes[i].deg - 2
			}
			for i := foundC; i < M; i++ {
				outDeg[i] = nodes[i].deg - 1
			}

			for i := 0; i < foundC; i++ {
				fmt.Fprintln(writer, nodes[i].id, nodes[(i+1)%foundC].id)
			}

			layers := make([][]int, foundK)
			nodeIdx := 0
			for i, sz := range LSizes {
				layers[i] = make([]int, sz)
				for j := 0; j < sz; j++ {
					layers[i][j] = nodeIdx
					nodeIdx++
				}
			}

			for i := 0; i < foundK-1; i++ {
				prevLayer := layers[i]
				currLayer := layers[i+1]

				for j := 0; j < foundC; j++ {
					fmt.Fprintln(writer, nodes[prevLayer[j]].id, nodes[currLayer[j]].id)
					outDeg[prevLayer[j]] -= 1
				}

				parentIdx := 0
				for j := foundC; j < len(currLayer); j++ {
					for parentIdx < len(prevLayer) && outDeg[prevLayer[parentIdx]] == 0 {
						parentIdx++
					}
					p := prevLayer[parentIdx]
					fmt.Fprintln(writer, nodes[p].id, nodes[currLayer[j]].id)
					outDeg[p] -= 1
				}
			}

			leafIdx := M
			for i := 0; i < M; i++ {
				for outDeg[i] > 0 {
					fmt.Fprintln(writer, nodes[i].id, nodes[leafIdx].id)
					outDeg[i] -= 1
					leafIdx++
				}
			}
		}
	}

	writer.Flush()
	return strings.TrimSpace(buf.String())
}

type Case struct{ input string }

func genCases() []Case {
	rng := rand.New(rand.NewSource(1868_04))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(5) + 2
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", rng.Intn(4)))
		}
		sb.WriteByte('\n')
		cases[i] = Case{sb.String()}
	}
	return cases
}

func runBinary(bin, input string) (string, error) {
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

func runCase(bin string, c Case) error {
	expected := solveD(c.input)
	got, err := runBinary(bin, c.input)
	if err != nil {
		return err
	}
	if expected != strings.TrimSpace(got) {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	var bin string
	if len(os.Args) == 2 {
		bin = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		bin = os.Args[2]
	} else {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	cases := genCases()
	for i, c := range cases {
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
