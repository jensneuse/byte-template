package byte_template

import (
	"bytes"
	"io"
	"testing"
)

func TestTemplating(t *testing.T) {

	run := func(input, expectedOutput string, fetch Fetch, directiveDefinitions ...DirectiveDefinition) func(t *testing.T) {
		return func(t *testing.T) {
			template := New(directiveDefinitions...)
			buf := bytes.Buffer{}
			err := template.Execute(&buf, []byte(input), fetch)
			if err != nil {
				t.Fatal(err)
			}
			want := expectedOutput
			got := buf.String()

			if want != got {
				t.Fatalf("want: %s, got: %s", want, got)
			}
		}
	}

	t.Run("should not crash when first token is open", run("{", "{", func(w io.Writer, path []byte) (err error) {
		return
	}))
	t.Run("single open", run("/api/user/{", "/api/user/{", func(w io.Writer, path []byte) (err error) {
		return
	}))
	t.Run("simple id", run("/api/user/{{.id }}", "/api/user/1", func(w io.Writer, path []byte) (err error) {
		if string(path) == ".id" {
			_, err = w.Write([]byte("1"))
		}
		return
	}))
	t.Run("simple id", run("/api/user/{{.id}}", "/api/user/1", func(w io.Writer, path []byte) (err error) {
		if string(path) == ".id" {
			_, err = w.Write([]byte("1"))
		}
		return
	}))
	t.Run("simple id", run("/api/user/{{ .id }}", "/api/user/1", func(w io.Writer, path []byte) (err error) {
		if string(path) == ".id" {
			_, err = w.Write([]byte("1"))
		}
		return
	}))
	t.Run("simple directive with item", run("/api/user/{{ toLower .Name }}", "/api/user/sergey", func(w io.Writer, path []byte) (err error) {
		if string(path) == ".Name" {
			_, err = w.Write([]byte("Sergey"))
		}
		return
	}, DirectiveDefinition{
		Name: []byte("toLower"),
		Resolve: func(w io.Writer, arg []byte) error {
			_, err := w.Write(bytes.ToLower(arg))
			return err
		},
	}))
}

func BenchmarkTemplate_Execute(b *testing.B) {
	input := []byte("/api/user/{{ customDirective .Name }}")
	variable := []byte("Sergey")
	fetch := func(w io.Writer, path []byte) (err error) {
		_, err = w.Write(variable)
		return
	}
	template := New(DirectiveDefinition{
		Name: []byte("customDirective"),
		Resolve: func(w io.Writer, arg []byte) error {
			_, err := w.Write(arg)
			return err
		},
	})
	buf := bytes.Buffer{}

	b.SetBytes(int64(len(input) + len(variable)))
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		_ = template.Execute(&buf, input, fetch)
	}
}
