package nef_test

import (
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
	"github.com/snivilised/nefilim/test/luna"
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
		root = luna.Repo("test")
	})

	DescribeTable("fs: UniversalFS",
		func(entry fsTE[nef.UniversalFS]) {
			for _, overwrite := range []bool{false, true} {
				scratch(root)

				fS = nef.NewUniversalFS(nef.Rel{
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
				destination := filepath.Base(entry.to)
				Expect(fS.Change(entry.from, destination)).To(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
				Expect(luna.AsDirectory(lab.Static.FS.Change.Destination)).To(luna.ExistInFS(fS))
				Expect(luna.AsDirectory(destination)).NotTo(luna.ExistInFS(fS))
			},
		}),

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
				Expect(luna.AsFile(luna.Yoke(lab.Static.FS.Scratch, entry.to))).To(luna.ExistInFS(fS))
				Expect(luna.AsDirectory(entry.to)).NotTo(luna.ExistInFS(fS))
			},
		}),

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
				file := luna.Yoke(lab.Static.FS.Scratch, entry.to)
				Expect(luna.AsFile(file)).To(luna.ExistInFS(fS))
				Expect(luna.AsDirectory(entry.to)).NotTo(luna.ExistInFS(fS))
			},
		}),

		Entry(nil, fsTE[nef.UniversalFS]{
			given:   "[from] directory exists, [to] directory exists, [clash]",
			should:  "fail, overwrite/merge not supported",
			op:      "Change",
			require: lab.Static.FS.Scratch,
			from:    lab.Static.FS.Change.From.Directory,
			to:      lab.Static.FS.Change.To.Directory,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				Expect(require(root, entry.from)).To(Succeed())
				Expect(require(root, entry.to)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				destination := filepath.Base(entry.to)
				Expect(fS.Change(entry.from, destination)).NotTo(Succeed(),
					fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
				)
				luna.IsLinkError(fS.Change(entry.from, destination), entry.should)
				Expect(luna.AsDirectory(destination)).NotTo(luna.ExistInFS(fS))
			},
		}),

		Entry(nil, fsTE[nef.UniversalFS]{
			given:   "[from] file exists, [to] file exists [clash]",
			should:  "succeed, ignored",
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
					luna.Yoke(lab.Static.FS.Scratch, entry.to),
				)).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				err := fS.Change(entry.from, entry.to)
				if entry.overwrite {
					Expect(err).To(Succeed(),
						fmt.Sprintf("OVERWRITE: %v", entry.overwrite),
					)
					Expect(luna.AsFile(entry.from)).NotTo(luna.ExistInFS(fS))

					return
				}
				Expect(nef.IsBinaryFsOpError(err)).To(BeTrue())
				Expect(luna.AsFile(entry.from)).To(luna.ExistInFS(fS))
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
	)
})
