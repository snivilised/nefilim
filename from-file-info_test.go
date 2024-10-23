package nef_test

import (
	"io/fs"
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
)

// testFileInfo is a struct that implements the fs.FileInfo interface
type testFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
	sys     any
}

// Name returns the base name of the file
func (fi *testFileInfo) Name() string {
	return fi.name
}

// Size returns the length in bytes for regular files; system-dependent for others
func (fi *testFileInfo) Size() int64 {
	return fi.size
}

// Mode returns the file mode bits
func (fi *testFileInfo) Mode() fs.FileMode {
	return fi.mode
}

// ModTime returns the modification time
func (fi *testFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir returns true if it is a directory, false otherwise
func (fi *testFileInfo) IsDir() bool {
	return fi.isDir
}

// Sys returns underlying data source (can return nil)
func (fi *testFileInfo) Sys() any {
	return fi.sys
}

var _ = Describe("FromFileInfo", func() {
	Context("given: a FileInfo", func() {
		It("🧪 should: convert to DirEntry", func() {
			const (
				size   = 234
				layout = "2006-01-02 00:00:00"
			)

			released, _ := time.Parse(layout, "1986-09-29 00:00:00")

			entry := nef.FromFileInfo(&testFileInfo{
				name:    "heaven-can-wait.flac",
				size:    size,
				mode:    lab.Perms.File,
				modTime: released,
			})
			Expect(entry.Name()).To(Equal("heaven-can-wait.flac"))
			Expect(entry.IsDir()).To(BeFalse())
			Expect(entry.Type()).To(Equal(lab.Perms.File))
			Expect(entry.Info()).NotTo(BeNil())
		})
	})
})
