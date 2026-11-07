package pflag

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestEdgeCustomPrefixDisabledDelimiterEquals(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.SetPrefix('&')
	f.StringS("r", "r", "", "recurse")
	f.Lookup("r").NoOptDefVal = " "
	f.Lookup("r").OptargDelimiter = DelimiterDisabled

	got := []string{}
	store := func(flag *Flag, value string) error { got = append(got, flag.Name, value); return nil }
	if err := f.ParseAll([]string{"&rkey=value"}, store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"r", "key=value"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestEdgeZeroFlagSetPrefix(t *testing.T) {
	var f FlagSet
	if f.Prefix() != '-' {
		t.Errorf("expected default prefix '-', got %q", f.Prefix())
	}
}

func TestEdgeCustomPrefixShorthandOnly(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.SetPrefix('&')
	f.StringS("flag", "f", "default", "flag")

	got := []string{}
	store := func(flag *Flag, value string) error { got = append(got, flag.Name, value); return nil }
	if err := f.ParseAll([]string{"&f=val"}, store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"flag", "val"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	// long form is ignored for ShorthandOnly
	got = got[:0]
	if err := f.ParseAll([]string{"&&flag=val"}, store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected no sets, got %v", got)
	}
}

func TestEdgeCustomPrefixNameAsShorthandNargs(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.SetPrefix('&')
	f.StringSliceN("vals", "v", []string{}, "values")
	f.Lookup("vals").Nargs = -1
	f.IntP("other", "o", 0, "other")

	got := []string{}
	store := func(flag *Flag, value string) error { got = append(got, flag.Name, value); return nil }
	if err := f.ParseAll([]string{"&&vals", "a", "b", "&&other=3"}, store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"vals", "a,b", "other", "3"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestEdgeCustomPrefixInvalidSyntax(t *testing.T) {
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
	if !strings.Contains(err.Error(), "&&&flag") {
		t.Errorf("error %q should contain '&&&flag'", err.Error())
	}
}

func TestEdgeCustomPrefixNotExistError(t *testing.T) {
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
}

func TestEdgeArgumentStyleDisabledDelimiterAttachedRejected(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.StringS("r", "r", "", "recurse")
	f.Lookup("r").NoOptDefVal = " "
	f.Lookup("r").OptargDelimiter = DelimiterDisabled
	f.Lookup("r").ArgumentStyle = AcceptDelimited | AcceptNext // no AcceptAttached

	got := []string{}
	store := func(flag *Flag, value string) error { got = append(got, flag.Name, value); return nil }
	err := f.ParseAll([]string{"-rvalue"}, store)
	// attached should not be consumed; either error or bare flag used
	if err == nil {
		if len(got) >= 2 && got[1] == "value" {
			t.Errorf("attached value parsed despite AcceptsAttached=false")
		}
	}
}

func TestEdgeTopLevelPrefix(t *testing.T) {
	// ensure top-level SetPrefix/Parsed exist (compile-time)
	SetPrefix('-')
	if CommandLine.Prefix() != '-' {
		t.Errorf("expected '-', got %q", CommandLine.Prefix())
	}
}
