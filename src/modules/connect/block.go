package connect

import (
	"io"
)

type BlockType int

const BlockTypeText = BlockType(0)
const BlockTypeCommand = BlockType(1)
const BlockTypeCommandNoOpt = BlockType(2)
const BlockTypeSubNegotiation = BlockType(3)
const BlockTypeRaw = BlockType(4)

const (
	CR = byte('\r')
	LF = byte('\n')
)
const (
	IACCmdSE   = 240
	IACCmdNOP  = 241
	IACCmdData = 242

	IACCmdBreak = 243
	IACCmdGA    = 249
	IACCmdSB    = 250

	IACCmdWill = 251
	IACCmdWont = 252
	IACCmdDo   = 253
	IACCmdDont = 254

	IAC = 255
)

type Block struct {
	Type    BlockType
	Opt     byte
	Data    []byte
	Command byte
}

func EscapeIAC(data []byte) []byte {
	var result = make([]byte, 0, len(data))
	for i := range data {
		if data[i] == IAC {
			result = append(result, IAC, IAC)
		} else {
			result = append(result, data[i])
		}
	}
	return result
}
func (b *Block) Bytes() []byte {
	var data = []byte{}
	switch b.Type {
	case BlockTypeCommand:
		data = append(data, IAC, b.Command, b.Opt)
		return data
	case BlockTypeCommandNoOpt:
		data = append(data, IAC, b.Command)
		return data
	case BlockTypeSubNegotiation:
		data = append(data, IAC, IACCmdSB, b.Opt)
		data = append(data, EscapeIAC(b.Data)...)
		data = append(data, IAC, IACCmdSE)
		return data
	case BlockTypeRaw:
		return b.Data
	}
	return EscapeIAC(b.Data)
}
func (b *Block) WriteTo(w io.ByteWriter) error {
	var data = b.Bytes()
	for i := range data {
		err := w.WriteByte(data[i])
		if err != nil {
			return err
		}
	}
	return nil
}
func readSBDataBlock(r io.ByteReader, data []byte, iac bool) ([]byte, error) {
	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	if iac {
		switch b {
		case IACCmdSE:
			return data, nil
		default:
			data = append(data, b)
			return readSBDataBlock(r, data, false)
		}
	}
	if b == IAC {
		return readSBDataBlock(r, data, true)
	} else {
		data = append(data, b)
		return readSBDataBlock(r, data, false)
	}
}
func readSBOptBlock(r io.ByteReader) (*Block, error) {
	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	data, err := readSBDataBlock(r, []byte{}, false)
	if err != nil {
		return nil, err
	}
	return NewSubNegotiationBlock(b, data), nil
}
func readOptBlock(r io.ByteReader, cmd byte) (*Block, error) {
	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	return NewCommandBlock(cmd, b), nil
}
func readCmdBlock(r io.ByteReader) (*Block, error) {
	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	switch b {
	case IAC:
		return NewByteBlock(IAC), nil
	case IACCmdWill, IACCmdWont, IACCmdDo, IACCmdDont:
		return readOptBlock(r, b)
	case IACCmdSB:
		return readSBOptBlock(r)
	}
	return NewCommandNoOptBlock(b), nil
}
func ReadBlock(r io.ByteReader) (*Block, error) {
	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	if b != IAC {
		return NewByteBlock(b), nil
	}
	return readCmdBlock(r)
}
func NewByteBlock(data byte) *Block {
	return &Block{
		Type: BlockTypeText,
		Data: []byte{data},
	}
}
func NewTextBlock(data []byte) *Block {
	return &Block{
		Type: BlockTypeText,
		Data: data,
	}
}
func NewCommandBlock(cmd byte, opt byte) *Block {
	return &Block{
		Type:    BlockTypeCommand,
		Command: cmd,
		Opt:     opt,
	}
}
func NewCommandNoOptBlock(cmd byte) *Block {
	return &Block{
		Type:    BlockTypeCommandNoOpt,
		Command: cmd,
	}
}

func NewSubNegotiationBlock(opt byte, data []byte) *Block {
	return &Block{
		Type: BlockTypeSubNegotiation,
		Opt:  opt,
		Data: data,
	}
}
func NewRawBlock(data []byte) *Block {
	return &Block{
		Type: BlockTypeRaw,
		Data: data,
	}
}
