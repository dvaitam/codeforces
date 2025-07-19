package main

import (
   "bufio"
   "fmt"
 "os"
)

type Frac struct{ x, y int64 }

func (a Frac) add(b Frac) Frac { return Frac{a.x + b.x, a.y + b.y} }
func (a Frac) less(b Frac) bool { return a.x*b.y < b.x*a.y }
func (a Frac) eq(b Frac) bool { return a.x*b.y == b.x*a.y }

var (
   n, r, K int
   a        []int
   e        [][]int
   dfn, out, rv, fa []int
   op, ban []bool
   s, w   []Frac
   dfnCnt int
)

func dfs(u, p int) {
   dfnCnt++
   dfn[u] = dfnCnt
   rv[dfnCnt] = u
   fa[u] = p
   s[u] = Frac{int64(a[u]), 1}
   for _, v := range e[u] {
       if v == p {
           continue
       }
       dfs(v, u)
       s[u] = s[u].add(s[v])
   }
   out[u] = dfnCnt
}

func calc(x int) {
   f1 := 0
   for i := x; i != 0; i = fa[i] {
       if op[i] {
           f1 = i
       }
   }
   if f1 == 0 {
       w[x] = Frac{int64(a[x]), 1}
       if K > 0 {
           best := w[x]
           for i := x; i != 0; i = fa[i] {
               if ban[i] {
                   continue
               }
               if s[i].less(best) || (best.less(w[x]) && s[i].eq(best)) {
                   best = s[i]
                   f1 = i
               }
           }
       }
       if f1 != 0 {
           K--
           w[x] = s[f1]
           op[f1] = true
       }
       return
   }
   tmp := make([]int, 0)
   rs := K + 1
   for i := dfn[f1]; i <= out[f1]; i++ {
       v := rv[i]
       if v < x {
           // compare w[v] < {a[v],1}
           if w[v].less(Frac{int64(a[v]), 1}) {
               f2 := 0
               for o := v; !(out[o] >= dfn[x] && dfn[o] <= dfn[x]); o = fa[o] {
                   if !ban[o] {
                       f2 = o
                   }
               }
               if f2 == 0 {
                   rs = -1
                   break
               }
               tmp = append(tmp, f2)
               rs--
               i = out[f2]
           }
       }
   }
   if rs < 0 {
       w[x] = s[f1]
       return
   }
   w[x] = Frac{int64(a[x]), 1}
   f2 := 0
   mn := w[x]
   if rs > 0 {
       for i := x; i != 0; i = fa[i] {
           if ban[i] {
               continue
           }
           if s[i].less(mn) {
               mn = s[i]
               f2 = i
           }
       }
   }
   if mn.less(s[f1]) {
       op[f1] = false
       for _, v := range tmp {
           op[v] = true
       }
       K = rs
       w[x] = mn
       if f2 != 0 {
           K--
           op[f2] = true
       }
   } else {
       w[x] = s[f1]
   }
}

func solve() {
   reader := bufio.NewReader(os.Stdin)
   var T int
   fmt.Fscan(reader, &T)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for T > 0 {
       T--
       fmt.Fscan(reader, &n, &r, &K)
       a = make([]int, n+1)
       e = make([][]int, n+1)
       dfn = make([]int, n+1)
       out = make([]int, n+1)
       rv = make([]int, n+1)
       fa = make([]int, n+1)
       op = make([]bool, n+1)
       ban = make([]bool, n+1)
       s = make([]Frac, n+1)
       w = make([]Frac, n+1)
       for i := 1; i <= n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       for i := 2; i <= n; i++ {
           var x, y int
           fmt.Fscan(reader, &x, &y)
           e[x] = append(e[x], y)
           e[y] = append(e[y], x)
       }
       dfnCnt = 0
       dfs(r, 0)
       for i := 1; i <= n; i++ {
           calc(i)
           for j := i; j != 0; j = fa[j] {
               if w[i].less(s[j]) {
                   ban[j] = true
               }
           }
       }
       res := make([]int, 0)
       for i := 1; i <= n; i++ {
           if op[i] {
               res = append(res, i)
           }
       }
       fmt.Fprintln(writer, len(res))
       for _, v := range res {
           fmt.Fprint(writer, v, " ")
       }
       fmt.Fprintln(writer)
   }
}

func main() {
   solve()
}
