package main

import (
	"github.com/dveselov/mystem"
	"github.com/mfonda/simhash"
	"regexp"
)

type MystemFeatureSet struct {
	Ignores []int          // array of mystem constants to ignore (i.e. mystem.Adjective, mystem.Preposition, mystem.Abbreviation ...)
	WordReg *regexp.Regexp // regexp for word detection

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
	features := make([]simhash.Feature, len(words))

	for i, word := range words {
		features[i] = simhash.NewFeature(word)
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
	}

	return []byte(analyses.GetLemma(0).Text()), true // return canonical form of word
}
