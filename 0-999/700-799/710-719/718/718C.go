package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

// Matrix is 2x2 matrix
type Matrix [2][2]int64

// multiply two matrices
func (a Matrix) mul(b Matrix) Matrix {
   return Matrix{
       {(a[0][0]*b[0][0] + a[0][1]*b[1][0]) % MOD, (a[0][0]*b[0][1] + a[0][1]*b[1][1]) % MOD},
       {(a[1][0]*b[0][0] + a[1][1]*b[1][0]) % MOD, (a[1][0]*b[0][1] + a[1][1]*b[1][1]) % MOD},
   }
}

var prePow [32]Matrix

// fast power of base matrix M
func matPow(x int64) Matrix {
   // M^x = product of prePow bits
   var res = Matrix{{1, 0}, {0, 1}}
   for i := 0; x > 0; i++ {
       if x&1 == 1 {
           res = prePow[i].mul(res)
       }
       x >>= 1
   }
   return res
}

// Node for segment tree
type Node struct {
   sum0, sum1 int64  // sum f(a_i), sum f(a_i+1)
   tag        Matrix // pending matrix to apply
}

var n, m int
var a []int64
var tree []Node

// push tag to children
func push(u, l, r int) {
   tag := tree[u].tag
   if tag != (Matrix{{1, 0}, {0, 1}}) {
       // apply to children
       lc, rc := u*2, u*2+1
       for _, v := range []int{lc, rc} {
           // update sums
           s0 := (tag[0][0]*tree[v].sum1 + tag[0][1]*tree[v].sum0) % MOD
           s1 := (tag[1][0]*tree[v].sum1 + tag[1][1]*tree[v].sum0) % MOD
           tree[v].sum1, tree[v].sum0 = s0, s1
           // compose tags: new_tag = tag * old_tag
           tree[v].tag = tag.mul(tree[v].tag)
       }
       // reset
       tree[u].tag = Matrix{{1, 0}, {0, 1}}
   }
}

// pull up from children
func pull(u int) {
   lc, rc := u*2, u*2+1
   tree[u].sum0 = (tree[lc].sum0 + tree[rc].sum0) % MOD
   tree[u].sum1 = (tree[lc].sum1 + tree[rc].sum1) % MOD
}

// build segment tree
func build(u, l, r int) {
   tree[u].tag = Matrix{{1, 0}, {0, 1}}
   if l == r {
       // compute f(a[l]) and f(a[l]+1)
       mat := matPow(a[l])
       tree[u].sum1 = mat[0][0] // F(a+1)
       tree[u].sum0 = mat[0][1] // F(a)
       return
   }
   mid := (l + r) >> 1
   build(u*2, l, mid)
   build(u*2+1, mid+1, r)
   pull(u)
}

// update range [L,R] with add x: apply M^x
func update(u, l, r, L, R int, Mx Matrix) {
   if L > r || R < l {
       return
   }
   if L <= l && r <= R {
       // apply Mx
       s0 := (Mx[0][0]*tree[u].sum1 + Mx[0][1]*tree[u].sum0) % MOD
       s1 := (Mx[1][0]*tree[u].sum1 + Mx[1][1]*tree[u].sum0) % MOD
       tree[u].sum1, tree[u].sum0 = s0, s1
       tree[u].tag = Mx.mul(tree[u].tag)
       return
   }
   push(u, l, r)
   mid := (l + r) >> 1
   update(u*2, l, mid, L, R, Mx)
   update(u*2+1, mid+1, r, L, R, Mx)
   pull(u)
}

// query range sum f(a_i)
func query(u, l, r, L, R int) (int64, int64) {
   if L > r || R < l {
       return 0, 0
   }
   if L <= l && r <= R {
       return tree[u].sum0, tree[u].sum1
   }
   push(u, l, r)
   mid := (l + r) >> 1
   s0l, s1l := query(u*2, l, mid, L, R)
   s0r, s1r := query(u*2+1, mid+1, r, L, R)
   return (s0l + s0r) % MOD, (s1l + s1r) % MOD
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &m)
   a = make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // precompute base powers of M
   // M = [1 1;1 0]
   prePow[0] = Matrix{{1, 1}, {1, 0}}
   for i := 1; i < 32; i++ {
       prePow[i] = prePow[i-1].mul(prePow[i-1])
   }
   tree = make([]Node, 4*(n+5))
   build(1, 1, n)
   for i := 0; i < m; i++ {
       var t, l, r int
       fmt.Fscan(reader, &t, &l, &r)
       if t == 1 {
           var x int64
           fmt.Fscan(reader, &x)
           Mx := matPow(x)
           update(1, 1, n, l, r, Mx)
       } else {
           s0, _ := query(1, 1, n, l, r)
           fmt.Fprintln(writer, s0)
       }
   }
}
