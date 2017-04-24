package main

import (
	"io"
	"bufio"
	"strings"
	"strconv"
	"path/filepath"
	"os"
	"fmt"
)

func acceleration(in io.Reader, out io.Writer) error {
	bin := bufio.NewReader(in)
	bout := bufio.NewWriter(out)
	for {
		line, err := bin.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		elements := strings.Split(line, ",")
		if len(elements) >= 3 {
			var acc []string
			for i := 0; i < 3; i++ {
				a, err := strconv.ParseFloat(elements[0], 64)
				if err != nil {
					return err
				}
				sa := strconv.FormatFloat(a, 'f', 5, 64)
				acc = append(acc, sa)
			}

			if _, err := bout.WriteString(strings.Join(acc, ",") + "\n"); err != nil {
				return err
			}
		}
	}

	return nil
}

func accelerationWalk() filepath.WalkFunc {
	c := 0
	return func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			c++
			fileName := filepath.Join(filepath.Dir(path), fmt.Sprintf("%d.csv", c))
			if inf, err := os.Open(path); err == nil {
				defer inf.Close()
				if outf, err := os.Create(fileName); err == nil {
					defer func(){
						outf.Close()
						renamedFileName := filepath.Join(filepath.Dir(path), fmt.Sprintf("%d.src", c))
						os.Rename(path, renamedFileName)
					}()
					return acceleration(inf, outf)
				} else {
					return err
				}
			} else {
				return err
			}
		}

		return nil
	}
}

func main() {
	filepath.Walk("labelled", accelerationWalk())
}
