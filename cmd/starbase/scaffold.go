package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/monoptic-io/starbase/internal/scaffold"
)

func runVersion(args []string) int {
	fmt.Println(scaffold.Version())
	return 0
}

func runSkills(args []string) int {
	fs := flag.NewFlagSet("skills", flag.ExitOnError)
	out := fs.String("o", scaffold.SkillsDir, "directory to write the skills into")
	force := fs.Bool("force", false, "overwrite skill files you have edited locally")
	fs.Parse(reorder(args))

	dir, _ := filepath.Abs(*out)
	res, err := scaffold.EmitSkills(dir, scaffold.Version(), *force)
	if err != nil {
		fmt.Fprintf(os.Stderr, "starbase: %v\n", err)
		return 1
	}
	printEmit(res)
	return 0
}

func runInit(args []string) int {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	title := fs.String("title", "", "site title")
	force := fs.Bool("force", false, "overwrite skill files you have edited locally")
	fs.Parse(reorder(args))

	dir := fs.Arg(0)
	if dir == "" {
		dir = "."
	}
	abs, _ := filepath.Abs(dir)
	res, err := scaffold.Init(abs, *title, scaffold.Version(), *force)
	if err != nil {
		fmt.Fprintf(os.Stderr, "starbase: %v\n", err)
		return 1
	}
	for _, f := range res.Created {
		fmt.Printf("  create  %s\n", f)
	}
	for _, f := range res.Skipped {
		fmt.Printf("  exists  %s (left as-is)\n", f)
	}
	printEmit(res.Skills)
	fmt.Printf("\nInitialized a starbase KB in %s\n", dir)
	fmt.Printf("Next:\n  starbase check %s\n  starbase build %s -o _site\n", dir, dir)
	return 0
}

// printEmit reports what an EmitSkills/Init call changed under .claude/skills.
func printEmit(res scaffold.EmitResult) {
	for _, f := range res.Written {
		fmt.Printf("  skill   %s (new)\n", f)
	}
	for _, f := range res.Updated {
		fmt.Printf("  skill   %s (updated)\n", f)
	}
	for _, f := range res.Removed {
		fmt.Printf("  skill   %s (removed)\n", f)
	}
	for _, f := range res.Skipped {
		fmt.Printf("  skill   %s (kept your edits — rerun with -force to overwrite)\n", f)
	}
	if !res.Changed() {
		fmt.Printf("  skills up to date (%d files)\n", res.Unchanged)
	}
}

// noteSkillsDrift prints a one-line reminder when a repo's emitted skills were
// written by a different released version than the running binary.
func noteSkillsDrift(startDir string) {
	d := scaffold.DetectDrift(startDir, scaffold.Version())
	if d.Stale {
		fmt.Fprintf(os.Stderr,
			"note: skills in %s were emitted by starbase %s; you're on %s — run `starbase skills` to refresh\n",
			d.Dir, d.Emitted, d.Current)
	}
}
