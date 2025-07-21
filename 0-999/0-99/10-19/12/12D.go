package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Lady represents a participant with beauty B, intellect rank I, and richness R.
type Lady struct {
   B int
   I int
   R int
}

// BIT supports point update and prefix maximum query.
type BIT struct {
   n   int
   bit []int
}

// NewBIT creates a BIT for size n.
func NewBIT(n int) *BIT {
   return &BIT{n, make([]int, n+1)}
}

// Update sets position i to max(current, v).
func (b *BIT) Update(i int, v int) {
   for ; i <= b.n; i += i & -i {
       if b.bit[i] < v {
           b.bit[i] = v
       }
   }
}

// Query returns max over prefix [1..i].
func (b *BIT) Query(i int) int {
   res := 0
   for ; i > 0; i -= i & -i {
       if b.bit[i] > res {
           res = b.bit[i]
       }
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var N int
   fmt.Fscan(in, &N)
   B := make([]int, N)
   I := make([]int, N)
   R := make([]int, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(in, &B[i])
   }
   for i := 0; i < N; i++ {
       fmt.Fscan(in, &I[i])
   }
   for i := 0; i < N; i++ {
       fmt.Fscan(in, &R[i])
   }

   // Compress intellect values
   uniq := make([]int, N)
   copy(uniq, I)
   sort.Ints(uniq)
   rank := make(map[int]int, N)
   m := 1
   for _, v := range uniq {
       if _, ok := rank[v]; !ok {
           rank[v] = m
           m++
       }
   }
   M := m - 1

   ladies := make([]Lady, N)
   for i := 0; i < N; i++ {
       ladies[i] = Lady{B: B[i], I: rank[I[i]], R: R[i]}
   }

   // Sort by beauty descending
   sort.Slice(ladies, func(i, j int) bool {
       return ladies[i].B > ladies[j].B
   })

   bit := NewBIT(M)
   var cnt int64
   // Process in groups of equal beauty
   for i := 0; i < N; {
       j := i
       for j < N && ladies[j].B == ladies[i].B {
           j++
       }
       // Query dominance for this group
       for k := i; k < j; k++ {
           rnk := ladies[k].I
           // number of intellects strictly greater: ranks rnk+1..M
           t := M - rnk
           if t > 0 {
               if bit.Query(t) > ladies[k].R {
                   cnt++
               }
           }
       }
       // Update BIT with this group's richness
       for k := i; k < j; k++ {
           rnk := ladies[k].I
           idx := M - rnk + 1
           bit.Update(idx, ladies[k].R)
       }
       i = j
   }

   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprintln(out, cnt)
}
