package pflag

import (
	"reflect"
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
