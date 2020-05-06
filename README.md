# GoSpin

[![GoDoc](https://godoc.org/github.com/m1/gospin?status.svg)](https://godoc.org/github.com/m1/gospin)
[![Build Status](https://travis-ci.org/m1/gospin.svg?branch=master)](https://travis-ci.org/m1/gospin)
[![Go Report Card](https://goreportcard.com/badge/github.com/m1/gospin)](https://goreportcard.com/report/github.com/m1/gospin)
[![Release](https://img.shields.io/github/release/m1/gospin.svg)](https://github.com/m1/gospin/releases/latest)
[![Coverage Status](https://coveralls.io/repos/github/m1/gospin/badge.svg)](https://coveralls.io/github/m1/gospin)

__Article spinning and spintax/spinning syntax engine written in Go, useful for A/B, testing pieces of text/articles and creating more natural conversations. 
Use as a [library](#usage) or as a [CLI](#cli-usage).__

## Installation

Use go get to get the latest version
```text
go get github.com/m1/gospin
```

Then import it into your projects using the following:
```go
import (
	"github.com/m1/gospin"
)
```

## What is spintax?

Take this example:

```
{hello|hey} world
```

When spinning an article (the above sentence) each of the words/phrases contained within the curly brackets 
are randomly picked and substituted in the sentence. So for example, the above article could be spun to be:
`hey world` and `hello world`

You can also have nested spintax, take this example:

```
{hello|hey} world, {hope {you're|you are} {okay|good}|have a nice day}
```

A few examples of what the above could output:
- `hey world, hope you're okay`
- `hello world, hope you are good`
- `hello world, have a nice day`
- etc...

You can also have optional phrases, just don't specify a word after or before a pipe to make it optional:
```
{hello|hey}{ world|}, how are you today?
```

A few examples of what the above could output:
- `hey, how are you today?`
- `hello world, how are you today?`
- etc...

## Why use spintax?

Spintax can be used for several things. It used to be used a lot for spinning articles for SEO but is 
less useful for that these days. It's more used for A/B testing, testing pieces of text for efficiency/click through rate. 
Also it is used spinning content for users to keep things fresh, i.e home page text or ai/chat bots.

## Usage

To use as a library is pretty simple:

```go
spinner := gospin.New(nil)
simple := "The {slow|quick} {fox|deer} {gracefully |}jumps over the {sleeping|lazy} dog"

spin := spinner.Spin(simple) // The slow fox jumps over the sleeping dog
spins := spinner.SpinN(simple, 10)
// spins = [
// "The slow fox gracefully jumps over the lazy dog"
// "The slow deer jumps over the sleeping dog"
// "The quick fox jumps over the lazy dog"
// ...
// ]
```

You can also configure it to take custom syntax (the package uses Jet format as the default), e.g. if you 
wanted it to set the start and end characters (default are curly brackets) to square brackets:
```go
spinner := gospin.New(&gospin.Config{
        StartChar:     "[",
        EndChar:       "]",
        DelimiterChar: ";",
})
simple := "The [slow;quick] [fox;deer] [gracefully ;]jumps over the [sleeping;lazy] dog"
spin := spinner.Spin(simple) // The slow fox jumps over the sleeping dog
```

### Escaping

To escape, the default character to use is `\\`, e.g:

```
The \{slow|quick\} {fox|deer} {gracefully |}jumps over the {sleeping|lazy} dog
```

Would output something like:
```
The {slow|quick} fox jumps over the sleeping dog
```

You can customize the escape char in the config:
```go
spinner := gospin.New(&gospin.Config{
        EscapeChar:    "@",
})
simple := "The @{slow|quick@} {fox|deer} {gracefully |}jumps over the {sleeping|lazy} dog"
spin := spinner.Spin(simple) // The @{slow|quick@} fox jumps over the sleeping dog
```

### Random seeds

The `spin` by default generates a random seed each spin, to stop this and use your own global rand seed you can 
use the `UseGlobalRand` toggle in the config. This is useful for testing:

```go
spinner := gospin.New(&gospin.Config{
        UseGlobalRand: false,
})
```

## CLI usage
 
GoSpin can also be used on the cli, just install using: `go get github.com/m1/gospin/cmd/gospin`

To use:
```
➜  ~ gospin --help                    
GoSpin is a fast and configurable article spinning and spintax engine written in Go.

Usage:
  gospin [text] [flags]

Flags:
      --delimiter string   Delimiter char (default "|")
      --end string         End char for the spinning engine (default "}")
      --escape string      Escape char (default "\\")
  -h, --help               help for gospin
      --start string       Start char for the spinning engine (default "{")
      --times int          How many articles to generate (default 1)
```

For example: 
```
➜  ~ gospin "{hello|hey} friend"                     
hey friend
```

To spin multiple, use the `times` flag, this is outputted as json for easier parsing:
```
➜  ~ gospin "The {slow|quick} {fox|deer} {gracefully |}jumps over the {sleeping|lazy} dog" --times=5 | jq
[
  "The slow fox gracefully jumps over the sleeping dog",
  "The slow deer jumps over the lazy dog",
  "The quick deer jumps over the sleeping dog",
  "The slow fox gracefully jumps over the sleeping dog",
  "The quick fox jumps over the lazy dog"
]
```