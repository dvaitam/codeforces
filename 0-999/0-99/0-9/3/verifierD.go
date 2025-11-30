package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases (one per line: pattern followed by comma-separated cost pairs).
const embeddedTestcases = `?()))?() 2,5 7,4
))?) 1,8
?(((?)??() 2,6 9,1 4,7 1,9
?(?)()???)?)?)?(?? 4,6 4,8 2,3 5,2 6,9 4,5 8,9 1,8 7,7 3,6
?)?(?)?()?(?(?)?))?) 6,2 9,2 7,6 1,8 7,3 4,1 9,4 6,6 9,1
(?)())??)))(?)?????( 9,4 9,4 7,8 9,6 1,4 3,9 3,2 5,1 2,2
()()(?)((()???)(() 3,6 3,5 5,8 6,8
))()()??(()((??? 9,4 7,1 3,8 9,7 4,9
(??(?(()(())??(??? 1,7 6,7 5,3 3,7 5,3 1,4 8,3 9,1
()((?)(??)?((()? 7,4 7,5 8,1 7,5 3,6
()??)??(((?)?))) 2,7 6,9 4,2 1,2 4,5 9,5
)(??)( 8,3 9,2
()(()(???))?))(( 7,2 9,4 2,5 9,2
(?)((((?()(? 1,2 7,3 4,3
))?)?( 5,9 8,6
?(()?)))(( 6,1 6,8
?)(??(?(()()() 8,2 9,8 6,5 4,5
?))()( 6,4
?)(?((()()?(?) 5,4 2,4 2,2 1,5
)((??(())(?(((??() 6,2 3,3 9,5 1,6 9,4
??(??)?))())(?(? 3,1 4,5 8,7 5,9 8,1 7,1 6,3
(()))??)((??())???((?? 7,3 2,4 6,9 8,4 8,4 7,6 5,4 6,3 4,5
?)(?(?(())(? 5,9 8,2 9,7 1,8
)???()???()))())(?()?(?( 6,7 3,9 1,9 2,5 2,3 2,8 8,4 9,7 5,5
?(?(?(()(?)?(?)( 9,1 8,1 4,5 4,5 5,8 6,8
?)(((?)??) 7,4 1,9 3,2 6,5
?)(()?)???(??)??)?)? 6,9 8,6 7,6 8,2 7,7 1,5 9,4 9,7 5,3 9,4 1,7
)))????)((?)?)()))()?? 2,8 4,5 1,7 3,7 1,6 7,9 1,5 2,7
)(?(?) 8,1 3,2
?(?))((?)?()?)??))?)??)( 5,5 4,4 9,6 9,1 9,5 4,7 7,3 6,4 4,1 7,6
((??)?(?))(( 3,8 3,5 3,3 8,6
?(()))(((( 5,2
((??)??()()))((())?)?) 1,8 9,8 5,2 3,2 7,4 5,6
(((((()))?)(?()(??(()? 5,4 3,1 9,1 7,5 5,1
?(() 1,9
)((()? 3,5
(?)??(?)?()?(???((((()(( 8,7 5,5 4,4 3,6 9,9 7,9 9,7 2,5 2,5
(?)()((????())()()(??)(? 8,9 1,8 3,7 8,1 9,5 5,6 3,5 5,2
(?))?((??()?))?) 9,4 7,8 8,9 9,9 3,4 6,2
(?)?))))))?(?) 2,1 9,6 9,9 2,3
)?))))??(?(??) 2,8 7,2 1,3 8,5 6,6 6,7
??()?(()??()(?((?) 6,9 3,1 4,3 1,2 5,2 9,2 3,9
??(? 6,8 5,4 8,4
)?()?)((??)?())( 6,9 2,5 9,7 8,2 9,7
)?()????)) 1,9 6,9 9,5 7,9 7,5
????)))? 8,3 3,5 1,7 1,6 1,2
()))))
?))())()((())(?)??))(??( 8,6 5,7 5,7 5,7 8,7 7,2
(((?)()( 1,2
?()?(?(??)?)??(? 3,1 1,9 7,6 2,9 7,2 5,3 1,4 2,7 8,5
?(?()???)??((?(?((??(()( 8,7 8,2 4,3 5,7 9,5 9,4 6,8 6,5 8,5 5,5 4,7
??))))?)())?)(?? 1,9 9,3 5,5 8,4 3,8 3,1
)?)?)(((?(??(??)(()? 3,4 8,8 5,4 9,1 3,2 4,9 5,7 4,2
)??(())(() 6,5 9,7
?((?)()? 1,6 2,5 1,6
(())
??(?)?)(??))(?((?((())?? 4,2 6,5 6,2 3,5 9,8 7,9 5,9 9,2 1,5 5,9
)()?(???)?(? 2,6 9,6 9,9 1,5 3,3 3,4
)))(()))(?)?(??)() 3,5 1,9 3,7 9,2
?)))?()(?(??)??(?((??? 7,5 8,1 1,1 3,9 6,9 6,8 4,4 6,3 6,7 6,5 1,7
)))))?(()?(??)() 4,9 2,9 8,6 2,1
()?())??)?))?)((())((?)) 4,2 4,8 8,6 4,7 2,8 2,3
?(?())???(???()?()((??() 5,3 2,5 4,9 7,1 4,7 3,6 4,2 9,3 1,9 4,9 9,2
(?(?(?()))????)??) 1,7 2,8 4,1 7,3 3,9 6,9 8,9 3,7 7,4
)()??)())))) 5,3 2,6
?(()(((???)) 8,1 9,4 7,4 3,9
(???)??(() 2,3 1,1 7,7 3,3 9,9
)(??)?)? 7,6 3,4 3,6 5,2
)(?)?))?)((?)??(???? 8,1 2,6 1,1 2,2 4,4 9,7 4,7 9,3 5,1 2,4
)??(??)?))(??()()?(((? 8,9 4,5 3,9 4,7 7,7 3,3 1,8 9,6 1,7
((??)(()?()(?))(?? 9,4 7,9 2,6 3,3 9,2 7,2
?)))(?))??(?)?(? 5,7 8,1 6,3 9,5 6,7 1,2 9,1
)))??? 1,6 2,8 2,9
)?())((?))?))??))? 9,1 3,2 7,6 6,1 2,4 5,2
?)(()? 3,5 5,4
()??(?)))? 8,1 9,5 2,9 3,1
)?))?(?((((?(?(()) 1,1 3,9 9,7 3,9 5,8
??)(()?((??))( 6,4 1,1 1,1 3,1 4,4
?)????()((?)(???)? 6,6 2,4 3,4 5,7 8,6 7,6 7,4 9,8 9,9 8,8
)??)))(? 9,5 7,9 1,8
)(???)(?()?)()))?? 8,4 8,6 3,7 2,6 1,5 1,4 2,6
???(?? 2,3 5,7 6,4 3,9 6,5
))?()(???(?? 8,2 1,4 4,4 7,7 3,5 6,1
)))(?))?))??()(??))(((?? 3,1 6,9 2,5 5,7 8,2 4,5 9,9 3,3
((((??(((?(? 1,7 1,2 6,6 1,9
)())?)((() 9,9
(?)) 8,1
((??((()(()?)((( 3,9 2,3 3,2
)(?((?(((??? 2,8 3,2 2,2 4,5 5,9
(?())))?()())(?( 1,5 7,7 2,4
)?))))))))?()()? 5,6 1,6 3,9
()?((?)? 3,3 5,4 6,3
))()(?((???? 6,8 6,2 3,8 1,1 4,6
???)()?())?)(( 6,9 6,6 2,3 4,1 7,4
?((()?((()?? 1,4 5,3 9,9 9,7
?(?)()((((((())(() 9,8 6,4
?()?(? 4,4 4,6 1,4
)?))?)(??)))() 1,1 6,3 6,5 6,2
?()()) 6,1
)))()(?())?) 5,1 4,5`

