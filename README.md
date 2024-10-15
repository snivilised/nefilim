# 📁 nefilim: ___file system as used by traverse___

[![A B](https://img.shields.io/badge/branching-commonflow-informational?style=flat)](https://commonflow.org)
[![A B](https://img.shields.io/badge/merge-rebase-informational?style=flat)](https://git-scm.com/book/en/v2/Git-Branching-Rebasing)
[![A B](https://img.shields.io/badge/branch%20history-linear-blue?style=flat)](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/managing-a-branch-protection-rule)
[![Go Reference](https://pkg.go.dev/badge/github.com/snivilised/nefilim.svg)](https://pkg.go.dev/github.com/snivilised/nefilim)
[![Go report](https://goreportcard.com/badge/github.com/snivilised/nefilim)](https://goreportcard.com/report/github.com/snivilised/nefilim)
[![Coverage Status](https://coveralls.io/repos/github/snivilised/nefilim/badge.svg?branch=main)](https://coveralls.io/github/snivilised/nefilim?branch=main&kill_cache=1)
[![Astrolib Continuous Integration](https://github.com/snivilised/nefilim/actions/workflows/ci-workflow.yml/badge.svg)](https://github.com/snivilised/nefilim/actions/workflows/ci-workflow.yml)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)
[![A B](https://img.shields.io/badge/commit-conventional-commits?style=flat)](https://www.conventionalcommits.org/)

<!-- MD013/Line Length -->
<!-- MarkDownLint-disable MD013 -->

<!-- MD014/commands-show-output: Dollar signs used before commands without showing output mark down lint -->
<!-- MarkDownLint-disable MD014 -->

<!-- MD033/no-inline-html: Inline HTML -->
<!-- MarkDownLint-disable MD033 -->

<!-- MD040/fenced-code-language: Fenced code blocks should have a language specified -->
<!-- MarkDownLint-disable MD040 -->

<!-- MD028/no-blanks-blockquote: Blank line inside blockquote -->
<!-- MarkDownLint-disable MD028 -->

<p align="left">
  <a href="https://go.dev"><img src="resources/images/go-logo-light-blue.png" width="50" alt="go.dev" /></a>
</p>

## Introduction

___Nefilim___ is a file system abstraction used internally with snivilised packages for file system operations. In particular, it is the file system used by the directory walker as implemented by the ___traverse___ package.

An important note has to be acknowledged about usage of the file systems defined here and their comparison to the ones as defined in the Go standard library.

There are 2 ways of interacting with the file system in Go. The primary way that seems most intuitive would be to use those functions as defined within the `os` package, eg, we can open a file using os.Open:

> os.Open("~/foo.txt")

... or we can use the `os.DirFS`. However, this method uses a different concept. We can't directly open a file. We first need to create a new file system and to do so to access the local file system, we would use `os.DirFS` first:

> localFS := os.DirFS("/foo")

where ___\/foo___ represents an absolute path we have access to. The result is a file system instance as represented by the `fs.FS` interface. We can now open a file via this instance, but the crucial difference now is that we can now only use relative paths, where the path we specify is relative to the rooted path specified when we created the file system earlier:

> localFS.Open("foo.txt")

The file system defined in Nefilim, provides access to the file system in the latter case, but not yet the former (there are plans to create another abstraction that enables the more traditional way of accessing the file system, as denoted by the first example above).

Another rationale for this repo was to fill the gap left by the standard library in that there are no writer file system interfaces, so they are defined here, primarily for the purposes of snivilised projects, but also for the benefit of third parties. Contained within is an abstraction tat defines a file system as required by ___traverse___, but this particular instance only requires a subset of the full, set of operations one would expect of an file system, but there is also a __Universal File System__ which will contain the full set of operations, such as Copy, which is currently not required by ___traverse___.

There are also a few minor adjustments and additions that should be noted, such as:

* a slightly different name for creating new directories, `Mkdir` as defined in the standard packages is replaced by a more user friendly `MakeDir`.

* a new `Move` operation, which is similar to `Rename` but is defined to separate out the move semantics from rename; ie, Move will only move an item to a different directory. If a same directory move is detected, then this will be rejected with an appropriate error and the client is guided to use Rename instead.

* a new `Change` operation is defined, that is like `Move`, but is stricter in that it enforces the use of a name as the __to__ parameter denoting the destination; ie, it is prohibited to specify another relative directory as the Change operation assumes the destination should reside in the same directory as the source.

The semantics of `Rename` has not been changed so clients can expect consistent behaviour when compared to the standard package.

Other than these changes, the functionality in Nefilim aims to mirror the standard package as much as possible.

## 🎀 Features

<p align="left">
  <a href="https://onsi.github.io/ginkgo/"><img src="https://onsi.github.io/ginkgo/images/ginkgo.png" width="100" alt="ginkgo" /></a>
  <a href="https://onsi.github.io/gomega/"><img src="https://onsi.github.io/gomega/images/gomega.png" width="100" alt="gomega" /></a>
</p>

* unit testing with [Ginkgo](https://onsi.github.io/ginkgo/)/[Gomega](https://onsi.github.io/gomega/)
* linting configuration and pre-commit hooks, (see: [linting-golang](https://freshman.tech/linting-golang/)).
* uses [💥 lo](https://github.com/samber/lo)

## 📚 Usage

coming soon ...
