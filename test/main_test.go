package test

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/rogpeppe/go-internal/gotooltest"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	if os.Getenv("TESTSCRIPT_ON") == "" {
		flag.Parse()
		// Don't put the binaries in a temporary directory to delete, as that
		// means we have to re-link them every single time. That's quite
		// expensive, at around half a second per 'go test' invocation.
		binDir, err := filepath.Abs(".cache")
		if err != nil {
			panic(err)
		}
		os.Setenv("GOBIN", binDir)
		os.Setenv("PATH", binDir+string(filepath.ListSeparator)+os.Getenv("PATH"))
		cmd := exec.Command("go", "install", "-ldflags=-w -s",
			"github.com/gunk/gunk",
			"github.com/gunk/docgen",
		)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("%+v\n", os.Environ())

			panic(err)
		}
	}
	os.Exit(testscript.RunMain(m, map[string]func() int{}))
}

func TestScripts(t *testing.T) {
	t.Parallel()
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Fatal(err)
	}

	goCache := filepath.Join(os.TempDir(), "gunk-test-go-cache")

	p := testscript.Params{
		Dir: filepath.Join("testdata", "scripts"),
		Setup: func(e *testscript.Env) error {
			cmd := exec.Command("cp",
				"testdata/prepare-go-sum/go.mod",
				"testdata/prepare-go-sum/go.sum",
				e.WorkDir)
			cmd.Stderr = os.Stderr
			if _, err := cmd.Output(); err != nil {
				return fmt.Errorf("failed to copy go.sum: %w", err)
			}

			e.Vars = append(e.Vars, "GONOSUMDB=*")
			e.Vars = append(e.Vars, "GUNK_CACHE_DIR="+cacheDir)
			e.Vars = append(e.Vars, "TESTSCRIPT_ON=on")

			e.Vars = append(e.Vars, "HOME="+goCache)
			return nil
		},
	}
	if err := gotooltest.Setup(&p); err != nil {
		t.Fatal(err)
	}
	testscript.Run(t, p)
}
