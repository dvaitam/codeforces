package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Constraint struct {
   c byte
   l, r int
   p int // prefix count at U
   a, b int // current interval [a,b]
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var s string
   fmt.Fscan(in, &s)
   n := len(s)
   var k, L, R int
   fmt.Fscan(in, &k, &L, &R)
   cons := make([]*Constraint, k)
   byChar := make(map[byte][]int)
   for i := 0; i < k; i++ {
       var ci byte
       var li, ri int
       fmt.Fscan(in, &ci, &li, &ri)
       cons[i] = &Constraint{c: ci, l: li, r: ri}
       byChar[ci] = append(byChar[ci], i)
   }
   // occurrences positions 1-based
   occ := make(map[byte][]int)
   for i := 0; i < n; i++ {
       ch := s[i]
       occ[ch] = append(occ[ch], i+1)
   }
   // init intervals at U=0
   for i, c := range cons {
       arr := occ[c.c]
       tot := len(arr)
       // p=0
       if c.l == 0 {
           c.a = 1
       } else if c.l <= tot {
           c.a = arr[c.l-1]
       } else {
           c.a = n+1
       }
       if c.r < tot {
           c.b = arr[c.r] - 1
       } else {
           c.b = n
       }
   }
   var total int64
   // process U from 0 to n-1
   for U := 0; U < n; U++ {
       // compute ansU for substrings [U+1, V]
       // build events
       ev := make([]struct{pos, d int}, 0, 512)
       start := U + 1
       for _, c := range cons {
           a := c.a
           if a < start {
               a = start
           }
           b := c.b
           if b < start {
               continue
           }
           if a <= b {
               ev = append(ev, struct{pos, d int}{a, +1})
               if b+1 <= n {
                   ev = append(ev, struct{pos, d int}{b + 1, -1})
               }
           }
       }
       // add sentinel
       ev = append(ev, struct{pos, d int}{n+1, 0})
       sort.Slice(ev, func(i, j int) bool { return ev[i].pos < ev[j].pos })
       cur := 0
       prev := start
       var ansU int
       for _, e := range ev {
           if e.pos > prev {
               if cur >= L && cur <= R {
                   ansU += e.pos - prev
               }
               prev = e.pos
           }
           cur += e.d
       }
       total += int64(ansU)
       // update for next U: slide U->U+1, update constraints for char s[U]
       ch := s[U]
       for _, idx := range byChar[ch] {
           c := cons[idx]
           // increment p
           c.p++
           arr := occ[ch]
           tot := len(arr)
           // new a and b
           if c.l == 0 {
               c.a = 1
           } else if c.p + c.l <= tot {
               c.a = arr[c.p + c.l - 1]
           } else {
               c.a = n+1
           }
           if c.p + c.r < tot {
               c.b = arr[c.p + c.r] - 1
           } else {
               c.b = n
           }
       }
   }
   fmt.Fprintln(out, total)
}
