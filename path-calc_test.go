package nef_test

import (
	"fmt"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	nef "github.com/snivilised/nefilim"
)

type (
	PathCalcs map[CalcType]nef.PathCalc
)

var (
	static = struct {
		foo,
		bar,
		baz,
		foobar,
		foobarbaz,
		root string
	}{
		foo:       "foo.txt",
		bar:       "bar.txt",
		baz:       "baz.txt",
		foobar:    "foo/bar",
		foobarbaz: "foo/bar/baz.txt",
		root:      "/home/root",
	}
)

var _ = Describe("PathCalc", func() {
	DescribeTable("Base",
		func(entry *genericCalcTE[string, string]) {
			calcs := PathCalcs{
				CalcTypeAbsolute: &nef.AbsoluteCalc{},
				CalcTypeRelative: &nef.RelativeCalc{
					Root: static.root,
				},
			}

			for ct, calc := range calcs {
				Expect(calc.Base(entry.input)).To(Equal(entry.expect[ct]),
					fmt.Sprintf("💥 'Base' failed for input: '%v' (CALC:%v)", entry.input, ct),
				)
			}
		},
		func(entry *genericCalcTE[string, string]) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'",
				entry.given, entry.should,
			)
		},
		Entry(nil, &genericCalcTE[string, string]{
			calcTE: calcTE{
				given:  "path is empty",
				should: "return .",
			},
			input: "",
			expect: map[CalcType]string{
				CalcTypeAbsolute: ".",
				CalcTypeRelative: ".",
			},
		}),
		Entry(nil, &genericCalcTE[string, string]{
			calcTE: calcTE{
				given:  "path is single element",
				should: "return original path",
			},
			input: static.foo,
			expect: map[CalcType]string{
				CalcTypeAbsolute: static.foo,
				CalcTypeRelative: static.foo,
			},
		}),
		Entry(nil, &genericCalcTE[string, string]{
			calcTE: calcTE{
				given:  "path is multi element",
				should: "return last element",
			},
			input: static.foobarbaz,
			expect: map[CalcType]string{
				CalcTypeAbsolute: static.baz,
				CalcTypeRelative: static.baz,
			},
		}),
	)

	if runtime.GOOS != "windows" {
		DescribeTable("Clean",
			func(entry *genericCalcTE[string, string]) {
				calcs := PathCalcs{
					CalcTypeAbsolute: &nef.AbsoluteCalc{},
					CalcTypeRelative: &nef.RelativeCalc{
						Root: static.root,
					},
				}

				for ct, calc := range calcs {
					Expect(calc.Clean(entry.input)).To(Equal(entry.expect[ct]),
						fmt.Sprintf("💥 'Clean' failed for input: '%v' (CALC:%v)", entry.input, ct),
					)
				}
			},
			func(entry *genericCalcTE[string, string]) string {
				return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'",
					entry.given, entry.should,
				)
			},
			Entry(nil, &genericCalcTE[string, string]{
				calcTE: calcTE{
					given:  "path is empty",
					should: "return .",
				},
				input: "",
				expect: map[CalcType]string{
					CalcTypeAbsolute: ".",
					CalcTypeRelative: ".",
				},
			}),
			Entry(nil, &genericCalcTE[string, string]{
				calcTE: calcTE{
					given:  "path is single element",
					should: "return original path",
				},
				input: static.foo,
				expect: map[CalcType]string{
					CalcTypeAbsolute: static.foo,
					CalcTypeRelative: static.foo,
				},
			}),
			Entry(nil, &genericCalcTE[string, string]{
				calcTE: calcTE{
					given:  "path is multi element",
					should: "return last element",
				},
				input: static.foobarbaz,
				expect: map[CalcType]string{
					CalcTypeAbsolute: static.foobarbaz,
					CalcTypeRelative: static.foobarbaz,
				},
			}),
			Entry(nil, &genericCalcTE[string, string]{
				calcTE: calcTE{
					given:  "path contains consecutive separators",
					should: "remove consecutive separators",
				},
				input: "foo//bar///baz.txt",
				expect: map[CalcType]string{
					CalcTypeAbsolute: static.foobarbaz,
					CalcTypeRelative: static.foobarbaz,
				},
			}),

			Entry(nil, &genericCalcTE[string, string]{
				calcTE: calcTE{
					given:  "path represents root",
					should: "return appropriate root path",
				},
				input: "/",
				expect: map[CalcType]string{
					CalcTypeAbsolute: "/",
					CalcTypeRelative: ".",
				},
			}),

			Entry(nil, &genericCalcTE[string, string]{
				calcTE: calcTE{
					given:  "path ends with /",
					should: "remove trailing /",
				},
				input: "foo/bar/",
				expect: map[CalcType]string{
					CalcTypeAbsolute: static.foobar,
					CalcTypeRelative: static.foobar,
				},
			}),

			Entry(nil, &genericCalcTE[string, string]{
				calcTE: calcTE{
					given:  "path starts with /",
					should: "clean as appropriate",
				},
				input: "/foo/bar",
				expect: map[CalcType]string{
					CalcTypeAbsolute: "/foo/bar",
					CalcTypeRelative: static.foobar,
				},
			}),
		)
	}

	DescribeTable("Dir",
		func(entry *genericCalcTE[string, string]) {
			calcs := PathCalcs{
				CalcTypeAbsolute: &nef.AbsoluteCalc{},
				CalcTypeRelative: &nef.RelativeCalc{
					Root: static.root,
				},
			}

			for ct, calc := range calcs {
				Expect(calc.Dir(entry.input)).To(Equal(entry.expect[ct]),
					fmt.Sprintf("💥 'Dir' failed for input: '%v' (CALC:%v)", entry.input, ct),
				)
			}
		},
		func(entry *genericCalcTE[string, string]) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'",
				entry.given, entry.should,
			)
		},
		Entry(nil, &genericCalcTE[string, string]{
			calcTE: calcTE{
				given:  "path is empty",
				should: "return .",
			},
			input: "",
			expect: map[CalcType]string{
				CalcTypeAbsolute: ".",
				CalcTypeRelative: ".",
			},
		}),
		Entry(nil, &genericCalcTE[string, string]{
			calcTE: calcTE{
				given:  "path is single element",
				should: "return .",
			},
			input: static.foo,
			expect: map[CalcType]string{
				CalcTypeAbsolute: ".",
				CalcTypeRelative: ".",
			},
		}),
		Entry(nil, &genericCalcTE[string, string]{
			calcTE: calcTE{
				given:  "path is multi element",
				should: "return last element",
			},
			input: static.foobarbaz,
			expect: map[CalcType]string{
				CalcTypeAbsolute: static.foobar,
				CalcTypeRelative: static.foobar,
			},
		}),
	)

	DescribeTable("Elements",
		func(entry *genericCalcTE[string, []string]) {
			calcs := PathCalcs{
				CalcTypeAbsolute: &nef.AbsoluteCalc{},
				CalcTypeRelative: &nef.RelativeCalc{
					Root: static.root,
				},
			}

			for ct, calc := range calcs {
				Expect(calc.Elements(entry.input)).To(HaveExactElements(entry.expect[ct]),
					fmt.Sprintf("💥 'Elements' failed for input: '%v' (CALC:%v)", entry.input, ct),
				)
			}
		},
		func(entry *genericCalcTE[string, []string]) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'",
				entry.given, entry.should,
			)
		},
		Entry(nil, &genericCalcTE[string, []string]{
			calcTE: calcTE{
				given:  "path is empty",
				should: "return empty slice",
			},
			input: "",
			expect: map[CalcType][]string{
				CalcTypeAbsolute: {},
				CalcTypeRelative: {},
			},
		}),
		Entry(nil, &genericCalcTE[string, []string]{
			calcTE: calcTE{
				given:  "path is single element",
				should: "single element slice",
			},
			input: static.foo,
			expect: map[CalcType][]string{
				CalcTypeAbsolute: {static.foo},
				CalcTypeRelative: {static.foo},
			},
		}),
		Entry(nil, &genericCalcTE[string, []string]{
			calcTE: calcTE{
				given:  "path is multi element",
				should: "return last element",
			},
			input: static.foobarbaz,
			expect: map[CalcType][]string{
				CalcTypeAbsolute: {"foo", "bar", "baz.txt"},
				CalcTypeRelative: {"foo", "bar", "baz.txt"},
			},
		}),
	)

	DescribeTable("Join",
		func(entry *calcVariadicToOneTE) {
			calcs := PathCalcs{
				CalcTypeAbsolute: &nef.AbsoluteCalc{},
				CalcTypeRelative: &nef.RelativeCalc{
					Root: static.root,
				},
			}

			for ct, calc := range calcs {
				Expect(calc.Join(entry.input...)).To(Equal(entry.expect[ct]),
					fmt.Sprintf("💥 'Join' failed for input: '%v' (CALC:%v)", entry.input, ct),
				)
			}
		},
		func(entry *calcVariadicToOneTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'",
				entry.given, entry.should,
			)
		},
		Entry(nil, &calcVariadicToOneTE{
			calcTE: calcTE{
				given:  "path is empty",
				should: "return .",
			},
			input: []string{},
			expect: map[CalcType]string{
				CalcTypeAbsolute: "",
				CalcTypeRelative: "",
			},
		}),
		Entry(nil, &calcVariadicToOneTE{
			calcTE: calcTE{
				given:  "path is single element",
				should: "return original path",
			},
			input: []string{static.foo},
			expect: map[CalcType]string{
				CalcTypeAbsolute: static.foo,
				CalcTypeRelative: static.foo,
			},
		}),
		Entry(nil, &calcVariadicToOneTE{
			calcTE: calcTE{
				given:  "path is multi element",
				should: "return last element",
			},
			input: []string{"foo", "bar", "baz.txt"},
			expect: map[CalcType]string{
				CalcTypeAbsolute: filepath.Join("foo", "bar", "baz.txt"),
				CalcTypeRelative: static.foobarbaz,
			},
		}),
	)

	DescribeTable("Split",
		func(entry *calcOneToPairTE) {
			calcs := PathCalcs{
				CalcTypeAbsolute: &nef.AbsoluteCalc{},
				CalcTypeRelative: &nef.RelativeCalc{
					Root: static.root,
				},
			}

			for ct, calc := range calcs {
				dir, file := calc.Split(entry.input)
				Expect(pair{dir: dir, file: file}).To(Equal(entry.expect[ct]),
					fmt.Sprintf("💥 'Split' failed for input: '%v' (CALC:%v)", entry.input, ct),
				)
			}
		},
		func(entry *calcOneToPairTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'",
				entry.given, entry.should,
			)
		},
		Entry(nil, &calcOneToPairTE{
			calcTE: calcTE{
				given:  "path is empty",
				should: "return nothing",
			},
			input: "",
			expect: map[CalcType]pair{
				CalcTypeAbsolute: {dir: "", file: ""},
				CalcTypeRelative: {dir: "", file: ""},
			},
		}),
		Entry(nil, &calcOneToPairTE{
			calcTE: calcTE{
				given:  "path is single element",
				should: "return .",
			},
			input: static.foo,
			expect: map[CalcType]pair{
				CalcTypeAbsolute: {dir: "", file: static.foo},
				CalcTypeRelative: {dir: "", file: static.foo},
			},
		}),
		Entry(nil, &calcOneToPairTE{
			calcTE: calcTE{
				given:  "path is multi element",
				should: "return last element",
			},
			input: static.foobarbaz,
			expect: map[CalcType]pair{
				CalcTypeAbsolute: {dir: "foo/bar/", file: "baz.txt"},
				CalcTypeRelative: {dir: static.foobar, file: "baz.txt"},
			},
		}),
	)

	DescribeTable("Truncate",
		func(entry *genericCalcTE[string, string]) {
			calcs := PathCalcs{
				CalcTypeAbsolute: &nef.AbsoluteCalc{},
				CalcTypeRelative: &nef.RelativeCalc{
					Root: static.root,
				},
			}

			for ct, calc := range calcs {
				Expect(calc.Truncate(entry.input)).To(Equal(entry.expect[ct]),
					fmt.Sprintf("💥 'Truncate' failed for input: '%v' (CALC:%v)", entry.input, ct),
				)
			}
		},
		func(entry *genericCalcTE[string, string]) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'",
				entry.given, entry.should,
			)
		},
		Entry(nil, &genericCalcTE[string, string]{
			calcTE: calcTE{
				given:  "path is empty",
				should: "return .",
			},
			input: "",
			expect: map[CalcType]string{
				CalcTypeAbsolute: ".",
				CalcTypeRelative: ".",
			},
		}),
		Entry(nil, &genericCalcTE[string, string]{
			calcTE: calcTE{
				given:  "path does not contain a trailing separator",
				should: "return original path",
			},
			input: static.foobar,
			expect: map[CalcType]string{
				CalcTypeAbsolute: static.foobar,
				CalcTypeRelative: static.foobar,
			},
		}),
		Entry(nil, &genericCalcTE[string, string]{
			calcTE: calcTE{
				given:  "path contains trailing separator",
				should: "return truncate the separator",
			},
			input: "foo/bar/",
			expect: map[CalcType]string{
				CalcTypeAbsolute: static.foobar,
				CalcTypeRelative: static.foobar,
			},
		}),
	)
})
