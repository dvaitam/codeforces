package main

import (
   "bufio"
   "fmt"
   "os"
)

// Suffix automaton state
type state struct {
   next map[byte]int
   link int
   len  int
   cnt  int
}

var st []state
var last, sz int

// initialize suffix automaton with capacity for maximum states
func saInit(n int) {
   st = make([]state, 2*n)
   st[0].next = make(map[byte]int)
   st[0].link = -1
   st[0].len = 0
   st[0].cnt = 0
   last = 0
   sz = 1
}

// extend automaton with character c
func saExtend(c byte) {
   cur := sz
   sz++
   st[cur].len = st[last].len + 1
   st[cur].cnt = 1
   st[cur].next = make(map[byte]int)
   p := last
   for p != -1 && st[p].next[c] == 0 {
       st[p].next[c] = cur
       p = st[p].link
   }
   if p == -1 {
       st[cur].link = 0
   } else {
       q := st[p].next[c]
       if st[p].len+1 == st[q].len {
           st[cur].link = q
       } else {
           clone := sz
           sz++
           st[clone].len = st[p].len + 1
           st[clone].next = make(map[byte]int)
           for k, v := range st[q].next {
               st[clone].next[k] = v
           }
           st[clone].link = st[q].link
           st[clone].cnt = 0
           for p != -1 && st[p].next[c] == q {
               st[p].next[c] = clone
               p = st[p].link
           }
           st[q].link = clone
           st[cur].link = clone
       }
   }
   last = cur
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   nS := len(s)
   saInit(nS)
   for i := 0; i < nS; i++ {
       saExtend(s[i])
   }
   // accumulate end position counts
   maxLen := 0
   for i := 0; i < sz; i++ {
       if st[i].len > maxLen {
           maxLen = st[i].len
       }
   }
   cntLen := make([]int, maxLen+1)
   for i := 0; i < sz; i++ {
       cntLen[st[i].len]++
   }
   for i := 1; i <= maxLen; i++ {
       cntLen[i] += cntLen[i-1]
   }
   order := make([]int, sz)
   for i := sz - 1; i >= 0; i-- {
       l := st[i].len
       cntLen[l]--
       order[cntLen[l]] = i
   }
   for i := sz - 1; i > 0; i-- {
       v := order[i]
       if st[v].link >= 0 {
           st[st[v].link].cnt += st[v].cnt
       }
   }
   // prepare stamp array for marking visited states per query
   stamp := make([]int, sz)
   curStamp := 0

   var q int
   fmt.Fscan(reader, &q)
   for qi := 0; qi < q; qi++ {
       var x string
       fmt.Fscan(reader, &x)
       m := len(x)
       if m > nS {
           fmt.Fprintln(writer, 0)
           continue
       }
       curStamp++
       // build doubled string for rotations
       P := x + x[:m-1]
       v := 0
       l := 0
       var ans int64
       for i := 0; i < len(P); i++ {
           c := P[i]
           for v != -1 && st[v].next[c] == 0 {
               v = st[v].link
               if v >= 0 {
                   l = st[v].len
               } else {
                   l = 0
               }
           }
           if v == -1 {
               v = 0
               l = 0
           }
           if st[v].next[c] != 0 {
               v = st[v].next[c]
               l++
           }
           if l >= m {
               // shrink to length m
               u := v
               for st[st[u].link].len >= m {
                   u = st[u].link
               }
               if stamp[u] != curStamp {
                   stamp[u] = curStamp
                   ans += int64(st[u].cnt)
               }
           }
       }
       fmt.Fprintln(writer, ans)
   }
}
