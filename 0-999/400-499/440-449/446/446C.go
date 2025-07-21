package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000009

var (
   n, m int
   a    []int64
   fib  []int64
   sumF []int64 // prefix sum of fib
   tree []int64
   tagA []int64
   tagB []int64
)

func sumAdd(length int, A, B int64) int64 {
   if length == 1 {
       return A % mod
   }
   // sum = A*(1 + sum fib[1..length-2]) + B*(sum fib[1..length-1])
   sA := (1 + sumF[length-2]) % mod
   sB := sumF[length-1]
   return (A*sA + B*sB) % mod
}

func push(v, l, r int) {
   A := tagA[v]
   B := tagB[v]
   if A == 0 && B == 0 {
       return
   }
   mid := (l + r) >> 1
   left, right := v<<1, v<<1|1
   // left child covers length lenL
   lenL := mid - l + 1
   // apply to left
   addL := sumAdd(lenL, A, B)
   tree[left] = (tree[left] + addL) % mod
   tagA[left] = (tagA[left] + A) % mod
   tagB[left] = (tagB[left] + B) % mod
   // prepare for right child
   lenR := r - mid
   // offset = position in sequence for mid+1: offset = lenL+1
   offset := lenL + 1
   var A2, B2 int64
   switch {
   case offset == 1:
       A2, B2 = A, B
   case offset == 2:
       A2 = B
       B2 = (A + B) % mod
   default:
       // offset >=3
       A2 = (A*fib[offset-2] + B*fib[offset-1]) % mod
       B2 = (A*fib[offset-1] + B*fib[offset]) % mod
   }
   // apply to right
   addR := sumAdd(lenR, A2, B2)
   tree[right] = (tree[right] + addR) % mod
   tagA[right] = (tagA[right] + A2) % mod
   tagB[right] = (tagB[right] + B2) % mod
   // clear
   tagA[v], tagB[v] = 0, 0
}

func build(v, l, r int) {
   if l == r {
       tree[v] = a[l] % mod
       return
   }
   mid := (l + r) >> 1
   build(v<<1, l, mid)
   build(v<<1|1, mid+1, r)
   tree[v] = (tree[v<<1] + tree[v<<1|1]) % mod
}

func update(v, l, r, ql, qr int) {
   if qr < l || r < ql {
       return
   }
   if ql <= l && r <= qr {
       // offset for this segment
       k := l - ql + 1
       A := fib[k]
       B := fib[k+1]
       added := sumAdd(r-l+1, A, B)
       tree[v] = (tree[v] + added) % mod
       tagA[v] = (tagA[v] + A) % mod
       tagB[v] = (tagB[v] + B) % mod
       return
   }
   push(v, l, r)
   mid := (l + r) >> 1
   update(v<<1, l, mid, ql, qr)
   update(v<<1|1, mid+1, r, ql, qr)
   tree[v] = (tree[v<<1] + tree[v<<1|1]) % mod
}

func query(v, l, r, ql, qr int) int64 {
   if qr < l || r < ql {
       return 0
   }
   if ql <= l && r <= qr {
       return tree[v]
   }
   push(v, l, r)
   mid := (l + r) >> 1
   res := query(v<<1, l, mid, ql, qr) + query(v<<1|1, mid+1, r, ql, qr)
   if res >= mod {
       res -= mod
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n, &m)
   a = make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // precompute fib and prefix sums
   fib = make([]int64, n+5)
   sumF = make([]int64, n+5)
   fib[0], fib[1], fib[2] = 0, 1, 1
   sumF[0] = 0
   sumF[1] = fib[1]
   sumF[2] = (fib[1] + fib[2]) % mod
   for i := 3; i <= n+2; i++ {
       fib[i] = (fib[i-1] + fib[i-2]) % mod
       sumF[i] = (sumF[i-1] + fib[i]) % mod
   }
   tree = make([]int64, 4*(n+1))
   tagA = make([]int64, 4*(n+1))
   tagB = make([]int64, 4*(n+1))
   build(1, 1, n)
   for i := 0; i < m; i++ {
       typ, l, r := 0, 0, 0
       fmt.Fscan(in, &typ, &l, &r)
       if typ == 1 {
           update(1, 1, n, l, r)
       } else if typ == 2 {
           ans := query(1, 1, n, l, r)
           fmt.Fprintln(out, ans)
       }
   }
}
