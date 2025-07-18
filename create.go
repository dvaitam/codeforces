package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(os.Args[2])
	s := os.Args[2]
	// letter := s[len(s)-1:]
	n, _ := strconv.Atoi(s)
	sn := n
	thousands := (n / 1000) * 1000
	t_dir := fmt.Sprintf("%d-%d", thousands, thousands+999)
	n = n % 1000
	hundreds := (n / 100) * 100
	h_dir := fmt.Sprintf("%d-%d", thousands+hundreds, thousands+hundreds+99)
	n = n % 100
	tens := (n / 10) * 10
	te_dir := fmt.Sprintf("%d-%d", thousands+hundreds+tens, thousands+hundreds+tens+9)

	c_dir := fmt.Sprintf("%d", sn)
	path := strings.Join([]string{t_dir, h_dir, te_dir, c_dir}, "/")
	newpath := filepath.Join(".", path)
	err := os.MkdirAll(newpath, os.ModePerm)
	fmt.Println(err)
	final_path := "./" + path + "/" + s + os.Args[3] + ".go"
	fmt.Println(final_path)
	//os.OpenFile(final_path, os.O_RDONLY|os.O_CREATE, 0666)

	Copy("template.go", final_path, os.Args[3])

}
func Copy(src, dst, cmd string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
