package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "sort"
   "time"
)

// BIT implements a Binary Indexed Tree for prefix sums of int64
type BIT struct {
   n int
   t []int64
}

func NewBIT(n int) *BIT {
   return &BIT{n: n, t: make([]int64, n+2)}
}

// Add adds v at index i (0-based)
func (b *BIT) Add(i int, v int64) {
   for i++; i <= b.n; i += i & -i {
       b.t[i] += v
   }
}

// Sum returns sum of [0..i], returns 0 for i<0
func (b *BIT) Sum(i int) int64 {
   if i < 0 {
       return 0
   }
   var s int64
   for i++; i > 0; i -= i & -i {
       s += b.t[i]
   }
   return s
}

// Treap node for ordered set of ints
type treap struct {
   key        int
   prio       int
   left, right *treap
}

// split splits tree t into l (< key) and r (>= key)
func split(t *treap, key int) (l, r *treap) {
   if t == nil {
       return nil, nil
   }
   if t.key < key {
       var sr *treap
       t.right, sr = split(t.right, key)
       return t, sr
   }
   var sl *treap
   sl, t.left = split(t.left, key)
   return sl, t
}

// merge merges two treaps a and b, all keys in a < keys in b
func merge(a, b *treap) *treap {
   if a == nil {
       return b
   }
   if b == nil {
       return a
   }
   if a.prio < b.prio {
       a.right = merge(a.right, b)
       return a
   }
   b.left = merge(a, b.left)
   return b
}

// Set is a wrapper for treap-based ordered set of ints
type Set struct{ root *treap }

// Insert key into set (no-op if exists)
func (s *Set) Insert(key int) {
   var l, r = split(s.root, key)
   var m, _ = split(r, key+1)
   if m == nil {
       m = &treap{key: key, prio: rand.Int()}
   }
   s.root = merge(merge(l, m), r)
}

// Erase key from set (no-op if not exists)
func (s *Set) Erase(key int) {
   var l, r = split(s.root, key)
   var m, rr = split(r, key+1)
   // drop m
   s.root = merge(l, rr)
}

// Next returns smallest element >= key, or -1 if none
func (s *Set) Next(key int) int {
   var ans = -1
   for t := s.root; t != nil; {
       if t.key >= key {
           ans = t.key
           t = t.left
       } else {
           t = t.right
       }
   }
   return ans
}

// Prev returns largest element <= key, or -1 if none
func (s *Set) Prev(key int) int {
   var ans = -1
   for t := s.root; t != nil; {
       if t.key <= key {
           ans = t.key
           t = t.right
       } else {
           t = t.left
       }
   }
   return ans
}

// Xs handles coordinate compression
type Xs struct{ a []int }

func (x *Xs) Add(v int) { x.a = append(x.a, v) }
func (x *Xs) Init() {
   sort.Ints(x.a)
   x.a = unique(x.a)
}
func unique(a []int) []int {
   n := 0
   for i, v := range a {
       if i == 0 || v != a[i-1] {
           a[n] = v
           n++
       }
   }
   return a[:n]
}
func (x *Xs) Id(v int) int {
   // index of v in x.a
   i := sort.SearchInts(x.a, v)
   return i
}

func main() {
   rand.Seed(time.Now().UnixNano())
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var qn int
   fmt.Fscan(in, &qn)
   ops := make([]int, qn)
   xs := new(Xs)
   for i := 0; i < qn; i++ {
       var op string
       var v int
       fmt.Fscan(in, &op, &v)
       xs.Add(v)
       if op == "-" {
           ops[i] = -v
       } else {
           ops[i] = v
       }
   }
   xs.Init()
   n := len(xs.a)
   bit := NewBIT(n + 2)
   cs := make([]int, n)
   w := make([]int64, n)
   s := new(Set)
   bs := new(Set)
   cnt := 0

   check := func(i int) {
       j := s.Prev(i - 1)
       if j == -1 || int64(xs.a[j])*2 >= int64(xs.a[i]) {
           bs.Erase(i)
       } else {
           bs.Insert(i)
           w[i] = bit.Sum(i - 1)
       }
   }

   for _, a := range ops {
       idx := abs(a)
       i := xs.Id(idx)
       if a < 0 {
           cnt--
       } else {
           cnt++
       }
       if cs[i] == 0 {
           s.Insert(i)
           check(i)
       }
       if a < 0 {
           cs[i]--
       } else {
           cs[i]++
       }
       if cs[i] == 0 {
           bs.Erase(i)
           s.Erase(i)
       }
       if j := s.Next(i + 1); j != -1 {
           check(j)
       }
       bit.Add(i, int64(a))
       // update w for bad indices >= i
       j := i
       for {
           j = bs.Next(j + 1)
           if j == -1 {
               break
           }
           w[j] += int64(a)
       }
       ans := cnt - 1
       j = -1
       for {
           j = bs.Next(j + 1)
           if j == -1 {
               break
           }
           if w[j]*2 < int64(xs.a[j]) {
               ans--
           }
       }
       if ans < 0 {
           ans = 0
       }
       fmt.Fprintln(out, ans)
   }
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}
