package store

import (
	"context"
	"os"
	"path/filepath"

	"github.com/mholt/archiver/v4"
	"github.com/spf13/cobra"

	"github.com/rancherfederal/hauler/pkg/log"
)

type SaveOpts struct {
	*RootOpts
	FileName string
}

func (o *SaveOpts) AddArgs(cmd *cobra.Command) {
	f := cmd.Flags()

	f.StringVarP(&o.FileName, "filename", "f", "haul.tar.zst", "Name of archive")
}

// SaveCmd
// TODO: Just use mholt/archiver for now, even though we don't need most of it
func SaveCmd(ctx context.Context, o *SaveOpts, outputFile string) error {
	l := log.FromContext(ctx)

	absOutputfile, err := filepath.Abs(outputFile)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(cwd)
	if err := os.Chdir(o.StoreDir); err != nil {
		return err
	}

	files, err := archiver.FilesFromDisk(nil, map[string]string{
		".": "", // contents added recursively
	})
	if err != nil {
		return err
	}

	out, err := os.Create(absOutputfile)
	if err != nil {
		return err
	}
	defer out.Close()

	// we can use the CompressedArchive type to gzip a tarball
	// (compression is not required; you could use Tar directly)
	format := archiver.CompressedArchive{
		Compression: archiver.Xz{},
		Archival:    archiver.Tar{},
	}

	// create the archive
	err = format.Archive(context.Background(), out, files)
	if err != nil {
		return err
	}

	l.Infof("saved store [%s] -> [%s]", o.StoreDir, absOutputfile)
	return nil
}
