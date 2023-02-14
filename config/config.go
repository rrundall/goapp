package config

import (
	"path"
	"path/filepath"
)

// Define config variables
var DBFile, _ = filepath.Abs(path.Join("pkg/db", "book.db"))
var LogFile, _ = filepath.Abs(path.Join("./", "api.log"))
