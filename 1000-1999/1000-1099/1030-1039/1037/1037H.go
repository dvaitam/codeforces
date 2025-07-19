package main

import (
   "bufio"
   "fmt"
   "os"
)

const N = 223456

// segment tree node
type tnode struct {
   l, r *tnode
}

// automaton node
type node struct {
   goCh [26]*node
   p    *node
   r    *tnode
   val  int
   son  []*node
}

var (
   tpool [N * 50]tnode
   tcur  int
   npool [N * 2]node
   cur   int
   rt    *node
   s     []rune
   ch    [N]int
   w     [N]int
   n, m, L, R int
)

func newTNode() *tnode {
   p := &tpool[tcur]
   tcur++
   return p
}

func newNode() *node {
   p := &npool[cur]
   cur++
   return p
}

// extend automaton
func extend(p *node, c int) *node {
   np := newNode()
   np.val = p.val + 1
   for ; p != nil && p.goCh[c] == nil; p = p.p {
       p.goCh[c] = np
   }
   if p == nil {
       np.p = rt
   } else {
       q := p.goCh[c]
       if q.val == p.val+1 {
           np.p = q
       } else {
           nq := newNode()
           *nq = *q
           nq.val = p.val + 1
           q.p = nq
           np.p = nq
           for ; p != nil && p.goCh[c] == q; p = p.p {
               p.goCh[c] = nq
           }
       }
   }
   return np
}

// build a segment tree with single 1 at x
func build(l, r, x int) *tnode {
   p := newTNode()
   if l != r {
       mid := (l + r) >> 1
       if x <= mid {
           p.l = build(l, mid, x)
       } else {
           p.r = build(mid+1, r, x)
       }
   }
   return p
}

// merge two segment trees
func merge(a, b *tnode) *tnode {
   if a == nil {
       return b
   }
   if b == nil {
       return a
   }
   p := newTNode()
   p.l = merge(a.l, b.l)
   p.r = merge(a.r, b.r)
   return p
}

// exists any mark in [L,R]
func exists(p *tnode, l, r, L, R int) bool {
   if p == nil {
       return false
   }
   if L == l && R == r {
       return true
   }
   mid := (l + r) >> 1
   if R <= mid {
       return exists(p.l, l, mid, L, R)
   } else if L > mid {
       return exists(p.r, mid+1, r, L, R)
   }
   return exists(p.l, l, mid, L, mid) || exists(p.r, mid+1, r, mid+1, R)
}

// dfs to merge children segment trees
func dfsBuild(p *node) {
   for _, c := range p.son {
       dfsBuild(c)
       p.r = merge(p.r, c.r)
   }
}

// initialize automaton and trees
func initAutomaton() {
   rt = newNode()
   cur = 1
   tcur = 0
   rt.val = 0
   nodePtr := rt
   for i := 1; i <= n; i++ {
       nodePtr = extend(nodePtr, int(s[i]-'a'))
       nodePtr.r = build(1, n, i)
   }
   // build suffix link tree
   for i := 0; i < cur; i++ {
       p := &npool[i]
       if p.p != nil {
           p.p.son = append(p.p.son, p)
       }
   }
   dfsBuild(rt)
}

// check validity for next char depth d at node p
func valid(d int, p *node) bool {
   if p != nil && L+d-1 <= R {
       return exists(p.r, 1, n, L+d-1, R)
   }
   return false
}

// search answer
func dfsSearch(dep int, ty bool, p *node, out *bufio.Writer) bool {
   if ty {
       for i := 0; i < dep; i++ {
           out.WriteByte(byte(ch[i] + 'a'))
       }
       out.WriteByte('\n')
       return true
   }
   var start int
   if ty {
       start = 0
   } else {
       if w[dep] > 0 {
           start = w[dep]
       } else {
           start = 0
       }
   }
   for c := start; c < 26; c++ {
       if valid(dep+1, p.goCh[c]) {
           ch[dep] = c
           nextTy := ty || (c > w[dep])
           if dfsSearch(dep+1, nextTy, p.goCh[c], out) {
               return true
           }
       }
   }
   return false
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var str string
   fmt.Fscan(in, &str)
   n = len(str)
   s = make([]rune, n+1)
   for i, r := range str {
       s[i+1] = r
   }
   initAutomaton()
   var q int
   fmt.Fscan(in, &q)
   for i := 0; i < q; i++ {
       var tstr string
       fmt.Fscan(in, &L, &R, &tstr)
       m = len(tstr)
       for j, r := range tstr {
           w[j] = int(r - 'a')
       }
       w[m] = -1
       if !dfsSearch(0, false, rt, out) {
           out.WriteString("-1\n")
       }
   }
}
