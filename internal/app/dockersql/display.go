package dockersql

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func DisplayResults(rows *sql.Rows) error {
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	var (
		w      = tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)
		num    = len(columns)
		format = fmt.Sprintf("%s\n", getNumTabs(num))

		values = make([]interface{}, num)
		strs   = make([]string, num)
	)

	for i := 0; i < num; i++ {
		values[i] = &strs[i]
	}

	fmt.Fprintf(w, format, toInterface(columns)...)

	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return err
		}
		fmt.Fprintf(w, format, toInterface(strs)...)
	}
	return w.Flush()
}

func getNumTabs(num int) string {
	tabs := []string{}
	for i := 0; i < num; i++ {
		tabs = append(tabs, "%s\t")
	}
	return strings.Join(tabs, "")
}

func toInterface(strs []string) []interface{} {
	out := []interface{}{}
	for _, s := range strs {
		out = append(out, s)
	}
	return out
}
