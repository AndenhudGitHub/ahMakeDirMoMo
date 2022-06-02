package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	//尺寸表路徑
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Print(err)
		os.Exit(3)
	}

	//讀取config json 參數
	type config struct {
		WorkPath      string                 `json:"WorkPath"`
		NewWorkPath   string                 `json:"NewWorkPath"`
		SizeTablePath string                 `json:"SizeTablePath"`
		TryTablePath  string                 `json:"TryTablePath"`
		Leve3Dir      string                 `json:"Leve3Dir"`
		Logo          string                 `json:"Logo"`
		Story         string                 `json:"Story"`
		BeginCount    string                 `json:"BeginCount"`
		MaxCount      string                 `json:"MaxCount"`
		Copy1         string                 `json:"copy1"`
		Copy2         string                 `json:"Copy2"`
		Copy3         string                 `json:"Copy3"`
		Copy1Max      int                    `json:"Copy1Max"`
		GroupDir      string                 `json:"GroupDir"`
		XXLDir        string                 `json:"XXLDir"`
		XXLFile       string                 `json:"XXLFile"`
		TryMapping    map[string]interface{} `json:"TryMapping"`
	}

	//將config 寫入obj
	var obj config
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(3)
	}

	// fmt.Print(obj)
	// os.Exit(3)

	DirPath := ""
	mkOutOrg := ""
	XXLFilePath := ""

	//讀取目前要做的資料夾路徑
	DirPath = strings.Replace(obj.WorkPath, "\\", "\\\\", -1)
	//複製到的資料夾
	mkOutOrg = obj.WorkPath + string(os.PathSeparator) + obj.Leve3Dir
	mkDir(mkOutOrg)
	//XXL txt 路徑
	XXLFilePath = XXLFilePath + obj.WorkPath + string(os.PathSeparator) + obj.XXLFile

	//XXL txt內文字
	XXLArray := readXXlFile(XXLFilePath)
	//尺寸表的路徑
	SpecPath := strings.Replace(obj.SizeTablePath, "\\", "\\\\", -1)
	//試穿表路徑 (用料號前兩個字)
	TryTablePath := strings.Replace(obj.TryTablePath, "\\", "\\\\", -1)
	//掃描DIR
	dirArr := scandir(DirPath)

	//LOGO路徑
	//LogoPath := obj.SizeTablePath + string(os.PathSeparator) + obj.Logo
	//圖片開始編碼數字
	BeginCount, _ := strconv.Atoi(obj.BeginCount)
	//最多編到幾張 為一個group 資料夾
	MaxCount, _ := strconv.Atoi(obj.MaxCount)

	//專推 第二個 _ 後面永遠是1 例如: 10001_M_1_6
	imageIndex := 1
	//分70個為一組資料夾 group1 group2
	group := 1
	//放抓不到尺寸表的 error msg
	var failSizeTable []string
	//放抓不到試穿表的 error msg
	var failTryTable []string
	//需要壓浮水印的圖片紀錄 無使用
	// var needFillArr []string
	var needToSize1000T1040 []string
	//試穿表對應圖片
	var TryMap = obj.TryMapping

	for _, fileDir := range dirArr {

		//圖片DIR層
		level1Dir := obj.WorkPath + string(os.PathSeparator) + fileDir

		//掠過TXT
		if strings.Index(fileDir, ".txt") != -1 {
			continue
		}

		//建立資料夾
		mkOut := mkOutOrg + string(os.PathSeparator) + obj.GroupDir + strconv.Itoa(group)
		mkDir(mkOut)

		//掃描DIR 圖片
		inSideImageArr := scandir(level1Dir)

		// fmt.Println(inSideImageArr)
		// os.Exit(3)

		copyTo1 := ""
		copyTo2 := ""
		copyTo3 := ""
		copyToStory := ""
		imageIndex2 := 1
		level2ImgDir := ""
		indexNext := 0
		hasA := false
		//檢查圖片有 ABC 標記起來 下面跑 index 從 2 開始
		for _, fileImage := range inSideImageArr {
			if fileImage == "A.jpg" || fileImage == "B.jpg" || fileImage == "C.jpg" || fileImage == "a.jpg" || fileImage == "b.jpg" || fileImage == "c.jpg" {
				hasA = true
			}
		}

		//開始跑 料號資料夾
		for _, fileImage := range inSideImageArr {
			//檔名是 JPG 才做事
			if strings.Index(fileImage, ".jpg") > -1 {

				//組合圖片需要丟入的DIR
				level2ImgDir = level1Dir + string(os.PathSeparator) + fileImage

				//有ABC INDEX 從2開始
				if hasA {
					indexNext = imageIndex2 + 1
				} else {
					indexNext = imageIndex2
				}

				//圖片本身是 ABC index 直接從 1
				if inSideImageArr[0] == "A.jpg" || inSideImageArr[0] == "B.jpg" ||
					inSideImageArr[0] == "C.jpg" || inSideImageArr[0] == "a.jpg" ||
					inSideImageArr[0] == "b.jpg" || inSideImageArr[0] == "c.jpg" {
					indexNext = imageIndex2
				}

				/* copyTo1 為SKU 六張 copyTo2 專推圖抓完整全部 copy3 O檔案全圖 */
				if fileImage == "A.jpg" || fileImage == "B.jpg" || fileImage == "C.jpg" || fileImage == "a.jpg" || fileImage == "b.jpg" || fileImage == "c.jpg" {
					copyTo1 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy1 + "1.jpg"
					copyTo2 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) + "_1.jpg"
					copyTo3 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy3 + ".jpg"
					// needFillArr = append(needFillArr, copyTo1)
					//needFillArr = append(needFillArr, copyTo2)
				} else {

					copyTo1 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy1 + strconv.Itoa(indexNext) + ".jpg"
					copyTo2 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 +
						"_" + strconv.Itoa(imageIndex) +
						"_" + strconv.Itoa(indexNext) + ".jpg"
					copyTo3 = ""
				}
				needToSize1000T1040 = append(needToSize1000T1040, copyTo2)
				if (indexNext) <= obj.Copy1Max || fileImage == "A.jpg" || fileImage == "B.jpg" ||
					fileImage == "C.jpg" || fileImage == "a.jpg" || fileImage == "b.jpg" || fileImage == "c.jpg" {
					CopyFile(level2ImgDir, copyTo1)
					fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", level2ImgDir, copyTo1)
				}
				CopyFile(level2ImgDir, copyTo2)
				fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", level2ImgDir, copyTo2)
				if copyTo3 != "" {
					CopyFile(level2ImgDir, copyTo3)
					fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", level2ImgDir, copyTo3)
				}
				imageIndex2++
			}
		}

		//上面基礎圖片做完再來做
		if strings.Index(fileDir, "OUT") == -1 {

			tryJpg := fmt.Sprintf("%v", TryMap[fileDir[0:2]])
			//切割資料夾變成陣列
			jpgCutArray := strings.Split(tryJpg, ",")

			//試穿表路徑
			tryPicPath := TryTablePath + string(os.PathSeparator)

			// fmt.Println(jpgCutArray)
			// os.Exit(3)

			//複製尺寸表
			styleNoPath := SpecPath + string(os.PathSeparator) + fileDir[0:2] + fileDir[4:8] + ".jpg"

			//複製品排故事
			storyPath := SpecPath + string(os.PathSeparator) + obj.Story
			//B 開頭 SKU 補尺寸表
			copyTo1 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy1 + strconv.Itoa(imageIndex2) + ".jpg"
			//M 開頭 推文 繼續往後補尺寸表
			copyTo2 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) +
				"_" + strconv.Itoa(imageIndex2) + ".jpg"

			//沒滿6 SKU 還要再補尺寸表
			if (imageIndex2) <= obj.Copy1Max {
				err3 := CopyFile(styleNoPath, copyTo1)
				if err3 != nil {
					fmt.Printf("複製圖片 %s ， 到 %s 失敗!!\n", styleNoPath, copyTo1)
				} else {
					fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", styleNoPath, copyTo1)
				}
			}

			//複製尺寸表
			err4 := CopyFile(styleNoPath, copyTo2)

			//試穿表 有找到MAP
			if len(jpgCutArray) >= 1 {
				for _, tryOneJpg := range jpgCutArray {

					if tryOneJpg == "" {
						continue
					}

					imageIndex2 = imageIndex2 + 1

					if imageIndex2 <= obj.Copy1Max {
						copyTo4 := mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy1 + strconv.Itoa(imageIndex2) + ".jpg"
						newTryPath := tryPicPath + tryOneJpg
						err6 := CopyFile(newTryPath, copyTo4)
						if err6 != nil {
							failTryTable = append(failTryTable, "試穿表:"+newTryPath+"失敗\n")
						}
						// needToSize1000T1040 = append(needToSize1000T1040, copyTo4)
					}

					copyTo5 := mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) +
						"_" + strconv.Itoa(imageIndex2) + ".jpg"
					newTryPath := tryPicPath + tryOneJpg
					err6 := CopyFile(newTryPath, copyTo5)
					if err6 != nil {
						failTryTable = append(failTryTable, "試穿表:"+newTryPath+"失敗\n")
					}
					needToSize1000T1040 = append(needToSize1000T1040, copyTo5)

				}
			}

			//品排故事 排再最後
			copyToStory = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 +
				"_" + strconv.Itoa(imageIndex) + "_" + strconv.Itoa(imageIndex2+1) + ".jpg"
			//複製品排故事
			err5 := CopyFile(storyPath, copyToStory)

			//紀錄需要縮放圖的 品牌故事和 推文圖
			needToSize1000T1040 = append(needToSize1000T1040, copyToStory)
			needToSize1000T1040 = append(needToSize1000T1040, copyTo2)

			//抓不到 記錄起來到 txt
			if err4 != nil || err5 != nil {
				fmt.Printf("複製圖片 %s ， 到 %s 失敗!!\n", styleNoPath, copyTo2)
				fmt.Printf("複製圖片 %s ， 到 %s 失敗!!\n", storyPath, copyToStory)
				if err5 != nil {
					failSizeTable = append(failSizeTable, "品牌故事:"+storyPath+"失敗\n")
				} else {
					failSizeTable = append(failSizeTable, "複製尺寸表:"+styleNoPath+"失敗\n")
				}

			} else {
				fmt.Printf("複製尺寸表 %s\n", styleNoPath+".jpg")
				fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", styleNoPath, copyTo2)
				fmt.Printf("品牌故事 %s\n", storyPath+".jpg")
				fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", storyPath, copyToStory)
			}

			//起始++
			BeginCount++

			//超過一個資料夾的上限 資料夾群組重新 group +1
			if BeginCount > MaxCount {
				newbegin, _ := strconv.Atoi(obj.BeginCount)
				BeginCount = newbegin
				//imageIndex = 1
				group++
			}
		}
	}

	if len(XXLArray) > 0 {

		BeginCount, _ = strconv.Atoi(obj.BeginCount)
		MaxCount, _ = strconv.Atoi(obj.MaxCount)
		imageIndex = 1
		xxl := 1
		for _, pathName := range XXLArray {
			level1Dir := obj.WorkPath + string(os.PathSeparator) + pathName
			if _, err := os.Stat(level1Dir); os.IsNotExist(err) {
				continue
			}
			mkOut := mkOutOrg + string(os.PathSeparator) + obj.XXLDir + strconv.Itoa(xxl)
			mkDir(mkOut)
			inSideImageArr := scandir(level1Dir)
			//SKU圖 最多六張
			copyTo1 := ""
			//推文用途 無上限抓滿
			copyTo2 := ""
			//O大圖原圖
			copyTo3 := ""
			imageIndex2 := 1
			level2ImgDir := ""
			copyToStory := ""
			indexNext := 0
			hasA := false
			for _, fileImage := range inSideImageArr {
				if fileImage == "A.jpg" || fileImage == "B.jpg" || fileImage == "C.jpg" || fileImage == "a.jpg" || fileImage == "b.jpg" || fileImage == "c.jpg" {
					hasA = true
				}
			}

			for _, fileImage := range inSideImageArr {
				if strings.Index(fileImage, ".jpg") > -1 {
					level2ImgDir = level1Dir + string(os.PathSeparator) + fileImage

					if hasA {
						indexNext = imageIndex2 + 1
					} else {
						indexNext = imageIndex2
					}
					if inSideImageArr[0] == "A.jpg" || inSideImageArr[0] == "B.jpg" || inSideImageArr[0] == "C.jpg" || inSideImageArr[0] == "a.jpg" || inSideImageArr[0] == "b.jpg" || inSideImageArr[0] == "c.jpg" {
						indexNext = imageIndex2
					}

					if fileImage == "A.jpg" || fileImage == "B.jpg" || fileImage == "C.jpg" || fileImage == "a.jpg" || fileImage == "b.jpg" || fileImage == "c.jpg" {
						copyTo1 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy1 + "1.jpg"
						copyTo2 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) + "_1.jpg"
						copyTo3 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy3 + ".jpg"
						// needFillArr = append(needFillArr, copyTo1)
						//needFillArr = append(needFillArr, copyTo2)
					} else {
						copyTo1 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy1 + strconv.Itoa(indexNext) + ".jpg"
						copyTo2 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) + "_" + strconv.Itoa(indexNext) + ".jpg"
						copyTo3 = ""
					}
					//需要改size 為 1000X1040的推圖
					needToSize1000T1040 = append(needToSize1000T1040, copyTo2)

					//複製SKU 需要圖
					if (indexNext) <= obj.Copy1Max || fileImage == "A.jpg" || fileImage == "B.jpg" || fileImage == "C.jpg" || fileImage == "a.jpg" || fileImage == "b.jpg" || fileImage == "c.jpg" {
						CopyFile(level2ImgDir, copyTo1)
						fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", level2ImgDir, copyTo1)
					}

					//複製推文圖
					CopyFile(level2ImgDir, copyTo2)
					fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", level2ImgDir, copyTo2)

					//複製原圖O開頭
					if copyTo3 != "" {
						CopyFile(level2ImgDir, copyTo3)
						fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", level2ImgDir, copyTo3)
					}
					imageIndex2++
				}
			}

			//檢查OUT 資料夾是否存在
			if strings.Index(pathName, "OUT") == -1 {

				tryJpg := fmt.Sprintf("%v", TryMap[pathName[0:2]])
				//切割試穿表jpg 有逗號表示要跑兩張
				jpgCutArray := strings.Split(tryJpg, ",")

				//試穿表路徑
				tryPicPath := TryTablePath + string(os.PathSeparator)

				styleNoPath := SpecPath + string(os.PathSeparator) + pathName[0:2] + pathName[4:8] + ".jpg"
				storyPath := SpecPath + string(os.PathSeparator) + obj.Story

				copyTo2 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 +
					"_" + strconv.Itoa(imageIndex) + "_" + strconv.Itoa(imageIndex2) + ".jpg"

				if (imageIndex2) <= obj.Copy1Max {
					copyTo1 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy1 + strconv.Itoa(imageIndex2) + ".jpg"
					err3 := CopyFile(styleNoPath, copyTo1)
					if err3 != nil {
						fmt.Printf("複製圖片 %s ， 到 %s 失敗!!\n", styleNoPath, copyTo1)
					} else {
						fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", styleNoPath, copyTo1)
					}
				}

				//copy 尺寸表
				err4 := CopyFile(styleNoPath, copyTo2)

				//試穿表 有找到MAP
				if len(jpgCutArray) >= 1 {
					for _, tryOneJpg := range jpgCutArray {

						if tryOneJpg == "" {
							continue
						}

						imageIndex2 = imageIndex2 + 1

						if imageIndex2 <= obj.Copy1Max {
							copyTo4 := mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy1 + strconv.Itoa(imageIndex2) + ".jpg"
							newTryPath := tryPicPath + tryOneJpg
							err6 := CopyFile(newTryPath, copyTo4)
							if err6 != nil {
								failTryTable = append(failTryTable, "試穿表:"+newTryPath+"失敗\n")
							}
						}

						copyTo5 := mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) +
							"_" + strconv.Itoa(imageIndex2) + ".jpg"
						newTryPath := tryPicPath + tryOneJpg
						err6 := CopyFile(newTryPath, copyTo5)
						if err6 != nil {
							failTryTable = append(failTryTable, "試穿表:"+newTryPath+"失敗\n")
						}
						needToSize1000T1040 = append(needToSize1000T1040, copyTo5)
					}
				}

				copyToStory = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) + "_" + strconv.Itoa(imageIndex2+1) + ".jpg"
				err5 := CopyFile(storyPath, copyToStory)
				needToSize1000T1040 = append(needToSize1000T1040, copyTo2)
				needToSize1000T1040 = append(needToSize1000T1040, copyToStory)
				if err4 != nil || err5 != nil {
					fmt.Printf("複製圖片 %s ， 到 %s 失敗!!\n", styleNoPath, copyTo1)
					fmt.Printf("複製圖片 %s ， 到 %s 失敗!!\n", styleNoPath, copyTo2)
					fmt.Printf("複製圖片 %s ， 到 %s 失敗!!\n", storyPath, copyToStory)
					if err5 != nil {
						failSizeTable = append(failSizeTable, "品牌故事:"+storyPath+"失敗\n")
					} else {
						failSizeTable = append(failSizeTable, "複製尺寸表:"+styleNoPath+"失敗\n")
					}

				} else {
					fmt.Printf("複製尺寸表 %s\n", styleNoPath+".jpg")
					fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", styleNoPath, copyTo1)
					fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", styleNoPath, copyTo2)
					fmt.Printf("品牌故事 %s\n", storyPath+".jpg")
					fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", storyPath, copyToStory)
				}
				BeginCount++
				//imageIndex++
				if BeginCount > MaxCount {
					newbegin, _ := strconv.Atoi(obj.BeginCount)
					BeginCount = newbegin
					imageIndex = 1
					xxl++
				}
			}
		}
	}

	var txtString string
	if len(failSizeTable) > 0 {
		for _, errorMsg := range failSizeTable {
			txtString = txtString + errorMsg
		}
		content := []byte(txtString)
		err := ioutil.WriteFile("log.txt", content, 0666)
		if err != nil {
			fmt.Println("ioutil WriteFile error: ", err)
		}
	}

	txtString = ""
	if len(failTryTable) > 0 {
		for _, errorMsg := range failTryTable {
			txtString = txtString + errorMsg
		}
		content := []byte(txtString)
		err := ioutil.WriteFile("try-log.txt", content, 0666)
		if err != nil {
			fmt.Println("ioutil WriteFile error: ", err)
		}
	}

	// txtString = ""

	// if len(needFillArr) > 0 {
	// 	for _, fillPath := range needFillArr {
	// 		txtString = txtString + fillPath + ";\n"
	// 	}
	// 	fmt.Println(txtString)
	// 	content := []byte(txtString)
	// 	err := ioutil.WriteFile(".\\needFill.txt", content, 0777)
	// 	if err != nil {
	// 		fmt.Println("ioutil WriteFile error: ", err)
	// 	}
	// }
	txtString = ""

	if len(needToSize1000T1040) > 0 {
		for _, fillPath := range needToSize1000T1040 {
			txtString = txtString + fillPath + ";\n"
		}
		fmt.Println(txtString)
		content := []byte(txtString)
		err := ioutil.WriteFile(".\\needToSize1000T1040.txt", content, 0777)
		if err != nil {
			fmt.Println("ioutil WriteFile error: ", err)
		}
	}

	fmt.Println("執行完成")
	fmt.Scanln()
}

//掃描資料夾底下檔案
func scandir(dir string) []string {
	var files []string
	filelist, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range filelist {
		files = append(files, f.Name())
	}
	return files
}

func readXXlFile(XXLFilePath string) []string {
	var content []string
	f, err := os.Open(XXLFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}
	return content
}

//byte 轉 string
func BytesToString(data []byte) string {
	return string(data[:])
}

func moveFile(orgPath string, movePath string) {

	fmt.Println(movePath)
	path := filepath.Dir(movePath)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}
	err := os.Rename(orgPath, movePath)
	if err != nil {
		fmt.Println("移動檔案失敗!!")
	}
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func mkDir(src string) (err error) {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		os.Mkdir(src, 0755)
		fmt.Println("建立資料夾:" + src)
	}
	if err != nil {
		fmt.Print("建立錯誤:")
		fmt.Print(err)
	}
	return
}

func InStringSlice(haystack []string, needle string) bool {
	for _, e := range haystack {
		if e == needle {
			return true
		}
	}

	return false
}

func dd(data string) (err error) {
	fmt.Println(data)
	os.Exit(3)
	return
}
