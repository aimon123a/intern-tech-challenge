package main

import (
	"context"
	"fmt"

	"strconv"
	"strings"
	"sort"

	"os"
	"bufio"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)
var m map[string]string

// get key value from str
// e.g 1.8.0 --> key = "1.8"
func getKey(str string)string{
	var s = ""
	for i := len(str)-1; i > -1; i--{
		if str[i] == '.'{
			s = str[:i]
			break
		}
	}
	return s
}

//check existence of key
//and compare current key value vs new value
//e.g key[1.8] v = 0  : num3 = 1
// then update key[1.8] v = 1
func checkKey(key string, num3 int){
	if m[key] == ""{
		m[key] = strconv.Itoa(num3)
		return
	}
	intKey,err := strconv.Atoi(m[key])
	if err != nil{
		panic(err)
	}
	if intKey < num3{
		m[key] = strconv.Itoa(num3)
	}
}
// sort the keys in order to print in reversed format
func sortKey()[]string{
	var keysArr []int
	keysArr = make([]int, len(m))
	strKey := make([]string, len(m))
	i := 0
	// for key append to int array
	// e.g 1.8 = 108 & 1.10 = 110
	for k := range m { 
		splitkey := strings.Split(k,".")
		num1, err := strconv.Atoi(splitkey[0])
		num2, err := strconv.Atoi(splitkey[1])
		num1 = num1 * 100
		if err != nil {return strKey}
		keysArr[i] = num1+num2
		i++
	}
	sort.Ints(keysArr)
	// conver int array to string array
	for j := 0;j < len(keysArr); j++ {
		s := []string{strconv.Itoa(keysArr[j]/100),strconv.Itoa(keysArr[j]%100)}
		strKey[j] = strings.Join(s,".")
	}
	return strKey
}
//for test case purposes
func mapVersions(releases []*semver.Version, versionSlice []*semver.Version, minVersion *semver.Version) []*semver.Version {
	versionSlice = LatestVersions(releases, minVersion)
	if len(versionSlice) == 0 {
		versionSlice = make([]*semver.Version, 1)
		versionSlice[0] = minVersion
	}
	return versionSlice
}

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	versionSlice := make([]*semver.Version, len(releases))
	var str = fmt.Sprintf("%v",minVersion)
	m = make(map[string]string)
	version := strings.Split(str,"-")
	splitV := strings.Split(version[0],".")
	firstMinV, secondMinV, thirdMinV := returnVersions(splitV)

	key := ""
	//Position := 0
	// This is just an example structure of the code, if you implement this interface, the test cases in main_test.go are very easy to run
	for i := range releases{
		str = fmt.Sprintf("%v", releases[i])

		version := strings.Split(str,"-")
		splitV := strings.Split(version[0],".")
		//compare version by splitting into 3 digits
		num1, num2, num3 := returnVersions(splitV)
		if num1 < firstMinV {continue}
		if num1 > firstMinV {
			key = getKey(version[0])
			checkKey(key,num3)
			continue
		}
		if num2 < secondMinV {continue}
		if num2 > secondMinV {
			key = getKey(version[0])
			checkKey(key,num3)
			continue
		}
		if num3 < thirdMinV || num3 == thirdMinV {continue}
		key = getKey(version[0])
		checkKey(key,num3)
	}
	//get sorted keys
	sortedKey := sortKey()
	j := 0
	//put versions back into versionSlice from the sorted keyArray
	for i := len(sortedKey)-1; i > -1; i--{
		versionSlice[j] = semver.New(sortedKey[i] + "." + m[sortedKey[i]])
		j++
	}
	return versionSlice[:len(sortedKey)]
}
// for each version in test file or retrieved from github, split them into 3 numbers 
// to compare
func returnVersions(str []string)(int,int,int){
	firstMinV, err := strconv.Atoi(str[0])
	secondMinV, err := strconv.Atoi(str[1])
	thirdMinV, err := strconv.Atoi(str[2])
	if err != nil{
		panic(err)
	}
	return firstMinV,secondMinV,thirdMinV
}
// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	// Github
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}

	//os read in txt file
	file, err := os.Open("test.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		line := scanner.Text()
		splitLine := strings.Split(line,"/")
		repository := splitLine[0]
		version := strings.Split(splitLine[1],",")[1]
		releases, _, err := client.Repositories.ListReleases(ctx, repository, repository, opt)
		if err != nil {
			panic(err) // is this really a good way?
		}
		minVersion := semver.New(version)
		allReleases := make([]*semver.Version, len(releases))
		for i, release := range releases {
			versionString := *release.TagName
			
			if versionString[0] == 'v' {
				versionString = versionString[1:]
			}
			allReleases[i] = semver.New(versionString)
		}
		versionSlice := LatestVersions(allReleases, minVersion)
		fmt.Printf("latest versions of %s/%s: %s\n", repository, repository, versionSlice)
	}
	
}