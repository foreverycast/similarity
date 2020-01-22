package main


import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/agnivade/levenshtein"
)

var NumberOfAccounts int
var NumberOfSalesforceAccounts int

type Configuration struct {
	FileToCompare      string
	IsAddFileToCheck   bool
	AddFileNameToCheck string
	ResultFileName     string
	MinSimilarRanking  float64
	StopWordsFile      string
}

type Account struct {
	Id                string
	Name              string
	Street            string
	Compare           string
	CompareName       string
	CompareStreet     string
	Benchmark         float64
	BenchmarkName     float64
	BenchmarkStreet   float64
	RankingCombined   float64
	DuplicateId       string
	DuplicateNameId   string
	DuplicateStreetId string
	CompareString     string
}

type StopWord struct {
	replace_to   string
	replace_with string
}

func floattostr(fv float64) string {
	return strconv.FormatFloat(fv, 'f', 2, 64)
}

func loadingProgress(i int) float64 {
	return ((float64(i) * (2*float64(NumberOfAccounts) + 2*float64(NumberOfSalesforceAccounts) - float64(i) - 1) / 2) / (float64(NumberOfAccounts)*float64(NumberOfSalesforceAccounts) + float64(NumberOfAccounts)*(float64(NumberOfAccounts)-1)/2) * 100)
}

func calculateRating(object string, toCompareObject string) float64 {
	return ((math.Max(float64(len(object)), float64(len(toCompareObject))) - float64(levenshtein.ComputeDistance(object, toCompareObject))) / math.Max(float64(len(object)), float64(len(toCompareObject))))
}

func replaceStopWords(str string, stopwords []StopWord) string {
	//str = strings.ToLower(str)
	for index := range stopwords {
		strings.ReplaceAll(str, stopwords[index].replace_to, stopwords[index].replace_with)
		// fmt.Println(stopwords[index].replace_to, "index")
	}
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(str), "工業股份有限公司", ""), "股份有限公司", ""), "有限公司", ""), "科技", ""), ".", " "), ",", " "), " inc", ""), " corporation", ""), " gmbh", ""), " ltd", ""), "-", " "), "#", " "), " ", ""), "company", ""), "電子", ""), "电子", ""), "宁波", ""), "aktiengesellschaft", "")
}