type item struct {
	diff int
	idx  int
}

type minHeap []item

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].diff < h[j].diff }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(item)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

func solve(pattern string, costs [][2]int) (bool, int64, string) {
	n := len(pattern)
	if n%2 != 0 {
		return false, 0, ""
	}
	res := []byte(pattern)
	h := &minHeap{}
	heap.Init(h)
	var totalCost int64
	balance := 0
	qidx := 0
	for i := 0; i < n; i++ {
		ch := res[i]
		switch ch {
		case '(':
			balance++
		case ')':
			balance--
		case '?':
			if qidx >= len(costs) {
				return false, 0, ""
			}
			a := costs[qidx][0]
			b := costs[qidx][1]
			qidx++
			res[i] = ')'
			balance--
			totalCost += int64(b)
			heap.Push(h, item{diff: a - b, idx: i})
		}
		if balance < 0 {
			if h.Len() == 0 {
				return false, 0, ""
			}
			it := heap.Pop(h).(item)
			res[it.idx] = '('
			totalCost += int64(it.diff)
			balance += 2
		}
	}
	if balance != 0 {
		return false, 0, ""
	}
	return true, totalCost, string(res)
}

func runCandidate(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcases), "\n")
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		pattern := fields[0]
		costs := make([][2]int, 0, strings.Count(pattern, "?"))
		for _, f := range fields[1:] {
			parts := strings.Split(f, ",")
			if len(parts) != 2 {
				fmt.Fprintf(os.Stderr, "case %d: bad cost token\n", idx+1)
				os.Exit(1)
			}
			a, _ := strconv.Atoi(parts[0])
			b, _ := strconv.Atoi(parts[1])
			costs = append(costs, [2]int{a, b})
		}
		ok, expCost, expStr := solve(pattern, costs)
		var input strings.Builder
		input.WriteString(pattern)
		input.WriteByte('\n')
		for i, c := range costs {
			input.WriteString(fmt.Sprintf("%d %d", c[0], c[1]))
			if i+1 < len(costs) {
				input.WriteByte('\n')
			}
		}
		input.WriteByte('\n')
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if !ok {
			if strings.TrimSpace(got) != "-1" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected -1 got %s\n", idx+1, got)
				os.Exit(1)
			}
			continue
		}
		linesOut := strings.Split(strings.TrimSpace(got), "\n")
		if len(linesOut) < 2 {
			fmt.Fprintf(os.Stderr, "case %d failed: output should have 2 lines\n", idx+1)
			os.Exit(1)
		}
		gotCost, err := strconv.ParseInt(strings.TrimSpace(linesOut[0]), 10, 64)
		if err != nil || gotCost != expCost || strings.TrimSpace(linesOut[1]) != expStr {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%d\n%s\n\ngot:\n%s\n", idx+1, expCost, expStr, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
