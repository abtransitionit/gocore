package kubectl2

import (
	"fmt"
	"strings"

	"github.com/abelgacem/lucg/luccore"
)

func PlayQueryKubectl(kubectlQuery string, opts KubectlOptions) (result string, customErr string, err error) {
	// Play request
	output, _, cerrSrc, errSrc := luccore.RunCLIOnVM(KubectlHost, kubectlQuery)
	// handle FAILURE
	if errSrc != nil {
		customErr = fmt.Sprintf("❌ Kubectl Command failed > %v > Output > %s", cerrSrc, output)
		return "", customErr, errSrc
	}
	// handle OPTION
	if opts.FormatOutput {
		// Apply formatting if enabled
		output = formatKubectlResultAddIndx(output) // Apply formatting if enabled
	}
	// handle SUCCESS
	return output, "", nil
}

// formats kubectl output by adding line numbers (except header)
func formatKubectlResultAddIndx(output string) string {
	output = strings.TrimSuffix(output, "\n") // Remove trailing newline
	lines := strings.Split(output, "\n")

	if len(lines) == 0 {
		return output // Edge case: empty output
	}

	var formatted strings.Builder
	// formatted.WriteString(lines[0] + "\n")
	formatted.WriteString(fmt.Sprintf("Id   %s\n", lines[0])) // First line (header)

	for i, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			formatted.WriteString("\n")
			continue
		} // skip empty lines but keep them in the output
		if line != "" { // Skip empty lines (optional)
			formatted.WriteString(fmt.Sprintf("%-2d - %s\n", i+1, line))
		}
	}

	return formatted.String()
}
