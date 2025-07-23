package main

import (
   "bufio"
   "fmt"
   "os"
)

// BIT implements a Fenwick tree for int64 values
type BIT struct {
   n    int
   tree []int64
}

func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int64, n+1)}
}

// Add adds v at position i (1-based)
func (b *BIT) Add(i int, v int64) {
   for x := i; x <= b.n; x += x & -x {
       b.tree[x] += v
   }
}

// Sum returns sum of [1..i]
func (b *BIT) Sum(i int) int64 {
   var s int64
   for x := i; x > 0; x -= x & -x {
       s += b.tree[x]
   }
   return s
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, q int
   fmt.Fscan(reader, &n, &q)
   s := make([]string, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &s[i])
   }
   // Build Aho-Corasick trie
   // next[state][c] -> state
   next := make([][]int, 1)
   next[0] = make([]int, 26)
   fail := make([]int, 1)
   endpoints := make([]int, n+1)
   // insert strings
   for i := 1; i <= n; i++ {
       u := 0
       for _, ch := range s[i] {
           c := int(ch - 'a')
           if next[u][c] == 0 {
               next = append(next, make([]int, 26))
               fail = append(fail, 0)
               next[u][c] = len(next) - 1
           }
           u = next[u][c]
       }
       endpoints[i] = u
   }
   // build failure links
   queue := make([]int, 0, len(next))
   for c := 0; c < 26; c++ {
       v := next[0][c]
       if v != 0 {
           queue = append(queue, v)
       }
   }
   for idx := 0; idx < len(queue); idx++ {
       u := queue[idx]
       for c := 0; c < 26; c++ {
           v := next[u][c]
           if v != 0 {
               f := fail[u]
               next[u][c] = v
               fail[v] = next[f][c]
               queue = append(queue, v)
           } else {
               next[u][c] = next[fail[u]][c]
           }
       }
   }
   // build failure tree children
   sz := len(next)
   children := make([][]int, sz)
   for u := 1; u < sz; u++ {
       p := fail[u]
       children[p] = append(children[p], u)
   }
   // Euler tour on failure tree (iterative)
   tin := make([]int, sz)
   tout := make([]int, sz)
   time := 0
   type frame struct{ u, idx int }
   stack := []frame{{0, 0}}
   for len(stack) > 0 {
       top := stack[len(stack)-1]
       u := top.u
       if top.idx == 0 {
           time++
           tin[u] = time
       }
       if top.idx < len(children[u]) {
           v := children[u][top.idx]
           // increment index of current frame
           stack[len(stack)-1].idx++
           // push child
           stack = append(stack, frame{v, 0})
       } else {
           tout[u] = time
           // pop frame
           stack = stack[:len(stack)-1]
       }
   }
   // record visits for each text
   visits := make([][]int, n+1)
   for i := 1; i <= n; i++ {
       u := 0
       for _, ch := range s[i] {
           c := int(ch - 'a')
           u = next[u][c]
           visits[i] = append(visits[i], u)
       }
   }
   // events: for each time t, list of queries {patternNode, queryID, sign}
   type Event struct{ p, id, sign int }
   events := make([][]Event, n+1)
   ans := make([]int64, q)
   for j := 0; j < q; j++ {
       var l, r, k int
       fmt.Fscan(reader, &l, &r, &k)
       p := endpoints[k]
       events[r] = append(events[r], Event{p, j, 1})
       if l > 1 {
           events[l-1] = append(events[l-1], Event{p, j, -1})
       }
   }
   // process events with BIT over Euler time
   bit := NewBIT(sz + 2)
   for t := 1; t <= n; t++ {
       for _, u := range visits[t] {
           bit.Add(tin[u], 1)
       }
       for _, ev := range events[t] {
           lo := tin[ev.p]
           hi := tout[ev.p]
           cnt := bit.Sum(hi) - bit.Sum(lo-1)
           ans[ev.id] += int64(ev.sign) * cnt
       }
   }
   // output answers
   for i := 0; i < q; i++ {
       fmt.Fprintln(writer, ans[i])
   }
