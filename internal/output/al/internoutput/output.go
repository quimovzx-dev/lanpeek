package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func JSON(w io.Writer, value any) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return e.Encode(value)
}

func Table(w io.Writer, headers []string, rows [][]string) error {
	// Rounded separators and a restrained bold header. Lipgloss automatically
	// strips colors when output is not a terminal.
	t := table.New().Headers(headers...).Rows(rows...).Border(lipgloss.NormalBorder()).
		StyleFunc(func(row, col int) lipgloss.Style {
			s := lipgloss.NewStyle().Padding(0, 1)
			if row == table.HeaderRow {
				return s.Bold(true)
			}
			return s
		})
	_, err := fmt.Fprintln(w, strings.TrimRight(t.String(), "\n"))
	return err
}
