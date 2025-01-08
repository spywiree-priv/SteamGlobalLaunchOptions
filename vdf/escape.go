package vdf

import "strings"

// https://developer.valvesoftware.com/wiki/KeyValues#About_KeyValues_Text_File_Format
var (
	escaper   = strings.NewReplacer("\n", `\n`, "\t", `\t`, `\\`, "\\", "\"", `\"`)
	unescaper = strings.NewReplacer(`\n`, "\n", `\t`, "\t", `\\`, "\\", `\"`, "\"")
)

func Quote(s string) string {
	return `"` + escaper.Replace(s) + `"`
}

func Unquote(s string) string {
	return unescaper.Replace(s[1 : len(s)-1])
}
