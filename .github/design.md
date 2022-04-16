# Product design

## Idea

Copy an interface definition to your clipboard, and run `mokku` to have your clipboard replaced with a mock implementation of the copied interface.

This implementation is heavily inspired by [Mat Ryer's `moq`](https://github.com/matryer/moq).
The implementation is a struct with function fields for each mocked method.
The struct's actual method will check for the presence of a mock implementation and panic if it is missing with a helpful message.
The name of the implementation for an interface called `InterfaceName` be `InterfaceNameMock`.

## Key features

- External tool, no imports required! (could theoretically be installed with a 3rd party package manager, but since this is a dev tool for gophers, it's pretty safe to that `go` is installed, and I don't see a use-case for this).
- Mock can be easily changed after it's been created (just run against the new interface).

## Wishes

- Nice README
  - Logo
  - Badges
    - build status
    - (go awesome)
    - go report
    - coverage
    - latest version
    - license
  - Asciicinema demo
  - Installation instructions
  - Link to moq
  - Contribution guidelines
- Initial release managed with a GitHub project
- Looks good on `pkg.go.dev`
- High test coverage
- Included in awesome go
- Using GitHub actions for CI

## Shortcomings

### Only Mac

For simplicity, (and in order to have no dependencies other than the standard library) only mac will be supported initially.
Platform differences will are exiled to the packages in the `cmd/` subdirectories, where they can be implemented when the time comes.

#### Action

Ended up swallowing my pride an importing https://github.com/atotto/clipboard to do the hard work of using the clipboard for me.

### Only support for basic interfaces

Not supported directly:

- Embedded interfaces
- Type aliases
- Defined types (e.g. if `A` is a struct, and `B` is defined as `type B A`, then creating a mock against `B` won't work)

# Engineering design

## Nomenclature

- `Mock(data []byte) ([]byte, error)` - the core top-level func of the `github.com/kinbiko/mokku` library. TODO: Consider renaming this func before the first version to something that indicates that this is a shortcut around creating a type with a config etc. which is perhaps the natural progression of this package.
- `mokku` - the command line tool for turning interfaces into mock implementations
- `targetInterface` - A struct representation of the interface to create a mock implementation for.

## Structure

The codebase will be structured as a top-level library (tested), with a `cmd/` subdirectory to hold main functions.
It's the `main()` func's responsibility to read input (initially just from the clipboard) and write (again, to the clipboard initially).
All code lives in the top directory.
Golden files live in `testdata/`.

## Components

### Main

Main is responsible for knowing its own version, and calling the library with the data found in the system clipboard.

- Read flags from the command line
  - `-version` display the version
  - `-debug` display internal log output
  - `-help` display usage information
- Read from the clipboard
- Pass contents of clipboard to library
- If err is returned, print error message is user is to blame, along with usage.
  - If error is internal, then return error message alone with instructions to create an issue on the github repo.
- Pass contents of returned data to the clipboard.

### Mock

- Calls a `parse` function that validates and tries to parse the data given as a `targetInterface` type.
- Return error if this failed. In this case it's probably an issue with what the user has provided.
- Calls a `asBytes` function that takes the given `targetInterface` and parses a template with the given data.
- Return the result of `asBytes`. If there was an error here, then it's probably an issue with our implementation.

### parse

Might be able to use [`parser#ParseExpr`](https://golang.org/pkg/go/parser/#ParseExpr) in order to create a `targetInterface` (defined below).

Use the [Go spec definition of what an interface may look like](https://golang.org/ref/spec#Interface_types) and base the parser off of this:

- Trim whitespace
- If name is missing, just call it 'Mock'
- Expect 'type' to be the first token
  - if so, then store `name = secondToken`
  - if `type` is not the first token, strip everything until th e `interface` token.
    - if it doesn't exist, return an error.
- Then look for {}
- If missing, return an error
- Then parse the contents within the {}.
  - Expect some string before (). Store this as the method name
  - Then store everything until the next string outside ()s as the parameters.
  - Do this for all methods on the interface.
    Return a struct that looks something like this, with the extracted information:

```go
type targetInterface struct {
  Name string
  Methods []method
}

type method struct {
  Name string
  Params string // copied verbatim
}
```

### templating

Generally follow whatever's done in moq, but try and run gofmt against the output code. (don't believe there will be any differences between gofmt and goimports for this project)

## Moq analysis

This document are some notes around the implementation and design of [`moq`, the tool that this project draws most of its inspiration from](https://github.com/matryer/moq).

## Differences

- Moq is targeted for use with `go:generate`.
- Has an option for `gofmt` vs `goimports`, as it will generate entire files (and thus the import order may differ)
- Is slightly confusing as it requires you to use packages/directories
- Stores the data in the calls by default. I think I'd like this to be an opt-in by the developer, who has to choose to store these by modifying the struct manually.
- Should include a 'feature' list (e.g. "does not rely on a cleanup function", "No 'DO NOT EDIT' warnings -- modification is encouraged!", "not a dependency in your app like gomock").
- Moq handles safety concerns around concurrency itself, because it stores each call to each method. I'm on the fence on whether to include this as a feature or not. Implementation details are not hard.
- Moq has to think about build tags, as it is file based.
- Moq has to think about ensuring that acronyms are uppercased in camelcase names as it uses a Mock prefix, instead of suffix. (WRONG ASSUMPTION!)
- Moq does not use go mod and doesn't declare its dependencies

## Good ideas worth stealing (lovingly)

Note: just because these are good ideas, it doesn't mean that they are usable in mokku.

- Create a top-level unnamed (`_`) variable with type of the interface that you want, an initialise it to an empty struct of the mock type. This would not compile if `mokku` failed.
- ~No external dependencies outside of the standard library~ this isn't true, but wasn't obvious due to the lack of a `go.mod` file or a `vendor/` directory
- Golden files for testing
- Might have to think about line endings in the copied text. Moq standardises on `LF`
- Use a text template for the desired output! (In hindsight this is obvious)

## Notes for eng design

- Although at least initially the mokku won't take any command line arguments, it should still support:

  - `-debug` for verbose output
  - `-help` for help
  - `-version` for version number

- Define error messages when given invalid input.
- Definitely design keeping in mind that there's quite probably value in having multiple `cmd` implementations.
- We should give a warning if the interface we're given is empty, and there's no point in writing a mock. Notably, the copied interface still remains in the clipboard, and the user should be made aware of this.
- Will probably struggle to support type aliases and named types if we don't utilise an AST. I think I can live with this weakness.
- Looks like an initial task will be to create a parser of an interface.

## Questions

- [ ] Q: Why does moq need to know about all the abbreviations that will flag as lint issues if there's an inconsistency when AFAICT it uses a suffixes, and not a prefixes in its naming convention.
  - A: I'm guessing it has to do with going above and beyond to not cause a lint failure in the generated code.

## Implementation details:

### Main

Start by reading user provided flags.

```go
type userFlags struct {
	outFile   string
	pkgName   string
	formatter string
	args      []string
}
```

Then pass these to a `run` func that might return an err. Exit with 1 if it returns an err, in addition to printing the error message and listing the usage.

### Run

Run starts off by validating the input.
Then, if no output file is given it will write to stdout.
Extract positional arguments for srcDir and other arguments.
Create a new instance of a `moq` by passing in a config with the srcDir, package name and formatter.
Call `Mock` on this `moq` type.
If there's nothing to write then return.
Create the directory if it's missing, and write to the specified file.

### Constructor

Extract some source package info using a third party library: `golang.org/x/tools/go/packages`.
Then it identifies the complete package path of the output package, defaulting to the source package if absent. Lots of package manipulation logic here that probably doesn't apply to mokku.
Then it creates a new template with `tmpl, err := template.New("moq").Funcs(templateFuncs).Parse(moqTemplate)`. Here, `templateFuncs` includes logic (named `Exports`) to turn lowercased words like http into HTTP. This is used for piping inside templates.
Sets up a default formatter and returns.

### Mock

Validates that at least one interface is defined.
Creates an initial `doc`.
The previous call to `golang.org/x/tools/go/packages` extracted all the types for that package.
... I think I've got what I need.
