package main

import (
   "bufio"
   "fmt"
   "os"
)

// Conveyor mapping kinds
const (
   K_UP   = 0
   K_OUT  = 1
   K_LOOP = 2
)

// CellRes holds the result of starting at a cell: up, exit, or loop
type CellRes struct {
   kind int
   row  int // for K_OUT: exit row
   col  int // 0-based col position; for K_UP: next col; for K_OUT: exit colPos (-1 or m)
}

func main() {
   rdr := bufio.NewReader(os.Stdin)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()

   var n, m, q int
   fmt.Fscan(rdr, &n, &m, &q)
   belts := make([][]byte, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(rdr, &s)
       belts[i] = []byte(s)
   }
   // build segment tree
   size := 1
   for size < n {
       size <<= 1
   }
   // tree nodes: each mapping slice length m
   tree := make([][]CellRes, 2*size)
   // identity mapping
   for i := range tree {
       tree[i] = make([]CellRes, m)
       for j := 0; j < m; j++ {
           tree[i][j] = CellRes{kind: K_UP, col: j}
       }
   }
   // build leaves
   for i := 1; i <= n; i++ {
       res := buildRow(belts, i, m)
       pos := size + i - 1
       for j := 0; j < m; j++ {
           tree[pos][j] = res[j]
       }
   }
   // build internal nodes
   for i := size - 1; i >= 1; i-- {
       left := tree[2*i]
       right := tree[2*i+1]
       tree[i] = combine(left, right, m)
   }

   // process queries
   for k := 0; k < q; k++ {
       var op string
       fmt.Fscan(rdr, &op)
       if op[0] == 'A' {
           var x, y int
           fmt.Fscan(rdr, &x, &y)
           // query prefix [1,x]
           cur := make([]CellRes, m)
           for j := 0; j < m; j++ {
               cur[j] = CellRes{kind: K_UP, col: j}
           }
           // segment tree indices
           l := size
           r := size + x - 1
           // collect nodes in bottom-up order
           var nodes []int
           for l <= r {
               if l&1 == 1 {
                   // should not happen for prefix l=size
                   nodes = append(nodes, l)
                   l++
               }
               if r&1 == 0 {
                   nodes = append(nodes, r)
                   r--
               }
               l >>= 1
               r >>= 1
           }
           // debug: print segment nodes
           fmt.Fprintln(w, "DEBUG nodes:", nodes)
           // combine in collected order
           for _, idx := range nodes {
               cur = combine(tree[idx], cur, m)
           }
           // result at y-1
           res := cur[y-1]
           if res.kind == K_LOOP {
               fmt.Fprintln(w, -1, -1)
           } else if res.kind == K_OUT {
               // output exit position
               colOut := res.col
               // convert 0-based exit col to 1-based or 0/m+1
               var cout int
               if colOut < 0 {
                   cout = 0
               } else if colOut >= m {
                   cout = m + 1
               } else {
                   cout = colOut + 1
               }
               fmt.Fprintln(w, res.row, cout)
           } else {
               // UP beyond row 1
               var cout = res.col + 1
               fmt.Fprintln(w, 0, cout)
           }
       } else {
           // change
           var x, y int
           var cc string
           fmt.Fscan(rdr, &x, &y, &cc)
           cb := cc[0]
           belts[x-1][y-1] = cb
           // rebuild leaf and path
           pos := size + x - 1
           rowMap := buildRow(belts, x, m)
           for j := 0; j < m; j++ {
               tree[pos][j] = rowMap[j]
           }
           pos >>= 1
           for pos >= 1 {
               left := tree[2*pos]
               right := tree[2*pos+1]
               tree[pos] = combine(left, right, m)
               pos >>= 1
           }
       }
   }
}

// buildRow computes mapping for a single row i (1-based)
func buildRow(belts [][]byte, i, m int) []CellRes {
   res := make([]CellRes, m)
   row := belts[i-1]
   for j := 0; j < m; j++ {
       visited := make([]bool, m)
       pos := j
       for {
           if pos < 0 || pos >= m {
               res[j] = CellRes{kind: K_OUT, row: i, col: pos}
               break
           }
           if visited[pos] {
               res[j] = CellRes{kind: K_LOOP}
               break
           }
           visited[pos] = true
           switch row[pos] {
           case '^':
               res[j] = CellRes{kind: K_UP, col: pos}
               pos = -100 // to exit
           case '<':
               pos--
           case '>':
               pos++
           }
           if res[j].kind == K_UP {
               break
           }
       }
   }
   return res
}

// combine composes upper (A) and lower (B) mappings
func combine(upper, lower []CellRes, m int) []CellRes {
   out := make([]CellRes, m)
   for j := 0; j < m; j++ {
       b := lower[j]
       switch b.kind {
       case K_OUT, K_LOOP:
           out[j] = b
       case K_UP:
           a := upper[b.col]
           out[j] = a
       }
   }
   return out
}
