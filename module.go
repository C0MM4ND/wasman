package wasman

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/c0mm4nd/wasman/segments"
	"github.com/c0mm4nd/wasman/types"
)

// errors on parsing module
var (
	ErrInvalidMagicNumber = errors.New("invalid magic number")
	ErrInvalidVersion     = errors.New("invalid version header")
)

// Module is a standard wasm module implement according to wasm v1, https://www.w3.org/TR/wasm-core-1/#syntax-module%E2%91%A0
type Module struct {
	*ModuleConfig

	// sections
	TypesSection     []*types.FuncType
	ImportsSection   []*segments.ImportSegment
	FunctionsSection []uint32
	TablesSection    []*types.TableType
	MemorySection    []*types.MemoryType
	GlobalsSection   []*segments.GlobalSegment
	ExportsSection   map[string]*segments.ExportSegment
	StartSection     []uint32
	ElementsSection  []*segments.ElemSegment
	CodesSection     []*segments.CodeSegment
	DataSection      []*segments.DataSegment

	// index spaces
	indexSpace *indexSpace
}

// index to the imports
type indexSpace struct {
	Functions []fn
	Globals   []*global
	Tables    [][]*uint32
	Memories  [][]byte
}

type global struct {
	Type *types.GlobalType
	Val  interface{}
}

// NewModule reads bytes from the io.Reader and read all sections, finally return a wasman.Module entity if no error
func NewModule(r io.Reader, config *ModuleConfig) (*Module, error) {
	// magic number
	buf := make([]byte, 4)
	if n, err := io.ReadFull(r, buf); err != nil || n != 4 || !bytes.Equal(buf, magic) {
		return nil, ErrInvalidMagicNumber
	}

	// version
	if n, err := io.ReadFull(r, buf); err != nil || n != 4 || !bytes.Equal(buf, version) {
		return nil, ErrInvalidVersion
	}

	module := &Module{}
	if err := module.readSections(r); err != nil {
		return nil, fmt.Errorf("readSections failed: %w", err)
	}

	if config != nil {
		module.ModuleConfig = config
	}

	return module, nil
}
