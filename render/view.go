package render

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	api "github.com/sovorem/sovorem/client"
)

const safeStepIcon = "🛡︎"

func renderTestHeader(header string, spinner spinner.Model, isFinished bool, isSubmit bool, passed *bool, noPenaltyOnFail bool) string {
	if noPenaltyOnFail {
		header = fmt.Sprintf("%s %s", header, white.Render(safeStepIcon))
	}
	cmdStr := renderTest(header, spinner.View(), isFinished, &isSubmit, passed)
	box := borderBox.Render(fmt.Sprintf(" %s ", cmdStr))
	sliced := strings.Split(box, "\n")
	sliced[2] = strings.Replace(sliced[2], "─", "┬", 1)
	return strings.Join(sliced, "\n") + "\n"
}

func renderTests(tests []testModel, spinner string) string {
	var str strings.Builder
	var edges strings.Builder

	for _, test := range tests {
		testStr := renderTest(test.text, spinner, test.finished, nil, test.passed)
		testStr = fmt.Sprintf("  %s", testStr)
		height := lipgloss.Height(testStr)

		edges.Reset()
		edges.WriteString(" ├─")
		for i := 1; i < height; i++ {
			edges.WriteString("\n │ ")
		}

		str.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, edges.String(), testStr))
		str.WriteByte('\n')
	}

	return str.String()
}

func renderTest(text string, spinner string, isFinished bool, isSubmit *bool, passed *bool) string {
	testStr := ""
	if !isFinished {
		testStr += fmt.Sprintf("%s %s", spinner, text)
	} else if isSubmit != nil && !*isSubmit {
		testStr += text
	} else if passed == nil {
		testStr += gray.Render(fmt.Sprintf("?  %s", text))
	} else if *passed {
		testStr += green.Render(fmt.Sprintf("✓  %s", text))
	} else {
		testStr += red.Render(fmt.Sprintf("X  %s", text))
	}
	return testStr
}

func renderJqOutputs(outputs []api.CLICommandJqOutput) string {
	if len(outputs) == 0 {
		return ""
	}

	var str strings.Builder
	str.WriteString("\n > jq-ի output-ը.\n\n")
	for _, output := range outputs {
		str.WriteString(gray.Render(fmt.Sprintf("Query: %s", output.Query)))
		str.WriteByte('\n')
		if output.Error != "" {
			str.WriteString(gray.Render(fmt.Sprintf("Error: %s", output.Error)))
			str.WriteByte('\n')
			str.WriteByte('\n')
			continue
		}
		if len(output.Results) == 0 {
			str.WriteString(gray.Render("Արդյունքներ. [չկան]"))
			str.WriteByte('\n')
			str.WriteByte('\n')
			continue
		}
		str.WriteString(gray.Render("Արդյունքներ."))
		str.WriteByte('\n')
		for _, line := range output.Results {
			str.WriteString(gray.Render("  - " + line))
			str.WriteByte('\n')
		}
		str.WriteByte('\n')
	}
	return str.String()
}

