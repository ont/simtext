package main

import (
	"fmt"
	"github.com/dveselov/mystem"
	"github.com/mfonda/simhash"
	"io/ioutil"
	"os"
	"strconv"
)

type Group struct {
	hash  uint64
	names []string
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s [threshold] [dir]\n", os.Args[0])
		os.Exit(1)
	}

	thresh, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Error: value %s is not integer\n", os.Args[1])
		os.Exit(1)
	}

	dir := os.Args[2]
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if os.Getenv("DEBUG") != "" {
		debugFiles(dir, files)
		os.Exit(0)
	}

	findGroups(dir, files, thresh)
}

func debugFiles(dir string, files []os.FileInfo) {
	for _, file := range files {
		data, err := ioutil.ReadFile(dir + "/" + file.Name())

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		hash := calcHash(data, true)

		fmt.Printf("%s --> %s", file.Name(), strconv.FormatInt(int64(hash), 2))
	}
}

func findGroups(dir string, files []os.FileInfo, thresh int) {
	groups := make([]*Group, 0)

	for _, file := range files {
		data, err := ioutil.ReadFile(dir + "/" + file.Name())

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		group := findGroup(groups, data, thresh)

		if group == nil {
			group = newGroup(data)
			groups = append(groups, group)
		}

		group.names = append(group.names, file.Name())
	}

	printGroups(groups)
}

func findGroup(groups []*Group, txt []byte, thresh int) *Group {
	hash := calcHash(txt, false)

	for _, group := range groups {
		if simhash.Compare(hash, group.hash) < uint8(thresh) {
			return group
		}
	}

	return nil
}

func newGroup(txt []byte) *Group {
	return &Group{
		hash:  calcHash(txt, false),
		names: make([]string, 0),
	}
}

func calcHash(txt []byte, debug bool) uint64 {
	fset := NewMystemFeatureSet(string(txt), []int{
		mystem.Interjunction,
		mystem.Preposition,
		mystem.Abbreviation,
		mystem.Adjective,
		mystem.Particle,
		mystem.AdjPronoun,
	})

	fset.Debug = debug

	return simhash.Simhash(fset)
}

func printGroups(groups []*Group) {
	var total int

	for _, group := range groups {
		if len(group.names) == 1 {
			continue
		}

		total++

		fmt.Println("-----------")
		for _, name := range group.names {
			fmt.Println(name)
		}
	}

	fmt.Println("===============")
	fmt.Printf("Total: %d groups\n", total)
}
