package shell

import "fmt"
import "strings"

func (bashFormatStrategy) FormatSetArg(name string, oneBasedIndex int) string {
	return fmt.Sprintf("%s=$%d", name, oneBasedIndex)
}

func (zshFormatStrategy) FormatSetArg(name string, oneBasedIndex int) string {
	return fmt.Sprintf("%s=$%d", name, oneBasedIndex)
}

func (fishFormatStrategy) FormatSetArg(name string, oneBasedIndex int) string {
	return fmt.Sprintf("set %s $argv[%d]", name, oneBasedIndex)
}

func (pwshFormatStrategy) FormatSetArg(name string, oneBasedIndex int) string {
	return fmt.Sprintf("$%s = $args[%d]", name, oneBasedIndex-1)
}

func (nuFormatStrategy) FormatSetArg(name string, oneBasedIndex int) string {
	return fmt.Sprintf("$%s = $args.%d", name, oneBasedIndex-1)
}

func (xonshFormatStrategy) FormatSetArg(name string, oneBasedIndex int) string {
	return fmt.Sprintf("%s=$argv[%d]", name, oneBasedIndex)
}

func (tcshFormatStrategy) FormatSetArg(name string, oneBasedIndex int) string {
	return fmt.Sprintf("set %s=$%d", name, oneBasedIndex)
}

func (cmdFormatStrategy) FormatSetArg(name string, oneBasedIndex int) string {
	return fmt.Sprintf("set %s=%%%d", name, oneBasedIndex)
}

func (bashFormatStrategy) FormatProgress(state, percentage int) string {
	return fmt.Sprintf("printf '\\033]9;4;%d;%d\\007'", state, percentage)
}

func (zshFormatStrategy) FormatProgress(state, percentage int) string {
	return fmt.Sprintf("printf '\\033]9;4;%d;%d\\007'", state, percentage)
}

func (fishFormatStrategy) FormatProgress(state, percentage int) string {
	return fmt.Sprintf("printf '\\033]9;4;%d;%d\\007'", state, percentage)
}

func (pwshFormatStrategy) FormatProgress(state, percentage int) string {
	return fmt.Sprintf(`[Console]::Out.Write("$([char]27)]9;4;%d;%d$([char]7)")`, state, percentage)
}

func (nuFormatStrategy) FormatProgress(state, percentage int) string {
	return fmt.Sprintf("printf '\\033]9;4;%d;%d\\007'", state, percentage)
}

func (xonshFormatStrategy) FormatProgress(state, percentage int) string {
	return fmt.Sprintf("printf '\\033]9;4;%d;%d\\007'", state, percentage)
}

func (tcshFormatStrategy) FormatProgress(state, percentage int) string {
	return fmt.Sprintf("printf '\\033]9;4;%d;%d\\007'", state, percentage)
}

func (cmdFormatStrategy) FormatProgress(state, percentage int) string {
	return fmt.Sprintf("printf '\\033]9;4;%d;%d\\007'", state, percentage)
}

func (bashFormatStrategy) EscapeString(value string) string {
	return defaultEscapedString(value)
}

func (zshFormatStrategy) EscapeString(value string) string {
	return defaultEscapedString(value)
}

func (fishFormatStrategy) EscapeString(value string) string {
	return defaultEscapedString(value)
}

func (pwshFormatStrategy) EscapeString(value string) string {
	return pwshEscapedString(value)
}

func (nuFormatStrategy) EscapeString(value string) string {
	return defaultEscapedString(value)
}

func (xonshFormatStrategy) EscapeString(value string) string {
	return defaultEscapedString(value)
}

func (tcshFormatStrategy) EscapeString(value string) string {
	return defaultEscapedString(value)
}

func (cmdFormatStrategy) EscapeString(value string) string {
	return defaultEscapedString(value)
}

func (bashFormatStrategy) FormatAliasScriptPrelude() string   { return "" }
func (bashFormatStrategy) FormatAliasScriptPostlude() string  { return "" }
func (zshFormatStrategy) FormatAliasScriptPrelude() string    { return "" }
func (zshFormatStrategy) FormatAliasScriptPostlude() string   { return "" }
func (fishFormatStrategy) FormatAliasScriptPrelude() string   { return "" }
func (fishFormatStrategy) FormatAliasScriptPostlude() string  { return "" }
func (pwshFormatStrategy) FormatAliasScriptPrelude() string   { return "" }
func (pwshFormatStrategy) FormatAliasScriptPostlude() string  { return "" }
func (nuFormatStrategy) FormatAliasScriptPrelude() string     { return "" }
func (nuFormatStrategy) FormatAliasScriptPostlude() string    { return "" }
func (xonshFormatStrategy) FormatAliasScriptPrelude() string  { return "" }
func (xonshFormatStrategy) FormatAliasScriptPostlude() string { return "" }
func (tcshFormatStrategy) FormatAliasScriptPrelude() string   { return "" }
func (tcshFormatStrategy) FormatAliasScriptPostlude() string  { return "" }
func (cmdFormatStrategy) FormatAliasScriptPrelude() string    { return cmdAliasPre() }
func (cmdFormatStrategy) FormatAliasScriptPostlude() string   { return cmdAliasPost() }

func defaultEscapedString(value string) string {
	return strings.NewReplacer(
		`\`, `\\`,
		`"`, `\"`,
	).Replace(value)
}

func pwshEscapedString(value string) string {
	return strings.NewReplacer(
		"`", "``",
		`"`, "`\"",
	).Replace(value)
}
