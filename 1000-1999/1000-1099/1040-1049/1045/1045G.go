package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Robot represents a robot with position x, range r, and intelligence iq
type Robot struct { x, r, iq int }

// Fenwick implements a Fenwick (BIT) tree for prefix sums
type Fenwick struct { n int; a []int }

// NewFenwick creates a Fenwick tree of size n
func NewFenwick(n int) Fenwick { return Fenwick{n, make([]int, n)} }

// Add adds val at index pos
func (f *Fenwick) Add(pos, val int) {
   for pos < f.n {
       f.a[pos] += val
       pos |= pos + 1
   }
}

// Sum returns the prefix sum [0..pos]
func (f *Fenwick) Sum(pos int) int {
   res := 0
   for pos >= 0 {
       res += f.a[pos]
       pos = (pos&(pos+1)) - 1
   }
   return res
}

// Query returns the sum in [l..r)
func (f *Fenwick) Query(l, r int) int {
   return f.Sum(r-1) - f.Sum(l-1)
}

// unique removes duplicates from sorted slice a
func unique(a []int) []int {
   if len(a) == 0 {
       return a
   }
   j := 1
   for i := 1; i < len(a); i++ {
       if a[i] != a[i-1] {
           a[j] = a[i]
           j++
       }
   }
   return a[:j]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   fmt.Fscan(reader, &n, &k)
   robots := make([]Robot, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &robots[i].x, &robots[i].r, &robots[i].iq)
   }

   // Compress intelligence values
   iqs := make([]int, n)
   for i, rob := range robots {
       iqs[i] = rob.iq
   }
   sort.Ints(iqs)
   iqs = unique(iqs)
   m := len(iqs)

   // Prepare positions per intelligence group
   positions := make([][]int, m)
   for _, rob := range robots {
       idx := sort.SearchInts(iqs, rob.iq)
       positions[idx] = append(positions[idx], rob.x)
   }
   for i := 0; i < m; i++ {
       sort.Ints(positions[i])
       positions[i] = unique(positions[i])
   }

   // Initialize Fenwick trees
   fenwicks := make([]Fenwick, m)
   for i := 0; i < m; i++ {
       fenwicks[i] = NewFenwick(len(positions[i]))
   }

   // Sort robots by decreasing range
   sort.Slice(robots, func(i, j int) bool {
       return robots[i].r > robots[j].r
   })

   var ans int64
   for _, rob := range robots {
       lowerIq := rob.iq - k
       upperIq := rob.iq + k
       // find IQ groups within [lowerIq..upperIq]
       it := sort.Search(m, func(i int) bool { return iqs[i] >= lowerIq })
       for it < m && iqs[it] <= upperIq {
           xs := positions[it]
           left := sort.Search(len(xs), func(i int) bool { return xs[i] >= rob.x-rob.r })
           right := sort.Search(len(xs), func(i int) bool { return xs[i] > rob.x+rob.r })
           ans += int64(fenwicks[it].Query(left, right))
           it++
       }
       // add current robot to its IQ group
       idx := sort.SearchInts(iqs, rob.iq)
       pos := sort.Search(len(positions[idx]), func(i int) bool { return positions[idx][i] >= rob.x })
       fenwicks[idx].Add(pos, 1)
   }

   fmt.Fprintln(writer, ans)
}
