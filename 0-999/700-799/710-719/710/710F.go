package main

import (
   "bufio"
   "fmt"
   "os"
)

// AhoCorasick automaton
type AC struct {
   next [][]int
   fail []int
   out  []int
}

// NewAC builds an AC automaton for given patterns
func NewAC(pats []string) *AC {
   // node 0 is root
   next := make([][]int, 1)
   next[0] = make([]int, 26)
   for i := range next[0] {
       next[0][i] = -1
   }
   out := make([]int, 1)
   // insert patterns
   for _, p := range pats {
       u := 0
       for _, ch := range p {
           c := int(ch - 'a')
           if next[u][c] == -1 {
               next[u][c] = len(next)
               next = append(next, make([]int, 26))
               for i := range next[len(next)-1] {
                   next[len(next)-1][i] = -1
               }
               out = append(out, 0)
           }
           u = next[u][c]
       }
       out[u]++
   }
   // build fail
   n := len(next)
   fail := make([]int, n)
   // init queue
   q := make([]int, 0, n)
   // first level
   for c := 0; c < 26; c++ {
       v := next[0][c]
       if v != -1 {
           fail[v] = 0
           q = append(q, v)
       } else {
           next[0][c] = 0
       }
   }
   // bfs
   for i := 0; i < len(q); i++ {
       u := q[i]
       for c := 0; c < 26; c++ {
           v := next[u][c]
           if v != -1 {
               f := fail[u]
               fail[v] = next[f][c]
               out[v] += out[fail[v]]
               q = append(q, v)
           } else {
               next[u][c] = next[fail[u]][c]
           }
       }
   }
   return &AC{next: next, fail: fail, out: out}
}

// Count occurrences of patterns in text s
func (ac *AC) Count(s string) int {
   res, u := 0, 0
   for _, ch := range s {
       c := int(ch - 'a')
       u = ac.next[u][c]
       res += ac.out[u]
   }
   return res
}

// Bucket holds patterns and its automaton
type Bucket struct {
   pats []string
   ac   *AC
}

// add pattern to buckets (binary merge)
func addPattern(buckets []*Bucket, s string) []*Bucket {
   nb := &Bucket{pats: []string{s}}
   for i := 0; ; i++ {
       if i >= len(buckets) {
           buckets = append(buckets, nb)
           break
       }
       if buckets[i] == nil {
           // place here
           buckets[i] = nb
           break
       }
       // merge
       merged := make([]string, 0, len(buckets[i].pats)+len(nb.pats))
       merged = append(merged, buckets[i].pats...)
       merged = append(merged, nb.pats...)
       nb = &Bucket{pats: merged}
       buckets[i] = nil
   }
   // rebuild only the last placed bucket's ac
   // find last i where nb placed
   // rebuild all non-nil buckets that have nil ac
   for _, b := range buckets {
       if b != nil && b.ac == nil {
           b.ac = NewAC(b.pats)
       }
   }
   return buckets
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var m int
   fmt.Fscan(reader, &m)
   addBuckets := make([]*Bucket, 0)
   delBuckets := make([]*Bucket, 0)
   for i := 0; i < m; i++ {
       var t int
       var s string
       fmt.Fscan(reader, &t, &s)
       if t == 1 {
           addBuckets = addPattern(addBuckets, s)
       } else if t == 2 {
           delBuckets = addPattern(delBuckets, s)
       } else if t == 3 {
           // query
           cnt := 0
           for _, b := range addBuckets {
               if b != nil {
                   cnt += b.ac.Count(s)
               }
           }
           for _, b := range delBuckets {
               if b != nil {
                   cnt -= b.ac.Count(s)
               }
           }
           fmt.Fprintln(writer, cnt)
           writer.Flush()
       }
   }
}
