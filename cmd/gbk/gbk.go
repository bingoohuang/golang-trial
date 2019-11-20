package main

// Reading a non UTF-8 text file in Go
// https://stackoverflow.com/a/31544542

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// The golang.org/x/text/encoding package defines an interface for generic character encodings
// that can convert to/from UTF-8.
// The golang.org/x/text/encoding/simplifiedchinese sub-package
// provides GB18030, GBK and HZ-GB2312 encoding implementations.

// Encoding to use. Since this implements the encoding.Encoding
// interface from golang.org/x/text/encoding you can trivially
// change this out for any of the other implemented encoders,
// e.g. `traditionalchinese.Big5`, `charmap.Windows1252`,
// `korean.EUCKR`, etc.
var gbk = simplifiedchinese.GBK

func main() {
	const filename = "example_GBK_file"
	exampleWriteGBK(filename)
	exampleReadGBK(filename)
}

func exampleReadGBK(filename string) {
	// Read UTF-8 from a GBK encoded file.
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	r := transform.NewReader(f, gbk.NewDecoder())

	// Read converted UTF-8 from `r` as needed.
	// As an example we'll read line-by-line showing what was read:
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		fmt.Printf("Read line: %s\n", sc.Bytes())
	}
	if err = sc.Err(); err != nil {
		log.Fatal(err)
	}

	if err = f.Close(); err != nil {
		log.Fatal(err)
	}
}

func exampleWriteGBK(filename string) {
	// Write UTF-8 to a GBK encoded file.
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	w := transform.NewWriter(f, gbk.NewEncoder())

	// Write UTF-8 to `w` as desired.
	// As an example we'll write some text from the Wikipedia
	// GBK page that includes Chinese.
	_, err = fmt.Fprintln(w,
		`In 1995, China National Information Technology Standardization
Technical Committee set down the Chinese Internal Code Specification
(Chinese: 汉字内码扩展规范（GBK）; pinyin: Hànzì Nèimǎ
Kuòzhǎn Guīfàn (GBK)), Version 1.0, known as GBK 1.0, which is a
slight extension of Codepage 936. The newly added 95 characters were not
found in GB 13000.1-1993, and were provisionally assigned Unicode PUA
code points.`)
	if err != nil {
		log.Fatal(err)
	}

	if err = f.Close(); err != nil {
		log.Fatal(err)
	}
}
