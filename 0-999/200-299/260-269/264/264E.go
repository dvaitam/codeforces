package main

import (
   "bufio"
   "fmt"
   "os"
)

const H = 10

type Node [H][H]uint8

var (
   n, m int
   tree []Node
   bit  []int
   outBuf *bufio.Writer
)

func bitAdd(i, v int) {
   for ; i <= n; i += i & -i {
       bit[i] += v
   }
}

// find smallest p such that sum[1..p] >= k
func bitKth(k int) int {
   pos := 0
   sum := 0
   // highest bit of n
   bitMax := 1
   for bitMax<<1 <= n {
       bitMax <<= 1
   }
   for d := bitMax; d > 0; d >>= 1 {
       if pos+d <= n && sum+bit[pos+d] < k {
           sum += bit[pos+d]
           pos += d
       }
   }
   return pos + 1
}

// merge left and right into dst
func merge(dst *Node, a, b *Node) {
   // for all l<=r: dst[l][r] = max over k in [l..r] a[l][k] + b[k][r]
   for l := 0; l < H; l++ {
       for r := l; r < H; r++ {
           maxv := uint8(0)
           for k := l; k <= r; k++ {
               v := a[l][k] + b[k][r]
               if v > maxv {
                   maxv = v
               }
           }
           dst[l][r] = maxv
       }
   }
}

// update position pos (0-indexed). h: 0..H-1 for insert, -1 for delete
func update(idx, l, r, pos, h int) {
   if l == r {
       // leaf
       var empty Node
       if h < 0 {
           tree[idx] = empty
       } else {
           // single height h
           var nd Node
           for i := 0; i < H; i++ {
               for j := i; j < H; j++ {
                   if i <= h && h <= j {
                       nd[i][j] = 1
                   } else {
                       nd[i][j] = 0
                   }
               }
           }
           tree[idx] = nd
       }
       return
   }
   mid := (l + r) >> 1
   if pos <= mid {
       update(idx<<1, l, mid, pos, h)
   } else {
       update(idx<<1|1, mid+1, r, pos, h)
   }
   // merge children
   merge(&tree[idx], &tree[idx<<1], &tree[idx<<1|1])
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   outBuf = bufio.NewWriter(os.Stdout)
   defer outBuf.Flush()
   fmt.Fscan(reader, &n, &m)
   // init structures
   size := 1
   for size < n {
       size <<= 1
   }
   tree = make([]Node, size<<1)
   bit = make([]int, n+1)
   // process queries
   for qi := 0; qi < m; qi++ {
       var t int
       fmt.Fscan(reader, &t)
       if t == 1 {
           var p, h int
           fmt.Fscan(reader, &p, &h)
           // insert
           bitAdd(p, 1)
           update(1, 0, size-1, p-1, h-1)
       } else {
           var x int
           fmt.Fscan(reader, &x)
           p := bitKth(x)
           bitAdd(p, -1)
           update(1, 0, size-1, p-1, -1)
       }
       // answer is dp[0][H-1] at root
       ans := tree[1][0][H-1]
       // print ans (0..10)
       if ans < 10 {
           outBuf.WriteByte(byte('0' + ans))
       } else {
           outBuf.WriteString("10")
       }
       if qi+1 < m {
           outBuf.WriteByte(' ')
       }
   }
}
