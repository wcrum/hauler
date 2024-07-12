package store

import (
	"context"
	"os"

	"github.com/mholt/archiver/v4"
	"github.com/rancherfederal/hauler/pkg/content"
	"github.com/rancherfederal/hauler/pkg/store"
	"github.com/spf13/cobra"

	"github.com/rancherfederal/hauler/pkg/log"
)

type LoadOpts struct {
	*RootOpts
}

func (o *LoadOpts) AddFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	_ = f
}

// LoadCmd
// TODO: Just use mholt/archiver for now, even though we don't need most of it
func LoadCmd(ctx context.Context, o *LoadOpts, archiveRefs ...string) error {
	l := log.FromContext(ctx)

	for _, archiveRef := range archiveRefs {
		l.Infof("loading content from [%s] to [%s]", archiveRef, o.StoreDir)
		err := unarchiveLayoutTo(ctx, archiveRef, o.StoreDir)
		if err != nil {
			return err
		}
	}

	return nil
}

// unarchiveLayoutTo accepts an archived oci layout and extracts the contents to an existing oci layout, preserving the index
func unarchiveLayoutTo(ctx context.Context, archivePath string, dest string) error {
	tmpdir, err := os.MkdirTemp("", "hauler")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpdir)

	format := archiver.CompressedArchive{
		Compression: archiver.Xz{},
		Archival:    archiver.Tar{},
	}

	fileList := []string{}

	handler := func(ctx context.Context, f archiver.File) error {
		// do something with the file
		return nil
	}

	archive, err := os.Open(archivePath)
	if err != nil {
		return err
	}

	err = format.Extract(ctx, archive, fileList, handler)
	if err != nil {
		return err
	}

	s, err := store.NewLayout(tmpdir)
	if err != nil {
		return err
	}

	ts, err := content.NewOCI(dest)
	if err != nil {
		return err
	}

	_, err = s.CopyAll(ctx, ts, nil)
	return err
}
