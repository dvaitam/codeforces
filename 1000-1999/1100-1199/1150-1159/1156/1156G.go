package main

import (
   "bufio"
   "fmt"
   "os"
)

type Node struct {
   op    byte
   left  *Node
   right *Node
   // for leaf
   varName string
   // list of variable names that held this node
   names []string
   id int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   fmt.Fscan(reader, &n)
   // maps
   leafMap := make(map[string]*Node)
   varMap := make(map[string]*Node)
   opMap := make(map[string]*Node)
   var nodes []*Node
   idCnt := 1
   // helper to get current node of var
   getVar := func(v string) *Node {
       if nd, ok := varMap[v]; ok {
           return nd
       }
       // leaf
       if nd, ok := leafMap[v]; ok {
           varMap[v] = nd
           return nd
       }
       nd := &Node{op: 0, varName: v, names: []string{v}, id: idCnt}
       idCnt++
       leafMap[v] = nd
       varMap[v] = nd
       nodes = append(nodes, nd)
       return nd
   }
   // process lines
   for i := 0; i < n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       // split at '='
       var lhs, rhs string
       for j := 0; j < len(line); j++ {
           if line[j] == '=' {
               lhs = line[:j]
               rhs = line[j+1:]
               break
           }
       }
       // find operation in rhs
       var nd *Node
       var nodeR *Node
       for j := 0; j < len(rhs); j++ {
           c := rhs[j]
           if c == '$' || c == '^' || c == '#' || c == '&' {
               // binary op
               a1 := rhs[:j]
               a2 := rhs[j+1:]
               n1 := getVar(a1)
               n2 := getVar(a2)
               key := string(c) + "," + fmt.Sprint(n1.id) + "," + fmt.Sprint(n2.id)
               if ex, ok := opMap[key]; ok {
                   nd = ex
               } else {
                   ex := &Node{op: c, left: n1, right: n2, names: []string{}, id: idCnt}
                   idCnt++
                   opMap[key] = ex
                   nodes = append(nodes, ex)
                   nd = ex
               }
               break
           }
       }
       if nd == nil {
           // copy assignment
           nodeR = getVar(rhs)
           nd = nodeR
       }
       // record name
       nd.names = append(nd.names, lhs)
       varMap[lhs] = nd
   }
   // determine result node
   resNode, ok := varMap["res"]
   if !ok {
       // res unchanged
       fmt.Fprintln(writer, 0)
       return
   }
   // if leaf
   if resNode.op == 0 {
       // if res maps to res, no need
       if resNode.varName == "res" {
           fmt.Fprintln(writer, 0)
       } else {
           fmt.Fprintln(writer, 1)
           fmt.Fprintf(writer, "res=%s\n", resNode.varName)
       }
       return
   }
   // collect reachable nodes
   visited := make(map[int]bool)
   var order []*Node
   var dfs func(nd *Node)
   dfs = func(nd *Node) {
       if visited[nd.id] {
           return
       }
       visited[nd.id] = true
       if nd.op != 0 {
           dfs(nd.left)
           dfs(nd.right)
       }
       order = append(order, nd)
   }
   dfs(resNode)
   // gather internal nodes in order
   var internals []*Node
   for _, nd := range order {
       if nd.op != 0 {
           internals = append(internals, nd)
       }
   }
   k := len(internals)
   fmt.Fprintln(writer, k)
   // assign varNames
   nameMap := make(map[int]string)
   for _, nd := range internals {
       if nd == resNode {
           nameMap[nd.id] = "res"
           continue
       }
       // pick last non-res if possible
       sel := ""
       for i := len(nd.names) - 1; i >= 0; i-- {
           if nd.names[i] != "res" {
               sel = nd.names[i]
               break
           }
       }
       if sel == "" {
           sel = nd.names[0]
       }
       nameMap[nd.id] = sel
   }
   // output lines
   for _, nd := range internals {
       tgt := nameMap[nd.id]
       // children
       var lvar, rvar string
       if nd.left.op == 0 {
           lvar = nd.left.varName
       } else {
           lvar = nameMap[nd.left.id]
       }
       if nd.right.op == 0 {
           rvar = nd.right.varName
       } else {
           rvar = nameMap[nd.right.id]
       }
       fmt.Fprintf(writer, "%s=%s%c%s\n", tgt, lvar, nd.op, rvar)
   }
}
