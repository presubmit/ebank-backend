package utils

import "database/sql"

func ContainsString(arr []string, v string) bool {
	for _, x := range arr {
		if x == v {
			return true
		}
	}
	return false
}

func NullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
