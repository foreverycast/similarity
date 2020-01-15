# similarity
Search for duplicates

## INTRODUCTION
Similarity script is for searching duplicates in master data. The script use Levenshtein distance to find similar names. 

## INSTALLATION
 * Create csv files for compare (replace the names in os.Open("duplicates_" + country + ".csv") for file to load and os.Create("result_" + country + ".csv") for result csv)
 * If neccessary, the additional file can be added, important for intergrating data set isAddFileToCheck = True and replace "additional.csv" for file name 
 * Go Lang is required
 * The package is still in development
 * Just execute the  go run duplicatesearch.go

## ABOUT THIS RELEASE
* Tests are not ready
* stopwords.csv is not working
