# carapace-pflag

[![PkgGoDev](https://pkg.go.dev/badge/github.com/carapace-sh/carapace-pflag)](https://pkg.go.dev/github.com/carapace-sh/carapace-pflag)
[![GoReportCard](https://goreportcard.com/badge/github.com/carapace-sh/carapace-pflag)](https://goreportcard.com/report/github.com/carapace-sh/carapace-pflag)
[![Coverage Status](https://coveralls.io/repos/github/carapace-sh/carapace-pflag/badge.svg?branch=master)](https://coveralls.io/github/carapace-sh/carapace-pflag?branch=master)

Fork of [spf13/pflag](https://github.com/spf13/pflag) aimed to provide support for non-posix variants in [carapace](https://github.com/carapace-sh/carapace) standalone mode (e.g. [carapace-bin](https://github.com/carapace-sh/carapace-bin)).

## Customizations

### Shorthand-Only

Support shorthand-only flags (e.g. `-h` without `--help`) using `S` suffix for flag functions:

```go
pflag.BoolS("help", "h", false, "show help") // -h
```

### Long Shorthand

Support shorthand flags that are more than one character long (e.g. `pkill -<sig>`).
> This implicitly disables posix shorthand chaining (e.g. `ls -lha`):

```go
pflag.BoolS("STOP", "STOP", false, "Stop process, unblockable") // -STOP
pflag.BoolN("help", "h", false, "show help") // -h, -help
```

### Custom Optarg Delimiter

Support custom optarg delimiter (e.g. `java -agentlib:jdwp`).

```go
rootCmd.Flags().StringS("agentlib", "agentlib", "", "load native agent library")
rootCmd.Flag("agentlib").NoOptDefVal = " "
rootCmd.Flag("agentlib").OptargDelimiter = ':'
```

Set `OptargDelimiter` to `pflag.DelimiterDisabled` to disable delimiter-based
parsing and accept directly-attached values (e.g. `-rvalue` → flag `r`, value
`value`):

```go
rootCmd.Flag("agentlib").OptargDelimiter = pflag.DelimiterDisabled
rootCmd.Flag("agentlib").NoOptDefVal = " "
```

### Argument Style

Control which argument forms a flag accepts using the `ArgumentStyle` bitmask.
The zero value accepts all forms (backward compatible).

```go
rootCmd.Flag("verbose").ArgumentStyle = pflag.AcceptDelimited | pflag.AcceptNext
// accepts: --verbose=val, --verbose val
// rejects: --verboseval (attached)
```

### Custom Flag Prefix

Change the flag prefix character (e.g. `&` for elvish) on a `FlagSet`:

```go
rootCmd.Flags().SetPrefix('&') // &flag, &&flag, && terminates
```

### Nargs

Support for flags consuming multiple arguments.
> Use with `Slice` and `Array` flag types.

```go
rootCmd.Flags().StringSlice("nargs-any", []string{}, "Nargs")
rootCmd.Flags().StringSlice("nargs-two", []string{}, "Nargs")

rootCmd.Flag("nargs-any").Nargs = -1 // consumes at least one argument until one starts with `-`
rootCmd.Flag("nargs-two").Nargs = 2 // consumes exactly 2 arguments
```
