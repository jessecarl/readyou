package main

import (
	"flag"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/gohugoio/hugo/parser/metadecoders"
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
	content, err := parseFrontMatterAndContent(in)
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

func parseFrontMatterAndContent(r io.Reader) (pageparser.ContentFrontMatter, error) {
	var cf pageparser.ContentFrontMatter

	psr, err := pageparser.Parse(r, pageparser.Config{})
	if err != nil {
		return cf, err
	}

	var frontMatterSource []byte

	iter := psr.Iterator()

	walkFn := func(item pageparser.Item) bool {
		switch {
		case item.Type == pageparser.TypeIgnore:
		case item.IsFrontMatter():
			cf.FrontMatterFormat = pageparser.FormatFromFrontMatterType(item.Type)
			frontMatterSource = item.Val
		case item.IsDone():
			cf.Content = psr.Input()[:]
		case frontMatterSource != nil:
			// The rest is content.
			cf.Content = psr.Input()[item.Pos:]
			// Done
			return false
		}
		return true
	}

	iter.PeekWalk(walkFn)

	cf.FrontMatter, err = metadecoders.Default.UnmarshalToMap(frontMatterSource, cf.FrontMatterFormat)
	return cf, err
}
