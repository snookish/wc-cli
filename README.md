# wc-cli

A small command-line tool written in Go to count words, lines, and bytes/characters from files or standard input.

## Features

- Count words, lines, and bytes.
- Read from files or STDIN.

## Installation

Build locally:

```sh
make build
```

## Usage

Basic usage:

```sh
wc-cli [flags] [file...]
```

If no files are provided, wc-cli reads from standard input.

Common flags:

- -w Count words
- -l Count lines
- -c Count bytes

Combine flags to display multiple metrics.

## Examples

Count words in a file:

```sh
wc-cli -w file.txt
```

Count lines and bytes for multiple files:

```sh
wc-cli -l -c file1.txt file2.txt
```

Pipe input:

```sh
cat file.txt | wc-cli -w
```
