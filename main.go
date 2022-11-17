package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"pocGoby2Xray/pkg/runner"
	"strings"
)

var (
	gobyPocFile string
	gobyPocPath string
	outputPath  string
)

func main() {
	flag.StringVar(&gobyPocFile, "f", "", "transform one goby poc file")
	flag.StringVar(&gobyPocPath, "p", "", "transform poc file for path")
	flag.StringVar(&outputPath, "o", "", "the xray poc output path")
	flag.Parse()

	if gobyPocFile == "" && gobyPocPath == "" {
		flag.Usage()
		return
	}
	if len(gobyPocFile) > 0 {
		err := transformOnePoc(gobyPocFile)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if len(gobyPocPath) > 0 {
		err := transformPathPoc(gobyPocPath)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func transformPathPoc(filePath string) error {
	files, err := filepath.Glob(filepath.Join(filePath, "*.json"))
	if err != nil {
		return err
	}
	var total int
	for _, file := range files {
		err = transformOnePoc(file)
		if err != nil {
			pocFile := filepath.Base(file)
			log.Printf("%s:%v\n", pocFile, err)
		} else {
			total++
		}
	}
	log.Printf("[+]Finish: %d/%d", total, len(files))
	return nil
}

func transformOnePoc(filePathName string) error {
	r := runner.NewPocRunner()
	gpoc, err := r.LoadGobyJsonPocFile(filePathName)
	if err != nil {
		return err
	}
	xpoc, err := r.Run(gpoc)
	if err != nil {
		return err
	}
	//获取文件名称带后缀
	fileNameWithSuffix := filepath.Base(filePathName)
	//获取文件的后缀(文件类型)
	fileType := filepath.Ext(fileNameWithSuffix)
	//获取文件名称(不带后缀)
	fileNameOnly := strings.TrimSuffix(fileNameWithSuffix, fileType)
	//转换后的POC文件路径
	dir := outputPath
	if len(dir) == 0 {
		dir = filepath.Dir(filePathName)
	}
	//拼接最终的文件名
	newPathFileName := filepath.Join(dir, fmt.Sprintf("%s.yml", fileNameOnly))
	err = r.SaveXrayPocFile(newPathFileName, xpoc)
	if err != nil {
		return err
	}
	return nil
}
