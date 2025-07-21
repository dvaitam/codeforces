package main

import (
   "bufio"
   "io"
   "os"
   "sort"
   "strconv"
)

// BIT implements a Fenwick Tree for sum queries
type BIT struct {
   n    int
   tree []int
}

// NewBIT creates a Fenwick Tree of size n
func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int, n+1)}
}

// Add increases element at index i by v (1-based)
func (b *BIT) Add(i, v int) {
   for ; i <= b.n; i += i & -i {
       b.tree[i] += v
   }
}

// Sum returns the prefix sum up to index i (1-based)
func (b *BIT) Sum(i int) int {
   s := 0
   for ; i > 0; i -= i & -i {
       s += b.tree[i]
   }
   return s
}

func main() {
   data, _ := io.ReadAll(os.Stdin)
   // fast integer parsing
   idx := 0
   nextInt := func() int {
       for idx < len(data) && (data[idx] < '0' || data[idx] > '9') {
           idx++
       }
       x := 0
       for idx < len(data) && data[idx] >= '0' && data[idx] <= '9' {
           x = x*10 + int(data[idx]-'0')
           idx++
       }
       return x
   }
   n := nextInt()
   vals := make([]int, n)
   for i := 0; i < n; i++ {
       vals[i] = nextInt()
   }
   // coordinate compression
   type pair struct{ val, idx int }
   arr := make([]pair, n)
   for i, v := range vals {
       arr[i] = pair{v, i}
   }
   sort.Slice(arr, func(i, j int) bool { return arr[i].val < arr[j].val })
   ranks := make([]int, n)
   for i, p := range arr {
       ranks[p.idx] = i + 1
   }
   // compute rightLess: number of k>j with a[k] < a[j]
   rightLess := make([]int, n)
   bit := NewBIT(n)
   for j := n - 1; j >= 0; j-- {
       r := ranks[j]
       if r > 1 {
           rightLess[j] = bit.Sum(r - 1)
       }
       bit.Add(r, 1)
   }
   // compute result using leftGreater * rightLess
   bit = NewBIT(n)
   var ans int64
   for j := 0; j < n; j++ {
       r := ranks[j]
       // number of i<j with vals[i] > vals[j] is j - sum(r)
       leftGt := j - bit.Sum(r)
       ans += int64(leftGt) * int64(rightLess[j])
       bit.Add(r, 1)
   }
   // output
   w := bufio.NewWriter(os.Stdout)
   w.WriteString(strconv.FormatInt(ans, 10))
   w.Flush()
}
