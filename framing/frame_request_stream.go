package framing

import (
	"encoding/binary"
	"fmt"
	"github.com/rsocket/rsocket-go/common"
)

const (
	minRequestStreamFrameLen = initReqLen
)

// FrameRequestStream is frame for requesting a completable stream.
type FrameRequestStream struct {
	*BaseFrame
}

// Validate returns error if frame is invalid.
func (p *FrameRequestStream) Validate() (err error) {
	if p.body.Len() < minRequestStreamFrameLen {
		err = errIncompleteFrame
	}
	return
}

func (p *FrameRequestStream) String() string {
	return fmt.Sprintf("FrameRequestStream{%s,data=%s,metadata=%s,initialRequestN=%d}", p.header, p.Data(), p.Metadata(), p.InitialRequestN())
}

// InitialRequestN returns initial request N.
func (p *FrameRequestStream) InitialRequestN() uint32 {
	return binary.BigEndian.Uint32(p.body.Bytes())
}

// Metadata returns metadata bytes.
func (p *FrameRequestStream) Metadata() []byte {
	return p.trySliceMetadata(4)
}

// Data returns data bytes.
func (p *FrameRequestStream) Data() []byte {
	return p.trySliceData(4)
}

// NewFrameRequestStream returns a new request stream frame.
func NewFrameRequestStream(id uint32, n uint32, data, metadata []byte, flags ...FrameFlag) *FrameRequestStream {
	fg := newFlags(flags...)
	bf := common.BorrowByteBuffer()
	for range [4]struct{}{} {
		_ = bf.WriteByte(0)
	}
	binary.BigEndian.PutUint32(bf.Bytes(), n)
	if len(metadata) > 0 {
		fg |= FlagMetadata
		_ = bf.WriteUint24(len(metadata))
		_, _ = bf.Write(metadata)
	}
	if len(data) > 0 {
		_, _ = bf.Write(data)
	}
	return &FrameRequestStream{
		NewBaseFrame(NewFrameHeader(id, FrameTypeRequestStream, fg), bf),
	}
}
