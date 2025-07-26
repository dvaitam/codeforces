package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
)

const MOD = 1000000007

var (
   inv2   int64
   pow2   []int64
   inv2p  []int64
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

type Node struct {
   pVal       int64
   id         int
   prio       int
   left, right *Node
   size       int
   sumA, sumB, C int64
}

func size(n *Node) int {
   if n == nil {
       return 0
   }
   return n.size
}

func upd(n *Node) {
   if n == nil {
       return
   }
   // children
   L, R := n.left, n.right
   szL := size(L)
   szR := size(R)
   n.size = szL + 1 + szR
   // compute A_node and B_node
   // A_i = p * pow2[i-1], B_i = p * inv2p[i]
   // node position = szL+1
   A_node := n.pVal * pow2[szL] % MOD
   B_node := n.pVal * inv2p[szL+1] % MOD
   // sumA
   sumA := int64(0)
   sumA = (sumA + getSumA(L)) % MOD
   sumA = (sumA + A_node) % MOD
   if R != nil {
       sumA = (sumA + getSumA(R) * pow2[szL+1]) % MOD
   }
   n.sumA = sumA
   // sumB
   sumB := int64(0)
   if L != nil {
       sumB = (sumB + getSumB(L)) % MOD
   }
   sumB = (sumB + B_node) % MOD
   if R != nil {
       sumB = (sumB + getSumB(R) * inv2p[szL+1]) % MOD
   }
   n.sumB = sumB
   // C
   c := int64(0)
   if L != nil {
       c = (c + L.C) % MOD
   }
   if R != nil {
       c = (c + R.C) % MOD
   }
   // cross terms
   // sumA_L * B_node
   if L != nil {
       c = (c + getSumA(L) * B_node) % MOD
   }
   // sumA_L * sumB_R_shift
   if L != nil && R != nil {
       tmp := getSumB(R) * inv2p[szL+1] % MOD
       c = (c + getSumA(L) * tmp) % MOD
   }
   // A_node * sumB_R_shift
   if R != nil {
       tmp := getSumB(R) * inv2p[szL+1] % MOD
       c = (c + A_node * tmp) % MOD
   }
   n.C = c
}

func getSumA(n *Node) int64 {
   if n == nil {
       return 0
   }
   return n.sumA
}

func getSumB(n *Node) int64 {
   if n == nil {
       return 0
   }
   return n.sumB
}

// compare node key < (p,id)
func lessNode(n *Node, p int64, id int) bool {
   if n.pVal != p {
       return n.pVal < p
   }
   return n.id < id
}

func split(n *Node, p int64, id int) (l, r *Node) {
   if n == nil {
       return nil, nil
   }
   if lessNode(n, p, id) {
       ll, rr := split(n.right, p, id)
       n.right = ll
       upd(n)
       return n, rr
   } else {
       ll, rr := split(n.left, p, id)
       n.left = rr
       upd(n)
       return ll, n
   }
}

func merge(a, b *Node) *Node {
   if a == nil {
       return b
   }
   if b == nil {
       return a
   }
   if a.prio < b.prio {
       a.right = merge(a.right, b)
       upd(a)
       return a
   } else {
       b.left = merge(a, b.left)
       upd(b)
       return b
   }
}

func insert(root *Node, node *Node) *Node {
   if root == nil {
       return node
   }
   if node.prio < root.prio {
       l, r := split(root, node.pVal, node.id)
       node.left = l
       node.right = r
       upd(node)
       return node
   }
   if lessNode(root, node.pVal, node.id) {
       root.right = insert(root.right, node)
   } else {
       root.left = insert(root.left, node)
   }
   upd(root)
   return root
}

func erase(root *Node, p int64, id int) *Node {
   if root == nil {
       return nil
   }
   if root.pVal == p && root.id == id {
       // remove root
       return merge(root.left, root.right)
   }
   if lessNode(root, p, id) {
       root.right = erase(root.right, p, id)
   } else {
       root.left = erase(root.left, p, id)
   }
   upd(root)
   return root
}

func main() {
   defer writer.Flush()
   var n, q int
   fmt.Fscan(reader, &n)
   pArr := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &pArr[i])
   }
   fmt.Fscan(reader, &q)
   inv2 = (MOD + 1) / 2
   maxN := n + q + 5
   pow2 = make([]int64, maxN)
   inv2p = make([]int64, maxN)
   pow2[0] = 1
   inv2p[0] = 1
   for i := 1; i < maxN; i++ {
       pow2[i] = pow2[i-1] * 2 % MOD
       inv2p[i] = inv2p[i-1] * inv2 % MOD
   }
   var root *Node
   // build initial treap
   for i := 1; i <= n; i++ {
       node := &Node{pVal: pArr[i], id: i, prio: rand.Int()}
       upd(node)
       root = insert(root, node)
   }
   // print initial
   var ans int64
   if root != nil {
       ans = root.C % MOD
   }
   fmt.Fprintln(writer, ans)
   // process queries
   for k := 0; k < q; k++ {
       var idx int
       var x int64
       fmt.Fscan(reader, &idx, &x)
       // erase old
       root = erase(root, pArr[idx], idx)
       // insert new
       pArr[idx] = x
       node := &Node{pVal: x, id: idx, prio: rand.Int()}
       upd(node)
       root = insert(root, node)
       // output
       ans = root.C % MOD
       fmt.Fprintln(writer, ans)
   }
}
