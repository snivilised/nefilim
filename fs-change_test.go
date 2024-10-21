package nef_test

import (
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
)

// Note [clash/no-clash]: when an item is moved to the destination, clash
// refers to the scenario where the source already exists in the destination,
// whereas no-clash is the opposite.
//
// Eg, if moving src/foo.txt to dest/foo.txt, a clash scenario means foo.txt
// already exists in dest/ and no clash where it doesn't. The success of the
// operation depends on wether the overwrite flag has been specified in the
// filesystem.

var _ = Describe("op: change", Ordered, func() {
	var (
		root string
		fS   nef.UniversalFS
	)

	BeforeAll(func() {
		root = Repo("test")
	})

	DescribeTable("fs: UniversalFS",
		func(entry fsTE[nef.UniversalFS]) {
			for _, overwrite := range []bool{false, true} {
				scratch(root)

				fS = nef.NewUniversalFS(nef.At{
					Root:      root,
					Overwrite: overwrite,
				})
				entry.overwrite = overwrite

				if entry.arrange != nil {
					entry.arrange(entry, fS)
				}
				entry.action(entry, fS)
			}
		},
		func(entry fsTE[nef.UniversalFS]) string {
			return fmt.Sprintf("🧪 ===> given: target is '%v', %v should: '%v'",
				entry.given, entry.op, entry.should,
			)
		},

		// + SUCCESS CASES

		// FROM IS DIRECTORY, TO DOES NOT EXIST (same for both tentative & overwrite)
		Entry(nil, fsTE[nef.UniversalFS]{
			given:   "[from] directory exists, [to] directory missing, [no-clash]",
			should:  "succeed",
			op:      "Change",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Change.From.Directory,
			to:      lab.Static.FS.Change.To.Directory,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root, entry.require)).To(Succeed())
				Expect(require(root, entry.from)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				_, destination := nef.SplitParent(entry.to)
				Expect(fS.Change(entry.from, destination)).To(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
				Expect(AsDirectory(lab.Static.FS.Change.Destination)).To(ExistInFS(fS))
			},
		}),

		// FROM IS FILE, TO DOES NOT EXIST (same for both tentative & overwrite)
		// FROM IS DIRECTORY, TO DOES NOT EXIST (same for both tentative & overwrite)
		Entry(nil, fsTE[nef.UniversalFS]{
			given:   "[from] file exists, [to] file missing, [no-clash]",
			should:  "succeed",
			op:      "Change",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Change.From.File,
			to:      lab.Static.FS.Change.To.File,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root, entry.require, entry.from)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				Expect(fS.Change(entry.from, entry.to)).To(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
				Expect(AsFile(Join(lab.Static.FS.Scratch, entry.to))).To(ExistInFS(fS))
			},
		}),

		// FROM IS FILE, TO EXISTS: (overwrite only)
		Entry(nil, fsTE[nef.UniversalFS]{
			given:   "[from] file exists, [to] file missing, [no-clash]",
			should:  "succeed",
			op:      "Change",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Change.From.File,
			to:      lab.Static.FS.Change.To.File,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root, entry.require, entry.from)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				Expect(fS.Change(entry.from, entry.to)).To(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
				file := Join(lab.Static.FS.Scratch, entry.to)
				Expect(AsFile(file)).To(ExistInFS(fS))
			},
		}),

		// + FAILURE CASES (other than not in same dir as they have already been covered)

		// FROM IS DIRECTORY, TO EXISTS: (overwrite) -> [MERGE] os.rename, (tentative) -> PROHIBITED

		Entry(nil, fsTE[nef.UniversalFS]{
			given:   "[from] directory exists, [to] directory exists, [clash]",
			should:  "succeed, same file ignored",
			op:      "Change",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Change.From.Directory,
			to:      lab.Static.FS.Change.To.Directory,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root, entry.from)).To(Succeed())
				Expect(require(root, entry.to)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				_, destination := nef.SplitParent(entry.to)
				Expect(fS.Change(entry.from, destination)).NotTo(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
			},
		}),

		// FROM IS FILE, TO EXISTS: (tentative) -> PROHIBITED

		Entry(nil, fsTE[nef.UniversalFS]{
			given:   "[from] file exists, [to] directory exists, [clash]",
			should:  "succeed, same file ignored",
			op:      "Change",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Change.From.File,
			to:      lab.Static.FS.Change.To.File,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root,
					entry.require,
					entry.from,
				)).To(Succeed())
				Expect(require(root,
					entry.require,
					Join(lab.Static.FS.Scratch, entry.to),
				)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				Expect(fS.Change(entry.from, entry.to)).To(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
				Expect(AsFile(entry.from)).NotTo(ExistInFS(fS))
			},
		}),

		// 🔆
		XEntry(nil, Label("INVALID:DIFF DIR"), fsTE[nef.UniversalFS]{
			given:   "[from] file exists, [to] file missing, [no-clash]",
			should:  "succeed",
			note:    "filename not included in the destination path (from/file.txt => to)", // !!! WRONG
			op:      "Change",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Move.From.File,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root, entry.require, entry.from)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				Expect(fS.Change(entry.from, lab.Static.FS.Move.Destination)).To(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
				Expect(AsFile(lab.Static.FS.Move.To.File)).To(ExistInFS(fS))
			},
		}),

		XEntry(nil, Label("INVALID:DIFF DIR"), fsTE[nef.UniversalFS]{
			given:   "[from] file exists, [to] directory exists, [clash]",
			should:  "succeed, only if override",
			note:    "filename not included in the destination path (from/file.txt => to)",
			op:      "Change",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Move.From.File,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root,
					entry.require,
					entry.from,
				)).To(Succeed())
				Expect(require(root,
					lab.Static.FS.Move.Destination,
					lab.Static.FS.Move.To.File,
				)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				if entry.overwrite {
					Expect(fS.Change(entry.from, lab.Static.FS.Move.Destination)).To(Succeed(),
						fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
					)
					Expect(AsFile(lab.Static.FS.Move.To.File)).To(ExistInFS(fS))
					return
				}
				Expect(fS.Change(entry.from, lab.Static.FS.Move.Destination)).NotTo(Succeed())
			},
		}),

		XEntry(nil, Label("INVALID:DIFF DIR"), fsTE[nef.UniversalFS]{
			given:   "[from] file exists, [to] directory exists, [no-clash]",
			should:  "succeed",
			note:    "filename IS included in the destination path (from/file.txt => to/file.txt)",
			op:      "Change",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Move.From.File,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root, entry.require, entry.from)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				destination := lab.Static.FS.Move.To.File
				Expect(fS.Change(entry.from, destination)).To(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
				Expect(AsFile(destination)).To(ExistInFS(fS))
			},
		}),

		XEntry(nil, Label("INVALID:DIFF DIR"), fsTE[nef.UniversalFS]{
			given:   "[from] file exists, [to] directory exists, [clash]",
			should:  "succeed, only if override",
			note:    "filename IS included in the destination path (from/file.txt => to/file.txt)",
			op:      "Change",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Move.From.File,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
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
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				destination := lab.Static.FS.Move.To.File

				if entry.overwrite {
					Expect(fS.Change(entry.from, destination)).To(Succeed(),
						fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
					)
					Expect(AsFile(destination)).To(ExistInFS(fS))
					return
				}
				Expect(fS.Change(entry.from, destination)).NotTo(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
			},
		}),

		XEntry(nil, Label("INVALID:DIFF DIR"), fsTE[nef.UniversalFS]{
			given:   "[from] directory exists, [to] directory exists, [no clash]",
			should:  "succeed",
			note:    "directory not included in the destination path (from/dir => to)",
			op:      "Change",
			require: lab.Static.FS.Move.From.Directory,
			from:    lab.Static.FS.Move.From.Directory,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root, entry.require)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				Expect(fS.Change(entry.from, lab.Static.FS.Move.Destination)).To(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
				Expect(AsDirectory(lab.Static.FS.Move.To.Directory)).To(ExistInFS(fS))
			},
		}),

		XEntry(nil, Label("INVALID:DIFF DIR"), fsTE[nef.UniversalFS]{
			given:   "[from] directory exists, [to] directory exists, [clash]",
			should:  "fail",
			note:    "directory not included in the destination path (from/dir => to)",
			op:      "Change",
			require: lab.Static.FS.Move.From.Directory,
			from:    lab.Static.FS.Move.From.Directory,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root, entry.require)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.To.Directory)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				if entry.overwrite {
					Expect(fS.Change(entry.from, lab.Static.FS.Move.Destination)).NotTo(Succeed(),
						fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
					)
					return
				}

				Expect(fS.Change(entry.from, lab.Static.FS.Move.Destination)).NotTo(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
			},
		}),

		XEntry(nil, Label("INVALID:DIFF DIR"), fsTE[nef.UniversalFS]{
			given:   "[from] directory exists, [to] directory exists, [no clash]",
			should:  "succeed",
			note:    "directory IS included in the destination path (from/dir => to/dir)",
			op:      "Change",
			require: lab.Static.FS.Move.From.Directory,
			from:    lab.Static.FS.Move.From.Directory,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root, entry.require)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				destination := lab.Static.FS.Move.To.Directory
				Expect(fS.Change(entry.from, destination)).To(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
				Expect(AsDirectory(destination)).To(ExistInFS(fS))
			},
		}),

		XEntry(nil, Label("INVALID:DIFF DIR"), fsTE[nef.UniversalFS]{
			given:   "[from] directory exists, [to] directory exists, [clash]",
			should:  "fail",
			note:    "directory IS included in the destination path (from/dir => to/dir)",
			op:      "Change",
			require: lab.Static.FS.Move.From.Directory,
			from:    lab.Static.FS.Move.From.Directory,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root, entry.require)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.Destination)).To(Succeed())
				Expect(require(root, lab.Static.FS.Move.To.Directory)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				destination := lab.Static.FS.Move.To.Directory
				Expect(fS.Change(entry.from, destination)).NotTo(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
			},
		}),

		Entry(nil, fsTE[nef.UniversalFS]{
			given:   "[from] path missing",
			should:  "fail",
			op:      "Change",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.Foo,
			to:      lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root, entry.require)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				err := fS.Change(entry.from, entry.to)
				Expect(err).NotTo(Succeed(), fmt.Sprintf("OVERWRITE: %v", entry.overwrite))
				Expect(nef.IsBinaryFsOpError(err)).To(BeTrue())
			},
		}),

		XEntry(nil, Label("INVALID:DIFF DIR"), fsTE[nef.UniversalFS]{
			given:   "[from] file exists, [to] directory missing",
			should:  "fail",
			op:      "Change",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Move.From.File,
			to:      lab.Static.FS.Move.Destination,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root, entry.require, entry.from)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				destination := filepath.Join(entry.to, lab.Static.Foo)
				Expect(fS.Change(entry.from, destination)).NotTo(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
			},
		}),

		// The following tests are a duplicate of those defined for the rename
		// operation 💠, where the target item is being renamed into the same directory,
		// but these should be rejected with an error, because this amounts to a
		// Move which is not the intended use of Change.
		Entry(nil, fsTE[nef.UniversalFS]{
			given:   "[from] file exists, [to] name does not exist, [no-clash]",
			should:  "fail, [to] path should not include directory path",
			op:      "Change",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Rename.From.File,
			to:      lab.Static.FS.Rename.To.File,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root, entry.require, entry.from)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				IsInvalidPathError(
					fS.Change(entry.from, entry.to), entry.should,
				)
			},
		}),

		Entry(nil, fsTE[nef.UniversalFS]{
			given:   "[from] file exists, [to] equal to [from], [clash]",
			should:  "fail, [to] path should not include directory path",
			op:      "Change",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Rename.From.File,
			to:      lab.Static.FS.Rename.From.File,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root, entry.require, entry.from)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				IsInvalidPathError(
					fS.Change(entry.from, entry.to), entry.should,
				)
			},
		}),

		Entry(nil, fsTE[nef.UniversalFS]{
			given:   "[from] directory exists, [to] name does not exist, [no-clash]",
			should:  "fail, [to] path should not include directory path",
			op:      "Change",
			require: lab.Static.FS.Rename.From.Directory,
			from:    lab.Static.FS.Rename.From.Directory,
			to:      lab.Static.FS.Rename.To.Directory,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root, entry.require)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				IsInvalidPathError(
					fS.Change(entry.from, entry.to), entry.should,
				)
			},
		}),

		Entry(nil, fsTE[nef.UniversalFS]{
			given:   "[from] directory exists, [to] equal to [from], [clash]",
			should:  "fail, [to] path should not include directory path",
			op:      "Change",
			require: lab.Static.FS.Rename.From.Directory,
			from:    lab.Static.FS.Rename.From.Directory,
			to:      lab.Static.FS.Rename.From.Directory,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root, entry.require)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				IsInvalidPathError(
					fS.Change(entry.from, entry.to), entry.should,
				)
			},
		}),

		// ===> may have to repeat the above to cover successful cases where to
		// does not contains a different directory, but these may already be covered.
	)
})
