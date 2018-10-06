package p2p

import (
	"encoding/binary"
	"hash/crc32"
	"testing"

	"github.com/golang/snappy"
	"github.com/stretchr/testify/assert"
)

var (
	testChainID     uint32 = 1
	testMessageType        = Ping
	testVerion      uint16 = 1
	testData               = []byte("test")
)

func TestP2PMessageContent(t *testing.T) {
	m := newP2PMessage(testChainID, testMessageType, testVerion, testData)
	assert.NotNil(t, m)
	content := m.content()
	assert.NotNil(t, content)
	assert.True(t, len(content) > dataBegin)
	assert.Equal(t, testChainID, binary.BigEndian.Uint32(content[chainIDBegin:chainIDEnd]))
	assert.Equal(t, testMessageType, MessageType(binary.BigEndian.Uint16(content[messageTypeBegin:messageTypeEnd])))
	assert.Equal(t, testVerion, binary.BigEndian.Uint16(content[versionBegin:versionEnd]))
	assert.EqualValues(t, len(snappy.Encode(nil, testData)), binary.BigEndian.Uint32(content[dataLengthBegin:dataLengthEnd]))
	assert.Equal(t, crc32.ChecksumIEEE(snappy.Encode(nil, testData)), binary.BigEndian.Uint32(content[dataChecksumBegin:dataChecksumEnd]))
	assert.Equal(t, snappy.Encode(nil, testData), content[dataBegin:])
}

func TestP2PChanID(t *testing.T) {
	m := newP2PMessage(testChainID, testMessageType, testVerion, testData)
	assert.Equal(t, testChainID, m.chainID())
}

func TestP2PMessageType(t *testing.T) {
	m := newP2PMessage(testChainID, testMessageType, testVerion, testData)
	assert.Equal(t, testMessageType, m.messageType())
}

func TestP2PVersion(t *testing.T) {
	m := newP2PMessage(testChainID, testMessageType, testVerion, testData)
	assert.Equal(t, testVerion, m.version())
}

func TestP2PDataLength(t *testing.T) {
	m := newP2PMessage(testChainID, testMessageType, testVerion, testData)
	assert.EqualValues(t, len(snappy.Encode(nil, testData)), m.dataLength())
}

func TestP2PChecksum(t *testing.T) {
	m := newP2PMessage(testChainID, testMessageType, testVerion, testData)
	assert.Equal(t, crc32.ChecksumIEEE(snappy.Encode(nil, testData)), m.checksum())
}

func TestP2PRawData(t *testing.T) {
	m := newP2PMessage(testChainID, testMessageType, testVerion, testData)
	assert.Equal(t, snappy.Encode(nil, testData), m.rawData())
}

func TestP2PData(t *testing.T) {
	m := newP2PMessage(testChainID, testMessageType, testVerion, testData)
	data, err := m.data()
	assert.NotNil(t, data)
	assert.Nil(t, err)
	assert.Equal(t, testData, data)
}

func TestP2PMessageParse(t *testing.T) {
	m := newP2PMessage(testChainID, testMessageType, testVerion, testData)
	newM, err := parseP2PMessage(m.content())
	assert.Nil(t, err)
	assert.Equal(t, m, newM)
}
