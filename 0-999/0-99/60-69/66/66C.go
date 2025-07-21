package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

// Node represents a folder or disk node in the file system tree
type Node struct {
   children map[string]*Node
   fileCount int      // number of files directly in this folder
   // computed counts
   subfolders int     // total number of subfolders (descendants)
   files int          // total number of files in subtree
   isDisk bool        // true for disk root nodes (C:, D:, etc.)
}

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   // build tree with root having disk children
   root := &Node{children: make(map[string]*Node)}
   for scanner.Scan() {
       line := scanner.Text()
       if line == "" {
           continue
       }
       parts := strings.Split(line, "\\")
       // parts[0] is disk like "C:"
       disk := parts[0]
       node, ok := root.children[disk]
       if !ok {
           node = &Node{children: make(map[string]*Node), isDisk: true}
           root.children[disk] = node
       }
       // traverse folders
       curr := node
       // folders are parts[1] .. parts[len-2]
       for i := 1; i < len(parts)-1; i++ {
           name := parts[i]
           child, exists := curr.children[name]
           if !exists {
               child = &Node{children: make(map[string]*Node)}
               curr.children[name] = child
           }
           curr = child
       }
       // last part is file name
       curr.fileCount++
   }
   // compute counts
   var maxSub, maxFiles int
   var dfs func(n *Node)
   dfs = func(n *Node) {
       sf := 0
       ff := n.fileCount
       for _, c := range n.children {
           dfs(c)
           // add child's subfolders and itself if folder
           sf += c.subfolders
           if !c.isDisk {
               sf++
           }
           // add files in child's subtree
           ff += c.files
       }
       n.subfolders = sf
       n.files = ff
       // skip disk nodes for max
       if !n.isDisk {
           if sf > maxSub {
               maxSub = sf
           }
           if ff > maxFiles {
               maxFiles = ff
           }
       }
   }
   // run dfs from each disk child
   for _, d := range root.children {
       dfs(d)
   }
   // output results
   fmt.Printf("%d %d\n", maxSub, maxFiles)
}
