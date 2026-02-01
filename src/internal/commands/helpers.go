package commands

// filterUsedArgs filters out already-used arguments from suggestions
// Used by ValidArgsFunction in cobra commands for shell completion
func filterUsedArgs(suggestions []string, usedArgs []string) []string {
	usedSet := make(map[string]bool)
	for _, arg := range usedArgs {
		usedSet[arg] = true
	}

	var filtered []string
	for _, s := range suggestions {
		if !usedSet[s] {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
