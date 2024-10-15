package nef_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
)

var _ = Describe("op: remove", Ordered, func() {
	var (
		root string
		fS   nef.UniversalFS
	)

	BeforeAll(func() {
		root = Repo("test")
	})

	BeforeEach(func() {
		scratchPath := filepath.Join(root, lab.Static.FS.Scratch)

		if _, err := os.Stat(scratchPath); err == nil {
			Expect(os.RemoveAll(scratchPath)).To(Succeed(),
				fmt.Sprintf("failed to delete existing directory %q", scratchPath),
			)
		}
	})

	DescribeTable("removal",
		func(entry fsTE[nef.UniversalFS]) {
			for _, overwrite := range []bool{false, true} {
				fS = nef.NewUniversalFS(nef.At{
					Root:      root,
					Overwrite: entry.overwrite,
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
			given:   "file and exists",
			should:  "succeed",
			op:      "Remove",
			require: lab.Static.FS.Scratch,
			target:  lab.Static.FS.Remove.File,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				err := require(root, entry.require, entry.target)
				Expect(err).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				if entry.overwrite {
					// tbd
					return
				}
				Expect(fS.Remove(entry.target)).To(Succeed())
			},
		}),
		Entry(nil, fsTE[nef.UniversalFS]{
			given:   "path does not exist",
			should:  "fail",
			op:      "Remove",
			require: lab.Static.FS.Scratch,
			target:  lab.Static.Foo,
			action: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				if entry.overwrite {
					// tbd
					return
				}
				Expect(fS.Remove(entry.target)).To(MatchError(os.ErrNotExist))
			},
		}),
		Entry(nil, fsTE[nef.UniversalFS]{
			given:   "directory exists and not empty",
			should:  "fail",
			op:      "Remove",
			require: lab.Static.FS.Scratch,
			target:  lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				err := require(root, entry.require, lab.Static.FS.Remove.File)
				Expect(err).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				if entry.overwrite {
					// tbd
					return
				}
				Expect(errors.Unwrap(fS.Remove(entry.target))).To(
					MatchError("directory not empty"),
				)
			},
		}),
		//
		Entry(nil, fsTE[nef.UniversalFS]{
			given:   "path does not exist",
			should:  "succeed",
			op:      "RemoveAll",
			require: lab.Static.FS.Scratch,
			target:  lab.Static.Foo,
			action: func(entry fsTE[nef.UniversalFS], fS nef.UniversalFS) {
				if entry.overwrite {
					// tbd
					return
				}
				Expect(fS.RemoveAll(entry.target)).To(Succeed())
			},
		}),
		Entry(nil, fsTE[nef.UniversalFS]{
			given:   "directory exists and not empty",
			should:  "succeed",
			op:      "RemoveAll",
			require: lab.Static.FS.Scratch,
			target:  lab.Static.FS.Scratch,
			arrange: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				err := require(root, entry.require, lab.Static.FS.Remove.File)
				Expect(err).To(Succeed())
			},
			action: func(entry fsTE[nef.UniversalFS], _ nef.UniversalFS) {
				if entry.overwrite {
					// tbd
					return
				}
				Expect(fS.RemoveAll(entry.target)).To(Succeed())
			},
		}),
	)
})
