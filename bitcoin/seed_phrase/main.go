package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
	"unicode/utf8"
)

type bitString string

func main() {

	x := bitString(get128BitString())
	hash := NewSHA256(x.AsByteSlice())
	x += bitString(getFirst4Bits(hash))
	indexes := x.getIndexes()
	dict := readFromFile(`./bip-0039/english.txt`)
	seed := getSeedPhrase(indexes, dict)
	fmt.Println(seed)
}

func getSeedPhrase(indexes [12]int, dict []string) (out [12]string) {
	for i := 0; i < 12; i++ {
		out[i] = dict[indexes[i]]
	}
	return
}
func readFromFile(filePath string) (out []string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	// read line by line
	for fileScanner.Scan() {
		out = append(out, fileScanner.Text())
	}
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}
	return
}
func get128BitString() (bitString string) {
	for i := 0; i < 128; i++ {
		if eagle := flipCoin(); eagle {
			bitString += "1"
		} else {
			bitString += "0"
		}
		time.Sleep(time.Nanosecond)
	}
	return
}

func (b bitString) getIndexes() (out [12]int) {
	for i := range out {
		bin := string(b[:11])
		n, _ := strconv.ParseInt(bin, 2, 64)
		out[i] = int(n)
		b = b[11:]
	}
	return
}
func getFirst4Bits(sha256hash []byte) (bitString string) {
	for _, v := range sha256hash {
		s := fmt.Sprintf("%b", v)
		bitString += s
		if utf8.RuneCountInString(bitString) >= 4 {
			break
		}
	}
	return bitString[:4]
}

func flipCoin() bool {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int()
	return n%2 == 1
}

func NewSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func (b bitString) AsByteSlice() []byte {
	var out []byte
	var str string

	for i := len(b); i > 0; i -= 8 {
		if i-8 < 0 {
			str = string(b[0:i])
		} else {
			str = string(b[i-8 : i])
		}
		v, err := strconv.ParseUint(str, 2, 8)
		if err != nil {
			panic(err)
		}
		out = append([]byte{byte(v)}, out...)
	}
	return out
}
