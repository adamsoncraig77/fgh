package repos

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Matt-Gleich/fgh/pkg/utils"
	"github.com/Matt-Gleich/statuser/v2"
	"github.com/briandowns/spinner"
)

// Get all repos in the directory and all subdirectories
func Repos(rootPath string, fatalErr bool) []LocalRepo {
	spin := spinner.New(utils.SpinnerCharSet, utils.SpinnerSpeed)
	spin.Suffix = " Getting list of repos"
	spin.Start()

	oldRepos := []LocalRepo{}
	err := filepath.Walk(
		rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && IsGitRepo(path) {

				absPath, err := filepath.Abs(path)
				if err != nil {
					statuser.Error("Failed to get absolute path for "+path, err, 1)
				}

				owner, name, err := OwnerAndNameFromRemote(path)
				if err != nil {
					msg := "Failed to get owner and name from remote in " + path
					if fatalErr {
						statuser.Error(msg, err, 1)
					}
					statuser.Warning(msg + fmt.Sprintln(err))
					return nil
				}

				oldRepos = append(oldRepos, LocalRepo{
					Owner: owner,
					Name:  name,
					Path:  absPath,
				})
			}
			return nil
		},
	)

	spin.Stop()
	if err != nil {
		statuser.Error("Failed to get list of repos", err, 1)
	}

	if len(oldRepos) == 0 {
		statuser.Warning("0 repos found inside " + rootPath)
	}
	return oldRepos
}
