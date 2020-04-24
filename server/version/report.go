package version

import "fmt"

func Report() {
	fmt.Printf("Version : %d.%d-%d\n", VersionMajor, VersionMinor, VersionBuild)
	fmt.Printf("Built : %s\n", BuildTime)
	fmt.Printf("Repo URL : %s\n", RepoURL)
	fmt.Printf("Branch : %s\n", BranchName)
	fmt.Printf("Commit Id : Short : %s Long %s\n", ShortCommitId, LongCommitId)
	fmt.Printf("Assigned Tags : %s\n", AssignedTags)
}
