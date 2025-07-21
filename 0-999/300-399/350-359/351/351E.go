package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Fenwick tree for int counts
type Fenwick struct {
   n    int
   tree []int
}

func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, tree: make([]int, n+1)}
}

// Add v at position i (1-indexed)
func (f *Fenwick) Update(i, v int) {
   for x := i; x <= f.n; x += x & -x {
       f.tree[x] += v
   }
}

// Query sum [1..i]
func (f *Fenwick) Query(i int) int {
   s := 0
   for x := i; x > 0; x -= x & -x {
       s += f.tree[x]
   }
   return s
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   p := make([]int, n)
   absv := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &p[i])
       if p[i] < 0 {
           absv[i] = -p[i]
       } else {
           absv[i] = p[i]
       }
   }
   // Coordinate compress abs values
   vals := make([]int, n)
   copy(vals, absv)
   sort.Ints(vals)
   uniq := vals[:0]
   for _, v := range vals {
       if len(uniq) == 0 || uniq[len(uniq)-1] != v {
           uniq = append(uniq, v)
       }
   }
   m := len(uniq)
   rank := make(map[int]int, m)
   for i, v := range uniq {
       rank[v] = i + 1 // 1-indexed
   }
   // Compute cnt_less_before a[i]
   a := make([]int64, n)
   bit := NewFenwick(m)
   for i := 0; i < n; i++ {
       r := rank[absv[i]]
       if r > 1 {
           a[i] = int64(bit.Query(r - 1))
       }
       bit.Update(r, 1)
   }
   // Compute cnt_smaller_after b[i]
   b := make([]int64, n)
   bit = NewFenwick(m)
   for i := n - 1; i >= 0; i-- {
       r := rank[absv[i]]
       if r > 1 {
           b[i] = int64(bit.Query(r - 1))
       }
       bit.Update(r, 1)
   }
   // Group indices by equal abs
   groups := make(map[int][]int, m)
   for i := 0; i < n; i++ {
       r := rank[absv[i]]
       groups[r] = append(groups[r], i)
   }
   var ans int64
   const INF = int64(4e18)
   // Process each group independently
   for _, idxs := range groups {
       // DP over positions in group: dp[ones] = min cost
       dp := make([]int64, 1)
       dp[0] = 0
       // iterate through group in order
       for _, idx := range idxs {
           ci0 := a[idx]
           ci1 := b[idx]
           j := len(dp)
           dp2 := make([]int64, j+1)
           for k := range dp2 {
               dp2[k] = INF
           }
           // for each possible ones count so far
           for ones, cur := range dp {
               if cur >= INF {
                   continue
               }
               // assign L=0 at this position
               // cost = cur + c0 + ones
               c0 := cur + ci0 + int64(ones)
               if c0 < dp2[ones] {
                   dp2[ones] = c0
               }
               // assign L=1
               c1 := cur + ci1
               if c1 < dp2[ones+1] {
                   dp2[ones+1] = c1
               }
           }
           dp = dp2
       }
       // take min over all ones counts
       best := INF
       for _, v := range dp {
           if v < best {
               best = v
           }
       }
       ans += best
   }
   // Output result
   fmt.Println(ans)
}
