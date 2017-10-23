package main

import (
	"fmt"
	"github.com/dveselov/mystem"
	"github.com/mfonda/simhash"
	"golang.org/x/text/unicode/norm"
	"io/ioutil"
	"regexp"
	"strings"
)

func GetHash(fname string) uint64 {
	dat, _ := ioutil.ReadFile(fname)
	//txt := strip.StripTags(string(dat))
	txt := string(dat)

	r := regexp.MustCompile("<[^>]+>|[^а-яА-Я]+")
	txt = r.ReplaceAllString(txt, " ")

	println("---------------")
	println(txt)
	println("---------------")

	//txt := "проверяется доступное место на сервере"
	words := strings.Fields(txt)

	wordsCanonical := make([]string, 0, len(words))

	for _, word := range words {
		analyses := mystem.NewAnalyses(word)

		bad := false
		arr := make([]int, 0)
		for i := 0; i < analyses.Count(); i++ {
			lemma := analyses.GetLemma(i)
			grammemes := lemma.StemGram()

			arr = append(arr, grammemes...)
			if len(grammemes) == 0 {
				bad = true
				break
			}
			for _, g := range grammemes {
				if /*g == mystem.Conjunction ||*/ g == mystem.Interjunction || g == mystem.Preposition || g == mystem.Abbreviation || g == mystem.Adjective {
					//fmt.Println("--->", lemma.Text())
					bad = true
					break
				}
			}

			if bad {
				break
			}

			//fmt.Printf("%s - %s %v (%d)\n", word, lemma.Text(), grammemes, analyses.Count())
		}
		if !bad {
			fmt.Print(analyses.GetLemma(0).Text(), " ", arr, " ")
			wordsCanonical = append(wordsCanonical, analyses.GetLemma(0).Text())
			//fmt.Print(analyses.GetLemma(0).Text(), " ")
		}

		analyses.Close()
	}

	txtCanonical := strings.Join(wordsCanonical, " ")
	println("\n---------------")
	fmt.Printf("%v %v\n", txtCanonical, []byte(txtCanonical))
	println("---------------")
	hash := simhash.Simhash(simhash.NewUnicodeWordFeatureSet([]byte(txtCanonical), norm.NFC))

	return hash
}

func main() {
	//hash1 := GetHash("test.txt")
	//hash2 := GetHash("test2.txt")
	//println(simhash.Compare(hash1, hash2))

	m := NewMystemFeatureSet("<b>Тестовая, строка... ? слов123!! 123 bestмыыыыхх", []int{})
	for _, word := range m.GetWords() {
		println(string(word))
	}
	println("---------------")
	for _, word := range m.GetNormalizedWords() {
		println(string(word))
	}

	//fmt.Println(fmt.Sprintf("Analyze of '%s':", "маша"))
	//for i := 0; i < analyses.Count(); i++ {
	//    lemma := analyses.GetLemma(i)
	//    grammemes := lemma.StemGram()
	//    fmt.Println(fmt.Sprintf("%d. %s - %v", i+1, lemma.Text(), grammemes))
	//}
}
