package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 1000000000

type deque struct {
   data       []int
   head, tail int
}

func newDeque(n int) *deque {
   size := 2*n + 10
   d := &deque{data: make([]int, size)}
   d.head = size / 2
   d.tail = d.head
   return d
}

func (d *deque) empty() bool {
   return d.head == d.tail
}

func (d *deque) pushFront(x int) {
   d.head--
   d.data[d.head] = x
}

func (d *deque) pushBack(x int) {
   d.data[d.tail] = x
   d.tail++
}

func (d *deque) popFront() int {
   x := d.data[d.head]
   d.head++
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n, m int
       fmt.Fscan(reader, &n, &m)
       s := make([][]byte, n)
       for i := 0; i < n; i++ {
           var row string
           fmt.Fscan(reader, &row)
           s[i] = []byte(row)
       }
       nm := n * m
       dis := make([]int, nm)
       par := make([]int, nm)
       mark := make([]bool, nm)
       for i := 0; i < nm; i++ {
           dis[i] = INF
       }
       dq := newDeque(nm)
       // initial states: left column
       for i := 0; i < n; i++ {
           idx := i*m
           dis[idx] = 0
           if isok(i, 0, n, m, s) {
               if s[i][0] != '#' {
                   dis[idx] = 1
                   dq.pushBack(idx)
               } else {
                   dq.pushFront(idx)
               }
           }
       }
       // 0-1 BFS
       for !dq.empty() {
           x := dq.popFront()
           if x%m == m-1 || mark[x] {
               if mark[x] {
                   continue
               }
               // if reached right border, skip marking further
               continue
           }
           mark[x] = true
           r := x / m
           c := x % m
           // diagonal moves
           for _, dxy := range [][2]int{{-1, 1}, {1, 1}, {1, -1}, {-1, -1}} {
               nr := r + dxy[0]
               nc := c + dxy[1]
               if nr < 0 || nr >= n || nc < 0 || nc >= m {
                   continue
               }
               ni := nr*m + nc
               if mark[ni] {
                   continue
               }
               if !isok(nr, nc, n, m, s) {
                   continue
               }
               w := 0
               if s[nr][nc] != '#' {
                   w = 1
               }
               nd := dis[x] + w
               if nd < dis[ni] {
                   dis[ni] = nd
                   par[ni] = x
                   if w == 0 {
                       dq.pushFront(ni)
                   } else {
                       dq.pushBack(ni)
                   }
               }
           }
       }
       // find best end
       ans := INF
       stp := -1
       for i := 0; i < n; i++ {
           idx := i*m + (m - 1)
           if dis[idx] < ans {
               ans = dis[idx]
               stp = idx
           }
       }
       if ans == INF {
           fmt.Fprintln(writer, "NO")
           continue
       }
       // reconstruct path
       fmt.Fprintln(writer, "YES")
       x := stp
       for x%m != 0 {
           r := x / m
           c := x % m
           s[r][c] = '#'
           x = par[x]
       }
       s[x/m][0] = '#'
       // output grid
       for i := 0; i < n; i++ {
           writer.Write(s[i])
           writer.WriteByte('\n')
       }
   }
}

func isok(r, c, n, m int, s [][]byte) bool {
   // no adjacent cardinal neighbor is '#'
   for _, dxy := range [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
       nr := r + dxy[0]
       nc := c + dxy[1]
       if nr >= 0 && nr < n && nc >= 0 && nc < m && s[nr][nc] == '#' {
           return false
       }
   }
   return true
}
