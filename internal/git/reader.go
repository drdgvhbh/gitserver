package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/src-d/go-billy.v4"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

const dotGitFolderName = ".git"

type Reader interface {
	Open(path string) (Repository, error)
}

type StorageReader struct {
	fileSystem billy.Filesystem
}

func NewReader(fileSystem billy.Filesystem) Reader {
	return &StorageReader{
		fileSystem: fileSystem,
	}
}

// Open opens a git repository from the given path
func (reader *StorageReader) Open(path string) (Repository, error) {
	folder := findDotGitFolder(path)

	dotgitFolder, err := reader.fileSystem.Chroot(folder)
	if err != nil {
		return nil, err
	}

	repoRoot, err := reader.fileSystem.Chroot(path)
	if err != nil {
		return nil, err
	}

	storage := filesystem.NewStorageWithOptions(
		dotgitFolder,
		cache.NewObjectLRUDefault(),
		filesystem.Options{KeepDescriptors: true})
	repo, err := gogit.Open(storage, repoRoot)
	if err != nil {
		return nil, err
	}

	return &GitRepository{Wrapee: repo}, nil
}

func findDotGitFolder(path string) string {
	dotGitPath := fmt.Sprintf("%s/.git", path)
	dotGitStat, err := os.Stat(dotGitPath)
	if err != nil {
		return ""
	}

	repositoryIsNotASubmodule := dotGitStat.IsDir()
	if repositoryIsNotASubmodule {
		return dotGitPath
	}

	dotGitData, err := ioutil.ReadFile(dotGitPath)
	if err != nil {
		return ""
	}

	relativeDotGitPath := strings.Trim(
		fmt.Sprintf("%s/%s",
			path,
			strings.Split(string(dotGitData), " ")[1]), "\n")

	return relativeDotGitPath
}
