package git

import (
	"gopkg.in/src-d/go-billy.v4"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

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
	repoRoot, err := reader.fileSystem.Chroot(path)
	if err != nil {
		return nil, err
	}

	dotgitFolder, err := repoRoot.Chroot(".git")
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
