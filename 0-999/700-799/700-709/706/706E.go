package main

import (
   "bufio"
   "fmt"
   "os"
)

// Node represents a cell in the matrix with links to the right and down
type Node struct {
   right, down *Node
   val         int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, q int
   fmt.Fscan(reader, &n, &m, &q)
   // create nodes with dummy border
   nodes := make([][]*Node, n+2)
   for i := 0; i < n+2; i++ {
       nodes[i] = make([]*Node, m+2)
       for j := 0; j < m+2; j++ {
           nodes[i][j] = &Node{}
       }
   }
   // assign values for inner grid
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           fmt.Fscan(reader, &nodes[i][j].val)
       }
   }
   // link pointers
   for i := 0; i < n+2; i++ {
       for j := 0; j < m+2; j++ {
           if j+1 < m+2 {
               nodes[i][j].right = nodes[i][j+1]
           }
           if i+1 < n+2 {
               nodes[i][j].down = nodes[i+1][j]
           }
       }
   }
   // process queries
   for qi := 0; qi < q; qi++ {
       var a, b, c, d, h, w int
       fmt.Fscan(reader, &a, &b, &c, &d, &h, &w)
       // find predecessors (node before top-left corners)
       p1 := nodes[0][0]
       for i := 0; i < a-1; i++ {
           p1 = p1.down
       }
       for j := 0; j < b-1; j++ {
           p1 = p1.right
       }
       p2 := nodes[0][0]
       for i := 0; i < c-1; i++ {
           p2 = p2.down
       }
       for j := 0; j < d-1; j++ {
           p2 = p2.right
       }
       // swap top edges (down pointers)
       x1, x2 := p1, p2
       for i := 0; i < w; i++ {
           x1 = x1.right
           x2 = x2.right
           x1.down, x2.down = x2.down, x1.down
       }
       // swap right edges (right pointers)
       y1, y2 := p1, p2
       for i := 0; i < h; i++ {
           y1 = y1.down
           y2 = y2.down
           y1.right, y2.right = y2.right, y1.right
       }
       // swap bottom edges (down pointers)
       z1, z2 := p1, p2
       for i := 0; i < h; i++ {
           z1 = z1.down
           z2 = z2.down
       }
       for i := 0; i < w; i++ {
           z1 = z1.right
           z2 = z2.right
           z1.down, z2.down = z2.down, z1.down
       }
       // swap left edges (right pointers)
       w1, w2 := p1, p2
       for i := 0; i < w; i++ {
           w1 = w1.right
           w2 = w2.right
       }
       for i := 0; i < h; i++ {
           w1 = w1.down
           w2 = w2.down
           w1.right, w2.right = w2.right, w1.right
       }
   }
   // output resulting matrix
   // start at dummy head
   for i := 1; i <= n; i++ {
       cur := nodes[0][0]
       // move to row i
       for k := 0; k < i; k++ {
           cur = cur.down
       }
       // now cur is at start of row (at col 0)
       for j := 1; j <= m; j++ {
           cur = cur.right
           fmt.Fprintf(writer, "%d", cur.val)
           if j < m {
               writer.WriteByte(' ')
           }
       }
       writer.WriteByte('\n')
   }
}
