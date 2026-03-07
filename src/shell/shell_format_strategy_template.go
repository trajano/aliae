package shell

import "fmt"

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
