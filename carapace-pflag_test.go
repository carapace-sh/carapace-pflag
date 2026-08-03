package pflag

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestLongShorthand(t *testing.T) {
	f := NewFlagSet("longShorthand", ContinueOnError)
	f.BoolP("boola", "a", false, "bool value")
	f.BoolP("boolb", "ab", false, "bool2 value")
	f.BoolP("boolc", "c", false, "bool value")
	f.StringP("stringa", "s", "0", "string value")
	f.StringP("stringx", "sx", "0", "string value")
	f.Lookup("stringx").NoOptDefVal = "1"
	args := []string{
		"-ab",
		"-sx=something",
	}
	want := []string{
		"boolb", "true",
		"stringx", "something",
	}
	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name)
		if len(value) > 0 {
			got = append(got, value)
		}
		return nil
	}
	if err := f.ParseAll(args, store); err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	if !f.Parsed() {
		t.Errorf("f.Parse() = false after Parse")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("f.TestLongShorthand() fail to restore the args")
		t.Errorf("Got:  %v", got)
		t.Errorf("Want: %v", want)
	}
}

func TestNonPosix(t *testing.T) {
	f := NewFlagSet("nonPosix", ContinueOnError)
	f.StringN("stringa", "sa", "0", "string value")
	f.StringN("stringx", "sx", "0", "string value")
	f.Lookup("stringx").NoOptDefVal = "1"
	args := []string{
		"-sa", "somearg",
		"-stringx=something",
	}
	want := []string{
		"stringa", "somearg",
		"stringx", "something",
	}
	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name)
		if len(value) > 0 {
			got = append(got, value)
		}
		return nil
	}
	if err := f.ParseAll(args, store); err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	if !f.Parsed() {
		t.Errorf("f.Parse() = false after Parse")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("f.TestLongShorthand() fail to restore the args")
		t.Errorf("Got:  %v", got)
		t.Errorf("Want: %v", want)
	}
}

func TestOptargDelimiter(t *testing.T) {
	f := NewFlagSet("optargdelimiter", ContinueOnError)
	f.StringN("stringa", "a", "0", "string value")
	f.StringN("stringx", "x", "0", "string value")
	f.Lookup("stringa").NoOptDefVal = "1"
	f.Lookup("stringa").OptargDelimiter = '/'
	f.Lookup("stringx").NoOptDefVal = "2"
	f.Lookup("stringx").OptargDelimiter = ':'

	args := []string{
		"-stringa/somearg",
		"-stringx:something",
	}
	want := []string{
		"stringa", "somearg",
		"stringx", "something",
	}
	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name)
		if len(value) > 0 {
			got = append(got, value)
		}
		return nil
	}
	if err := f.ParseAll(args, store); err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	if !f.Parsed() {
		t.Errorf("f.Parse() = false after Parse")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("f.TestLongShorthand() fail to restore the args")
		t.Errorf("Got:  %v", got)
		t.Errorf("Want: %v", want)
	}
}

func TestNargs(t *testing.T) {
	f := NewFlagSet("nargs", ContinueOnError)
	f.StringSlice("stringa", []string{}, "string value")
	f.StringSlice("stringx", []string{}, "string value")
	f.Lookup("stringa").Nargs = 2
	f.Lookup("stringx").Nargs = -1

	args := []string{
		"--stringa", "one", "two", "three",
		"--stringx", "four", "five",
	}
	want := []string{
		"stringa", "one,two",
		"stringx", "four,five",
	}
	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name)
		if len(value) > 0 {
			got = append(got, value)
		}
		return nil
	}
	if err := f.ParseAll(args, store); err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	if !f.Parsed() {
		t.Errorf("f.Parse() = false after Parse")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("f.TestLongShorthand() fail to restore the args")
		t.Errorf("Got:  %v", got)
		t.Errorf("Want: %v", want)
	}

	// ensure slice is correctly set
	f.Parse(args)
	got, err := f.GetStringSlice("stringa")
	if err != nil {
		t.Error(err.Error())
	}
	want = []string{"one", "two"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got:  %v", got)
		t.Errorf("Want: %v", want)
	}

}

func TestShorthandNameOverlap(t *testing.T) {
	f := NewFlagSet("shorthandNameOverlap", ContinueOnError)
	f.StringN("overlapping", "o", "", "overlapping flag")

	args := []string{
		"-overlapping", "value",
	}
	want := []string{
		"overlapping", "value",
	}
	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name)
		if len(value) > 0 {
			got = append(got, value)
		}
		return nil
	}
	if err := f.ParseAll(args, store); err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("f.TestShorthandNameOverlap() fail to restore the args")
		t.Errorf("Got:  %v", got)
		t.Errorf("Want: %v", want)
	}

	// ensure args are correctly parsed
	f.Parse(args)
	got2, err := f.GetString("overlapping")
	if err != nil {
		t.Error(err.Error())
	}
	want2 := "value"
	if got2 != want2 {
		t.Errorf("Got:  %v", got2)
		t.Errorf("Want: %v", want2)
	}
}

