package main

import (
   "bufio"
   "fmt"
   "os"
)

// Actually implement with explicit size
type SAM2 struct {
   next []map[byte]int
   link []int
   len  []int
   cnt  []int64
   last int
   sz   int
}

func NewSAM2(maxLen int) *SAM2 {
   size := 2*maxLen + 1
   sam := &SAM2{
       next: make([]map[byte]int, size),
       link: make([]int, size),
       len:  make([]int, size),
       cnt:  make([]int64, size),
       last: 0,
       sz:   1,
   }
   sam.next[0] = make(map[byte]int)
   sam.link[0] = -1
   return sam
}

// Extend adds character c
func (sam *SAM2) Extend(c byte) {
   p := sam.last
   cur := sam.sz
   sam.sz++
   sam.len[cur] = sam.len[p] + 1
   sam.next[cur] = make(map[byte]int)
   sam.cnt[cur] = 1
   // propagate transitions
   for p >= 0 && sam.next[p][c] == 0 {
       sam.next[p][c] = cur
       p = sam.link[p]
   }
   if p == -1 {
       sam.link[cur] = 0
   } else {
       q := sam.next[p][c]
       if sam.len[p]+1 == sam.len[q] {
           sam.link[cur] = q
       } else {
           clone := sam.sz
           sam.sz++
           sam.len[clone] = sam.len[p] + 1
           // copy transitions
           sam.next[clone] = make(map[byte]int)
           for k, v := range sam.next[q] {
               sam.next[clone][k] = v
           }
           sam.link[clone] = sam.link[q]
           // no cnt for clone
           sam.cnt[clone] = 0
           for p >= 0 && sam.next[p][c] == q {
               sam.next[p][c] = clone
               p = sam.link[p]
           }
           sam.link[q] = clone
           sam.link[cur] = clone
       }
   }
   sam.last = cur
}

// Build automaton from string
func BuildSAM(s string) *SAM2 {
   sam := NewSAM2(len(s))
   for i := 0; i < len(s); i++ {
       sam.Extend(s[i])
   }
   // topo order by len
   maxLen := 0
   for i := 0; i < sam.sz; i++ {
       if sam.len[i] > maxLen {
           maxLen = sam.len[i]
       }
   }
   bucket := make([]int, maxLen+1)
   for i := 0; i < sam.sz; i++ {
       bucket[sam.len[i]]++
   }
   for i := 1; i <= maxLen; i++ {
       bucket[i] += bucket[i-1]
   }
   order := make([]int, sam.sz)
   for i := sam.sz - 1; i >= 0; i-- {
       l := sam.len[i]
       bucket[l]--
       order[bucket[l]] = i
   }
   // accumulate counts
   for i := sam.sz - 1; i > 0; i-- {
       v := order[i]
       p := sam.link[v]
       if p >= 0 {
           sam.cnt[p] += sam.cnt[v]
       }
   }
   return sam
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   fmt.Fscan(reader, &s)
   var n int
   fmt.Fscan(reader, &n)
   type Rule struct {
       sam *SAM2
       l   int64
       r   int64
   }
   rules := make([]Rule, n)
   for i := 0; i < n; i++ {
       var p string
       var l, r int64
       fmt.Fscan(reader, &p, &l, &r)
       sam := BuildSAM(p)
       rules[i] = Rule{sam: sam, l: l, r: r}
   }
   good := make(map[string]struct{})
   // iterate all substrings
   for i := 0; i < len(s); i++ {
       // current state in each SAM
       cur := make([]int, n)
       for k := 0; k < n; k++ {
           cur[k] = 0
       }
       for j := i; j < len(s); j++ {
           c := s[j]
           ok := true
           for k := 0; k < n; k++ {
               if cur[k] >= 0 {
                   nxt, has := rules[k].sam.next[cur[k]][c]
                   if !has {
                       cur[k] = -1
                   } else {
                       cur[k] = nxt
                   }
               }
               // get occurrences
               var occ int64
               if cur[k] >= 0 {
                   occ = rules[k].sam.cnt[cur[k]]
               } else {
                   occ = 0
               }
               if occ < rules[k].l || occ > rules[k].r {
                   ok = false
                   break
               }
           }
           if ok {
               sub := s[i : j+1]
               good[sub] = struct{}{}
           }
       }
   }
   fmt.Println(len(good))
}
