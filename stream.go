package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func readStream(s io.Reader) (string, error) {
	r := bufio.NewReader(s)

	var inp string
	for {
		var err error
		inp, err = r.ReadString('\n')
		if err != io.EOF {
			if err != nil {
				return "", err
			}
		}

		inp = strings.TrimSpace(inp)
		if len(inp) > 0 {
			break
		}
	}

	fmt.Println()
	return inp, nil
}