func TestArgumentStyle(t *testing.T) {
	tests := []struct {
		name        string
		style       ArgumentStyle
		args        []string
		wantErr     bool
		wantValue   string
		description string
	}{
		// Zero value (default) accepts all
		{
			name:        "zero accepts delimiter",
			style:       0,
			args:        []string{"--stringa=value"},
			wantErr:     false,
			wantValue:   "value",
			description: "delimiter style should work",
		},
		{
			name:        "zero accepts next arg",
			style:       0,
			args:        []string{"--stringa", "value"},
			wantErr:     false,
			wantValue:   "value",
			description: "next arg style should work",
		},
		{
			name:        "zero accepts posix attached",
			style:       0,
			args:        []string{"-svalue"},
			wantErr:     false,
			wantValue:   "value",
			description: "posix attached style should work",
		},
		// AcceptNext only
		{
			name:        "AcceptNext accepts next arg",
			style:       AcceptNext,
			args:        []string{"--stringa", "value"},
			wantErr:     false,
			wantValue:   "value",
			description: "next arg style should work",
		},
		{
			name:        "AcceptNext rejects delimiter",
			style:       AcceptNext,
			args:        []string{"--stringa=value"},
			wantErr:     true,
			description: "delimiter style should fail",
		},
		// AcceptDelimited only
		{
			name:        "AcceptDelimited accepts delimiter",
			style:       AcceptDelimited,
			args:        []string{"--stringa=value"},
			wantErr:     false,
			wantValue:   "value",
			description: "delimiter style should work",
		},
		{
			name:        "AcceptDelimited rejects next arg",
			style:       AcceptDelimited,
			args:        []string{"--stringa", "value"},
			wantErr:     true,
			description: "next arg style should fail",
		},
		// AcceptAttached only
		{
			name:        "AcceptAttached accepts posix attached",
			style:       AcceptAttached,
			args:        []string{"-svalue"},
			wantErr:     false,
			wantValue:   "value",
			description: "posix attached style should work",
		},
		{
			name:        "AcceptAttached rejects next arg",
			style:       AcceptAttached,
			args:        []string{"-s", "value"},
			wantErr:     true,
			description: "next arg style should fail",
		},
		// Combined AcceptDelimited | AcceptNext
		{
			name:        "DelimitedOrNext accepts delimiter",
			style:       AcceptDelimited | AcceptNext,
			args:        []string{"--stringa=value"},
			wantErr:     false,
			wantValue:   "value",
			description: "delimiter style should work",
		},
		{
			name:        "DelimitedOrNext accepts next arg",
			style:       AcceptDelimited | AcceptNext,
			args:        []string{"--stringa", "value"},
			wantErr:     false,
			wantValue:   "value",
			description: "next arg style should work",
		},
		{
			name:        "DelimitedOrNext rejects attached",
			style:       AcceptDelimited | AcceptNext,
			args:        []string{"-svalue"},
			wantErr:     true,
			description: "attached style should fail",
		},
		// Combined AcceptDelimited | AcceptAttached
		{
			name:        "DelimitedOrAttached accepts delimiter",
			style:       AcceptDelimited | AcceptAttached,
			args:        []string{"--stringa=value"},
			wantErr:     false,
			wantValue:   "value",
			description: "delimiter style should work",
		},
		{
			name:        "DelimitedOrAttached accepts attached",
			style:       AcceptDelimited | AcceptAttached,
			args:        []string{"-svalue"},
			wantErr:     false,
			wantValue:   "value",
			description: "attached style should work",
		},
		{
			name:        "DelimitedOrAttached rejects next arg",
			style:       AcceptDelimited | AcceptAttached,
			args:        []string{"--stringa", "value"},
			wantErr:     true,
			description: "next arg style should fail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFlagSet("argumentstyle", ContinueOnError)
			f.StringP("stringa", "s", "default", "string value")
			f.Lookup("stringa").ArgumentStyle = tt.style

			got := []string{}
			store := func(flag *Flag, value string) error {
				got = append(got, flag.Name)
				if len(value) > 0 {
					got = append(got, value)
				}
				return nil
			}
			err := f.ParseAll(tt.args, store)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for %s, got nil", tt.description)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for %s: %s", tt.description, err)
				return
			}
			if !f.Parsed() {
				t.Errorf("f.Parse() = false after Parse")
			}
			if tt.wantValue != "" && len(got) >= 2 && got[1] != tt.wantValue {
				t.Errorf("got value %q, want %q", got[1], tt.wantValue)
			}
		})
	}
}

