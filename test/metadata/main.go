package main

import (
	"flag"
	"log"

	metadata "github.com/niakr1s/moee/src/lib/metatdata"
	"github.com/niakr1s/moee/src/lib/util"
)

var clean *bool = flag.Bool("clean", false, "if set, cleans test files; if not - applies test metadata to it")
var dir *string = flag.String("dir", "test_files", "test files dir path")

func init() {
	flag.Parse()
}

func main() {
	dir := *dir
	dir, _ = util.ToAbs(dir)
	clean := *clean

	log.Printf("got dir=%s, clean=%v", dir, clean)

	musicFiles, err := util.GetAllMusicFilesFromDir(dir, false)
	if err != nil {
		log.Fatal(err)
	}

	for _, filePath := range musicFiles {
		if clean {
			err = metadata.WriteClean(filePath)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err = metadata.WriteTitle(filePath, "title")
			if err != nil {
				log.Fatal(err)
			}
			err = metadata.WriteArtist(filePath, "artist")
			if err != nil {
				log.Fatal(err)
			}
			err = metadata.WriteLyrics(filePath, "lyrics")
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}
