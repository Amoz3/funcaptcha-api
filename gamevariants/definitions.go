package gamevariants

type AudioSolver struct {
	// all these may be returning , error in the future unsure rn
	// might be justified to have a pre processing function for ensuring the file format is correct but provider and consumer are the same party so idc

	// takes path to audio file, splits and returns 3 paths, 3 paths should be in chronological order
	// file to split, dir to store splits in.
	Splitter func(string, string) []string
	// should take varargs string (output from Splitter) returns guess, 1, 2, 3
	Classifier func(...string) int
}

func (a *AudioSolver) Solve(filePath, dirPath string) int {
	splits := a.Splitter(filePath, dirPath)
	return a.Classifier(splits...)
}
