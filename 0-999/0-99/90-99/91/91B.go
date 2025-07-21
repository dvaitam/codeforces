package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Fenwick tree for range maximum query on prefix [1..i]
type Fenwick struct {
   n    int
   data []int
}

// NewFenwick creates a Fenwick tree of size n (1-based)
func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, data: make([]int, n+1)}
}

// Update position i to max(current, v)
func (f *Fenwick) Update(i, v int) {
   for ; i <= f.n; i += i & -i {
       if f.data[i] < v {
           f.data[i] = v
       }
   }
}

// Query returns max over data[1..i]
func (f *Fenwick) Query(i int) int {
   res := 0
   for ; i > 0; i -= i & -i {
       if f.data[i] > res {
           res = f.data[i]
       }
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // coordinate compression
   vals := make([]int, n)
   copy(vals, a)
   sort.Ints(vals)
   m := 0
   for i := 0; i < n; i++ {
       if i == 0 || vals[i] != vals[i-1] {
           vals[m] = vals[i]
           m++
       }
   }
   vals = vals[:m]
   // Fenwick tree stores max(position+1) among seen
   bit := NewFenwick(m)
   ans := make([]int, n)
   // sweep from right to left
   for i := n - 1; i >= 0; i-- {
       // rank of a[i]
       r := sort.SearchInts(vals, a[i]) // 0-based
       // values strictly less than a[i] have ranks [0..r-1]
       if r > 0 {
           j1 := bit.Query(r) // max of pos+1
           if j1 > 0 {
               j := j1 - 1
               ans[i] = j - i - 1
           } else {
               ans[i] = -1
           }
       } else {
           ans[i] = -1
       }
       // update current position
       bit.Update(r+1, i+1)
   }
   // output
   for i, v := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       if v < 0 {
           writer.WriteString("-1")
       } else {
           writer.WriteString(fmt.Sprint(v))
       }
   }
   writer.WriteByte('\n')
}
