package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/http-live-streaming/m3u8-downloader/dl"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

var (
	url      string
	output   string
	chanSize int
)

func init() {
	flag.StringVar(&url, "u", "", "M3U8 URL, required")
	flag.IntVar(&chanSize, "c", 25, "Maximum number of occurrences")
	flag.StringVar(&output, "o", "", "Output folder, required")
}

func main() {
	flag.Parse()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[error]", r)
			os.Exit(-1)
		}
	}()
	if url == "" {
		panicParameter("u")
	}
	if output == "" {
		panicParameter("o")
	}
	if chanSize <= 0 {
		panic("parameter 'c' must be greater than 0")
	}
	downloader, err := dl.NewTask(output, url)
	if err != nil {
		panic(err)
	}
	if err := downloader.Start(chanSize); err != nil {
		panic(err)
	}

	err = ffmpeg_go.Input(output+"/main.ts").
		Output(output+"/main.mp4", ffmpeg_go.KwArgs{"c": "copy", "f": "mp4"}).
		OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		panic("cover to mp4 failed")
	}

	fmt.Println("removing ts file....")
	err = os.Remove(output + "/main.ts")
	if err != nil {
		fmt.Println("remove ts file failed", err)
	}

	fmt.Println("Done!")
}

func panicParameter(name string) {
	panic("parameter '" + name + "' is required")
}
