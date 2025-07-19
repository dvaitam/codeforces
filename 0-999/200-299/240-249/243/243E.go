package main

import (
   "fmt"
   "os"
)

const maxNodes = 5000

type node struct {
   ch   []int
   tag  bool
   col  int
   leaf int
}

var n, root, dn int
var bo []bool
var ss [][]byte
var d [maxNodes]node
var ans []int

func gameover() {
   fmt.Println("NO")
   os.Exit(0)
}

func link(u, v int) {
   if v != 0 {
       d[u].ch = append(d[u].ch, v)
   }
}

func merge(u *[]int, v []int) {
   for _, x := range v {
       if x != 0 {
           *u = append(*u, x)
       }
   }
}

func chk(u int) bool {
   a := make([]int, len(d[u].ch))
   for j, v := range d[u].ch {
       a[j] = d[v].col
   }
   o := 0
   for j := range a {
       if a[j] > a[o] {
           o = j
       }
   }
   for j := 0; j < o; j++ {
       if a[j] > a[j+1] {
           return false
       }
   }
   for j := o; j < len(a)-1; j++ {
       if a[j] < a[j+1] {
           return false
       }
   }
   return true
}

func srt(u int) bool {
   a := make([]int, len(d[u].ch))
   z0, z1 := true, true
   for j, v := range d[u].ch {
       a[j] = d[v].col
       if j > 0 && a[j] < a[j-1] {
           z0 = false
       }
       if j > 0 && a[j] > a[j-1] {
           z1 = false
       }
   }
   if z0 {
       return true
   }
   if z1 {
       // reverse children order
       for i, j := 0, len(d[u].ch)-1; i < j; i, j = i+1, j-1 {
           d[u].ch[i], d[u].ch[j] = d[u].ch[j], d[u].ch[i]
       }
       return true
   }
   return false
}

func newNode() int {
   dn++
   // initialize fields
   d[dn].ch = d[dn].ch[:0]
   d[dn].tag = false
   d[dn].col = 0
   d[dn].leaf = 0
   return dn
}

func build() {
   root = newNode()
   for i := 1; i <= n; i++ {
       nd := newNode()
       link(root, nd)
       d[nd].leaf = i
   }
}

func paint(u int) {
   if d[u].leaf != 0 {
       id := d[u].leaf
       if bo[id] {
           d[u].col = 2
       } else {
           d[u].col = 0
       }
       return
   }
   cnt := [3]int{}
   for _, v := range d[u].ch {
       paint(v)
       cnt[d[v].col]++
   }
   d[u].col = 1
   if cnt[0] == 0 && cnt[1] == 0 {
       d[u].col = 2
   }
   if cnt[1] == 0 && cnt[2] == 0 {
       d[u].col = 0
   }
}

func bud(v []int) int {
   if len(v) == 0 {
       return 0
   }
   if len(v) == 1 {
       return v[0]
   }
   nd := newNode()
   for _, x := range v {
       link(nd, x)
   }
   return nd
}

func pushrig(u int) []int {
   if u == 0 {
       return nil
   }
   if d[u].tag && !srt(u) {
       gameover()
   }
   cnt := [3]int{}
   ls := [3][]int{}
   for _, v := range d[u].ch {
       c := d[v].col
       cnt[c]++
       ls[c] = append(ls[c], v)
   }
   if cnt[1] > 1 {
       gameover()
   }
   var res []int
   if !d[u].tag {
       res = append(res, bud(ls[0]))
       if cnt[1] == 1 {
           res = append(res, pushrig(ls[1][0])...)
       }
       res = append(res, bud(ls[2]))
   } else {
       merge(&res, ls[0])
       if cnt[1] == 1 {
           merge(&res, pushrig(ls[1][0]))
       }
       merge(&res, ls[2])
   }
   return res
}

func setcond(u int) {
   if d[u].col != 1 {
       return
   }
   sc, scn := 0, 0
   cnt := [3]int{}
   for _, v := range d[u].ch {
       c := d[v].col
       cnt[c]++
       if c != 0 {
           sc = v
           scn++
       }
   }
   if scn == 1 {
       setcond(sc)
       return
   }
   if cnt[1] > 2 {
       gameover()
   }
   var ls [3][]int
   for _, v := range d[u].ch {
       ls[d[v].col] = append(ls[d[v].col], v)
   }
   // pad ls[1] with two zeros
   ls[1] = append(ls[1], 0, 0)
   if !d[u].tag {
       // rebuild children
       d[u].ch = d[u].ch[:0]
       a := bud(ls[2])
       if a != 0 {
           d[a].col = 2
       }
       if len(ls[1]) > 2 {
           b := newNode()
           link(b, ls[1][0])
           link(b, a)
           link(b, ls[1][1])
           d[b].col = 1
           d[b].tag = true
           link(u, b)
           merge(&d[u].ch, ls[0])
           setcond(b)
           return
       } else {
           link(u, a)
           merge(&d[u].ch, ls[0])
           return
       }
   }
   // now d[u].tag must be true
   var newch []int
   flag := false
   for _, v := range d[u].ch {
       c := d[v].col
       if c == 1 {
           res := pushrig(v)
           if flag {
               // reverse res
               for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
                   res[i], res[j] = res[j], res[i]
               }
           }
           merge(&newch, res)
           flag = true
       } else {
           flag = flag || (c == 2)
           newch = append(newch, v)
       }
   }
   d[u].ch = newch
}

func workans(u int) {
   if d[u].leaf != 0 {
       ans = append(ans, d[u].leaf)
       return
   }
   for _, v := range d[u].ch {
       workans(v)
   }
}

func main() {
   if _, err := fmt.Fscan(os.Stdin, &n); err != nil {
       return
   }
   bo = make([]bool, n+1)
   ss = make([][]byte, n+1)
   for i := 1; i <= n; i++ {
       var line string
       fmt.Fscan(os.Stdin, &line)
       // 1-based indexing
       ss[i] = make([]byte, len(line)+1)
       copy(ss[i][1:], line)
   }
   build()
   for i := 1; i <= n; i++ {
       for j := 1; j <= n; j++ {
           bo[j] = (ss[i][j] == '1')
       }
       paint(root)
       setcond(root)
   }
   ans = make([]int, 0, n)
   workans(root)
   fmt.Println("YES")
   for i := 1; i <= n; i++ {
       for j := 0; j < n; j++ {
           fmt.Printf("%c", ss[i][ans[j]])
       }
       fmt.Println()
   }
}
