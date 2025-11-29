package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded source for the reference solution (was 1054C.go).
const solutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

type node struct {
   w, id int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   L := make([]int, n)
   R := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &L[i])
       if L[i] > i {
           fmt.Fprintln(writer, "NO")
           return
       }
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &R[i])
       if R[i] > n-1-i {
           fmt.Fprintln(writer, "NO")
           return
       }
   }
   seq := make([]int, n)
   cnt := n
   var Q []node
   // initial positions with zero constraints
   for i := 0; i < n; i++ {
       if L[i] == 0 && R[i] == 0 {
           seq[i] = cnt
           Q = append(Q, node{cnt, i})
       }
   }
   // assign values
   for len(Q) > 0 {
       // process current layer
       for len(Q) > 0 {
           x := Q[0]
           Q = Q[1:]
           for j := x.id + 1; j < n; j++ {
               if seq[j] != 0 {
                   continue
               }
               L[j]--
           }
           for j := 0; j < x.id; j++ {
               if seq[j] != 0 {
                   continue
               }
               R[j]--
           }
       }
       cnt--
       // find new zeros
       for i := 0; i < n; i++ {
           if seq[i] == 0 && L[i] == 0 && R[i] == 0 {
               seq[i] = cnt
               Q = append(Q, node{cnt, i})
           }
       }
   }
   // check
   for i := 0; i < n; i++ {
       if seq[i] == 0 {
           fmt.Fprintln(writer, "NO")
           return
       }
   }
   // output
   fmt.Fprintln(writer, "YES")
   for i := 0; i < n; i++ {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, seq[i])
   }
   writer.WriteByte('\n')
}
`

const testcasesRaw = `1 0 0
6 0 0 2 2 1 5 5 4 2 0 0 0
5 0 1 0 2 3 2 1 0 0 0
3 0 0 1 0 0 0
2 0 1 0 0
6 0 0 1 2 4 4 1 3 0 2 0 0
4 0 1 0 1 1 2 1 0
5 0 0 1 1 4 3 1 1 1 0
1 0 0
3 0 1 1 2 0 0
3 0 1 0 1 0 0
3 0 1 2 0 0 0
1 0 0
6 0 0 0 2 3 1 1 3 3 0 1 0
2 0 1 1 0
4 0 1 0 0 3 1 0 0
5 0 0 0 2 1 1 1 0 0 0
6 0 1 2 0 4 5 3 1 0 1 1 0
1 0 0
1 0 0
6 0 1 0 2 1 5 3 2 2 1 1 0
3 0 1 0 2 1 0
4 0 0 1 2 3 0 0 0
6 0 1 2 2 1 4 4 2 3 1 0 0
5 0 0 2 2 1 1 2 2 0 0
5 0 1 1 2 2 2 3 1 1 0
1 0 0
4 0 1 0 3 2 2 0 0
5 0 1 0 1 0 1 0 2 1 0
3 0 1 0 1 1 0
5 0 1 0 3 1 0 1 2 1 0
4 0 1 2 1 0 1 0 0
6 0 1 0 1 4 1 2 1 3 0 1 0
6 0 0 0 1 3 0 3 3 2 1 1 0
5 0 0 1 0 3 0 2 0 1 0
4 0 0 2 3 2 2 0 0
3 0 1 1 2 0 0
1 0 0
1 0 0
6 0 1 0 2 0 2 3 1 2 1 0 0
4 0 0 0 1 1 1 0 0
5 0 0 1 0 0 4 1 1 0 0
3 0 1 2 1 0 0
5 0 0 2 1 4 0 2 0 0 0
5 0 1 2 1 2 1 3 2 0 0
5 0 0 0 0 3 4 2 2 1 0
4 0 0 0 2 0 0 0 0
3 0 0 1 2 0 0
4 0 0 1 0 3 0 1 0
5 0 1 2 3 3 1 0 1 1 0
4 0 1 0 0 0 1 0 0
4 0 1 1 3 3 0 0 0
3 0 0 0 1 1 0
4 0 0 2 2 0 0 1 0
2 0 0 0 0
6 0 1 0 3 3 0 1 1 2 1 0 0
3 0 0 1 1 1 0
4 0 0 2 0 3 2 0 0
1 0 0
5 0 0 0 0 4 4 1 0 1 0
4 0 1 2 1 1 1 0 0
1 0 0
5 0 0 1 3 1 3 3 0 1 0
6 0 1 2 0 2 2 4 2 0 2 1 0
3 0 1 1 2 0 0
2 0 0 1 0
4 0 1 2 3 1 0 1 0
4 0 1 2 0 2 1 0 0
4 0 0 1 1 2 1 1 0
3 0 0 2 1 1 0
6 0 1 0 1 0 4 2 1 0 1 1 0
1 0 0
5 0 0 0 1 1 1 0 2 0 0
3 0 1 2 0 1 0
2 0 1 0 0
6 0 0 1 2 3 3 4 2 3 1 0 0
2 0 0 0 0
3 0 1 0 0 0 0
3 0 1 0 2 1 0
5 0 1 0 1 4 4 0 2 0 0
6 0 1 2 1 3 5 5 3 2 0 0 0
2 0 0 1 0
6 0 1 0 3 0 3 4 0 2 2 1 0
4 0 0 0 2 1 1 0 0
1 0 0
4 0 0 1 3 1 0 1 0
3 0 0 2 0 0 0
1 0 0
5 0 0 1 2 2 0 0 0 0 0
3 0 0 0 1 1 0
1 0 0
5 0 0 2 2 3 1 1 1 0 0
6 0 1 1 1 0 3 2 4 3 0 1 0
4 0 1 1 2 0 1 1 0
4 0 1 0 1 2 1 0 0
6 0 1 1 3 4 4 1 3 2 1 0 0
5 0 0 1 2 3 1 3 0 0 0
1 0 0
6 0 0 0 1 2 0 1 0 0 1 0 0
1 0 0`

var _ = solutionSource

func computeExpected(n int, L, R []int) string {
	for i := 0; i < n; i++ {
		if L[i] > i || R[i] > n-1-i {
			return "NO"
		}
	}
	seq := make([]int, n)
	cnt := n
	var queue []int
	for i := 0; i < n; i++ {
		if L[i] == 0 && R[i] == 0 {
			seq[i] = cnt
			queue = append(queue, i)
		}
	}
	for len(queue) > 0 {
		for len(queue) > 0 {
			id := queue[0]
			queue = queue[1:]
			for j := id + 1; j < n; j++ {
				if seq[j] == 0 {
					L[j]--
				}
			}
			for j := 0; j < id; j++ {
				if seq[j] == 0 {
					R[j]--
				}
			}
		}
		cnt--
		for i := 0; i < n; i++ {
			if seq[i] == 0 && L[i] == 0 && R[i] == 0 {
				seq[i] = cnt
				queue = append(queue, i)
			}
		}
	}
	for i := 0; i < n; i++ {
		if seq[i] == 0 {
			return "NO"
		}
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(seq[i]))
	}
	return sb.String()
}

func main() {
	var _ = solutionSource
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil || len(fields) != 1+2*n {
			fmt.Fprintf(os.Stderr, "invalid test format at case %d\n", idx)
			os.Exit(1)
		}
		L := make([]int, n)
		R := make([]int, n)
		for i := 0; i < n; i++ {
			L[i], _ = strconv.Atoi(fields[1+i])
		}
		for i := 0; i < n; i++ {
			R[i], _ = strconv.Atoi(fields[1+n+i])
		}
		expected := strings.TrimSpace(computeExpected(n, append([]int(nil), L...), append([]int(nil), R...)))
		input := line + "\n"

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed. Expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
