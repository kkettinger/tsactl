package main

func containsBinary(s string) bool {
	for _, r := range s {
		if (r < 32 || r > 126) && r != '\n' && r != '\r' && r != '\t' {
			return true
		}
	}
	return false
}
