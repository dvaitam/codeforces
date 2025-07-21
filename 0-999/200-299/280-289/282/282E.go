package main

import (
   "bufio"
   "fmt"
   "os"
)

// Trie node for binary representation
type node struct {
   child [2]int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([]uint64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // prefix xor P[0..n]
   P := make([]uint64, n+1)
   for i := 1; i <= n; i++ {
       P[i] = P[i-1] ^ a[i]
   }
   // suffix xor S[1..n+1]
   S := make([]uint64, n+2)
   for i := n; i >= 1; i-- {
       S[i] = S[i+1] ^ a[i]
   }

   // initialize trie
   // pre-allocate nodes: root at index 0
   nodes := make([]node, 1, (n+1)*41)
   // insert initial P[0]
   insertTrie(&nodes, P[0])
   var ans uint64
   // for each possible suffix starting at j
   for j := 1; j <= n+1; j++ {
       // query best with S[j]
       val := queryTrie(nodes, S[j])
       if val > ans {
           ans = val
       }
       // insert next prefix P[j]
       if j <= n {
           insertTrie(&nodes, P[j])
       }
   }
   fmt.Fprintln(writer, ans)
}

// insert x into trie
func insertTrie(nodes *[]node, x uint64) {
   cur := 0
   for bit := 39; bit >= 0; bit-- {
       b := int((x >> uint(bit)) & 1)
       if (*nodes)[cur].child[b] == 0 {
           *nodes = append(*nodes, node{})
           (*nodes)[cur].child[b] = len(*nodes) - 1
       }
       cur = (*nodes)[cur].child[b]
   }
}

// query max xor with x
func queryTrie(nodes []node, x uint64) uint64 {
   cur := 0
   var res uint64
   for bit := 39; bit >= 0; bit-- {
       b := int((x >> uint(bit)) & 1)
       opp := b ^ 1
       if nodes[cur].child[opp] != 0 {
           res |= uint64(1) << uint(bit)
           cur = nodes[cur].child[opp]
       } else {
           cur = nodes[cur].child[b]
       }
   }
   return res
}
