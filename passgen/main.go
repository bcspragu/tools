package main

import (
	"bufio"
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("must specify a dictionary path as an arg")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("failed to open word list: %v", err)
	}
	defer f.Close()

	dict := make([]string, 0, 370000 /* roughly the size of the word list */)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		dict = append(dict, sc.Text())
	}

	r := rand.New(&cryptoSource{})
	pwd := wordLoop(r, dict)

	fmt.Println()
	for i, w := range pwd {
		fmt.Print(w)
		if i < len(pwd)-1 {
			fmt.Print(r.Intn(100))
		}
	}
}

func wordLoop(r *rand.Rand, dict []string) []string {
	var pwd []string
	for {
		word := strings.Title(dict[r.Intn(len(dict))])

		fmt.Printf("Proposed word: %q\n", word)

		fmt.Print("Command (a - accept, n - next, q - quit: ")
		var input string
		fmt.Scanln(&input)
		switch input {
		case "a":
			pwd = append(pwd, word)
		case "n":
			continue
		case "q":
			return pwd
		}
	}
}

type cryptoSource struct{}

func (c *cryptoSource) Int63() int64 {
	return int64(c.Uint64() & 0x7FFFFFFFFFFFFFFF)
}

func (*cryptoSource) Uint64() uint64 {
	var buf [8]byte
	n, err := cryptorand.Read(buf[:])
	if err != nil {
		panic(err)
	}
	if n != 8 {
		panic("read incorrect number of bytes")
	}
	return binary.BigEndian.Uint64(buf[:])
}

func (*cryptoSource) Seed(_ int64) {}
