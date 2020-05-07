package main

import (
	"flag"
	"log"
	"os"
	"text/template"

	"github.com/gohugoio/hugo/parser/pageparser"
)

func main() {
	infile := flag.String("input-file", "", "optional path to input file (otherwise uses stdin)")
	flag.StringVar(infile, "i", "", "optional path to input file (otherwise uses stdin)")
	outfile := flag.String("output-file", "", "optional path to output file (otherwise uses stdout)")
	flag.StringVar(outfile, "o", "", "optional path to output file (otherwise uses stdout)")

	flag.Parse()

	in := os.Stdin
	if *infile != "" {
		f, err := os.Open(*infile)
		if err != nil {
			log.Panic(err)
		}
		defer f.Close()
		in = f
	}

	out := os.Stdout
	if *outfile != "" {
		f, err := os.Create(*outfile)
		if err != nil {
			log.Panic(err)
		}
		defer f.Close()
		in = f
	}

	// parse input for data and raw template
	content, err := pageparser.ParseFrontMatterAndContent(in)
	if err != nil {
		log.Panic(err)
	}

	// create template and parse raw
	tmpl, err := template.New("input").Parse(string(content.Content))
	if err != nil {
		log.Panic(err)
	}

	err = tmpl.Execute(out, content.FrontMatter)
	if err != nil {
		log.Panic(err)
	}
}
