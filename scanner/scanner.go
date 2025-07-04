package scanner

type ScanOptions struct {
	IncludeLockfiles bool
	IncludeVendor    bool
	Verbosity        int
}

type DepMap = map[string][]string
type AllDeps = map[string]DepMap
type LanguageScanner func(projectPath string, opts ScanOptions) (map[string][]string, error)

func mergeDeps(dest, src DepMap) DepMap {
	if dest == nil {
		dest = make(DepMap)
	}
	for pkg, paths := range src {
		dest[pkg] = append(dest[pkg], paths...)
	}
	return dest
}

// Main scan entrypoint
func ScanDependencies(path string, opts ScanOptions) (AllDeps, error) {
	results := make(AllDeps)

	if pDeps, err := ScanPython(path, opts.IncludeLockfiles, opts.IncludeVendor, opts.Verbosity); err == nil {
		results["pypi"] = mergeDeps(results["pypi"], pDeps)
	}

	if nDeps, err := ScanJavaScript(path, opts.IncludeLockfiles, opts.IncludeVendor, opts.Verbosity); err == nil {
		results["npm"] = mergeDeps(results["npm"], nDeps)
	}

	if gDeps, err := ScanGo(path, opts.IncludeLockfiles, opts.IncludeVendor, opts.Verbosity); err == nil {
		results["go"] = mergeDeps(results["go"], gDeps)
	}

	if phDeps, err := ScanPHP(path, opts.IncludeLockfiles, opts.IncludeVendor, opts.Verbosity); err == nil {
		results["php"] = mergeDeps(results["php"], phDeps)
	}

	if rDeps, err := ScanRuby(path, opts.IncludeLockfiles, opts.IncludeVendor, opts.Verbosity); err == nil {
		results["ruby"] = mergeDeps(results["ruby"], rDeps)
	}

	if ruDeps, err := ScanRust(path, opts.IncludeLockfiles, opts.IncludeVendor, opts.Verbosity); err == nil {
		results["rust"] = mergeDeps(results["rust"], ruDeps)
	}

	return results, nil
}
