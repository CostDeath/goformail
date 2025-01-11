# Go project root location relative to execution location
$root_location = "."
# Linter executable location relative to Go project root.
$exec_location = "../../scripts/golangci-lint.exe"

cd $root_location
& $exec_location run --out-format tab
#../../scripts/golangci-lint.exe run --out-format tab