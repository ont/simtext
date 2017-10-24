package main

type Group struct {
	hash  uint64
	files []*HashedFile
}

/*
func findGroups(dir string, files []os.FileInfo, thresh int) {
	groups := make([]*Group, 0)

	for _, file := range files {

		hashedFile := NewHashedFile(file.Name(), data)

		group := findGroup(groups, hashedFile, thresh)

		if group == nil {
			group = newGroup()
			groups = append(groups, group)
		}

		addToGroup(group, hashedFile)
	}

	printGroups(groups)
}

func findGroup(groups []*Group, file *HashedFile, thresh int) *Group {
	for _, group := range groups {
		if simhash.Compare(file.Hash, group.hash) < uint8(thresh) {
			return group
		}
	}

	return nil
}

func newGroup() *Group {
	return &Group{
		files: make([]*HashedFile, 0),
	}
}

func addToGroup(group *Group, file *HashedFile) {
	if len(group.files) == 0 {
		group.hash = file.Hash
	}

	group.files = append(group.files, file)
}

func printGroups(groups []*Group) {
	var total int

	for _, group := range groups {
		if len(group.files) == 1 {
			continue
		}

		total++

		fmt.Println("-----------")
		for _, file := range group.files {
			fmt.Println(file.Name, fmt.Sprintf("%064b", file.Hash))
		}
	}

	fmt.Println("===============")
	fmt.Printf("Total: %d groups\n", total)
}
*/
