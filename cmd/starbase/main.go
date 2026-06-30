// Command starbase builds a highly interactive but fully static knowledge base
// from a tree of markdown files. It is designed for agents: every problem it
// finds (dead links, missing template arguments, broken templates) is reported
// precisely so the content can be iterated to completion.
//
// Usage:
//
//	starbase build [flags] <content-dir>   generate the site
//	starbase check [flags] <content-dir>   fast validation only (no rendering)
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/monoptic-io/starbase/internal/build"
	"github.com/monoptic-io/starbase/internal/model"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "build":
		os.Exit(runBuild(os.Args[2:]))
	case "check":
		os.Exit(runCheck(os.Args[2:]))
	case "verify":
		os.Exit(runVerify(os.Args[2:]))
	case "templates":
		os.Exit(runTemplates(os.Args[2:]))
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "starbase: unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `starbase — a static knowledge-base generator for agents

Commands:
  build <content-dir>   Generate the static site (incremental).
  check <content-dir>   Fast validation: report dead links and bad template calls.
  verify <content-dir>  Re-run the evidence/ program and diff every claim's
                        embedded result against the freshly computed output.
  templates [dir]       List available shortcode templates and their arguments.

Run "starbase build -h" or "starbase check -h" for flags.
`)
}

func runTemplates(args []string) int {
	fs := flag.NewFlagSet("templates", flag.ExitOnError)
	fs.Parse(reorder(args))
	cats, err := build.Catalog(build.Config{ContentDir: contentDir(fs), OutDir: "_site"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "starbase: %v\n", err)
		return 1
	}
	for _, c := range cats {
		fmt.Printf("\n{{< %s … >}}", c.Name)
		if c.Doc != "" {
			fmt.Printf("\n  %s", c.Doc)
		}
		if c.Open {
			fmt.Printf("\n  (accepts arbitrary extra arguments)")
		}
		for _, p := range c.Params {
			req := "optional"
			switch {
			case p.Required:
				req = "REQUIRED"
			case p.HasDef:
				req = fmt.Sprintf("default %q", p.Default)
			}
			fmt.Printf("\n    %-12s %s", p.Name, req)
		}
		fmt.Println()
	}
	return 0
}

func runBuild(args []string) int {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	out := fs.String("o", "_site", "output directory")
	title := fs.String("title", "Knowledge Base", "site title")
	baseURL := fs.String("base", "", "base URL (optional)")
	drafts := fs.Bool("drafts", false, "include draft topics")
	force := fs.Bool("force", false, "ignore cache and rebuild everything")
	strict := fs.Bool("strict", false, "exit non-zero on warnings too")
	quiet := fs.Bool("quiet", false, "only print diagnostics and summary")
	vendorAssets := fs.Bool("vendor", false, "download and bundle third-party assets (e.g. KaTeX) locally for a self-contained, offline site instead of linking a CDN")
	offline := fs.Bool("offline", false, "with -vendor, use only previously cached downloads (never hit the network)")
	fs.Parse(reorder(args))

	cfg := build.Config{
		ContentDir: contentDir(fs),
		OutDir:     *out,
		SiteTitle:  *title,
		BaseURL:    *baseURL,
		Drafts:     *drafts,
		Force:      *force,
		Vendor:     *vendorAssets,
		Offline:    *offline,
	}
	res, err := build.Build(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "starbase: %v\n", err)
		return 1
	}
	printDiagnostics(res.Diagnostics)
	if !*quiet {
		fmt.Printf("built %d topics → %s  (%d rendered, %d unchanged)\n",
			res.Topics, *out, res.Rendered, res.Skipped)
	}
	fmt.Printf("%d error(s), %d warning(s)\n", res.Errors(), res.Warnings())
	if res.Errors() > 0 || (*strict && res.Warnings() > 0) {
		return 1
	}
	return 0
}

func runCheck(args []string) int {
	fs := flag.NewFlagSet("check", flag.ExitOnError)
	drafts := fs.Bool("drafts", false, "include draft topics")
	strict := fs.Bool("strict", false, "exit non-zero on warnings too")
	fs.Parse(reorder(args))

	cfg := build.Config{ContentDir: contentDir(fs), Drafts: *drafts, OutDir: "_site"}
	res, err := build.Check(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "starbase: %v\n", err)
		return 1
	}
	printDiagnostics(res.Diagnostics)
	fmt.Printf("checked %d topics: %d error(s), %d warning(s)\n",
		res.Topics, res.Errors(), res.Warnings())
	if res.Errors() > 0 || (*strict && res.Warnings() > 0) {
		return 1
	}
	return 0
}

// reorder moves positional arguments after flags so the content directory can
// appear before or after flags (agents shouldn't have to remember the order).
func reorder(args []string) []string {
	valueFlags := map[string]bool{
		"-o": true, "--o": true, "-title": true, "--title": true,
		"-base": true, "--base": true,
	}
	var flags, pos []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if strings.HasPrefix(a, "-") {
			flags = append(flags, a)
			if !strings.Contains(a, "=") && valueFlags[a] && i+1 < len(args) {
				flags = append(flags, args[i+1])
				i++
			}
			continue
		}
		pos = append(pos, a)
	}
	return append(flags, pos...)
}

func runVerify(args []string) int {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	drafts := fs.Bool("drafts", false, "include draft topics")
	force := fs.Bool("force", false, "re-run all evidence units, ignoring the cache")
	fs.Parse(reorder(args))

	cfg := build.Config{ContentDir: contentDir(fs), Drafts: *drafts, Force: *force, OutDir: "_site"}
	res, err := build.Verify(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "starbase: %v\n", err)
		return 1
	}
	printDiagnostics(res.Diagnostics)
	fmt.Printf("evidence: %d unit(s) run, %d cached\n", res.UnitsRun, res.UnitsCached)
	fmt.Printf("verified %d claim(s), %d attested (not re-run): %d error(s)\n",
		res.Verified, res.Attested, res.Errors())
	if res.Errors() > 0 {
		return 1
	}
	return 0
}

func contentDir(fs *flag.FlagSet) string {
	dir := fs.Arg(0)
	if dir == "" {
		dir = "."
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return dir
	}
	return abs
}

func printDiagnostics(diags []model.Diagnostic) {
	sorted := append([]model.Diagnostic(nil), diags...)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].File != sorted[j].File {
			return sorted[i].File < sorted[j].File
		}
		return sorted[i].Line < sorted[j].Line
	})
	for _, d := range sorted {
		loc := d.File
		if d.Line > 0 {
			loc = fmt.Sprintf("%s:%d", d.File, d.Line)
		}
		fmt.Fprintf(os.Stderr, "%s  %s  %s\n", strings.ToUpper(d.Severity.String()), loc, d.Message)
	}
}
