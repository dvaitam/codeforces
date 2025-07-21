package main

import (
   "bufio"
   "fmt"
   "os"
)

// BIT supports prefix max query and point update
type BIT struct {
   n    int
   tree []int
}

// NewBIT creates a BIT for indices 1..n
func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int, n+1)}
}

// Update updates index i to value v if v is greater
func (b *BIT) Update(i, v int) {
   for ; i <= b.n; i += i & -i {
       if b.tree[i] < v {
           b.tree[i] = v
       }
   }
}

// Query returns max value in prefix 1..i
func (b *BIT) Query(i int) int {
   m := 0
   for ; i > 0; i -= i & -i {
       if b.tree[i] > m {
           m = b.tree[i]
       }
   }
   return m
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // read n and sequence
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   maxVal := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] > maxVal {
           maxVal = a[i]
       }
   }

   // dp1: LIS ending at i
   bit1 := NewBIT(maxVal + 2)
   dp1 := make([]int, n)
   L := 0
   for i := 0; i < n; i++ {
       v := a[i]
       // query max for values < v
       best := 0
       if v > 1 {
           best = bit1.Query(v - 1)
       }
       dp1[i] = best + 1
       if dp1[i] > L {
           L = dp1[i]
       }
       bit1.Update(v, dp1[i])
   }

   // dp2: LIS starting at i
   bit2 := NewBIT(maxVal + 2)
   dp2 := make([]int, n)
   for i := n - 1; i >= 0; i-- {
       v := a[i]
       ra := maxVal - v + 1
       best := 0
       if ra > 1 {
           best = bit2.Query(ra - 1)
       }
       dp2[i] = best + 1
       bit2.Update(ra, dp2[i])
   }

   // count positions at each layer
   cnt := make([]int, L+1)
   for i := 0; i < n; i++ {
       if dp1[i]+dp2[i]-1 == L {
           cnt[dp1[i]]++
       }
   }

   // build result
   res := make([]byte, n)
   for i := 0; i < n; i++ {
       if dp1[i]+dp2[i]-1 < L {
           res[i] = '1'
       } else if cnt[dp1[i]] == 1 {
           res[i] = '3'
       } else {
           res[i] = '2'
       }
   }
   fmt.Fprint(writer, string(res))
}
