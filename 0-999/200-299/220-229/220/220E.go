package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Fenwick tree for int64
type Fenwick struct {
   n    int
   tree []int64
}

func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, tree: make([]int64, n+1)}
}

// add v at index i (1-based)
func (f *Fenwick) Add(i int, v int64) {
   for x := i; x <= f.n; x += x & -x {
       f.tree[x] += v
   }
}

// sum of [1..i]
func (f *Fenwick) Sum(i int) int64 {
   var s int64
   for x := i; x > 0; x -= x & -x {
       s += f.tree[x]
   }
   return s
}

// sum of [l..r]
func (f *Fenwick) Range(l, r int) int64 {
   if r < l {
       return 0
   }
   return f.Sum(r) - f.Sum(l-1)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var k int64
   fmt.Fscan(reader, &n, &k)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // compress values
   vals := make([]int64, n)
   copy(vals, a)
   sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
   vals = unique(vals)
   m := len(vals)
   comp := make([]int, n)
   for i := 0; i < n; i++ {
       // 1-based index
       idx := sort.Search(len(vals), func(j int) bool { return vals[j] >= a[i] })
       comp[i] = idx + 1
   }
   // compute pr: prefix inversions
   pr := make([]int64, n)
   bit := NewFenwick(m)
   for i := 0; i < n; i++ {
       ci := comp[i]
       // count previous > a[i]
       pr[i] = bit.Range(ci+1, m)
       if i > 0 {
           pr[i] += pr[i-1]
       }
       bit.Add(ci, 1)
   }
   // compute su: suffix inversions
   su := make([]int64, n+1)
   bit2 := NewFenwick(m)
   su[n] = 0
   for i := n - 1; i >= 0; i-- {
       ci := comp[i]
       // count later < a[i]
       cnt := bit2.Range(1, ci-1)
       su[i] = cnt + su[i+1]
       bit2.Add(ci, 1)
   }
   // two pointers l, r
   bitP := NewFenwick(m)
   bitS := NewFenwick(m)
   // initial prefix a[0]
   bitP.Add(comp[0], 1)
   // initial suffix a[1..]
   for i := 1; i < n; i++ {
       bitS.Add(comp[i], 1)
   }
   var ans int64
   var cnt int64
   // initial l=0 (index), r=1
   l := 0
   r := 1
   // cnt = number of pairs i<=l, j>=r where a[i]>a[j]
   cnt = bitS.Sum(comp[0] - 1)
   for l = 0; l < n-1; l++ {
       if l > 0 {
           // prefix include a[l]
           // update cnt for new i=l
           cnt += bitS.Sum(comp[l] - 1)
           bitP.Add(comp[l], 1)
       }
       // ensure r > l
       for r <= l {
           // remove a[r] from suffix
           // decrease cnt by (# i<=l where a[i] > a[r])
           cnt -= bitP.Range(comp[r]+1, m)
           bitS.Add(comp[r], -1)
           r++
       }
       // move r until f(l,r)<=k or r>n
       for r < n {
           // f = pr[l] + su[r] + cnt
           fval := pr[l] + su[r] + cnt
           if fval <= k {
               break
           }
           // remove a[r] from suffix
           cnt -= bitP.Range(comp[r]+1, m)
           bitS.Add(comp[r], -1)
           r++
       }
       if r >= n {
           // no valid r for this and further l
           break
       }
       // r in [l+1..n-1]
       ans += int64(n - r)
   }
   fmt.Fprintln(writer, ans)
}

func unique(a []int64) []int64 {
   j := 0
   for i := 0; i < len(a); i++ {
       if i == 0 || a[i] != a[i-1] {
           a[j] = a[i]
           j++
       }
   }
   return a[:j]
}
