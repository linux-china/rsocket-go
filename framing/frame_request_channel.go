package framing

import (
	"encoding/binary"
	"fmt"
	"github.com/rsocket/rsocket-go/common"
)

const (
	initReqLen                = 4
	minRequestChannelFrameLen = initReqLen
)

// FrameRequestChannel is frame for RequestChannel.
type FrameRequestChannel struct {
	*BaseFrame
}

// Validate returns error if frame is invalid.
func (p *FrameRequestChannel) Validate() (err error) {
	if p.body.Len() < minRequestChannelFrameLen {
		err = errIncompleteFrame
	}
	return
}

func (p *FrameRequestChannel) String() string {
	return fmt.Sprintf("FrameRequestChannel{%s,data=%s,metadata=%s,initialRequestN=%d}", p.header, p.Data(), p.Metadata(), p.InitialRequestN())
}

// InitialRequestN returns initial N.
func (p *FrameRequestChannel) InitialRequestN() uint32 {
	return binary.BigEndian.Uint32(p.body.Bytes())
}

// Metadata returns metadata bytes.
func (p *FrameRequestChannel) Metadata() []byte {
	return p.trySliceMetadata(initReqLen)
}

// Data returns data bytes.
func (p *FrameRequestChannel) Data() []byte {
	return p.trySliceData(initReqLen)
}

// NewFrameRequestChannel returns a new RequestChannel frame.
func NewFrameRequestChannel(sid uint32, n uint32, data, metadata []byte, flags ...FrameFlag) *FrameRequestChannel {
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
	return &FrameRequestChannel{
		&BaseFrame{
			header: NewFrameHeader(sid, FrameTypeRequestChannel, fg),
			body:   bf,
		},
	}
}
