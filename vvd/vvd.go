package vvd

import (
	"encoding/binary"
	"io"
)

type VVD struct {
	Header   *Header
	Fixups   []*Fixup
	Vertexes []*Vertex
	Tangents [][4]float32
}

func (d *Decoder) decodeVVD() (*VVD, error) {
	var err error
	vvd := new(VVD)
	vvd.Header, err = d.decodeHeader()
	if err != nil {
		return nil, err
	}

	if _, err := d.r.Seek(int64(vvd.Header.FixupTableStart), io.SeekStart); err != nil {
		return nil, err
	}
	vvd.Fixups = make([]*Fixup, vvd.Header.NumFixups)
	for i := range vvd.Fixups {
		vvd.Fixups[i], err = d.decodeFixup()
		if err != nil {
			return nil, err
		}
	}

	if _, err := d.r.Seek(int64(vvd.Header.VertexDataStart), io.SeekStart); err != nil {
		return nil, err
	}
	vvd.Vertexes = make([]*Vertex, vvd.Header.NumLODVertexes[0])
	for i := range vvd.Vertexes {
		vvd.Vertexes[i], err = d.decodeVertex()
		if err != nil {
			return nil, err
		}
	}

	if _, err := d.r.Seek(int64(vvd.Header.TangentDataStart), io.SeekStart); err != nil {
		return nil, err
	}
	vvd.Tangents = make([][4]float32, vvd.Header.NumLODVertexes[0])
	err = binary.Read(d.r, binary.LittleEndian, &vvd.Tangents)
	if err != nil {
		return nil, err
	}
	return vvd, nil
}
