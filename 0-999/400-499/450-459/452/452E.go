package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

const mod = 1000000007
const alpha = 29

type state struct {
   next  [alpha]int
   link, length int
   c1, c2, c3 int64
}

func getId(c rune) int {
   if c >= 'a' && c <= 'z' {
       return int(c - 'a')
   }
   switch c {
   case '#':
       return 26
   case '$':
       return 27
   case '%':
       return 28
   }
   return -1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   s1, _ := reader.ReadString('\n')
   s2, _ := reader.ReadString('\n')
   s3, _ := reader.ReadString('\n')
   s1 = strings.TrimSpace(s1)
   s2 = strings.TrimSpace(s2)
   s3 = strings.TrimSpace(s3)
   n1, n2, n3 := len(s1), len(s2), len(s3)
   minL := n1
   if n2 < minL {
       minL = n2
   }
   if n3 < minL {
       minL = n3
   }
   // build suffix automaton
   totalLen := n1 + n2 + n3 + 3
   maxStates := totalLen*2 + 5
   st := make([]state, maxStates)
   sz := 1
   last := 1
   st[1].link = 0
   st[1].length = 0
   // extend function
   extend := func(id int, which int) {
       sz++
       cur := sz
       st[cur].length = st[last].length + 1
       if which == 1 {
           st[cur].c1 = 1
       } else if which == 2 {
           st[cur].c2 = 1
       } else if which == 3 {
           st[cur].c3 = 1
       }
       p := last
       for p > 0 && st[p].next[id] == 0 {
           st[p].next[id] = cur
           p = st[p].link
       }
       if p == 0 {
           st[cur].link = 1
       } else {
           q := st[p].next[id]
           if st[p].length+1 == st[q].length {
               st[cur].link = q
           } else {
               sz++
               clone := sz
               st[clone] = st[q]
               st[clone].length = st[p].length + 1
               // reset counts on clone
               st[clone].c1, st[clone].c2, st[clone].c3 = 0, 0, 0
               for p > 0 && st[p].next[id] == q {
                   st[p].next[id] = clone
                   p = st[p].link
               }
               st[q].link = clone
               st[cur].link = clone
           }
       }
       last = cur
   }
   // add s1
   for _, ch := range s1 {
       extend(getId(ch), 1)
   }
   extend(getId('#'), 0)
   // add s2
   for _, ch := range s2 {
       extend(getId(ch), 2)
   }
   extend(getId('$'), 0)
   // add s3
   for _, ch := range s3 {
       extend(getId(ch), 3)
   }
   extend(getId('%'), 0)
   // bucket sort by length
   maxLen := st[last].length
   cntLen := make([]int, maxLen+1)
   for i := 1; i <= sz; i++ {
       l := st[i].length
       cntLen[l]++
   }
   for i := 1; i <= maxLen; i++ {
       cntLen[i] += cntLen[i-1]
   }
   order := make([]int, sz+1)
   for i := sz; i >= 1; i-- {
       l := st[i].length
       order[cntLen[l]] = i
       cntLen[l]--
   }
   // propagate counts
   for k := sz; k > 1; k-- {
       v := order[k]
       p := st[v].link
       st[p].c1 += st[v].c1
       st[p].c2 += st[v].c2
       st[p].c3 += st[v].c3
   }
   // difference array
   diff := make([]int64, minL+2)
   for i := 2; i <= sz; i++ {
       c1 := st[i].c1 % mod
       c2 := st[i].c2 % mod
       c3 := st[i].c3 % mod
       v := c1 * c2 % mod * c3 % mod
       if v == 0 {
           continue
       }
       link := st[st[i].link].length
       lbound := link + 1
       if lbound > minL {
           continue
       }
       rbound := st[i].length
       if rbound > minL {
           rbound = minL
       }
       if lbound > rbound {
           continue
       }
       diff[lbound] = (diff[lbound] + v) % mod
       diff[rbound+1] = (diff[rbound+1] - v + mod) % mod
   }
   // compute answers
   builder := strings.Builder{}
   builder.Grow(minL * 10)
   var cur int64
   for i := 1; i <= minL; i++ {
       cur = (cur + diff[i]) % mod
       if i > 1 {
           builder.WriteByte(' ')
       }
       builder.WriteString(strconv.FormatInt(cur, 10))
   }
   builder.WriteByte('\n')
   fmt.Print(builder.String())
}
