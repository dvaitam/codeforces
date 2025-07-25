package main

import (
   "bufio"
   "fmt"
   "os"
)

// Suffix Automaton for lowercase letters
type state struct {
   next map[byte]int
   link int
   len  int
}

type suffixAutomaton struct {
   st   []state
   last int
}

func newSuffixAutomaton(n int) *suffixAutomaton {
   sa := &suffixAutomaton{
       st:   make([]state, 0, 2*n),
       last: 0,
   }
   sa.st = append(sa.st, state{next: make(map[byte]int), link: -1, len: 0})
   return sa
}

func (sa *suffixAutomaton) extend(c byte) {
   p := sa.last
   cur := len(sa.st)
   sa.st = append(sa.st, state{next: make(map[byte]int), len: sa.st[p].len + 1})
   for p != -1 && sa.st[p].next[c] == 0 {
       sa.st[p].next[c] = cur
       p = sa.st[p].link
   }
   if p == -1 {
       sa.st[cur].link = 0
   } else {
       q := sa.st[p].next[c]
       if sa.st[p].len+1 == sa.st[q].len {
           sa.st[cur].link = q
       } else {
           // clone
           clone := len(sa.st)
           // copy q
           mp := make(map[byte]int)
           for k, v := range sa.st[q].next {
               mp[k] = v
           }
           sa.st = append(sa.st, state{next: mp, len: sa.st[p].len + 1, link: sa.st[q].link})
           for p != -1 && sa.st[p].next[c] == q {
               sa.st[p].next[c] = clone
               p = sa.st[p].link
           }
           sa.st[q].link = sa.st[cur].link = clone
       }
   }
   sa.last = cur
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, a, b int
   fmt.Fscan(reader, &n, &a, &b)
   var s string
   fmt.Fscan(reader, &s)

   dp := make([]int, n+1)
   const inf = 1e18
   for i := 1; i <= n; i++ {
       dp[i] = int(1e18)
   }

   sa := newSuffixAutomaton(n)
   // process positions
   for i := 0; i < n; i++ {
       // cost of single char
       if dp[i]+a < dp[i+1] {
           dp[i+1] = dp[i] + a
       }
       // longest substring from i in prefix
       u, length := 0, 0
       for j := i; j < n; j++ {
           c := s[j]
           nxt, ok := sa.st[u].next[c]
           if !ok {
               break
           }
           u = nxt
           length++
           // can encode s[i:j+1] with cost b
           if dp[i]+b < dp[j+1] {
               dp[j+1] = dp[i] + b
           }
       }
       // extend automaton with s[i]
       sa.extend(s[i])
   }
   fmt.Println(dp[n])
}
