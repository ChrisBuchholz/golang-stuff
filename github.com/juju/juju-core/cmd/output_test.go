// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package cmd_test

import (
	"launchpad.net/gnuflag"
	. "launchpad.net/gocheck"

	"launchpad.net/juju-core/cmd"
	"launchpad.net/juju-core/testing"
)

// OutputCommand is a command that uses the output.go formatters.
type OutputCommand struct {
	cmd.CommandBase
	out   cmd.Output
	value interface{}
}

func (c *OutputCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "output",
		Args:    "<something>",
		Purpose: "I like to output",
		Doc:     "output",
	}
}

func (c *OutputCommand) SetFlags(f *gnuflag.FlagSet) {
	c.out.AddFlags(f, "smart", cmd.DefaultFormatters)
}

func (c *OutputCommand) Init(args []string) error {
	return cmd.CheckEmpty(args)
}

func (c *OutputCommand) Run(ctx *cmd.Context) error {
	return c.out.Write(ctx, c.value)
}

// use a struct to control field ordering.
var defaultValue = struct {
	Juju   int
	Puppet bool
}{1, false}

var outputTests = map[string][]struct {
	value  interface{}
	output string
}{
	"": {
		{1, "1\n"},
		{-1, "-1\n"},
		{1.1, "1.1\n"},
		{true, "True\n"},
		{false, "False\n"},
		{"hello", "hello\n"},
		{"\n\n\n", "\n\n\n\n"},
		{"foo: bar", "foo: bar\n"},
		{[]string{"blam", "dink"}, "blam\ndink\n"},
		{map[interface{}]interface{}{"foo": "bar"}, "foo: bar\n"},
	},
	"smart": {
		{1, "1\n"},
		{-1, "-1\n"},
		{1.1, "1.1\n"},
		{true, "True\n"},
		{false, "False\n"},
		{"hello", "hello\n"},
		{"\n\n\n", "\n\n\n\n"},
		{"foo: bar", "foo: bar\n"},
		{[]string{"blam", "dink"}, "blam\ndink\n"},
		{map[interface{}]interface{}{"foo": "bar"}, "foo: bar\n"},
	},
	"json": {
		{1, "1\n"},
		{-1, "-1\n"},
		{1.1, "1.1\n"},
		{true, "true\n"},
		{false, "false\n"},
		{"hello", `"hello"` + "\n"},
		{"\n\n\n", `"\n\n\n"` + "\n"},
		{"foo: bar", `"foo: bar"` + "\n"},
		{[]string{}, `[]` + "\n"},
		{[]string{"blam", "dink"}, `["blam","dink"]` + "\n"},
		{defaultValue, `{"Juju":1,"Puppet":false}` + "\n"},
	},
	"yaml": {
		{1, "1\n"},
		{-1, "-1\n"},
		{1.1, "1.1\n"},
		{true, "true\n"},
		{false, "false\n"},
		{"hello", "hello\n"},
		{"\n\n\n", "'\n\n\n\n'\n"},
		{"foo: bar", "'foo: bar'\n"},
		{[]string{"blam", "dink"}, "- blam\n- dink\n"},
		{defaultValue, "juju: 1\npuppet: false\n"},
	},
}

func (s *CmdSuite) TestOutputFormat(c *C) {
	for format, tests := range outputTests {
		c.Logf("format %s", format)
		var args []string
		if format != "" {
			args = []string{"--format", format}
		}
		for i, t := range tests {
			c.Logf("  test %d", i)
			ctx := testing.Context(c)
			result := cmd.Main(&OutputCommand{value: t.value}, ctx, args)
			c.Assert(result, Equals, 0)
			c.Assert(bufferString(ctx.Stdout), Equals, t.output)
			c.Assert(bufferString(ctx.Stderr), Equals, "")
		}
	}

	ctx := testing.Context(c)
	result := cmd.Main(&OutputCommand{}, ctx, []string{"--format", "cuneiform"})
	c.Assert(result, Equals, 2)
	c.Assert(bufferString(ctx.Stdout), Equals, "")
	c.Assert(bufferString(ctx.Stderr), Matches, ".*: unknown format \"cuneiform\"\n")
}

// Py juju allowed both --format json and --format=json. This test verifies that juju is
// being built against a version of the gnuflag library (rev 14 or above) that supports
// this argument format.
// LP #1059921
func (s *CmdSuite) TestFormatAlternativeSyntax(c *C) {
	ctx := testing.Context(c)
	result := cmd.Main(&OutputCommand{}, ctx, []string{"--format=json"})
	c.Assert(result, Equals, 0)
	c.Assert(bufferString(ctx.Stdout), Equals, "null\n")
}
