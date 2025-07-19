package main

import (
   "bufio"
   "os"
   "strconv"
)

var (
   rd = bufio.NewReader(os.Stdin)
   wr = bufio.NewWriter(os.Stdout)
)

func readInt() int {
   neg := false
   b, _ := rd.ReadByte()
   for (b < '0' || b > '9') && b != '-' {
       b, _ = rd.ReadByte()
   }
   if b == '-' {
       neg = true
       b, _ = rd.ReadByte()
   }
   x := int(b - '0')
   for {
       b, _ = rd.ReadByte()
       if b < '0' || b > '9' {
           break
       }
       x = x*10 + int(b-'0')
   }
   if neg {
       return -x
   }
   return x
}

type edge struct{ to, nxt int }
type nodeOp struct{ a, b, tp int }

func main() {
   defer wr.Flush()
   T := readInt()
   for T > 0 {
       T--
       solve()
   }
}

func solve() {
   n := readInt()
   // allocate per test
   head := make([]int, n+1)
   deg := make([]int, n+1)
   nx := make([]int, n+1)
   fa := make([]int, n+1)
   v := make([]bool, n+1)
   edges := make([]edge, 2*(n+1))
   aops := make([]nodeOp, 0, n)
   et, at := 0, 0

   addEdge := func(x, y int) {
       et++
       edges[et] = edge{to: y, nxt: head[x]}
       head[x] = et
       deg[x]++
   }
   // read edges
   for i := 1; i < n; i++ {
       x := readInt(); y := readInt()
       addEdge(x, y)
       addEdge(y, x)
   }
   // find root
   rt := 1
   for i := 1; i <= n; i++ {
       if deg[i] >= 2 {
           rt = i
           break
       }
   }
   // get representative
   var getRep func(int) int
   getRep = func(x int) int {
       for nx[x] != 0 {
           x = nx[x]
       }
       return x
   }
   // dfs recursive
   var dfs func(int, int)
   dfs = func(x, p int) {
       fa[x] = p
       var a, b int = -1, -1
       for i := head[x]; i != 0; i = edges[i].nxt {
           y := edges[i].to
           if y == p {
               continue
           }
           dfs(y, x)
           if v[y] {
               continue
           }
           if a == -1 {
               a = y
           } else if b == -1 {
               b = y
           } else {
               // extra child
               at++
               aops = append(aops, nodeOp{a: y, b: getRep(y), tp: y})
           }
       }
       if b != -1 {
           at++
           aops = append(aops, nodeOp{a: getRep(a), b: getRep(b), tp: x})
           v[x] = true
       } else if a != -1 {
           nx[x] = a
       } else {
           nx[x] = 0
       }
   }
   dfs(rt, 0)
   // ensure root op
   if at == 0 || aops[at-1].tp != rt {
       // append root
       at++
       aops = append(aops, nodeOp{a: rt, b: getRep(rt), tp: rt})
   }
   // output
   // number of operations = at-1
   wr.WriteString(strconv.Itoa(at - 1))
   wr.WriteByte('\n')
   ls := aops[at-1].b
   // print in reverse order
   for i := at - 2; i >= 0; i-- {
       op := aops[i]
       // op.tp, fa[op.tp], ls, op.b
       wr.WriteString(strconv.Itoa(op.tp))
       wr.WriteByte(' ')
       wr.WriteString(strconv.Itoa(fa[op.tp]))
       wr.WriteByte(' ')
       wr.WriteString(strconv.Itoa(ls))
       wr.WriteByte(' ')
       wr.WriteString(strconv.Itoa(op.b))
       wr.WriteByte('\n')
       ls = op.a
   }
}
