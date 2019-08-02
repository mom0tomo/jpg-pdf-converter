package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gographics/imagick/imagick"
	unicommon "github.com/unidoc/unidoc/common"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func init() {
	// Debug log level.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	imagick.Initialize()
	defer imagick.Terminate()
}

func main() {
	srcDir := "./images/"
	dstDir := "./pdfs/"

	if err := imagesToPdf(srcDir, dstDir); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := mergePDFs(dstDir); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file in: %s\n", dstDir)
}

func imagesToPdf(srcDir string, dstDir string) error {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	files, _ := ioutil.ReadDir(srcDir)
	fileNum := len(files)

	for _, file := range files {
		for i := 1; i < fileNum; i++ {
			// 拡張子を抜いたファイル名を取得する
			fileName := file.Name()
			baseName := filepath.Base(fileName[:len(fileName)-len(filepath.Ext(fileName))])

			// jpg/pngを読む
			mw.ReadImage(srcDir + fileName)

			// pdfに変える
			mw.SetImageFormat("pdf")
			mw.WriteImage(dstDir + baseName + ".pdf")
		}
	}

	return nil
}

func mergePDFs(dstDir string) error {
	pdfWriter := pdf.NewPdfWriter()

	pdfFiles, _ := ioutil.ReadDir(dstDir)

	for _, pdfFile := range pdfFiles {
		f, err := os.Open(dstDir + pdfFile.Name())
		if err != nil {
			return err
		}
		defer f.Close()

		pdfReader, err := pdf.NewPdfReader(f)
		if err != nil {
			return err
		}

		numPages, err := pdfReader.GetNumPages()
		for j := 0; j < numPages; j++ {
			page, err := pdfReader.GetPage(j + 1)
			if err != nil {
				return err
			}

			err = pdfWriter.AddPage(page)
			if err != nil {
				return err
			}
		}
		if err := pdfWriter.Write(f); err != nil {
			return err
		}
	}

	fWrite, err := os.Create(dstDir + "output.pdf")
	if err != nil {
		return err
	}
	defer fWrite.Close()

	if err := pdfWriter.Write(fWrite); err != nil {
		return err
	}

	return nil
}
