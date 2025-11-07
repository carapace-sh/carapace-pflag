package pflag

// -- func Value
type funcValue func(string) error

func (f funcValue) Set(s string) error { return f(s) }

func (f funcValue) Type() string { return "func" }

func (f funcValue) String() string { return "" } // same behavior as stdlib 'flag' package

// Func defines a func flag with specified name, callback function and usage string.
//
// The callback function will be called every time "--{name}={value}" (or equivalent) is
// parsed on the command line, with "{value}" as an argument.
func (f *FlagSet) Func(name string, usage string, fn func(string) error) {
	f.FuncP(name, "", usage, fn)
}

// FuncN is like FuncP, but adds the name as shorthand (non-posix).
func (f *FlagSet) FuncN(name string, shorthand string, usage string, fn func(string) error) {
	var val Value = funcValue(fn)
	f.VarN(val, name, shorthand, usage)
}

// FuncP is like Func, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) FuncP(name string, shorthand string, usage string, fn func(string) error) {
	var val Value = funcValue(fn)
	f.VarP(val, name, shorthand, usage)
}

// FuncS is like Func, but accepts a shorthand letter to be used after a single dash, alone.
func (f *FlagSet) FuncS(name string, shorthand string, usage string, fn func(string) error) {
	var val Value = funcValue(fn)
	f.VarS(val, name, shorthand, usage)
}

// Func defines a func flag with specified name, callback function and usage string.
//
// The callback function will be called every time "--{name}={value}" (or equivalent) is
// parsed on the command line, with "{value}" as an argument.
func Func(name string, usage string, fn func(string) error) {
	CommandLine.FuncP(name, "", usage, fn)
}

// FuncN is like FuncP, but adds the name as shorthand (non-posix).
func FuncN(name, shorthand string, usage string, fn func(string) error) {
	CommandLine.FuncN(name, shorthand, usage, fn)
}

// FuncP is like Func, but accepts a shorthand letter that can be used after a single dash.
func FuncP(name, shorthand string, usage string, fn func(string) error) {
	CommandLine.FuncP(name, shorthand, usage, fn)
}

// FuncS is like Func, but accepts a shorthand letter to be used after a single dash, alone.
func FuncS(name, shorthand string, usage string, fn func(string) error) {
	CommandLine.FuncS(name, shorthand, usage, fn)
}