func (m rootModel) View() string {
	if m.clear {
		return ""
	}
	s := m.spinner.View()
	var str strings.Builder
	for _, step := range m.steps {
		str.WriteString(renderTestHeader(step.step, m.spinner, step.finished, m.isSubmit, step.passed, step.noPenaltyOnFail))
		str.WriteString(renderTests(step.tests, s))

		if step.sleepAfter != "" && step.finished {
			sleepBox := borderBox.Render(fmt.Sprintf(" %s ", step.sleepAfter))
			str.WriteString(sleepBox)
			str.WriteByte('\n')
		}

		if step.result == nil || !m.finalized {
			continue
		}

		if step.result.CLICommandResult != nil {
			for _, test := range step.tests {
				if strings.Contains(test.text, "exit code") {
					fmt.Fprintf(&str, "\n > Command-ի exit code-ը. %d\n", step.result.CLICommandResult.ExitCode)
					break
				}
			}
			str.WriteString(" > Command-ի stdout-ը.\n\n")
			sliced := strings.SplitSeq(step.result.CLICommandResult.Stdout, "\n")
			i := 0
			runeCount := 0
			const maxLines, maxRunes = 32, 5120
			for s := range sliced {
				if i >= maxLines || runeCount >= maxRunes {
					str.WriteString(gray.Render("... output-ը տեսողականորեն կրճատվել ա, բայց ամբողջական տարբերակը պահվել ա"))
					str.WriteByte('\n')
					break
				}
				runeCount += utf8.RuneCountInString(s)
				str.WriteString(gray.Render(s))
				str.WriteByte('\n')
				i++
			}
			str.WriteString(renderJqOutputs(step.result.CLICommandResult.JqOutputs))
			availableVariables, expectsVariables := availableVariablesForCLIResult(*step.result.CLICommandResult)
			if expectsVariables {
				str.WriteString(renderVariableSection("Հասանելի փոփոխականները", availableVariables))
			}
		}

		if step.result.HTTPRequestResult != nil {
			str.WriteString(printHTTPRequestResult(*step.result.HTTPRequestResult))
		}
	}

	if m.result == api.VerificationResultSlugSuccess && m.isSubmit {
		str.WriteByte('\n')
		str.WriteByte('\n')
		str.WriteString(green.Render("Բոլոր test-երը անցան! 🎉"))
		str.WriteByte('\n')
		if m.xpReward >= 0 {
			str.WriteByte('\n')
			str.WriteString(green.Bold(true).Render(fmt.Sprintf("Ստացար +%d XP", m.xpReward)))
			str.WriteByte('\n')
			for _, item := range m.xpBreakdown {
				if item.XP == 0 {
					continue
				}
				sign := "+"
				xp := item.XP
				if xp < 0 {
					sign = "-"
					xp = -xp
				}
				if item.Percent > 0 {
					str.WriteString(gray.Render(fmt.Sprintf("%s%3d XP (%-4s %s)", sign, xp, fmt.Sprintf("%.0f%%", item.Percent*100), item.Name)))
				} else {
					str.WriteString(gray.Render(fmt.Sprintf("%s%3d XP (%s)", sign, xp, item.Name)))
				}
				str.WriteByte('\n')
			}
		}
		str.WriteByte('\n')
		str.WriteString(green.Render("Վերադարձիր browser՝ հաջորդ դասին անցնելու համար։"))
		str.WriteByte('\n')
		str.WriteByte('\n')
	} else if m.result == api.VerificationResultSlugNoop {
		str.WriteString("\n\nTest-երը չանցան! ❌")
		fmt.Fprintf(&str, "\n\nՉանցած Step-ը. %v", m.failure.FailedStepIndex+1)
		str.WriteString("\nError. ")
		str.WriteString(m.failure.ErrorMessage)
		str.WriteByte('\n')
		str.WriteByte('\n')
		str.WriteString(white.Render(safeStepIcon))
		str.WriteString(" Սա safe step էր։\n")
		str.WriteString("Դեռ չես անցել, բայց Amulet-ի կամ Perfect Run-ի progress-ը չես կորցրել։\n\n")
	} else if m.result == api.VerificationResultSlugFailure {
		str.WriteByte('\n')
		str.WriteByte('\n')
		str.WriteString(red.Render("Test-երը չանցան! ❌"))
		if m.failure != nil {
			if m.failure.FailedStepIndex >= 0 && m.failure.FailedStepIndex < len(m.steps) {
				str.WriteString(red.Render(fmt.Sprintf("\n\nՉանցած Command-ը. %s", m.steps[m.failure.FailedStepIndex].step)))
			}
			str.WriteString(red.Render(fmt.Sprintf("\n\nՉանցած Step-ը. %v", m.failure.FailedStepIndex+1)))
			str.WriteString(red.Render(fmt.Sprintf("\nError. %s", m.failure.ErrorMessage)))
		} else {
			str.WriteString(red.Render("\n\nՉանցած Step-ը. unknown"))
			str.WriteString(red.Render("\nError. unknown"))
		}
		str.WriteByte('\n')
		str.WriteByte('\n')
		currentDate := time.Now().Format("2006-01-02")
		if strings.HasSuffix(currentDate, "04-01") {
			str.WriteString(magenta.Render(fmt.Sprintf("Էս դեպքի մասին զեկուցվել ա քո համակարգային ադմինիստրատորին։ [%s]\n", currentDate)))
		}
	}

	return str.String()
}
