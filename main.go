package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gographics/imagick/imagick"
	unicommon "github.com/unidoc/unidoc/common"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func init() {
	// Debug log level.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
}

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	srcDir := "./images/"
	dstDir := "./pdfs/"

	files, _ := ioutil.ReadDir(srcDir)
	fileNum := len(files)

	for _, file := range files {
		// 拡張子を抜いたファイル名を取得する
		fileName := file.Name()
		baseName := filepath.Base(fileName[:len(fileName) - len(filepath.Ext(fileName))])

		for i := 1; i < fileNum; i++ {
			// jpgを読む
			mw.ReadImage(srcDir + baseName + ".jpg")

			// pdfに変える
			mw.SetImageFormat("pdf")
			mw.WriteImage(dstDir + baseName + ".pdf")
		}
	}

	// pdfを結合する
	pdfWriter := pdf.NewPdfWriter()

	pdfFiles, _ := ioutil.ReadDir(dstDir)

	for _, pdfFile := range pdfFiles {
		f, err := os.Open(dstDir + pdfFile.Name())
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		pdfReader, err := pdf.NewPdfReader(f)
		if err != nil {
			log.Fatal(err)
		}

		numPages, err := pdfReader.GetNumPages()
		for j := 0; j < numPages; j++ {
			page, err := pdfReader.GetPage(j + 1)
			if err != nil {
				log.Fatal(err)
			}

			err = pdfWriter.AddPage(page)
			if err != nil {
				log.Fatal(err)
			}
		}
		if err := pdfWriter.Write(f); err != nil {
			log.Fatal(err)
		}
	}

	fWrite, err := os.Create(dstDir + "output.pdf")
	if err != nil {
		log.Fatal(err)
	}

	defer fWrite.Close()

	if err := pdfWriter.Write(fWrite); err != nil {
		log.Fatal(err)
	}
}