// TestArgumentStyleNoOptDefValLong tests the interaction between NoOptDefVal
// and ArgumentStyle on long flags. NoOptDefVal uses the AcceptsDelimited code
// path, so a bare --flag only works when AcceptDelimited is set (or style is 0).
func TestArgumentStyleNoOptDefValLong(t *testing.T) {
	tests := []struct {
		name    string
		style   ArgumentStyle
		args    []string
		wantErr bool
		wantVal string
	}{
		{
			name:    "zero accepts bare bool",
			style:   0,
			args:    []string{"--verbose"},
			wantErr: false,
			wantVal: "true",
		},
		{
			name:    "AcceptDelimited accepts bare bool",
			style:   AcceptDelimited,
			args:    []string{"--verbose"},
			wantErr: false,
			wantVal: "true",
		},
		{
			name:    "AcceptNext rejects bare bool",
			style:   AcceptNext,
			args:    []string{"--verbose"},
			wantErr: true,
		},
		{
			name:    "AcceptAttached rejects bare bool",
			style:   AcceptAttached,
			args:    []string{"--verbose"},
			wantErr: true,
		},
		{
			name:    "AcceptNext|AcceptAttached rejects bare bool",
			style:   AcceptNext | AcceptAttached,
			args:    []string{"--verbose"},
			wantErr: true,
		},
		{
			name:    "AcceptNext|AcceptDelimited accepts bare bool",
			style:   AcceptNext | AcceptDelimited,
			args:    []string{"--verbose"},
			wantErr: false,
			wantVal: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFlagSet("test", ContinueOnError)
			f.BoolP("verbose", "v", false, "verbose")
			f.Lookup("verbose").ArgumentStyle = tt.style

			got := []string{}
			store := func(flag *Flag, value string) error {
				got = append(got, flag.Name)
				if len(value) > 0 {
					got = append(got, value)
				}
				return nil
			}
			err := f.ParseAll(tt.args, store)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				var vre *ValueRequiredError
				if !errors.As(err, &vre) {
					t.Errorf("expected ValueRequiredError, got %T: %v", err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantVal != "" && len(got) >= 2 && got[1] != tt.wantVal {
				t.Errorf("got value %q, want %q", got[1], tt.wantVal)
			}
		})
	}
}

// TestArgumentStyleNoOptDefValShort tests the interaction between NoOptDefVal
// and ArgumentStyle on shorthand flags. The bare -v form uses AcceptsDelimited.
func TestArgumentStyleNoOptDefValShort(t *testing.T) {
	tests := []struct {
		name    string
		style   ArgumentStyle
		args    []string
		wantErr bool
		wantVal string
	}{
		{
			name:    "zero accepts bare -v",
			style:   0,
			args:    []string{"-v"},
			wantErr: false,
			wantVal: "true",
		},
		{
			name:    "AcceptDelimited accepts bare -v",
			style:   AcceptDelimited,
			args:    []string{"-v"},
			wantErr: false,
			wantVal: "true",
		},
		{
			name:    "AcceptNext rejects bare -v",
			style:   AcceptNext,
			args:    []string{"-v"},
			wantErr: true,
		},
		{
			name:    "AcceptAttached rejects bare -v",
			style:   AcceptAttached,
			args:    []string{"-v"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFlagSet("test", ContinueOnError)
			f.BoolP("verbose", "v", false, "verbose")
			f.Lookup("verbose").ArgumentStyle = tt.style

			got := []string{}
			store := func(flag *Flag, value string) error {
				got = append(got, flag.Name)
				if len(value) > 0 {
					got = append(got, value)
				}
				return nil
			}
			err := f.ParseAll(tt.args, store)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				var vre *ValueRequiredError
				if !errors.As(err, &vre) {
					t.Errorf("expected ValueRequiredError, got %T: %v", err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantVal != "" && len(got) >= 2 && got[1] != tt.wantVal {
				t.Errorf("got value %q, want %q", got[1], tt.wantVal)
			}
		})
	}
}

// TestArgumentStyleAllBitsSet tests that explicitly setting all three bits
// behaves identically to the zero value (accept all).
func TestArgumentStyleAllBitsSet(t *testing.T) {
	allBits := AcceptNext | AcceptDelimited | AcceptAttached

	tests := []struct {
		name      string
		args      []string
		wantValue string
	}{
		{"delimiter", []string{"--flag=value"}, "value"},
		{"next", []string{"--flag", "value"}, "value"},
		{"attached", []string{"-fvalue"}, "value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFlagSet("test", ContinueOnError)
			f.StringP("flag", "f", "default", "flag")
			f.Lookup("flag").ArgumentStyle = allBits

			got := []string{}
			store := func(flag *Flag, value string) error {
				got = append(got, flag.Name)
				if len(value) > 0 {
					got = append(got, value)
				}
				return nil
			}
			if err := f.ParseAll(tt.args, store); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) < 2 || got[1] != tt.wantValue {
				t.Errorf("got %v, want value %q", got, tt.wantValue)
			}
		})
	}
}

// TestNargsGreedyConsumesNothing tests that greedy Nargs (<0) consumes zero
// args when the next arg starts with '-' (assumed to be a flag).
func TestNargsGreedyConsumesNothing(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.StringSlice("vals", []string{}, "values")
	f.StringP("next", "n", "default", "next")
	f.Lookup("vals").Nargs = -1

	args := []string{"--vals", "--next=val"}
	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name)
		if len(value) > 0 {
			got = append(got, value)
		}
		return nil
	}
	if err := f.ParseAll(args, store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"vals", "next", "val"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestNargsGreedyStopsAtDashDash tests that greedy Nargs stops at '--'
// (which starts with '-') and that '--' acts as the positional terminator.
func TestNargsGreedyStopsAtDashDash(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.StringSlice("vals", []string{}, "values")
	f.Lookup("vals").Nargs = -1

	args := []string{"--vals", "a", "b", "--", "positional"}
	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name)
		if len(value) > 0 {
			got = append(got, value)
		}
		return nil
	}
	if err := f.ParseAll(args, store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// greedy Nargs should have consumed "a", "b" (stopping at '--')
	want := []string{"vals", "a,b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
	// '--' should be the terminator, so ArgsLenAtDash = 0 (no positionals before it)
	if f.ArgsLenAtDash() != 0 {
		t.Errorf("ArgsLenAtDash = %d, want 0", f.ArgsLenAtDash())
	}
	// 'positional' should be in f.Args()
	if !reflect.DeepEqual(f.Args(), []string{"positional"}) {
		t.Errorf("Args = %v, want [positional]", f.Args())
	}
}

// TestNargsGreaterThanAvailable tests that Nargs > available args consumes
// only what's available (no panic, no index-out-of-bounds).
func TestNargsGreaterThanAvailable(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.StringSlice("vals", []string{}, "values")
	f.Lookup("vals").Nargs = 5

	args := []string{"--vals", "a", "b"}
	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name)
		if len(value) > 0 {
			got = append(got, value)
		}
		return nil
	}
	if err := f.ParseAll(args, store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"vals", "a,b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestNargsZeroEqualsOne tests that Nargs=0 and Nargs=1 behave identically:
// both consume exactly one argument.
func TestNargsZeroEqualsOne(t *testing.T) {
	for _, nargs := range []int{0, 1} {
		t.Run("nargs_"+itoa(nargs), func(t *testing.T) {
			f := NewFlagSet("test", ContinueOnError)
			f.String("flag", "default", "flag")
			f.Lookup("flag").Nargs = nargs

			args := []string{"--flag", "val"}
			got := []string{}
			store := func(flag *Flag, value string) error {
				got = append(got, flag.Name)
				if len(value) > 0 {
					got = append(got, value)
				}
				return nil
			}
			if err := f.ParseAll(args, store); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			want := []string{"flag", "val"}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	return "1"
}

// TestOptargDelimiterWithAcceptNextOnly tests that a custom delimiter is
// rejected when only AcceptNext is set, and accepted when AcceptDelimited is set.
func TestOptargDelimiterWithAcceptNextOnly(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.StringP("flag", "f", "default", "flag")
	f.Lookup("flag").OptargDelimiter = ':'
	f.Lookup("flag").ArgumentStyle = AcceptNext

	// --flag:val should be rejected (AcceptDelimited not set)
	err := f.ParseAll([]string{"--flag:val"}, func(flag *Flag, value string) error { return nil })
	if err == nil {
		t.Fatal("expected error for delimited form with AcceptNext only, got nil")
	}
	var vre *ValueRequiredError
	if !errors.As(err, &vre) {
		t.Errorf("expected ValueRequiredError, got %T: %v", err, err)
	}
}

// TestOptargDelimiterWithAcceptDelimitedOnly tests that the next-arg form
// is rejected when only AcceptDelimited is set, even with a custom delimiter.
func TestOptargDelimiterWithAcceptDelimitedOnly(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.StringP("flag", "f", "default", "flag")
	f.Lookup("flag").OptargDelimiter = ':'
	f.Lookup("flag").ArgumentStyle = AcceptDelimited

	// --flag val should be rejected (AcceptNext not set)
	err := f.ParseAll([]string{"--flag", "val"}, func(flag *Flag, value string) error { return nil })
	if err == nil {
		t.Fatal("expected error for next-arg form with AcceptDelimited only, got nil")
	}
	var vre *ValueRequiredError
	if !errors.As(err, &vre) {
		t.Errorf("expected ValueRequiredError, got %T: %v", err, err)
	}

	// --flag:val should work
	f2 := NewFlagSet("test2", ContinueOnError)
	f2.StringP("flag", "f", "default", "flag")
	f2.Lookup("flag").OptargDelimiter = ':'
	f2.Lookup("flag").ArgumentStyle = AcceptDelimited
	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name, value)
		return nil
	}
	if err := f2.ParseAll([]string{"--flag:val"}, store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"flag", "val"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestShorthandOnlyIgnoresLongForm tests that --name is silently skipped
// for ShorthandOnly flags (treated like an unknown flag, no error).
func TestShorthandOnlyIgnoresLongForm(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.StringS("flag", "f", "default", "flag")

	// --flag=val is silently skipped for ShorthandOnly flags
	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name, value)
		return nil
	}
	if err := f.ParseAll([]string{"--flag=val"}, store); err != nil {
		t.Fatalf("expected no error for --flag with ShorthandOnly flag, got: %v", err)
	}
	// flag should not have been set (value stays at default)
	if len(got) != 0 {
		t.Errorf("flag should not have been set, got %v", got)
	}
	val, _ := f.GetString("flag")
	if val != "default" {
		t.Errorf("flag value = %q, want %q (unchanged default)", val, "default")
	}
}

// TestShorthandOnlyWithArgumentStyle tests that -f shorthand works with
// ArgumentStyle restrictions on ShorthandOnly flags.
func TestShorthandOnlyWithArgumentStyle(t *testing.T) {
	tests := []struct {
		name    string
		style   ArgumentStyle
		args    []string
		wantErr bool
		wantVal string
	}{
		{
			name:    "AcceptDelimited accepts -f=val",
			style:   AcceptDelimited,
			args:    []string{"-f=val"},
			wantErr: false,
			wantVal: "val",
		},
		{
			name:    "AcceptNext accepts -f val",
			style:   AcceptNext,
			args:    []string{"-f", "val"},
			wantErr: false,
			wantVal: "val",
		},
		{
			name:    "AcceptDelimited rejects -f val",
			style:   AcceptDelimited,
			args:    []string{"-f", "val"},
			wantErr: true,
		},
		{
			name:    "AcceptNext rejects -f=val",
			style:   AcceptNext,
			args:    []string{"-f=val"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFlagSet("test", ContinueOnError)
			f.StringS("flag", "f", "default", "flag")
			f.Lookup("flag").ArgumentStyle = tt.style

			got := []string{}
			store := func(flag *Flag, value string) error {
				got = append(got, flag.Name)
				if len(value) > 0 {
					got = append(got, value)
				}
				return nil
			}
			err := f.ParseAll(tt.args, store)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				var vre *ValueRequiredError
				if !errors.As(err, &vre) {
					t.Errorf("expected ValueRequiredError, got %T: %v", err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantVal != "" && len(got) >= 2 && got[1] != tt.wantVal {
				t.Errorf("got value %q, want %q", got[1], tt.wantVal)
			}
		})
	}
}

// TestNameAsShorthandWithArgumentStyle tests ArgumentStyle restrictions on
// NameAsShorthand flags (non-POSIX single-dash longhand).
func TestNameAsShorthandWithArgumentStyle(t *testing.T) {
	tests := []struct {
		name    string
		style   ArgumentStyle
		args    []string
		wantErr bool
		wantVal string
	}{
		{
			name:    "AcceptDelimited accepts -flag=val",
			style:   AcceptDelimited,
			args:    []string{"-flag=val"},
			wantErr: false,
			wantVal: "val",
		},
		{
			name:    "AcceptNext accepts -flag val",
			style:   AcceptNext,
			args:    []string{"-flag", "val"},
			wantErr: false,
			wantVal: "val",
		},
		{
			name:    "AcceptDelimited rejects -flag val",
			style:   AcceptDelimited,
			args:    []string{"-flag", "val"},
			wantErr: true,
		},
		{
			name:    "AcceptNext rejects -flag=val",
			style:   AcceptNext,
			args:    []string{"-flag=val"},
			wantErr: true,
		},
		{
			name:    "AcceptAttached rejects -flagval (non-posix)",
			style:   AcceptAttached,
			args:    []string{"-flagval"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFlagSet("test", ContinueOnError)
			f.StringN("flag", "f", "default", "flag")
			f.Lookup("flag").ArgumentStyle = tt.style

			got := []string{}
			store := func(flag *Flag, value string) error {
				got = append(got, flag.Name)
				if len(value) > 0 {
					got = append(got, value)
				}
				return nil
			}
			err := f.ParseAll(tt.args, store)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				var vre *ValueRequiredError
				if !errors.As(err, &vre) {
					// could be NotExistError for attached in non-posix
					var nee *NotExistError
					if !errors.As(err, &nee) {
						t.Errorf("expected ValueRequiredError or NotExistError, got %T: %v", err, err)
					}
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantVal != "" && len(got) >= 2 && got[1] != tt.wantVal {
				t.Errorf("got value %q, want %q", got[1], tt.wantVal)
			}
		})
	}
}

// TestAcceptsHelperMethods tests the Accepts* helper methods directly,
// including the zero-value-accepts-all semantic.
func TestAcceptsHelperMethods(t *testing.T) {
	// Zero value accepts all
	var zero ArgumentStyle
	if !zero.AcceptsNext() {
		t.Error("zero value should accept Next")
	}
	if !zero.AcceptsDelimited() {
		t.Error("zero value should accept Delimited")
	}
	if !zero.AcceptsAttached() {
		t.Error("zero value should accept Attached")
	}

	// Individual bits
	if !AcceptNext.AcceptsNext() {
		t.Error("AcceptNext should accept Next")
	}
	if AcceptNext.AcceptsDelimited() {
		t.Error("AcceptNext should not accept Delimited")
	}
	if AcceptNext.AcceptsAttached() {
		t.Error("AcceptNext should not accept Attached")
	}

	if AcceptDelimited.AcceptsNext() {
		t.Error("AcceptDelimited should not accept Next")
	}
	if !AcceptDelimited.AcceptsDelimited() {
		t.Error("AcceptDelimited should accept Delimited")
	}
	if AcceptDelimited.AcceptsAttached() {
		t.Error("AcceptDelimited should not accept Attached")
	}

	if AcceptAttached.AcceptsNext() {
		t.Error("AcceptAttached should not accept Next")
	}
	if AcceptAttached.AcceptsDelimited() {
		t.Error("AcceptAttached should not accept Delimited")
	}
	if !AcceptAttached.AcceptsAttached() {
		t.Error("AcceptAttached should accept Attached")
	}

	// Combined
	combined := AcceptNext | AcceptDelimited
	if !combined.AcceptsNext() {
		t.Error("AcceptNext|AcceptDelimited should accept Next")
	}
	if !combined.AcceptsDelimited() {
		t.Error("AcceptNext|AcceptDelimited should accept Delimited")
	}
	if combined.AcceptsAttached() {
		t.Error("AcceptNext|AcceptDelimited should not accept Attached")
	}
}

// TestArgumentStyleErrorType tests that rejection produces the correct
// error type (ValueRequiredError) with the correct flag reference.
func TestArgumentStyleErrorType(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.StringP("flag", "f", "default", "flag")
	f.Lookup("flag").ArgumentStyle = AcceptAttached

	err := f.ParseAll([]string{"--flag", "val"}, func(flag *Flag, value string) error { return nil })
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var vre *ValueRequiredError
	if !errors.As(err, &vre) {
		t.Fatalf("expected ValueRequiredError, got %T: %v", err, err)
	}
	if vre.GetFlag() == nil {
		t.Error("GetFlag() should return the flag")
	}
	if vre.GetFlag().Name != "flag" {
		t.Errorf("flag name = %q, want %q", vre.GetFlag().Name, "flag")
	}
	if vre.GetSpecifiedName() != "flag" {
		t.Errorf("specified name = %q, want %q", vre.GetSpecifiedName(), "flag")
	}
}

// TestInvalidValueErrorModeRendering tests that InvalidValueError renders
// the flag name correctly for different flag modes.
func TestInvalidValueErrorModeRendering(t *testing.T) {
	tests := []struct {
		name     string
		mode     mode
		setup    func(f *FlagSet)
		args     []string
		wantFrag string
	}{
		{
			name:     "Default mode shows -s, --name",
			mode:     Default,
			setup:    func(f *FlagSet) { f.IntP("flag", "f", 0, "flag") },
			args:     []string{"--flag=abc"},
			wantFrag: "-f, --flag",
		},
		{
			name:     "ShorthandOnly shows -s only",
			mode:     ShorthandOnly,
			setup:    func(f *FlagSet) { f.IntS("flag", "f", 0, "flag") },
			args:     []string{"-f=abc"},
			wantFrag: "-f",
		},
		{
			name:     "NameAsShorthand shows -s, -name",
			mode:     NameAsShorthand,
			setup:    func(f *FlagSet) { f.IntN("flag", "f", 0, "flag") },
			args:     []string{"-flag=abc"},
			wantFrag: "-f, -flag",
		},
		{
			name:     "Default no shorthand shows --name",
			mode:     Default,
			setup:    func(f *FlagSet) { f.Int("flag", 0, "flag") },
			args:     []string{"--flag=abc"},
			wantFrag: "--flag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFlagSet("test", ContinueOnError)
			tt.setup(f)

			err := f.ParseAll(tt.args, func(flag *Flag, value string) error {
				return &InvalidValueError{
					flag:  flag,
					value: value,
					cause: errors.New("strconv error"),
				}
			})
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			var ive *InvalidValueError
			if !errors.As(err, &ive) {
				t.Fatalf("expected InvalidValueError, got %T: %v", err, err)
			}
			if !stringsContains(err.Error(), tt.wantFrag) {
				t.Errorf("error %q should contain %q", err.Error(), tt.wantFrag)
			}
		})
	}
}

// TestNotExistErrorNonPosixShorthand tests the non-POSIX shorthand error
// message format (flagUnknownShorthandFlagMessageNonPosix).
func TestNotExistErrorNonPosixShorthand(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.StringN("known", "k", "default", "known flag")

	err := f.ParseAll([]string{"-unknown"}, func(flag *Flag, value string) error { return nil })
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var nee *NotExistError
	if !errors.As(err, &nee) {
		t.Fatalf("expected NotExistError, got %T: %v", err, err)
	}
	if nee.GetSpecifiedName() != "unknown" {
		t.Errorf("specified name = %q, want %q", nee.GetSpecifiedName(), "unknown")
	}
	// Non-POSIX error should not include the chain format
	if stringsContains(err.Error(), "in -") {
		t.Errorf("non-POSIX error should not contain chain format, got: %q", err.Error())
	}
}

// TestParseErrorsAllowlistWithArgumentStyle tests that unknown flags are
// silently skipped even when ArgumentStyle is set on known flags.
func TestParseErrorsAllowlistWithArgumentStyle(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.StringP("flag", "f", "default", "flag")
	f.Lookup("flag").ArgumentStyle = AcceptDelimited
	f.ParseErrorsAllowlist = ParseErrorsAllowlist{UnknownFlags: true}

	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name, value)
		return nil
	}
	args := []string{"--unknown=val", "--flag=val"}
	if err := f.ParseAll(args, store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"flag", "val"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestCountFlagWithArgumentStyle tests the interaction between count flags
// (which have NoOptDefVal="+1") and ArgumentStyle restrictions.
func TestCountFlagWithArgumentStyle(t *testing.T) {
	tests := []struct {
		name    string
		style   ArgumentStyle
		args    []string
		wantErr bool
		wantVal string
	}{
		{
			name:    "zero accepts bare -c",
			style:   0,
			args:    []string{"-c", "-c", "-c"},
			wantErr: false,
			wantVal: "+1",
		},
		{
			name:    "AcceptDelimited accepts bare -c",
			style:   AcceptDelimited,
			args:    []string{"-c"},
			wantErr: false,
			wantVal: "+1",
		},
		{
			name:    "AcceptNext rejects bare -c",
			style:   AcceptNext,
			args:    []string{"-c"},
			wantErr: true,
		},
		{
			name:    "AcceptDelimited accepts -c=3",
			style:   AcceptDelimited,
			args:    []string{"-c=3"},
			wantErr: false,
			wantVal: "3",
		},
		{
			name:    "AcceptNext rejects -c=3",
			style:   AcceptNext,
			args:    []string{"-c=3"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFlagSet("test", ContinueOnError)
			f.CountP("count", "c", "count")
			f.Lookup("count").ArgumentStyle = tt.style

			got := []string{}
			store := func(flag *Flag, value string) error {
				got = append(got, flag.Name, value)
				return nil
			}
			err := f.ParseAll(tt.args, store)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantVal != "" && len(got) >= 2 && got[1] != tt.wantVal {
				t.Errorf("got value %q, want %q", got[1], tt.wantVal)
			}
		})
	}
}

// TestBoolShorthandChainingWithArgumentStyle tests that bool shorthand
// chaining (-abc = -a -b -c) still works with ArgumentStyle when the
// zero value is used, and that NoOptDefVal is respected.
func TestBoolShorthandChainingWithArgumentStyle(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.BoolP("alpha", "a", false, "alpha")
	f.BoolP("beta", "b", false, "beta")
	f.BoolP("gamma", "c", false, "gamma")
	// zero ArgumentStyle = accept all, so bare bools should work via NoOptDefVal

	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name, value)
		return nil
	}
	if err := f.ParseAll([]string{"-abc"}, store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"alpha", "true", "beta", "true", "gamma", "true"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestArgumentStyleShorthandAttachedNonPosix tests that in non-POSIX mode
// (when IsPosix() is false), the AcceptAttached form is never used because
// the attached-arg code path requires f.IsPosix() to be true.
func TestArgumentStyleShorthandAttachedNonPosix(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	// NameAsShorthand makes IsPosix() return false
	f.StringN("flag", "f", "default", "flag")
	f.Lookup("flag").ArgumentStyle = AcceptAttached

	// -flagval in non-posix mode: "flagval" is looked up as a single shorthand
	// and since AcceptAttached is set but IsPosix() is false, the attached path
	// (which requires IsPosix()) is not taken. This should produce an error.
	err := f.ParseAll([]string{"-flagval"}, func(flag *Flag, value string) error { return nil })
	if err == nil {
		t.Fatal("expected error for attached form in non-posix mode, got nil")
	}
}

// TestOptargDelimiterDisabledNonPosix tests that an optarg flag with
// OptargDelimiter set to a control character (e.g. -1) in non-POSIX mode
// parses directly attached values: -rvalue -> flag "r" with value "value".
func TestOptargDelimiterDisabledNonPosix(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		wantVal string
	}{
		{
			name:    "attached value",
			args:    []string{"-rvalue"},
			wantErr: false,
			wantVal: "value",
		},
		{
			name:    "bare flag (no value, optarg)",
			args:    []string{"-r"},
			wantErr: false,
		},
		{
			name:    "next arg style",
			args:    []string{"-r", "value"},
			wantErr: false,
			// optarg: -r alone uses NoOptDefVal, "value" is a positional, not the flag value
		},
		{
			name:    "delimiter not recognized",
			args:    []string{"-r=value"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFlagSet("test", ContinueOnError)
			// StringS with same name and shorthand makes it ShorthandOnly;
			// multi-char shorthand makes IsPosix() return false
			f.StringS("r", "r", "", "recurse")
			f.Lookup("r").NoOptDefVal = " "
			f.Lookup("r").OptargDelimiter = DelimiterDisabled

			got := []string{}
			store := func(flag *Flag, value string) error {
				got = append(got, flag.Name, value)
				return nil
			}
			err := f.ParseAll(tt.args, store)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantVal != "" && len(got) >= 2 && got[1] != tt.wantVal {
				t.Errorf("got value %q, want %q", got[1], tt.wantVal)
			}
		})
	}
}

// TestInvalidSyntaxErrorForDashDashEquals tests that --=value produces
// an InvalidSyntaxError (name starts with '=').
func TestInvalidSyntaxErrorForDashDashEquals(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.String("flag", "default", "flag")

	err := f.ParseAll([]string{"--=value"}, func(flag *Flag, value string) error { return nil })
	if err == nil {
		t.Fatal("expected error for --=value, got nil")
	}
	var ise *InvalidSyntaxError
	if !errors.As(err, &ise) {
		t.Fatalf("expected InvalidSyntaxError, got %T: %v", err, err)
	}
	if ise.GetSpecifiedFlag() != "--=value" {
		t.Errorf("specified flag = %q, want %q", ise.GetSpecifiedFlag(), "--=value")
	}
}

// TestInvalidSyntaxErrorForDashDashDash tests that -- followed by another
// dash at the start of name produces InvalidSyntaxError.
func TestInvalidSyntaxErrorForDashDashDash(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.String("flag", "default", "flag")

	err := f.ParseAll([]string{"----flag"}, func(flag *Flag, value string) error { return nil })
	if err == nil {
		t.Fatal("expected error for ----flag, got nil")
	}
	var ise *InvalidSyntaxError
	if !errors.As(err, &ise) {
		t.Fatalf("expected InvalidSyntaxError, got %T: %v", err, err)
	}
}

// stringsContains is a helper to avoid importing strings just for one check.
func stringsContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func TestCustomPrefix(t *testing.T) {
	tests := []struct {
		name   string
		prefix rune
		setup  func(f *FlagSet)
		args   []string
		want   []string
	}{
		{
			name:   "ampersand long flag with next arg",
			prefix: '&',
			setup:  func(f *FlagSet) { f.IntP("flag", "f", 0, "flag") },
			args:   []string{"&&flag", "42"},
			want:   []string{"flag", "42"},
		},
		{
			name:   "ampersand long flag with delimiter",
			prefix: '&',
			setup:  func(f *FlagSet) { f.IntP("flag", "f", 0, "flag") },
			args:   []string{"&&flag=42"},
			want:   []string{"flag", "42"},
		},
		{
			name:   "ampersand shorthand with next arg",
			prefix: '&',
			setup:  func(f *FlagSet) { f.IntP("flag", "f", 0, "flag") },
			args:   []string{"&f", "42"},
			want:   []string{"flag", "42"},
		},
		{
			name:   "ampersand shorthand with delimiter",
			prefix: '&',
			setup:  func(f *FlagSet) { f.IntP("flag", "f", 0, "flag") },
			args:   []string{"&f=42"},
			want:   []string{"flag", "42"},
		},
		{
			name:   "ampersand NameAsShorthand single prefix",
			prefix: '&',
			setup:  func(f *FlagSet) { f.IntN("flag", "f", 0, "flag") },
			args:   []string{"&flag=42"},
			want:   []string{"flag", "42"},
		},
		{
			name:   "ampersand ShorthandOnly single prefix",
			prefix: '&',
			setup:  func(f *FlagSet) { f.IntS("flag", "f", 0, "flag") },
			args:   []string{"&f=42"},
			want:   []string{"flag", "42"},
		},
		{
			name:   "ampersand bool flag",
			prefix: '&',
			setup:  func(f *FlagSet) { f.BoolP("verbose", "v", false, "verbose") },
			args:   []string{"&v"},
			want:   []string{"verbose", "true"},
		},
		{
			name:   "ampersand terminator",
			prefix: '&',
			setup:  func(f *FlagSet) { f.IntP("flag", "f", 0, "flag") },
			args:   []string{"&f=1", "&&", "&f", "2"},
			want:   []string{"flag", "1"},
		},
		{
			name:   "ampersand positional args",
			prefix: '&',
			setup:  func(f *FlagSet) { f.IntP("flag", "f", 0, "flag") },
			args:   []string{"positional", "&f=42", "more"},
			want:   []string{"flag", "42"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFlagSet("test", ContinueOnError)
			f.SetPrefix(tt.prefix)
			tt.setup(f)

			got := []string{}
			store := func(flag *Flag, value string) error {
				got = append(got, flag.Name)
				if len(value) > 0 {
					got = append(got, value)
				}
				return nil
			}

			err := f.ParseAll(tt.args, store)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomPrefixErrors(t *testing.T) {
	// Unknown flag with ampersand prefix
	f := NewFlagSet("test", ContinueOnError)
	f.SetPrefix('&')
	f.IntP("flag", "f", 0, "flag")

	err := f.Parse([]string{"&&unknown=42"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
	var nee *NotExistError
	if !errors.As(err, &nee) {
		t.Fatalf("expected NotExistError, got %T: %v", err, err)
	}
	if !strings.Contains(err.Error(), "&&unknown") {
		t.Errorf("error %q should contain '&&unknown'", err.Error())
	}

	// Unknown shorthand with ampersand prefix
	err = f.Parse([]string{"&x=42"})
	if err == nil {
		t.Fatal("expected error for unknown shorthand")
	}
	if !strings.Contains(err.Error(), "&x") {
		t.Errorf("error %q should contain '&x'", err.Error())
	}
}

func TestCustomPrefixUsage(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.SetPrefix('&')
	f.IntP("flag", "f", 0, "a flag")
	f.Int("other", 0, "another flag")

	usage := f.FlagUsages()
	if !strings.Contains(usage, "&f") {
		t.Errorf("usage should contain '&f', got:\n%s", usage)
	}
	if !strings.Contains(usage, "&&flag") {
		t.Errorf("usage should contain '&&flag', got:\n%s", usage)
	}
	if !strings.Contains(usage, "&&other") {
		t.Errorf("usage should contain '&&other', got:\n%s", usage)
	}
}

func TestCustomPrefixValueRequiredError(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.SetPrefix('&')
	f.IntP("flag", "f", 0, "flag")

	err := f.Parse([]string{"&f"})
	if err == nil {
		t.Fatal("expected error for missing argument")
	}
	var vre *ValueRequiredError
	if !errors.As(err, &vre) {
		t.Fatalf("expected ValueRequiredError, got %T: %v", err, err)
	}
	if !strings.Contains(err.Error(), "&f") {
		t.Errorf("error %q should contain '&f'", err.Error())
	}
}

func TestCustomPrefixInvalidValueError(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.SetPrefix('&')
	f.IntP("flag", "f", 0, "flag")

	err := f.ParseAll([]string{"&f=abc"}, func(flag *Flag, value string) error {
		return &InvalidValueError{
			flag:   flag,
			value:  value,
			cause:  errors.New("strconv error"),
			prefix: f.Prefix(),
		}
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var ive *InvalidValueError
	if !errors.As(err, &ive) {
		t.Fatalf("expected InvalidValueError, got %T: %v", err, err)
	}
	if !strings.Contains(err.Error(), "&f") {
		t.Errorf("error %q should contain '&f'", err.Error())
	}
}

func TestDefaultPrefixIsDash(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	if f.Prefix() != '-' {
		t.Errorf("expected default prefix '-', got %q", f.Prefix())
	}
}

func TestCustomPrefixGreedyNargs(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.SetPrefix('&')
	var vals []string
	f.StringSliceVarP(&vals, "vals", "v", []string{}, "values")
	f.Lookup("vals").Nargs = -1
	f.IntP("other", "o", 0, "other flag")

	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name)
		if len(value) > 0 {
			got = append(got, value)
		}
		return nil
	}

	err := f.ParseAll([]string{"&&vals", "a", "b", "&&other=3"}, store)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"vals", "a,b", "other", "3"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCustomPrefixOptargDelimiter(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.SetPrefix('&')
	f.StringP("flag", "f", "0", "flag")
	f.Lookup("flag").OptargDelimiter = ':'

	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name, value)
		return nil
	}

	err := f.ParseAll([]string{"&f:val"}, store)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"flag", "val"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	got = nil
	err = f.ParseAll([]string{"&&flag:val"}, store)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want = []string{"flag", "val"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCustomPrefixCount(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want int
	}{
		{name: "bare shorthand", args: []string{"&c"}, want: 1},
		{name: "chained shorthand", args: []string{"&ccc"}, want: 3},
		{name: "bare long", args: []string{"&&count"}, want: 1},
		{name: "delimited long", args: []string{"&&count=5"}, want: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFlagSet("test", ContinueOnError)
			f.SetPrefix('&')
			var c int
			f.CountVarP(&c, "count", "c", "count")

			err := f.Parse(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if c != tt.want {
				t.Errorf("got %d, want %d", c, tt.want)
			}
		})
	}
}

func TestCustomPrefixHelp(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.SetPrefix('&')

	err := f.Parse([]string{"&&help"})
	if err != ErrHelp {
		t.Fatalf("expected ErrHelp, got %v", err)
	}

	err = f.Parse([]string{"&h"})
	if err != ErrHelp {
		t.Fatalf("expected ErrHelp, got %v", err)
	}
}

func TestCustomPrefixInvalidSyntax(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.SetPrefix('&')
	f.IntP("flag", "f", 0, "flag")

	err := f.Parse([]string{"&&&flag=42"})
	if err == nil {
		t.Fatal("expected error for triple prefix")
	}
	var ise *InvalidSyntaxError
	if !errors.As(err, &ise) {
		t.Fatalf("expected InvalidSyntaxError, got %T: %v", err, err)
	}
}

func TestCustomPrefixLonePrefixIsPositional(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.SetPrefix('&')
	f.IntP("flag", "f", 0, "flag")

	err := f.Parse([]string{"&"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Args()) != 1 || f.Args()[0] != "&" {
		t.Errorf("expected positional ['&'], got %v", f.Args())
	}
}

func TestCustomPrefixBoolChaining(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.SetPrefix('&')
	f.BoolP("a", "a", false, "a")
	f.BoolP("b", "b", false, "b")
	f.BoolP("c", "c", false, "c")

	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name, value)
		return nil
	}

	err := f.ParseAll([]string{"&abc"}, store)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"a", "true", "b", "true", "c", "true"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCustomPrefixInterspersedFalse(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.SetPrefix('&')
	f.SetInterspersed(false)
	f.IntP("flag", "f", 0, "flag")

	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name, value)
		return nil
	}

	err := f.ParseAll([]string{"&f=1", "positional", "&f=2"}, store)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"flag", "1"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
	if len(f.Args()) != 2 {
		t.Errorf("expected 2 positional args, got %d: %v", len(f.Args()), f.Args())
	}
}

func TestCustomPrefixArgumentStyleReject(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.SetPrefix('&')
	f.IntP("flag", "f", 0, "flag")
	f.Lookup("flag").ArgumentStyle = AcceptNext // only next-arg form, reject delimited

	err := f.Parse([]string{"&f=42"})
	if err == nil {
		t.Fatal("expected error for rejected delimited form")
	}
	var vre *ValueRequiredError
	if !errors.As(err, &vre) {
		t.Fatalf("expected ValueRequiredError, got %T: %v", err, err)
	}

	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name, value)
		return nil
	}
	err = f.ParseAll([]string{"&f", "42"}, store)
	if err != nil {
		t.Fatalf("unexpected error for next-arg form: %v", err)
	}
	want := []string{"flag", "42"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
