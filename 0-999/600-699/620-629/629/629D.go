package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

// Cake represents a cake with its key values.
type Cake struct {
   idx  int
   r2h  float64
}

// BIT implements a Fenwick tree for range maximum query.
type BIT struct {
   n    int
   tree []float64
}

// NewBIT creates a BIT of size n.
func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]float64, n+1)}
}

// Update sets tree[i] = max(tree[i], val) and propagates.
func (b *BIT) Update(i int, val float64) {
   for ; i <= b.n; i += i & -i {
       if b.tree[i] < val {
           b.tree[i] = val
       }
   }
}

// Query returns max value in range [1..i].
func (b *BIT) Query(i int) float64 {
   var ret float64
   for ; i > 0; i -= i & -i {
       if ret < b.tree[i] {
           ret = b.tree[i]
       }
   }
   return ret
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   cakes := make([]Cake, n)
   for i := 0; i < n; i++ {
       var r, h int
       fmt.Fscan(reader, &r, &h)
       cakes[i].idx = n - i
       cakes[i].r2h = float64(r)*float64(r)*float64(h)
   }
   // sort by descending r2h
   sort.Slice(cakes, func(i, j int) bool {
       return cakes[i].r2h > cakes[j].r2h
   })

   bit := NewBIT(n)
   v := make([]float64, n)
   var ans float64
   last := 0
   for i := 0; i < n; i++ {
       v[i] = cakes[i].r2h + bit.Query(cakes[i].idx-1)
       // delay updates for equal r2h
       if i < n-1 && cakes[i].r2h == cakes[i+1].r2h {
           // do nothing
       } else {
           for j := last; j <= i; j++ {
               bit.Update(cakes[j].idx, v[j])
           }
           last = i + 1
       }
       if v[i] > ans {
           ans = v[i]
       }
   }
   fmt.Fprintf(writer, "%.9f", ans*math.Pi)
}
