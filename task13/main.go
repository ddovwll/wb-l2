package main

import (
	"flag"
	"log"
	"os"

	"task13/cut"
)

func main() {
	fieldSpec := flag.String("f", "", "List of fields to extract (e.g., 1,3-5)")
	delimiter := flag.String("d", "\t", "Field delimiter (default: tab character)")
	separated := flag.Bool("s", false, "Suppress lines without delimiter (only process lines containing the delimiter)")

	flag.Parse()

	if *fieldSpec == "" {
		log.Fatal("you must specify -f option")
	}

	fields, err := cut.ParseFields(*fieldSpec)
	if err != nil {
		log.Fatal(err)
	}

	opts := cut.Options{
		Fields:    fields,
		Delimiter: *delimiter,
		Separated: *separated,
	}

	files := flag.Args()
	if len(files) == 0 {
		if err := cut.Run(os.Stdin, os.Stdout, opts); err != nil {
			log.Fatal(err)
		}
	} else {
		for _, fileName := range files {
			file, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(file)

			if err := cut.Run(file, os.Stdout, opts); err != nil {
				log.Fatal(err)
			}
		}
	}
}