func main() {
	// Read config file
	fileConf, _ := os.Open("conf.json")
	defer fileConf.Close()
	decoder := json.NewDecoder(fileConf)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	var fileToCompare string = configuration.FileToCompare
	var isAddFileToCheck bool = configuration.IsAddFileToCheck
	var addFileNameToCheck string = configuration.AddFileNameToCheck
	var resultFileName string = configuration.ResultFileName
	var wordsToTest []string
	var compare []string
	var accountsId []string
	var compare_name []string
	var compareStreet []string

	fmt.Println("Duplicates search started")

	fmt.Println("File to compare", fileToCompare)

	reg, err := regexp.Compile("[ ]{2,}")
	if err != nil {
		log.Fatal(err)
	}

	csvStopFile, _ := os.Open(configuration.StopWordsFile)
	readerStop := csv.NewReader(bufio.NewReader(csvStopFile))
	var stopWords []StopWord
	for {
		line, error := readerStop.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		stopWords = append(stopWords, StopWord{replace_to: line[0], replace_with: line[1]})
		fmt.Println(line)
	}

	fmt.Println(stopWords)

	csvFile, _ := os.Open(fileToCompare)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var accounts []Account
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		accounts = append(accounts, Account{
			Id:      line[0],
			Name:    replaceStopWords(line[1], stopWords),
			Street:  reg.ReplaceAllString(strings.ToLower(line[2]), ""),
			Compare: strings.ToLower(replaceStopWords(line[1], stopWords) + " " + line[2]),
		})
		wordsToTest = append(wordsToTest, strings.ToLower(replaceStopWords(line[1], stopWords)+" "+line[2]))
		compare = append(compare, strings.ToLower(replaceStopWords(line[1], stopWords)+" "+line[2]))
		compare_name = append(compare_name, replaceStopWords(line[1], stopWords))
		compareStreet = append(compareStreet, reg.ReplaceAllString(strings.ToLower(line[2]), ""))
		accountsId = append(accountsId, line[0])
		//fmt.Printf(s, line[2])
	}

	NumberOfAccounts = len(compare)
	if isAddFileToCheck {
		csvSfFile, _ := os.Open(addFileNameToCheck)
		readerSF := csv.NewReader(bufio.NewReader(csvSfFile))
		for {
			line, error := readerSF.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				log.Fatal(error)
			}
			compare = append(compare, strings.ToLower(replaceStopWords(line[1], stopWords)+" "+line[2]))
			accountsId = append(accountsId, line[0])
			compare_name = append(compare_name, replaceStopWords(line[1], stopWords))
			compareStreet = append(compareStreet, reg.ReplaceAllString(strings.ToLower(line[2]), ""))
		}
	}
	NumberOfSalesforceAccounts = len(compare) - NumberOfAccounts

	file, err := os.Create(resultFileName)
	defer file.Close()

	if err != nil {
		os.Exit(1)
	}

	var isDebug bool = false

	csvWriter := csv.NewWriter(file)
	for index, element := range accounts {

		for i, com := range compare {
			if com == element.Compare {
				compare = append(compare[:i], compare[i+1:]...)
				accountsId = append(accountsId[:i], accountsId[i+1:]...)
				compare_name = append(compare_name[:i], compare_name[i+1:]...)
				compareStreet = append(compareStreet[:i], compareStreet[i+1:]...)
				break
			}
		}

		for i, com := range compare {

			rankingName := calculateRating(element.Name, compare_name[i])

			rankingStreet := float64(0)
			ranking := float64(0)
			if rankingName > configuration.MinSimilarRanking {

				ranking = calculateRating(element.Compare, com)

				streetNameList := strings.Split(element.Street, " ")
				streetNameToCompare := strings.Split(compareStreet[i], " ")
				for j := range streetNameList {
					rankingStreetPart := float64(0)
					for n := range streetNameToCompare {
						rankingStreetPartTemp := calculateRating(streetNameList[j], streetNameToCompare[n])
						if rankingStreetPartTemp > rankingStreetPart {
							rankingStreetPart = rankingStreetPartTemp
							if rankingStreetPart == 1 {
								break
							}
						}
					}
					rankingStreet += rankingStreetPart
				}
				rankingStreet = rankingStreet / float64(len(streetNameList))
			}
			rankingCombined := (ranking + rankingName + rankingStreet) / 3

			if rankingCombined > element.RankingCombined {

				element.DuplicateId = accountsId[i]
				element.CompareString = com
				element.DuplicateStreetId = accountsId[i]
				element.CompareStreet = compareStreet[i]

				element.Benchmark = ranking
				element.BenchmarkName = rankingName
				element.BenchmarkStreet = rankingStreet

				element.RankingCombined = rankingCombined

				if ranking == 1.00 {
					element.RankingCombined = ranking
					break
				}
			}

		}

		// Loading progress
		if index%10 == 0 {
			fmt.Println()
			fmt.Println(floattostr(loadingProgress(index)), "% ready")
			fmt.Println()
		}

		fmt.Println("Similarity: " + floattostr((element.Benchmark * 100)) + "% between " + element.Compare + " and " + element.DuplicateId)
		var s []string
		s = append(s, element.Id, element.Name, element.Street, element.DuplicateId, element.CompareString, floattostr(element.Benchmark), floattostr(element.RankingCombined), element.DuplicateNameId, element.CompareName, floattostr(element.BenchmarkName), element.DuplicateStreetId, element.CompareStreet, floattostr(element.BenchmarkName))
		csvWriter.Write(s)
		if isDebug && index == 10 {
			break
		}
	}

	csvWriter.Flush()
}
