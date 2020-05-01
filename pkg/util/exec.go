package util

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
	"text/template"
)

// TemplateExec executes a command as described by the template string and the
// data map.
func TemplateExec(tmpl string, data map[string]string) ([]byte, error) {
	t := template.Must(template.New("").Parse(tmpl))

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, err
	}

	replaced := strings.ReplaceAll(buf.String(), "\n", " ")
	parts := strings.Split(replaced, " ")
	log.Println("Running:", parts)
	return exec.Command(parts[0], parts[1:]...).CombinedOutput()
}
