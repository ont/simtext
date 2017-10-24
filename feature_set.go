package main

import (
	"fmt"
	"github.com/dveselov/mystem"
	"github.com/mfonda/simhash"
	"regexp"
)

type MystemFeatureSet struct {
	Ignores []int          // array of mystem constants to ignore (i.e. mystem.Adjective, mystem.Preposition, mystem.Abbreviation ...)
	WordReg *regexp.Regexp // regexp for word detection
	Debug   bool           // if true then produce debug output

	TailPercent float32 // part of text: 0.25 is 1/4 of text
	TailLoss    int     // tail weight loss. Example: if TailLoss == 4 then all non-tail features will have weight 4 and tail features will have weight equal to one.

	txt []byte // original text
}

func NewMystemFeatureSet(txt string, ignores []int) *MystemFeatureSet {
	return &MystemFeatureSet{
		Ignores: ignores,
		WordReg: regexp.MustCompile(`[\pL-_]+`), // \pL - unicode character class L-"letters"
		txt:     []byte(txt),
	}
}

func (m *MystemFeatureSet) GetFeatures() []simhash.Feature {
	words := m.GetNormalizedWords()

	if m.Debug {
		fmt.Println("----------")
		fmt.Println(string(m.txt))
		fmt.Println("----------")
		for _, w := range words {
			print(string(w), " ")
		}
		fmt.Println("\n----------")
	}

	features := make([]simhash.Feature, len(words))

	for i, word := range words {
		if float32(i) <= float32(len(words))*(1-m.TailPercent) {
			features[i] = simhash.NewFeatureWithWeight(word, m.TailLoss) // non-tail
		} else {
			features[i] = simhash.NewFeatureWithWeight(word, 1) // tail
		}
	}

	return features
}

func (m *MystemFeatureSet) GetNormalizedWords() [][]byte {
	words := m.GetWords()
	res := make([][]byte, 0, len(words))

	for _, word := range words {
		form, good := m.normalize(word)

		if good {
			res = append(res, form)
		}
	}

	return res
}

func (m *MystemFeatureSet) GetWords() [][]byte {
	return m.WordReg.FindAll(m.txt, -1)
}

func (m *MystemFeatureSet) normalize(word []byte) (canonical []byte, good bool) {
	analyses := mystem.NewAnalyses(string(word))

	if analyses.Count() == 0 {
		return word, false // word is not detected, drop it
	}

	allGrammemes := make([]int, 0)

	for i := 0; i < analyses.Count(); i++ {
		lemma := analyses.GetLemma(i)
		grammemes := lemma.StemGram()

		if len(grammemes) == 0 {
			return word, false // word is not detected (in one of it's form), drop it
		}

		for _, g := range grammemes {
			for _, i := range m.Ignores {
				if g == i {
					return word, false
				}
			}
		}

		allGrammemes = append(allGrammemes, grammemes...)
	}

	if m.Debug {
		fmt.Println(" -->", string(word), fmt.Sprintf("%v", allGrammemes))
	}

	return []byte(analyses.GetLemma(0).Text()), true // return canonical form of word
}
