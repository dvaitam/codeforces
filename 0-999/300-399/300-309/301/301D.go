package main

import (
   "bufio"
   "os"
   "strconv"
   "strings"
)

// BIT implements a Fenwick tree for prefix sums on 1..n
type BIT struct {
   n    int
   tree []int
}

// NewBIT creates a BIT for size n
func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int, n+1)}
}

// Add adds value v at position i
func (b *BIT) Add(i, v int) {
   for ; i <= b.n; i += i & -i {
       b.tree[i] += v
   }
}

// Sum returns prefix sum up to i
func (b *BIT) Sum(i int) int {
   if i <= 0 {
       return 0
   }
   s := 0
   for ; i > 0; i -= i & -i {
       s += b.tree[i]
   }
   return s
}

// Query represents a sub-query for inclusion-exclusion
type Query struct {
   y     int
   id    int
   coeff int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   // read n and m
   line, _ := in.ReadBytes('\n')
   parts := bytesFields(line)
   n, _ := strconv.Atoi(parts[0])
   m, _ := strconv.Atoi(parts[1])

   // read permutation p
   line, _ = in.ReadBytes('\n')
   parts = bytesFields(line)
   pos := make([]int, n+1)
   for i := 1; i <= n; i++ {
       v, _ := strconv.Atoi(parts[i-1])
       pos[v] = i
   }

   // group edges by u
   edgesByU := make([][]int, n+1)
   for v := 1; v <= n; v++ {
       u := pos[v]
       for k := v + v; k <= n; k += v {
           edgesByU[u] = append(edgesByU[u], pos[k])
       }
   }

   // read queries and prepare sub-queries
   subsByX := make([][]Query, n+1)
   L := make([]int, m)
   R := make([]int, m)
   for i := 0; i < m; i++ {
       line, _ = in.ReadBytes('\n')
       parts = bytesFields(line)
       l, _ := strconv.Atoi(parts[0])
       r, _ := strconv.Atoi(parts[1])
       L[i], R[i] = l, r
       // inclusion-exclusion on rectangle [l,r]x[l,r]
       subsByX[r] = append(subsByX[r], Query{y: r, id: i, coeff: 1})
       if l > 0 {
           subsByX[l-1] = append(subsByX[l-1], Query{y: r, id: i, coeff: -1})
           subsByX[r] = append(subsByX[r], Query{y: l - 1, id: i, coeff: -1})
           subsByX[l-1] = append(subsByX[l-1], Query{y: l - 1, id: i, coeff: 1})
       }
   }

   bit := NewBIT(n)
   ansEdges := make([]int64, m)
   // sweep x from 0 to n
   for x := 0; x <= n; x++ {
       if x > 0 {
           for _, y := range edgesByU[x] {
               bit.Add(y, 1)
           }
       }
       for _, q := range subsByX[x] {
           cnt := bit.Sum(q.y)
           ansEdges[q.id] += int64(q.coeff) * int64(cnt)
       }
   }

   // output answers
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i := 0; i < m; i++ {
       // add self-pairs
       total := ansEdges[i] + int64(R[i]-L[i]+1)
       out.WriteString(strconv.FormatInt(total, 10))
       out.WriteByte('\n')
   }
}

// bytesFields splits a byte slice by spaces
func bytesFields(b []byte) []string {
   s := string(b)
   return strings.Fields(s)
}
