package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct { l, r int }

// mlist manages segments to remove based on adjacent equal chars
type mlist struct {
   n    int
   str  []byte
   nxt  []int
   pre  []int
   live []bool
   cnt  [26]int
   loc  [26][]int
   ans  []pair
}

func (l *mlist) init(n int, str []byte) {
   l.n = n
   l.str = str
   l.nxt = make([]int, n+2)
   l.pre = make([]int, n+2)
   l.live = make([]bool, n+2)
   for i := 0; i < 26; i++ {
       l.cnt[i] = 0
       l.loc[i] = l.loc[i][:0]
   }
   l.ans = l.ans[:0]
   for i := 0; i <= n+1; i++ {
       l.pre[i] = i - 1
       l.nxt[i] = i + 1
   }
   for i := 1; i <= n; i++ {
       l.live[i] = false
   }
   l.live[0], l.live[n+1] = true, true
}

func (l *mlist) insert(a int) {
   l.live[a] = true
   c := int(l.str[a] - 'a')
   l.cnt[c]++
   l.loc[c] = append(l.loc[c], a)
}

func (l *mlist) fin_nxt(a int) int {
   var path []int
   x := a
   for !l.live[x] {
       path = append(path, x)
       x = l.nxt[x]
   }
   for _, y := range path {
       l.nxt[y] = x
   }
   return x
}

func (l *mlist) fin_pre(a int) int {
   var path []int
   x := a
   for !l.live[x] {
       path = append(path, x)
       x = l.pre[x]
   }
   for _, y := range path {
       l.pre[y] = x
   }
   return x
}

func (l *mlist) ref(c int) {
   for len(l.loc[c]) > 0 && !l.live[l.loc[c][len(l.loc[c])-1]] {
       l.loc[c] = l.loc[c][:len(l.loc[c])-1]
   }
}

func (l *mlist) que(c int) {
   l.ref(c)
   if len(l.loc[c]) == 0 {
       return
   }
   t := l.loc[c][len(l.loc[c])-1]
   s := l.fin_nxt(l.nxt[t])
   tc := int(l.str[s] - 'a')
   if s == l.n+1 {
       maxi, ti := 0, -1
       for i := 0; i < 26; i++ {
           l.ref(i)
           if i == c || len(l.loc[i]) == 0 {
               continue
           }
           if maxi < l.loc[i][len(l.loc[i])-1] {
               maxi = l.loc[i][len(l.loc[i])-1]
               ti = i
           }
       }
       if ti < 0 {
           l.ans = append(l.ans, pair{t, t})
           l.live[t] = false
           l.cnt[c]--
           l.ref(c)
           return
       }
       t = maxi
       s = l.fin_nxt(l.nxt[t])
       tc = int(l.str[t] - 'a')
   }
   l.ans = append(l.ans, pair{t + 1, s})
   l.live[t], l.live[s] = false, false
   l.cnt[int(l.str[t]-'a')]--
   l.cnt[int(l.str[s]-'a')]--
   l.ref(c)
   l.ref(tc)
}

// slis manages deletion mapping with a BIT for current positions
type slis struct {
   n    int
   nxt  []int
   pre  []int
   live []bool
   bit  []int
}

func (s *slis) init(n int) {
   s.n = n
   s.nxt = make([]int, n+2)
   s.pre = make([]int, n+2)
   s.live = make([]bool, n+2)
   s.bit = make([]int, n+2)
   for i := 1; i <= n; i++ {
       s.bit[i] = 0
   }
   for i := 1; i <= n; i++ {
       s.upd(i, 1)
   }
   for i := 0; i <= n+1; i++ {
       s.pre[i] = i - 1
       s.nxt[i] = i + 1
       s.live[i] = true
   }
}

func (s *slis) upd(i, v int) {
   for i <= s.n {
       s.bit[i] += v
       i += i & -i
   }
}

func (s *slis) getv(i int) int {
   sum := 0
   for i > 0 {
       sum += s.bit[i]
       i &= i - 1
   }
   return sum
}

func (s *slis) fin_nxt(a int) int {
   var path []int
   x := a
   for !s.live[x] {
       path = append(path, x)
       x = s.nxt[x]
   }
   for _, y := range path {
       s.nxt[y] = x
   }
   return x
}

func (s *slis) que(l, r int) pair {
   left := s.getv(l)
   right := s.getv(r)
   t := l
   for t <= r {
       s.upd(t, -1)
       s.live[t] = false
       t = s.fin_nxt(s.nxt[t])
   }
   return pair{left, right}
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var s string
       fmt.Fscan(reader, &s)
       n := len(s)
       // 1-based indexing
       str := make([]byte, n+2)
       for i := 1; i <= n; i++ {
           str[i] = s[i-1]
       }

       var l mlist
       l.init(n, str)
       for i := 1; i < n; i++ {
           if str[i] == str[i+1] {
               l.insert(i)
           }
       }
       for {
           maxi, ci := 0, -1
           for i := 0; i < 26; i++ {
               if l.cnt[i] > maxi {
                   maxi = l.cnt[i]
                   ci = i
               }
           }
           if ci < 0 {
               break
           }
           l.que(ci)
       }

       var l2 slis
       l2.init(n)
       // output
       fmt.Fprintln(writer, len(l.ans)+1)
       remain := n
       for _, v := range l.ans {
           tmp := l2.que(v.l, v.r)
           fmt.Fprintln(writer, tmp.l, tmp.r)
           remain -= tmp.r - tmp.l + 1
       }
       fmt.Fprintln(writer, 1, remain)
   }
}
