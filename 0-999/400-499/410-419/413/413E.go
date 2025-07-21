package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

const INF = 1000000000

// Matrix represents distances from rows at left boundary to rows at right boundary
type Matrix [2][2]int

func combine(a, b Matrix) Matrix {
   var c Matrix
   for i := 0; i < 2; i++ {
       for j := 0; j < 2; j++ {
           mn := INF
           for k := 0; k < 2; k++ {
               d := a[i][k] + b[k][j]
               if d < mn {
                   mn = d
               }
           }
           c[i][j] = mn
       }
   }
   return c
}

func identity() Matrix {
   var m Matrix
   for i := 0; i < 2; i++ {
       for j := 0; j < 2; j++ {
           if i == j {
               m[i][j] = 0
           } else {
               m[i][j] = INF
           }
       }
   }
   return m
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // read ints
   nextInt := func() int {
       b, _ := reader.ReadBytes(' ')
       for len(b) > 0 && (b[len(b)-1] == ' ' || b[len(b)-1] == '\n' || b[len(b)-1] == '\r') {
           b = b[:len(b)-1]
       }
       x, _ := strconv.Atoi(string(b))
       return x
   }
   // read first line
   var n, m int
   // custom scan
   fmt.Fscan(reader, &n, &m)
   grid := make([]string, 2)
   fmt.Fscan(reader, &grid[0], &grid[1])
   // special case n == 1
   if n == 1 {
       for i := 0; i < m; i++ {
           var v, u int
           fmt.Fscan(reader, &v, &u)
           r1, r2 := 0, 0
           if v > 1 { r1 = 1 }
           if u > 1 { r2 = 1 }
           if v == u {
               fmt.Fprintln(writer, 0)
           } else {
               if grid[0][0] == '.' && grid[1][0] == '.' {
                   fmt.Fprintln(writer, 1)
               } else {
                   fmt.Fprintln(writer, -1)
               }
           }
       }
       return
   }
   // build leaf matrices for columns i->i+1, i=0..n-2
   size := 1
   for size < n-1 {
       size <<= 1
   }
   seg := make([]Matrix, 2*size)
   // initialize leaves
   for i := 0; i < n-1; i++ {
       // build matrix for column i to i+1
       // nodes: 0:(0,i),1:(1,i) -> 2:(0,i+1),3:(1,i+1)
       // adjacency list
       adj := [4][]int{}
       // vertical at i
       if grid[0][i] == '.' && grid[1][i] == '.' {
           adj[0] = append(adj[0], 1)
           adj[1] = append(adj[1], 0)
       }
       // vertical at i+1
       if grid[0][i+1] == '.' && grid[1][i+1] == '.' {
           adj[2] = append(adj[2], 3)
           adj[3] = append(adj[3], 2)
       }
       // horizontal row 0
       if grid[0][i] == '.' && grid[0][i+1] == '.' {
           adj[0] = append(adj[0], 2)
           adj[2] = append(adj[2], 0)
       }
       // horizontal row 1
       if grid[1][i] == '.' && grid[1][i+1] == '.' {
           adj[1] = append(adj[1], 3)
           adj[3] = append(adj[3], 1)
       }
       // compute distances from 0,1 to 2,3
       mat := Matrix{}
       for a := 0; a < 2; a++ {
           dist := [4]int{INF, INF, INF, INF}
           // BFS
           q := make([]int, 0, 4)
           dist[a] = 0
           q = append(q, a)
           for qi := 0; qi < len(q); qi++ {
               u := q[qi]
               for _, v := range adj[u] {
                   if dist[v] > dist[u]+1 {
                       dist[v] = dist[u] + 1
                       q = append(q, v)
                   }
               }
           }
           for b := 0; b < 2; b++ {
               mat[a][b] = dist[2+b]
           }
       }
       seg[size+i] = mat
   }
   // fill remaining leaves with identity
   id := identity()
   for i := n - 1; i < size; i++ {
       seg[size+i] = id
   }
   // build internal nodes
   for i := size - 1; i > 0; i-- {
       seg[i] = combine(seg[2*i], seg[2*i+1])
   }
   // process queries
   for qi := 0; qi < m; qi++ {
       var v, u int
       fmt.Fscan(reader, &v, &u)
       // decode row and col
       c1 := (v-1) % n
       r1 := 0
       if v-1 >= n {
           r1 = 1
       }
       c2 := (u-1) % n
       r2 := 0
       if u-1 >= n {
           r2 = 1
       }
       if c1 == c2 {
           if r1 == r2 {
               fmt.Fprintln(writer, 0)
           } else if grid[0][c1] == '.' && grid[1][c1] == '.' {
               fmt.Fprintln(writer, 1)
           } else {
               fmt.Fprintln(writer, -1)
           }
           continue
       }
       // determine query range
       var l, r int
       rev := false
       if c1 < c2 {
           l = c1
           r = c2 - 1
       } else {
           l = c2
           r = c1 - 1
           rev = true
       }
       // query segtree for [l..r-1] in leaves representation (i.e., indices l..r-1 are leaves)
       // but our leaves are at [0..n-2] representing i->i+1, so want [l..r-1] = [l..r-1]? Actually l..r-1 = c1..c2-1? no.
       // We need segments covering columns from lower to higher: [lower..higher-1] which is l..r-1? Actually above we set r=c2-1 if c1<c2, so leaves l..r-1? Mistake.
       // Actually want leaves indices from l to r-1 inclusive? Let's adjust: for c1<c2, want leaves from c1 to c2-1 inclusive: leaves at i=c1..c2-1. But c1 zero-based. So leaves index from c1 to c2-1.
       // With l=c1,r=c2-1, we use leaves [l..r-1]? No, we want [l..r]? So use [l..r].
       // Similarly for reverse.
       // Thus leaves index is [l..r].
       L := l
       R := r
       L += size
       R += size
       // initialize results
       resL := id
       resR := id
       for L <= R {
           if L&1 == 1 {
               resL = combine(resL, seg[L])
               L++
           }
           if R&1 == 0 {
               resR = combine(seg[R], resR)
               R--
           }
           L >>= 1
           R >>= 1
       }
       M := combine(resL, resR)
       var ans int
       if !rev {
           ans = M[r1][r2]
       } else {
           ans = M[r2][r1]
       }
       if ans >= INF/2 {
           fmt.Fprintln(writer, -1)
       } else {
           fmt.Fprintln(writer, ans)
       }
   }
}
