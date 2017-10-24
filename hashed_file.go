package main

import (
	"github.com/dveselov/mystem"
	"github.com/mfonda/simhash"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type HashedFile struct {
	Name string
	Hash uint64

	path string
}

func NewHashedFile(path string) (*HashedFile, error) {
	basename := filepath.Base(path)
	var hash = &HashedFile{
		Name: strings.TrimSuffix(basename, filepath.Ext(basename)),
		path: path,
	}

	err := hash.RefreshHash()
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func (f *HashedFile) RefreshHash() error {
	data, err := ioutil.ReadFile(f.path)

	if err != nil {
		return err
	}

	f.Hash = f.calcHash(data)
	return nil
}

func (f *HashedFile) calcHash(txt []byte) uint64 {
	fset := NewMystemFeatureSet(string(txt), []int{
		mystem.Interjunction,
		mystem.Preposition,
		mystem.Abbreviation,
		mystem.Adjective,
		mystem.Particle,
		mystem.AdjPronoun,
	})

	//fset.Debug = debug
	fset.TailPercent = 0.25
	fset.TailLoss = 3

	return simhash.Simhash(fset)
}
