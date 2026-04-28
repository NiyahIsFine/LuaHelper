package langserver

import (
	"context"
	"io/ioutil"
	lsp "luahelper-lsp/langserver/protocol"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAnnotateClassMemberReferencesAcrossFiles(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	paths, _ := filepath.Split(filename)

	strRootPath := paths + "../testdata/references"
	strRootPath, _ = filepath.Abs(strRootPath)
	strRootURI := "file://" + strRootPath
	lspServer := createLspTest(strRootPath, strRootURI)
	ctx := context.Background()

	for _, name := range []string{"class_def.lua", "class_use.lua"} {
		fileName := filepath.Join(strRootPath, name)
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			t.Fatalf("read file:%s err=%s", fileName, err.Error())
		}

		err = lspServer.TextDocumentDidOpen(ctx, lsp.DidOpenTextDocumentParams{
			TextDocument: lsp.TextDocumentItem{
				URI:  lsp.DocumentURI(fileName),
				Text: string(data),
			},
		})
		if err != nil {
			t.Fatalf("didopen file:%s err=%s", fileName, err.Error())
		}
	}

	defFile := filepath.Join(strRootPath, "class_def.lua")
	useFile := filepath.Join(strRootPath, "class_use.lua")
	useFileURI := lsp.DocumentURI("file://" + filepath.ToSlash(useFile))

	tests := []struct {
		name     string
		line     uint32
		char     uint32
		wantLine uint32
		wantChar uint32
	}{
		{
			name:     "field annotation",
			line:     1,
			char:     11,
			wantLine: 3,
			wantChar: 7,
		},
		{
			name:     "dynamic member function",
			line:     4,
			char:     33,
			wantLine: 4,
			wantChar: 7,
		},
		{
			name:     "dynamic member field",
			line:     5,
			char:     14,
			wantLine: 5,
			wantChar: 7,
		},
		{
			name:     "implicit alias variable",
			line:     4,
			char:     33,
			wantLine: 11,
			wantChar: 12,
		},
		{
			name:     "implicit function return",
			line:     4,
			char:     33,
			wantLine: 19,
			wantChar: 33,
		},
		{
			name:     "implicit assigned function return",
			line:     4,
			char:     33,
			wantLine: 22,
			wantChar: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locations, err := lspServer.TextDocumentReferences(ctx, lsp.ReferenceParams{
				TextDocumentPositionParams: lsp.TextDocumentPositionParams{
					TextDocument: lsp.TextDocumentIdentifier{URI: lsp.DocumentURI(defFile)},
					Position: lsp.Position{
						Line:      tt.line,
						Character: tt.char,
					},
				},
			})
			if err != nil {
				t.Fatalf("references error: %s", err.Error())
			}

			for _, location := range locations {
				if location.URI != useFileURI {
					continue
				}
				if location.Range.Start.Line == tt.wantLine && location.Range.Start.Character == tt.wantChar {
					return
				}
			}

			t.Fatalf("reference in class_use.lua at %d:%d not found, got %#v", tt.wantLine, tt.wantChar, locations)
		})
	}
}
