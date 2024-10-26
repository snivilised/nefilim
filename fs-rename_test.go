package nef_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
	"github.com/snivilised/nefilim/test/luna"
)

var _ = Describe("op: rename", Ordered, func() {
	var (
		root string
		fS   nef.RenameFS
	)

	BeforeAll(func() {
		root = luna.Repo("test")
	})

	BeforeEach(func() {
		fS = nef.NewUniversalFS(nef.Rel{
			Root:      root,
			Overwrite: false,
		})

		scratch(root)
	})

	DescribeTable("fs: RenameFS",
		func(entry fsTE[nef.RenameFS]) {
			if entry.arrange != nil {
				entry.arrange(entry, fS)
			}
			entry.action(entry, fS)
		},
		func(entry fsTE[nef.RenameFS]) string {
			return fmt.Sprintf("🧪 ===> given: target is '%v', %v should: '%v'",
				entry.given, entry.op, entry.should,
			)
		},
		// The following tests are a duplicate of those defined for the move
		// operation 🔆, but with appropriately adjusted expectations.
		Entry(nil, fsTE[nef.RenameFS]{
			given:   "[from] file exists, [to] directory exists, [no-clash]",
			should:  "fail, because filename is missing, from to path",
			note:    "filename not included in the destination path (from/file.txt => to)",
			op:      "Rename",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Move.From.File,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.RenameFS], _ nef.RenameFS) {
				Expect(require(root, entry.require, entry.from)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
			},
			action: func(entry fsTE[nef.RenameFS], fS nef.RenameFS) {
				luna.IsLinkError(fS.Rename(entry.from, lab.Static.FS.Move.Destination), entry.should)
			},
		}),

		Entry(nil, fsTE[nef.RenameFS]{
			given:   "[from] file exists, [to] directory exists, [clash]",
			should:  "succeed, only if override",
			note:    "filename not included in the destination path (from/file.txt => to)",
			op:      "Move",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Move.From.File,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.RenameFS], _ nef.RenameFS) {
				Expect(require(root,
					entry.require,
					entry.from,
				)).To(Succeed())
				Expect(require(root,
					lab.Static.FS.Move.Destination,
					lab.Static.FS.Move.To.File,
				)).To(Succeed())
			},
			action: func(entry fsTE[nef.RenameFS], fS nef.RenameFS) {
				luna.IsLinkError(fS.Rename(entry.from, lab.Static.FS.Move.Destination), entry.should)
			},
		}),

		Entry(nil, fsTE[nef.RenameFS]{
			given:   "[from] file exists, [to] directory exists, [no-clash]",
			should:  "succeed",
			note:    "filename IS included in the destination path (from/file.txt => to/file.txt)",
			op:      "Move",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Move.From.File,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.RenameFS], _ nef.RenameFS) {
				Expect(require(root, entry.require, entry.from)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
			},
			action: func(entry fsTE[nef.RenameFS], fS nef.RenameFS) {
				destination := lab.Static.FS.Move.To.File
				Expect(fS.Rename(entry.from, destination)).To(Succeed())
				Expect(luna.AsFile(destination)).To(luna.ExistInFS(fS))
			},
		}),

		Entry(nil, fsTE[nef.RenameFS]{
			given:   "[from] file exists, [to] directory exists, [clash]",
			should:  "succeed; move and overwrite",
			note:    "filename IS included in the destination path (from/file.txt => to/file.txt)",
			op:      "Move",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Move.From.File,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.RenameFS], _ nef.RenameFS) {
				Expect(require(root, entry.require, entry.from)).To(Succeed())

				if entry.overwrite {
					Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
					return
				}
				Expect(require(root,
					lab.Static.FS.Move.Destination,
					lab.Static.FS.Move.To.File,
				)).To(Succeed())
			},
			action: func(entry fsTE[nef.RenameFS], fS nef.RenameFS) {
				Expect(fS.Rename(entry.from, lab.Static.FS.Move.To.File)).To(Succeed())
				Expect(luna.AsFile(lab.Static.FS.Move.From.File)).NotTo(luna.ExistInFS(fS))
			},
		}),

		Entry(nil, fsTE[nef.RenameFS]{
			given:   "[from] directory exists, [to] directory exists, [no clash]",
			should:  "fail, because dir name is missing, from to path",
			note:    "directory not included in the destination path (from/dir => to)",
			op:      "Move",
			require: lab.Static.FS.Move.From.Directory,
			from:    lab.Static.FS.Move.From.Directory,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.RenameFS], _ nef.RenameFS) {
				Expect(require(root, entry.require)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
			},
			action: func(entry fsTE[nef.RenameFS], fS nef.RenameFS) {
				luna.IsLinkError(fS.Rename(entry.from, lab.Static.FS.Move.Destination), entry.should)
			},
		}),

		Entry(nil, fsTE[nef.RenameFS]{
			given:   "[from] directory exists, [to] directory exists, [clash]",
			should:  "fail",
			note:    "directory not included in the destination path (from/dir => to)",
			op:      "Move",
			require: lab.Static.FS.Move.From.Directory,
			from:    lab.Static.FS.Move.From.Directory,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.RenameFS], _ nef.RenameFS) {
				Expect(require(root, entry.require)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.To.Directory)).To(Succeed())
			},
			action: func(entry fsTE[nef.RenameFS], fS nef.RenameFS) {
				luna.IsLinkError(fS.Rename(entry.from, lab.Static.FS.Move.Destination), entry.should)
				Expect(luna.AsDirectory(lab.Static.FS.Move.From.Directory)).To(luna.ExistInFS(fS))
			},
		}),

		Entry(nil, fsTE[nef.RenameFS]{
			given:   "[from] directory exists, [to] directory exists, [no clash]",
			should:  "succeed",
			note:    "directory IS included in the destination path (from/dir => to/dir)",
			op:      "Move",
			require: lab.Static.FS.Move.From.Directory,
			from:    lab.Static.FS.Move.From.Directory,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.RenameFS], _ nef.RenameFS) {
				Expect(require(root, entry.require)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
			},
			action: func(entry fsTE[nef.RenameFS], fS nef.RenameFS) {
				destination := lab.Static.FS.Move.To.Directory
				Expect(fS.Rename(entry.from, destination)).To(Succeed())
				Expect(luna.AsDirectory(destination)).To(luna.ExistInFS(fS))
			},
		}),

		Entry(nil, fsTE[nef.RenameFS]{
			given:   "[from] directory exists, [to] directory exists, [clash]",
			should:  "fail",
			note:    "directory IS included in the destination path (from/dir => to/dir)",
			op:      "Move",
			require: lab.Static.FS.Move.From.Directory,
			from:    lab.Static.FS.Move.From.Directory,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.RenameFS], _ nef.RenameFS) {
				Expect(require(root, entry.require)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.To.Directory)).To(Succeed())
			},
			action: func(entry fsTE[nef.RenameFS], fS nef.RenameFS) {
				destination := lab.Static.FS.Move.To.Directory
				Expect(fS.Rename(entry.from, destination)).NotTo(Succeed())
				Expect(luna.AsDirectory(lab.Static.FS.Move.From.Directory)).To(luna.ExistInFS(fS))
			},
		}),

		// 💠 The tests in the follow section are defined for scenarios where the
		// target item is being renamed into the same directory; ie it is a rename,
		// without moving to a different directory.
		//
		Entry(nil, fsTE[nef.RenameFS]{
			given:   "[from] file exists, [to] name does not exist, [no-clash]",
			should:  "succeed",
			op:      "Rename",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Rename.From.File,
			to:      lab.Static.FS.Rename.To.File,
			arrange: func(entry fsTE[nef.RenameFS], _ nef.RenameFS) {
				Expect(require(root, entry.require, entry.from)).To(Succeed())
			},
			action: func(entry fsTE[nef.RenameFS], fS nef.RenameFS) {
				Expect(fS.Rename(entry.from, entry.to)).To(Succeed())
			},
		}),

		Entry(nil, fsTE[nef.RenameFS]{
			given:   "[from] file exists, [to] equal to [from], [clash]",
			should:  "succeed, ignored",
			op:      "Rename",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Rename.From.File,
			to:      lab.Static.FS.Rename.From.File,
			arrange: func(entry fsTE[nef.RenameFS], _ nef.RenameFS) {
				Expect(require(root, entry.require, entry.from)).To(Succeed())
			},
			action: func(entry fsTE[nef.RenameFS], fS nef.RenameFS) {
				Expect(fS.Rename(entry.from, entry.to)).To(Succeed())
			},
		}),

		Entry(nil, fsTE[nef.RenameFS]{
			given:   "[from] directory exists, [to] name does not exist, [no-clash]",
			should:  "succeed",
			op:      "Rename",
			require: lab.Static.FS.Rename.From.Directory,
			from:    lab.Static.FS.Rename.From.Directory,
			to:      lab.Static.FS.Rename.To.Directory,
			arrange: func(entry fsTE[nef.RenameFS], _ nef.RenameFS) {
				Expect(require(root, entry.require)).To(Succeed())
			},
			action: func(entry fsTE[nef.RenameFS], fS nef.RenameFS) {
				Expect(fS.Rename(entry.from, entry.to)).To(Succeed())
			},
		}),

		Entry(nil, fsTE[nef.RenameFS]{
			given:   "[from] directory exists, [to] equal to [from], [clash]",
			should:  "fail, directory names can't be same",
			op:      "Rename",
			require: lab.Static.FS.Rename.From.Directory,
			from:    lab.Static.FS.Rename.From.Directory,
			to:      lab.Static.FS.Rename.From.Directory,
			arrange: func(entry fsTE[nef.RenameFS], _ nef.RenameFS) {
				Expect(require(root, entry.require)).To(Succeed())
			},
			action: func(entry fsTE[nef.RenameFS], fS nef.RenameFS) {
				luna.IsLinkError(fS.Rename(entry.from, entry.to), entry.should)
			},
		}),
	)
})
