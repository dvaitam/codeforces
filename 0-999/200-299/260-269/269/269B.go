package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

// Fenwick tree (BIT) for prefix maximum
type Fenwick struct {
   n    int
   tree []int
}

func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, tree: make([]int, n+1)}
}

// Update position i with value v (keep max)
func (f *Fenwick) Update(i, v int) {
   for i <= f.n {
       if f.tree[i] < v {
           f.tree[i] = v
       }
       i += i & -i
   }
}

// Query returns max on prefix [1..i]
func (f *Fenwick) Query(i int) int {
   mx := 0
   for i > 0 {
       if f.tree[i] > mx {
           mx = f.tree[i]
       }
       i -= i & -i
   }
   return mx
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // fast scanning
   buf := bufio.NewScanner(reader)
   buf.Split(bufio.ScanWords)
   // read n, m
   buf.Scan()
   n, _ := strconv.Atoi(buf.Text())
   buf.Scan()
   m, _ := strconv.Atoi(buf.Text())

   fenw := NewFenwick(m)
   best := 0
   // read plants
   for i := 0; i < n; i++ {
       buf.Scan()
       s, _ := strconv.Atoi(buf.Text())
       buf.Scan() // skip position
       // compute longest non-decreasing subseq ending here
       cur := fenw.Query(s) + 1
       if cur > best {
           best = cur
       }
       fenw.Update(s, cur)
   }
   // minimal replant = total - longest kept
   fmt.Fprintln(writer, n-best)
}
