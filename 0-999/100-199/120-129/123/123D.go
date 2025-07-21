package main

import (
   "bufio"
   "fmt"
   "os"
)

const ALPHA = 26

// state for Suffix Automaton
type state struct {
   next [ALPHA]int // transitions: store stateIndex+1, 0 means none
   link int         // suffix link
   len  int         // length of longest string in this class
   occ  int         // number of end positions
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, _ := reader.ReadString('\n')
   if len(s) > 0 && s[len(s)-1] == '\n' {
       s = s[:len(s)-1]
   }
   n := len(s)
   // initialize SAM
   maxStates := 2 * n
   st := make([]state, maxStates)
   size, last := 1, 0
   st[0].len = 0
   st[0].link = -1
   // extend function
   extend := func(c int) {
       cur := size
       size++
       st[cur].len = st[last].len + 1
       st[cur].occ = 1
       p := last
       for p != -1 && st[p].next[c] == 0 {
           st[p].next[c] = cur + 1
           p = st[p].link
       }
       if p == -1 {
           st[cur].link = 0
       } else {
           q := st[p].next[c] - 1
           if st[p].len+1 == st[q].len {
               st[cur].link = q
           } else {
               clone := size
               size++
               st[clone] = st[q]
               st[clone].len = st[p].len + 1
               st[clone].occ = 0
               for p != -1 && st[p].next[c] == q+1 {
                   st[p].next[c] = clone + 1
                   p = st[p].link
               }
               st[q].link = clone
               st[cur].link = clone
           }
       }
       last = cur
   }
   // build SAM
   for i := 0; i < n; i++ {
       extend(int(s[i] - 'a'))
   }
   // count sort states by len
   maxLen := n
   cntLen := make([]int, maxLen+1)
   for i := 0; i < size; i++ {
       cntLen[st[i].len]++
   }
   pos := make([]int, maxLen+1)
   for l := 1; l <= maxLen; l++ {
       pos[l] = pos[l-1] + cntLen[l-1]
   }
   order := make([]int, size)
   for i := 0; i < size; i++ {
       l := st[i].len
       order[pos[l]] = i
       pos[l]++
   }
   // propagate occ in reverse order
   for i := size - 1; i >= 0; i-- {
       v := order[i]
       if st[v].link != -1 {
           st[st[v].link].occ += st[v].occ
       }
   }
   // compute answer
   var ans uint64
   for i := 0; i < size; i++ {
       linkLen := 0
       if st[i].link != -1 {
           linkLen = st[st[i].link].len
       }
       delta := st[i].len - linkLen
       if delta > 0 {
           c := uint64(st[i].occ)
           ans += uint64(delta) * c * (c + 1) / 2
       }
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   fmt.Fprint(writer, ans)
   writer.Flush()
}
