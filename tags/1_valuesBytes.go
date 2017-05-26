package tags

import "unsafe"

var sizeOfUint16 = (int)(unsafe.Sizeof(uint16(0)))

type valueBytes struct {
	buf []byte
	wIdx, rIdx int
}

func (vb *valueBytes) growIfRequired(expected int) {
	if len(vb.buf)-vb.wIdx < expected {
		tmp := make([]byte, 2*(len(vb.buf)+1)+expected)
		copy(tmp, vb.buf)
		vb.buf = tmp
	}
}

func (vb *valueBytes) writeValue(v []byte) {
	length := len(v)
	vb.growIfRequired(sizeOfUint16 + length)

/*	length := len(v)
	endIdx := vb.wIdx + sizeOfUint16 + int(length)
	vb.growIfRequired(endIdx)
*/	
	// writing length of v
	bytes := *(*[2]byte)(unsafe.Pointer(&length))
	vb.buf[vb.wIdx] = bytes[0]
	vb.wIdx++
	vb.buf[vb.wIdx] = bytes[1]
	vb.wIdx++

	if length == 0 {
		// No value was encoded for this key
		return
	}

	// writing v
	copy(vb.buf[vb.wIdx:], v)
	vb.wIdx += length
}

// readValue is the helper method to read the values when decoding valueBytes to a map[Key][]byte.
// It is meant to be used by toMap(...) only.
func (vb *valueBytes) readValue() []byte {
	// read length of v
	length := (int)(*(*uint16)(unsafe.Pointer(&vb.buf[vb.rIdx])))
	vb.rIdx += sizeOfUint16
	if length == 0 {
		// No value was encoded for this key
		return nil
	}

	// read value of v
	v := make([]byte, length)
	endIdx := vb.rIdx+length
	copy(v, vb.buf[vb.rIdx:endIdx])
	vb.rIdx += endIdx
	return v
}

func (vb *valueBytes) toMap(ks []Key) map[Key][]byte {
	m := make(map[Key][]byte, len(ks))
	for _, k := range ks {
		v := vb.readValue()
		if v != nil {
			m[k] = v
		}
	}
	return m	
}