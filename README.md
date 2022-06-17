# Ropes in Go

A value-oriented, immutable, functional rope in Go.

# Introduction

This module provides an implementation of the [rope data structure](https://en.wikipedia.org/wiki/Rope_(data_structure)) in Go.
Ropes allow for the efficient manipulation of long strings of text.

The rope structure implemented by this module is immutable, value-oriented, and
[fully persistent](https://en.wikipedia.org/wiki/Persistent_data_structure).

There is no way to modify an existing rope:
all operations on a rope return a new rope.

One notable use case that this makes particularly nice is undo/redo:
simply keep the old versions of the rope around.

# Example

A simple example:

```go

	rope := NewString("hello")
	rope = rope.AppendString(", world!")
	fmt.Printf("The value of the rope is: %v\n", rope.String())

```

# Documentation

Documentation is included inline and is also available at https://pkg.go.dev/github.com/deadpixi/rope

# License

rope - a persistent rope data structure
Copyright (C) 2022 Rob King

This library is free software; you can redistribute it and/or
modify it under the terms of the GNU Lesser General Public
License as published by the Free Software Foundation; either
version 2.1 of the License, or (at your option) any later version.

This library is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public
License along with this library; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301
USA
