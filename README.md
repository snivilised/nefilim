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

<!-- MD029/Ordered list item prefix -->
<!-- MarkDownLint-disable MD029 -->

<!-- vscode-markdown-toc -->
* 1. [Introduction](#Introduction)
* 2. [🎀 Features](#Features)
* 3. [🔰 Quick Start](#QuickStart)
* 4. [🎭 Relative vs Absolute](#RelativeVsAbsolute)
* 5. [📚 Usage](#Usage)
  * 5.1. [📂 File Systems](#FileSystems)
    * 5.1.1. [✨ Universal FS](#UniversalFS)
    * 5.1.2. [✨ Traverse FS](#TraverseFS)
    * 5.1.3. [✨ Exists In FS](#ExistsInFS)
    * 5.1.4. [✨ Read File FS](#ReadFileFS)
    * 5.1.5. [✨ Reader FS](#ReaderFS)
    * 5.1.6. [✨ Make Dir FS](#MakeDirFS)
    * 5.1.7. [✨ Move FS](#MoveFS)
    * 5.1.8. [✨ Change FS](#ChangeFS)
    * 5.1.9. [✨ Copy FS](#CopyFS)
    * 5.1.10. [✨ Remove FS](#RemoveFS)
    * 5.1.11. [✨ Rename FS](#RenameFS)
    * 5.1.12. [✨ Write File FS](#WriteFileFS)
    * 5.1.13. [✨ Writer FS](#WriterFS)
* 6. [Overwrite Flag](#OverwriteFlag)
* 7. [💔 Errors](#Errors)
  * 7.1. [⛔ Binary Fs Op Error](#BinaryFsOpError)
  * 7.2. [⛔ Invalid Path Error](#InvalidPathError)
  * 7.3. [⛔ Reject Same Directory Move Error](#RejectSameDirectoryMoveError)
  * 7.4. [⛔ Reject Different Directory Change Error](#RejectDifferentDirectoryChangeError)
* 8. [Utilities](#Utilities)
  * 8.1. [🛡️ EnsureAtPath](#EnsureAtPath)
  * 8.2. [🛡️ResolvePath](#ResolvePath)
* 9. [💥 Trouble Shooting](#TroubleShooting)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

<p align="left">
  <a href="https://go.dev"><img src="resources/images/go-logo-light-blue.png" width="50" alt="go.dev" title=" (THIS DOCUMENTATION IS STILL A WORK IN PROGRESS...)"/></a>
</p>

## 1. <a name='Introduction'></a>Introduction

___Nefilim___ is a file system abstraction used internally with ___snivilised___ packages for file system operations. In particular, it is the file system used by the directory walker as implemented by the ___traverse___ package.

An important note has to be acknowledged about usage of the file systems defined here and their comparison to the ones as defined in the Go standard library.

There are 2 ways of interacting with the file system in Go. The primary way that seems most intuitive would be to use those functions as defined within the `os` package, eg, we can open a file using os.Open:

> os.Open("~/foo.txt")

... or we can use the `os.DirFS`. However, this method uses a different concept. We can't directly open a file. We first need to create a new file system and to do so, to access the local file system, we would use `os.DirFS` first:

> localFS := os.DirFS("/foo")

where ___\/foo___ represents an absolute path we have access to. The result is a file system instance as represented by the `fs.FS` interface. We can now open a file via this instance, but the crucial difference now is that we can now only use relative paths, where the path we specify is relative to the rooted path specified when we created the file system earlier:

> localFS.Open("foo.txt")

The file system defined in Nefilim, provides access to the file system in the latter case (ie via a relative file system), but not yet the former, absolute file access (there are plans to create another abstraction that enables this more traditional way of accessing the file system, as denoted by the first example above).

Another rationale for this repo was to fill the gap left by the standard library, in that there are no writer file system interfaces, so they are defined here, primarily for the purposes of ___snivilised___ projects, but also for the benefit of third parties. Contained within is an abstraction that defines a file system as required by ___traverse___, but this particular instance only requires a subset of the full set of operations one would expect of a file system, but there is also a __Universal File System__ which will contain the full set of operations, such as Copy, which is currently not required by ___traverse___.

There are also a few minor adjustments and additions that should be noted, such as:

* a slightly different name for creating new directories, `Mkdir` as defined in the standard packages is replaced by a more user friendly `MakeDir`. (This is just a minor issue, but having to remember wether the 'd' in `Mkdir` was capitalised or not, is just friction that I would rather do without.)

* a new `Move` operation, which is similar to `Rename` but is defined to separate out the move semantics from rename; ie, Move will only move an item to a different directory. If a same directory move is detected, then this will be rejected with an appropriate error and the client is guided to use Rename instead.

* a new `Change` operation is defined, that is like `Move`, but is stricter in that it enforces the use of a name as the __to__ parameter denoting the destination to be in the same directory; ie, it is prohibited to specify another relative directory as the Change operation assumes the destination should reside in the same directory as the source.

The semantics of `Rename` remains unchanged, so clients can expect consistent behaviour when compared to the standard package.

Other than these changes, the functionality in Nefilim aims to mirror the standard package as much as possible.

## 2. <a name='Features'></a>🎀 Features

<p align="left">
  <a href="https://onsi.github.io/ginkgo/"><img src="https://onsi.github.io/ginkgo/images/ginkgo.png" width="100" alt="ginkgo" /></a>
  <a href="https://onsi.github.io/gomega/"><img src="https://onsi.github.io/gomega/images/gomega.png" width="100" alt="gomega" /></a>
</p>

* unit testing with [Ginkgo](https://onsi.github.io/ginkgo/)/[Gomega](https://onsi.github.io/gomega/)
* linting configuration and pre-commit hooks, (see: [linting-golang](https://freshman.tech/linting-golang/)).
* uses [💥 lo](https://github.com/samber/lo)

## 3. <a name='QuickStart'></a>🔰 Quick Start

* To create a universal file system that contains all reader and writer functions:

```go
  import (
    nef "github.com/snivilised/nefilim"
  )

  fS := nef.NewUniversalFS(nef.At{
    Root:      "/Users/marina/dev",
    Overwrite: false,
  })
```

... creates a file system whose root is _/Users/marina/dev_

Any operation invoked, should now be done with a path that is relative to this root, eg to open a file:

```go
  if file, err := fS.Open("foo.txt"); err!= nil {
    ...
  }
```

... will succeed if the file exists at _/Users/marina/dev/foo.txt_

When creating a file system with writer capabilities, the ___Overwrite___ flag can be set within the ___At___ struct, which will activate `overwrite` semantics, that are explained later for each writer operation.

## 4. <a name='RelativeVsAbsolute'></a>🎭 Relative vs Absolute

There are various file system constructor functions in the form NewXxxFS. Currently, these are all of the relative variety, whereby the client is required to invoke operations with paths that are relative to the root. ___Nefilim___ conforms to the semantics of ___io/fs___, so any paths that are not conformant with ___fs.ValidPath___ will be rejected.

The key rules for paths to confirm to are:

* must be unrooted
* must not start or end with a '/'
* must not contain '.' or '..' or the empty string, expect for the special case '.' which refers to root
* paths are forward '/' separated only, for all platforms
* characters such as backslash and colon are still valid, but should not be interpreted as path separators

## 5. <a name='Usage'></a>📚 Usage

### 5.1. <a name='FileSystems'></a>📂 File Systems

____Nefilim___ comes with predefined interfaces with different capabilities. The interfaces are as narrow as possible, most are single method interfaces. Some interfaces, contain more than 1 closely related methods. Clearly, ___Nefilim___ can't provide interfaces for all combination of methods, but the client is free to compose custom ones by combining those defined here by ___Nefilim___.

#### 5.1.1. <a name='UniversalFS'></a>✨ Universal FS

Capable of all read and write operations and can be used with ___traverse___ if so required. Actually, as previously indicated, traverse doesn't need a UniversalFS, it only requires a ___TraverseFS___, so why would we use a ___UniversalFS___ with ___traverse___? Well, within the callback of ___traverse___ navigation, we may need to invoke operations, not defined on ___TraverseFS___. But beware, do not invoke operations that would interfere with the currently running navigation session, without making required mitigating actions.

* interface: ___UniversalFS___
* Create: ___NewUniversalFS___

```go
  fS := nef.NewUniversalFS(nef.At{
    Root:      "/Users/marina/dev",
    Overwrite: false,
  })
```

* Composed of: ___ReaderFS___, ___WriterFS___

---

#### 5.1.2. <a name='TraverseFS'></a>✨ Traverse FS

A specialised file system as required for a ___traverse___ navigation.

* interface: ___TraverseFS___
* Create: ___NewTraverseFS___

```go
  fS := nef.NewTraverseFS(nef.At{
    Root:      "/Users/marina/dev",
    Overwrite: false,
  })

  result, err := tv.Walk().Configure().Extent(tv.Prime(
    &tv.Using{
      Tree:         "some-path-relative",
      Subscription: enums.SubscribeUniversal,
      Handler: func(node *core.Node) error {
        GinkgoWriter.Printf(
          "---> 🍯 EXAMPLE-REGEX-FILTER-CALLBACK: '%v'\n", node.Path,
        )
        return nil
      },
      GetTraverseFS: func(_ string) tv.TraverseFS {
        return fS
      },
    },
  )).Navigate(ctx)
```

In the above example, we create a universal file system rooted at _/Users/marina/dev_. The ___Tree___ path set in the ___tv.Using___ struct is relative to this root.

* Composed of: ___MakeDirFS___, ___ReaderFS___, ___WriteFileFS___

---

#### 5.1.3. <a name='ExistsInFS'></a>✨ Exists In FS

A file system that can determine the existence of a path and indicate if its a file or directory.

* interface: ___ExistsInFS___
* Create: ___NewExistsInFS___

```go
  fS := nef.NewExistsInFS(nef.At{
    Root:      "/Users/marina/dev",
  })
```

* Commands: FileExists, DirectoryExists

##### 💎 FileExists

> fS.FileExists("bar/baz/foo.txt")

returns true if _/Users/marina/dev/bar/baz/foo.txt_ exists as a file, false otherwise.

##### 💎 DirectoryExists

> fS.DirectoryExists("bar/baz")

returns true if _/Users/marina/dev/bar/baz/_ exists as a directory, false otherwise.

---

#### 5.1.4. <a name='ReadFileFS'></a>✨ Read File FS

* interface: ___ReadFileFS___

* Create: ___NewReadFileFS___

```go
  fS := nef.NewReadFileFS(nef.At{
    Root:      "/Users/marina/dev",
    Overwrite: false,
  })
```

* Composed of: ___fs.FS___
* Commands:

##### 💎 ReadFile

> fS.ReadFile("bar/baz/foo.txt")

returns no error if _/Users/marina/dev/bar/baz/foo.txt_ exists as a file, otherwise behaves as ___fs.ReadFile___.

---

#### 5.1.5. <a name='ReaderFS'></a>✨ Reader FS

Creates a read only file system.

* interface: ___ReaderFS___
* Create: ___NewReaderFS___

```go
  fS := nef.NewReaderFS(nef.At{
    Root:      "/Users/marina/dev",
  })
```

* Composed of: ___fs.StatFS___, ___fs.ReadDirFS___, ExistsInFS, ReadFileFS

---

#### 5.1.6. <a name='MakeDirFS'></a>✨ Make Dir FS

* interface: ___MakeDirFS___
* Create: ___NewMakeDirFS___

```go
  fS := nef.NewMakeDirFS(nef.At{
    Root:      "/Users/marina/dev",
    Overwrite: false,
  })
```

* Composed of: ___ExistsInFS___
* Commands: ___MakeDir___, ___MakeDirAll___

##### 💎 MakeDir

> fS.MakeDir("bar/baz")

behaves as ___os.Mkdir___.

##### 💎 MakeDirAll

> fS.MakeDir("bar/baz")

behaves as ___os.MkdirAll___.

---

#### 5.1.7. <a name='MoveFS'></a>✨ Move FS

Comes as part of the ___UniversalFS___ only. The ___Move___ command is a new operation, that does not exist in the standard library, created to isolate the `move` semantics of the ___os.Rename___ command. ___os.Rename___ implements both `move` and `rename` semantics combined.

Another problem with ___os.Rename___ occurs when moving a file eg:

when a file needs to be moved from _bar/file.txt_ to the directory _bar/baz/_, invoking ___so.Rename___ the intuitive way would be to do as follows:

> os.Rename("_bar/file.txt_", "_bar/baz/_")

but this will fail with a ___LinkError___. The correct way to achieve the desired result is

> os.Rename("_bar/file.txt_", "_bar/baz/file.txt_")

ie, the file name has to be replicated in the 'newpath' path.

The ___Move___ command challenges this requirement and allows the client to omit the file name from the second parameter and can be achieved as:

> fS.Move("_bar/file.txt_", "_bar/baz_")

In these examples, we have talked about the new path representing a file, but the same holds true for a directory.

Also, ___Move___ really does mean move, the new name is always retained, not renamed, and the new path always represents a different directory from the source.

* interface: ___MoveFS___
* Create: ___NewUniversalFS___

```go
  fS := nef.NewUniversalFS(nef.At{
    Root:      "/Users/marina/dev",
    Overwrite: false,
  })
```

* Commands: ___Move___

> fS.Move("_bar/file.txt_", "_bar/baz_")

As described previously, moves _/Users/marina/dev/bar/file.txt_ to _/Users/marina/dev/bar/baz/file.txt_. However, the behaviour differs depending on the prior existence of _bar/baz/file.txt_ and the value of the overwrite flag.

If the file already exists at the destination and _overwrite_ is _true_, then the existing file is overwritten, otherwise, an invalid file system operation ___NewInvalidBinaryFsOpError___ is returned. This denotes the from and to path and the name of the operation attempted, in this case _Move_.

#### 5.1.8. <a name='ChangeFS'></a>✨ Change FS

Comes as part of the ___UniversalFS___ only (implementation pending as of v0.1.2). The ___Change___ command is a new operation, that does not exist in the standard library, created to isolate the `rename` semantics of the ___os.Rename___ command.

The ___Change___ command imposes a further restriction to ___os.Rename___. In the same way that the ___Move___ command rejects setting a destination path that denotes the same directory as the source, ___Change___ will reject any attempt to move the item to a different directory. It is purely meant to `rename` the item in the same location; so ___Change___ is ___Rename___ but prevents `move` semantics.

* interface: ___ChangeFS___
* Create: ___NewUniversalFS___

```go
  fS := nef.NewUniversalFS(nef.At{
    Root:      "/Users/marina/dev",
    Overwrite: false,
  })
```

* Commands: ___Change___

> os.Change("_bar/from-file.txt_", "_bar/to-file.txt_")

renames _/Users/marina/dev/bar/from-file.txt_ to _/Users/marina/dev/bar/to-file.txt_. However, the behaviour differs depending on the prior existence of _bar/to-file_ and the value of the overwrite flag.

If the file already exists at the destination and _overwrite_ is _true_, then the existing file is overwritten, otherwise, an invalid file system operation ___NewInvalidBinaryFsOpError___ is returned. This denotes the `from` and `to` path and the name of the operation attempted, in this case _Change_.

---

#### 5.1.9. <a name='CopyFS'></a>✨ Copy FS

... pending

* interface: ___CopyFS___
* Create: ___tbd___

#### 5.1.10. <a name='RemoveFS'></a>✨ Remove FS

A file system interface that can delete files and directories

Comes as part of ___UniversalFS___ only.

* interface: ___RemoveFS___
* Commands: ___Remove___, ___RemoveAll___

##### 💎 Remove

> fS.Remove("bar/baz")

behaves as ___os.Remove___

##### 💎 RemoveAll

> fS.RemoveAll("bar/baz")

behaves as ___os.RemoveAll___

---

#### 5.1.11. <a name='RenameFS'></a>✨ Rename FS

A file system interface that can rename/move files and directories

Comes as part of the ___UniversalFS___ only.

* interface: ___RenameFS___
* Create: ___NewUniversalFS___

```go
  fS := nef.NewUniversalFS(nef.At{
    Root:      "/Users/marina/dev",
    Overwrite: false,
  })
```

* Commands: ___Rename___

> fs.Rename("from.txt", "bar/baz/to.txt")

behaves as ___os.Rename___, will move file from _/Users/marina/dev/from.txt_ to _/Users/marina/dev/bar/baz/to.txt_

> fs.Rename("from.txt", "to.txt")

behaves as ___os.Rename___, will move file from _/Users/marina/dev/from.txt_ to _/Users/marina/dev/to.txt_

📍 _Note_: the ___overwrite___ flag is ignored as it is not required by ___os.Rename___

---

#### 5.1.12. <a name='WriteFileFS'></a>✨ Write File FS

* interface: ___WriteFileFS___
* Create: ___NewWriteFileFS___

```go
  fS := nef.NewWriteFileFS(nef.At{
    Root:      "/Users/marina/dev",
    Overwrite: false,
  })
```

* Commands: ___Create___, ___WriteFile___

##### 💎 Create

> fS.Create("bar/baz")

behaves as ___os.Create___

##### 💎 WriteFile

> fS.WriteFile("bar/baz/file.txt")

behaves as ___os.WriteFile___

---

#### 5.1.13. <a name='WriterFS'></a>✨ Writer FS

* interface: ___WriterFS___
* Create: ___NewWriterFS___

```go
  fS := nef.NewNewWriterFS(nef.At{
    Root:      "/Users/marina/dev",
    Overwrite: false,
  })
```

* Composed of: ___CopyFS___, ___ExistsInFS___, ___MakeDirFS___, ___RemoveFS___, ___RenameFS___, ___WriteFileFS___

---

## 6. <a name='OverwriteFlag'></a>Overwrite Flag

The reader may have observed the presence of the overwrite flag at the construction site, being passed into the NewXxxFS functions and may have wondered why the flag is not passed into the command. This would be a valid observation, but it has been done this way in order to conform to the apis in the standard library. The ___overwrite___ flag is purely of the making of ___Nefilim___ and the only way to express it, is to pass it in at the time of creating the file system. This means that the client has to make an upfront decision as to what `overwrite` semantics are required, which is less than desirable, but necessary to avoid incompatibility with the standard packages.

## 7. <a name='Errors'></a>💔 Errors

As ___nefilim___ strives to conform to the standard library, commands contained within return the same errors as so defined, eg ___os.LinkError___, ___os.ErrExist___ and ___os.ErrNotExist___ to name but a few.

However, the custom commands, namely, ___Move___ and ___Change___ and to some extent, ___Rename___ can return new ___Nefilim___ defined errors, as described in the following sections.

The errors use idiomatic Go techniques for adding context to source errors by wrapping and also provided are convenience methods for identifying errors that typically invoke ___errors.Is/As___ on the client's behalf.

### 7.1. <a name='BinaryFsOpError'></a>⛔ Binary Fs Op Error

___IsBinaryFsOpError___ identifies an error that occurs as a result of a failed invoke of a command that take 2 parameters, typically `from` and `to` locations, representing either files or directories. The error also denotes the name of the command to which it relates.

### 7.2. <a name='InvalidPathError'></a>⛔ Invalid Path Error

___IsInvalidPathError___ identifies an error that occurs whenever a path fails validation using ___fs.ValidPath___.

### 7.3. <a name='RejectSameDirectoryMoveError'></a>⛔ Reject Same Directory Move Error

___IsRejectSameDirMoveError___ identifies an error that occurs as a result of a ___Move___ attempt to move an item to the same directory.

### 7.4. <a name='RejectDifferentDirectoryChangeError'></a>⛔ Reject Different Directory Change Error

___IsRejectDifferentDirChangeError___ an error that occurs as a result of a ___Change___ attempt to move an item to a different directory.

(not yet available)

## 8. <a name='Utilities'></a>Utilities

### 8.1. <a name='EnsureAtPath'></a>🛡️ EnsureAtPath

EnsurePathAt ensures that the specified path exists (including any non existing intermediate directories). Given a path and a default filename, the specified path is created in the following manner:

* If the path denotes a file (path does not end is a directory separator), then
the parent folder is created if it doesn't exist on the file-system provided.
* If the path denotes a directory, then that directory is created.

The returned string represents the file, so if the path specified was a directory path, then the defaultFilename provided is joined to the path and returned, otherwise the original path is returned un-modified.

Note: ___filepath.Join___ does not preserve a trailing separator, therefore to make sure a path is interpreted as a directory and not a file, then the separator has to be appended manually onto the end of the path. If vfs is not provided, then the path is ensured directly on the native file system.

[illustrative examples pending]

### 8.2. <a name='ResolvePath'></a>🛡️ResolvePath

ResolvePath performs 2 forms of path resolution. The first is resolving a home path reference, via the ~ character; ~ is replaced by the user's home path. The second resolves ./ or ../ relative path. (The overrides do not need to be provided.)

[illustrative examples pending]

## 9. <a name='TroubleShooting'></a>💥 Trouble Shooting

tbd...
