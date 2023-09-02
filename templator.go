package templator

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template" // TODO: add html support
)

const (
	DefaultDelimLeft  = "{{"
	DefaultDelimRight = "}}"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)

	Debug bool
)

func adjustTemplate(t *template.Template, left, right string) *template.Template {
	if left == "" {
		left = DefaultDelimLeft
	}
	if right == "" {
		right = DefaultDelimRight
	}
	if left != DefaultDelimLeft || right != DefaultDelimRight {
		t.Delims(left, right)
	}
	return t
}

func NewTemplate(name, contents, left, right string) (*template.Template, error) {
	t, err := template.New(name).Parse(contents)
	if err != nil {
		return nil, fmt.Errorf("parse fail: %w", err)
	}
	return adjustTemplate(t, left, right), nil
	// if left == "" {
	// 	left = DefaultDelimLeft
	// }
	// if right == "" {
	// 	right = DefaultDelimRight
	// }
	// if left != DefaultDelimLeft || right != DefaultDelimRight {
	// 	t.Dilims(left, right)
	// }
	// return t, nil
}

func NewTemplateFS(name, left, right string, fs fs.FS, patterns ...string) (*template.Template, error) {
	t, err := template.New(name).ParseFS(fs, patterns...)
	if err != nil {
		return nil, fmt.Errorf("parse fail: %w", err)
	}
	return adjustTemplate(t, left, right), nil
}

// func ApplyTemplateDefault(t *template.Template, w io.Writer, m map[string]any) error {
// 	return ApplyTemplate(t, w, m)
// }

func ApplyTemplate(t *template.Template, w io.Writer, m map[string]any) error {
	return t.Execute(w, m)
}

// func TemplateIO(t *template.Template, w io.Writer, m map[string]any) error {
// 	return t.Execute(w, m)
// }

func TemplateFilesDefaultEnv(src, dest string) error {
	m := EnvMap()
	return TemplateFiles(src, dest, "", "", m)
}

func TemplateFilesDefault(src, dest string, m map[string]any) error {
	return TemplateFiles(src, dest, "", "", m)
}

func TemplateFiles(src, dest, left, right string, m map[string]any) error {
	in, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source %q -- %w", src, err)
	}

	t, err := NewTemplate("default", string(in), left, right)
	if err != nil {
		return fmt.Errorf("failed to make template from %q -- %w", src, err)
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create dest %q -- %w", dest, err)
	}

	if err = t.Execute(f, m); err != nil {
		f.Close()
		return fmt.Errorf("template execute failed: %w", err)
	}
	return f.Close()
}

func TemplateIO(r io.Reader, w io.Writer, left, right string, m map[string]any) error {
	in, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read reader -- %w", err)
	}

	t, err := NewTemplate("default", string(in), left, right)
	if err != nil {
		return fmt.Errorf("failed to make template -- %w", err)
	}

	if err = t.Execute(w, m); err != nil {
		return fmt.Errorf("template execute failed: %w", err)
	}
	return nil
}

func TemplateIOEnv(r io.Reader, w io.Writer, left, right string) error {
	m := EnvMap()
	return TemplateIO(r, w, left, right, m)
}

func EnvMap() map[string]any {
	m := make(map[string]any)
	for _, pair := range os.Environ() {
		key, val, ok := strings.Cut(pair, "=")
		if !ok {
			log.Printf("failed to cut env pair value %q", pair)
			continue
		}
		if Debug {
			log.Printf("%-16s=%v", key, val)
		}
		m[key] = val
	}
	return m
}
