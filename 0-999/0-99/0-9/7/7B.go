package main

import (
   "fmt"
)

type block struct {
   id    int
   start int // start index (0-based)
   size  int
}

func main() {
   var t, m int
   if _, err := fmt.Scan(&t, &m); err != nil {
       return
   }
   blocks := make([]block, 0, t)
   nextID := 1
   for i := 0; i < t; i++ {
       var op string
       fmt.Scan(&op)
       switch op {
       case "alloc":
           var sz int
           fmt.Scan(&sz)
           // find first fit
           pos := -1
           // check before first block
           if len(blocks) == 0 {
               if sz <= m {
                   pos = 0
               }
           } else {
               // before first
               if blocks[0].start >= sz {
                   pos = 0
               } else {
                   // between blocks
                   for j := 0; j < len(blocks)-1; j++ {
                       freeStart := blocks[j].start + blocks[j].size
                       freeEnd := blocks[j+1].start
                       if freeEnd-freeStart >= sz {
                           pos = j + 1
                           break
                       }
                   }
                   // after last
                   if pos == -1 {
                       lastEnd := blocks[len(blocks)-1].start + blocks[len(blocks)-1].size
                       if m-lastEnd >= sz {
                           pos = len(blocks)
                       }
                   }
               }
           }
           if pos == -1 {
               fmt.Println("NULL")
           } else {
               // determine start position
               var start int
               if pos == 0 {
                   if len(blocks) > 0 {
                       start = 0
                   } else {
                       start = 0
                   }
               } else {
                   prev := blocks[pos-1]
                   start = prev.start + prev.size
               }
               // insert block
               newBlk := block{id: nextID, start: start, size: sz}
               blocks = append(blocks, block{})
               copy(blocks[pos+1:], blocks[pos:])
               blocks[pos] = newBlk
               fmt.Println(nextID)
               nextID++
           }
       case "erase":
           var x int
           fmt.Scan(&x)
           found := false
           for j, b := range blocks {
               if b.id == x {
                   // remove block
                   blocks = append(blocks[:j], blocks[j+1:]...)
                   found = true
                   break
               }
           }
           if !found {
               fmt.Println("ILLEGAL_ERASE_ARGUMENT")
           }
       case "defragment":
           // shift blocks to eliminate gaps
           cur := 0
           for j := range blocks {
               blocks[j].start = cur
               cur += blocks[j].size
           }
       }
   }
}
