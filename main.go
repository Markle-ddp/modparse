package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	_ "github.com/json-iterator/go"
)

func main() {
	ret, err := ParseModGraph()
	if err != nil {
		panic(err)
	}

	res := bytes.Split(ret, []byte("\n"))

	var out bytes.Buffer
	out.WriteString("graph TB")
	out.WriteString("\n")
	for _, item := range res {
		dep := strings.Split(string(item), " ")
		if len(dep) != 2 {
			break
		}
		out.WriteString("  ")
		out.WriteString(md5Sum(dep[0]))
		out.WriteByte('[')
		out.WriteString(strings.ReplaceAll(dep[0], "@", " "))
		out.WriteString("] --> ")
		out.WriteString(md5Sum(dep[1]))
		out.WriteByte('[')
		out.WriteString(strings.ReplaceAll(dep[1], "@", " "))
		out.WriteString("]\n")
	}
	err = os.WriteFile("dep.mmd", out.Bytes(), fs.ModePerm)
	if err != nil {
		panic(err)
	}
	fmt.Println("generate dep.mmd success")
}

func ParseModGraph() ([]byte, error) {
	cmd := exec.Command("go", "mod", "graph")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func md5Sum(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
