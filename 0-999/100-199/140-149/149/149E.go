package main

import (
   "bufio"
   "fmt"
   "os"
)

// Suffix automaton for uppercase letters
type state struct {
   next   [26]int
   link   int
   length int
   minEnd int
}

type SAM struct {
   nodes  []state
   tot    int
   last   int
   maxLen int
}

// NewSAM creates a suffix automaton that can handle strings up to maxLen
func NewSAM(maxLen int) *SAM {
   size := 2*maxLen + 5
   nodes := make([]state, size)
   inf := maxLen + 5
   for i := 0; i < size; i++ {
       nodes[i].minEnd = inf
       for c := 0; c < 26; c++ {
           nodes[i].next[c] = -1
       }
   }
   nodes[0].link = -1
   nodes[0].length = 0
   return &SAM{nodes: nodes, tot: 1, last: 0, maxLen: maxLen}
}

// Extend adds character c (0..25) at position pos (1-based) to the automaton
func (sam *SAM) Extend(c, pos int) {
   cur := sam.tot
   sam.tot++
   sam.nodes[cur].length = sam.nodes[sam.last].length + 1
   sam.nodes[cur].minEnd = pos
   p := sam.last
   for p >= 0 && sam.nodes[p].next[c] < 0 {
       sam.nodes[p].next[c] = cur
       p = sam.nodes[p].link
   }
   if p < 0 {
       sam.nodes[cur].link = 0
   } else {
       q := sam.nodes[p].next[c]
       if sam.nodes[p].length+1 == sam.nodes[q].length {
           sam.nodes[cur].link = q
       } else {
           clone := sam.tot
           sam.tot++
           // copy q
           sam.nodes[clone] = sam.nodes[q]
           sam.nodes[clone].length = sam.nodes[p].length + 1
           // reset minEnd to inf for clone
           sam.nodes[clone].minEnd = sam.maxLen + 5
           for p >= 0 && sam.nodes[p].next[c] == q {
               sam.nodes[p].next[c] = clone
               p = sam.nodes[p].link
           }
           sam.nodes[q].link = clone
           sam.nodes[cur].link = clone
       }
   }
   sam.last = cur
}

// Prepare propagates minEnd values through the automaton
func (sam *SAM) Prepare() {
   // bucket sort states by length
   maxL := 0
   for i := 0; i < sam.tot; i++ {
       if sam.nodes[i].length > maxL {
           maxL = sam.nodes[i].length
       }
   }
   cnt := make([]int, maxL+1)
   for i := 0; i < sam.tot; i++ {
       cnt[sam.nodes[i].length]++
   }
   for i := 1; i <= maxL; i++ {
       cnt[i] += cnt[i-1]
   }
   order := make([]int, sam.tot)
   for i := sam.tot - 1; i >= 0; i-- {
       l := sam.nodes[i].length
       cnt[l]--
       order[cnt[l]] = i
   }
   // propagate in descending length order
   for idx := sam.tot - 1; idx > 0; idx-- {
       v := order[idx]
       p := sam.nodes[v].link
       if p >= 0 && sam.nodes[p].minEnd > sam.nodes[v].minEnd {
           sam.nodes[p].minEnd = sam.nodes[v].minEnd
       }
   }
}

func reverseBytes(s []byte) []byte {
   n := len(s)
   r := make([]byte, n)
   for i := 0; i < n; i++ {
       r[i] = s[n-1-i]
   }
   return r
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   fmt.Fscan(reader, &s)
   n := len(s)
   // build SAM for s
   sam := NewSAM(n)
   for i := 0; i < n; i++ {
       sam.Extend(int(s[i]-'A'), i+1)
   }
   sam.Prepare()
   // build SAM for reversed s
   srev := reverseBytes([]byte(s))
   sams := NewSAM(n)
   for i := 0; i < n; i++ {
       sams.Extend(int(srev[i]-'A'), i+1)
   }
   sams.Prepare()
   var m int
   fmt.Fscan(reader, &m)
   ans := 0
   inf := n + 5
   for pi := 0; pi < m; pi++ {
       var p string
       fmt.Fscan(reader, &p)
       L := len(p)
       pre := make([]int, L+1)
       for i := 0; i <= L; i++ {
           pre[i] = inf
       }
       // forward prefixes
       cur := 0
       ok := true
       for j := 0; j < L; j++ {
           c := int(p[j] - 'A')
           if ok && sam.nodes[cur].next[c] >= 0 {
               cur = sam.nodes[cur].next[c]
               pre[j+1] = sam.nodes[cur].minEnd
           } else {
               ok = false
           }
       }
       // reversed prefixes -> original suffixes
       revp := reverseBytes([]byte(p))
       revpre := make([]int, L+1)
       for i := 0; i <= L; i++ {
           revpre[i] = inf
       }
       cur = 0
       ok = true
       for j := 0; j < L; j++ {
           c := int(revp[j] - 'A')
           if ok && sams.nodes[cur].next[c] >= 0 {
               cur = sams.nodes[cur].next[c]
               revpre[j+1] = sams.nodes[cur].minEnd
           } else {
               ok = false
           }
       }
       // check splits
       found := false
       for k := 1; k < L; k++ {
           if pre[k] < inf {
               vlen := L - k
               if revpre[vlen] < inf {
                   // max start of suffix v = n - minEndRev + 1
                   startV := n - revpre[vlen] + 1
                   if pre[k] < startV {
                       found = true
                       break
                   }
               }
           }
       }
       if found {
           ans++
       }
   }
   fmt.Println(ans)
}
