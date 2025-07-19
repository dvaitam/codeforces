package main

import (
   "bufio"
   "fmt"
   "os"
)

const N = 500500

type Edge1 struct { s, n int }
type Edge2 struct { s, n int }
type Edge struct { s, n, x, y int }
type CC struct { x, f int }
type Op struct { ax, ay, bx, by int }

var (
   n, En1, En2, En int
   fa, f, h1, h2, fg, out, Fa, d, h, tail []int
   E1 []Edge1
   E2 []Edge2
   E  []Edge
   A  []CC
   ops []Op
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func readInt() int {
   var x int
   var ch byte
   for {
       b, err := reader.ReadByte()
       if err != nil {
           return x
       }
       ch = b
       if ch >= '0' && ch <= '9' {
           x = int(ch - '0')
           break
       }
   }
   for {
       b, err := reader.ReadByte()
       if err != nil {
           break
       }
       ch = b
       if ch < '0' || ch > '9' {
           break
       }
       x = x*10 + int(ch-'0')
   }
   return x
}

func writeInt(x int) {
   if x >= 10 {
       writeInt(x / 10)
   }
   writer.WriteByte(byte('0' + x%10))
}

func getf(x int) int {
   if f[x] != x {
       f[x] = getf(f[x])
   }
   return f[x]
}

func E_add1(x, y int) {
   En1++
   E1[En1].s = y
   E1[En1].n = h1[x]
   h1[x] = En1
   En1++
   E1[En1].s = x
   E1[En1].n = h1[y]
   h1[y] = En1
}

func E_add2(x, y int) {
   En2++
   E2[En2].s = y
   E2[En2].n = h2[x]
   h2[x] = En2
}

func E_add(x, y, tx, ty int) {
   if h[x] == 0 {
       tail[x] = En + 1
   }
   En++
   E[En].s = y
   E[En].n = h[x]
   E[En].x = tx
   E[En].y = ty
   h[x] = En
   if h[y] == 0 {
       tail[y] = En + 1
   }
   En++
   E[En].s = x
   E[En].n = h[y]
   E[En].x = tx
   E[En].y = ty
   h[y] = En
}

func main() {
   defer writer.Flush()
   n = readInt()
   // allocate
   fa = make([]int, n+1)
   f = make([]int, n+1)
   h1 = make([]int, n+1)
   h2 = make([]int, n+1)
   fg = make([]int, n+1)
   out = make([]int, n+1)
   Fa = make([]int, n+1)
   d = make([]int, n+1)
   h = make([]int, n+1)
   tail = make([]int, n+1)
   E1 = make([]Edge1, 2*(n+1))
   E2 = make([]Edge2, 2*(n+1))
   E = make([]Edge, 2*(n+1))
   A = make([]CC, n+1)
   ops = make([]Op, 0, n)
   ans := n - 1
   // init DSU
   for i := 1; i <= n; i++ {
       f[i] = i
   }
   // read first tree
   for i := 1; i < n; i++ {
       x := readInt()
       y := readInt()
       E_add1(x, y)
   }
   // dfs1 iterative
   stack := make([][2]int, 0, n)
   stack = append(stack, [2]int{1, 0})
   for len(stack) > 0 {
       top := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       u, p := top[0], top[1]
       fa[u] = p
       for k := h1[u]; k != 0; k = E1[k].n {
           v := E1[k].s
           if v == p {
               continue
           }
           stack = append(stack, [2]int{v, u})
       }
   }
   // read second tree edges
   for i := 1; i < n; i++ {
       x := readInt()
       y := readInt()
       if fa[x] == y && getf(x) != getf(y) {
           f[getf(x)] = getf(y)
           ans--
           continue
       }
       if fa[y] == x && getf(x) != getf(y) {
           f[getf(y)] = getf(x)
           ans--
           continue
       }
       E_add2(x, y)
   }
   // prepare components
   for i := 1; i <= n; i++ {
       if getf(i) != getf(fa[i]) {
           fg[i] = 1
           if fa[i] == 0 {
               continue
           }
           fi := getf(i)
           ffi := getf(fa[i])
           A[fi].x = i
           A[fi].f = fa[i]
           out[ffi]++
           Fa[fi] = ffi
       }
   }
   // build component graph
   for i := 1; i <= n; i++ {
       for k := h2[i]; k != 0; k = E2[k].n {
           j := E2[k].s
           fi := getf(i)
           fj := getf(j)
           if fi != fj {
               E_add(fi, fj, i, j)
           }
       }
   }
   // output ans
   writeInt(ans)
   writer.WriteByte('\n')
   // queue
   q := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if fg[i] == 1 && out[i] == 0 {
           q = append(q, i)
       }
   }
   qi := 0
   for qi < len(q) {
       u := q[qi]; qi++
       if u == 1 {
           break
       }
       // remove edge from second tree and add from first tree
       ex, ey := 0, 0
       // find an available edge
       for kk := h[u]; kk != 0; kk = E[kk].n {
           v := E[kk].s
           if getf(u) != getf(v) {
               ex = E[kk].x
               ey = E[kk].y
               // merge lists
               nei := getf(v)
               headU := getf(u)
               // splice adjacency lists
               E[tail[nei]].n = h[headU]
               tail[nei] = tail[headU]
               f[headU] = nei
               break
           }
       }
       // record operation: add A[u] then remove ex,ey
       ops = append(ops, Op{A[u].x, A[u].f, ex, ey})
       // decrease out-degree
       parentComp := Fa[u]
       out[parentComp]--
       if out[parentComp] == 0 {
           q = append(q, parentComp)
       }
   }
   // print operations
   for _, o := range ops {
       writeInt(o.ax); writer.WriteByte(' ')
       writeInt(o.ay); writer.WriteByte(' ')
       writeInt(o.bx); writer.WriteByte(' ')
       writeInt(o.by); writer.WriteByte('\n')
   }
}
