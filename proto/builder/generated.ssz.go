// Code generated by fastssz. DO NOT EDIT.
// Hash: 437b6c535fb8770320d4e909148f105f0105004d5a471b241331da65d01b817e
package builder

import (
	ssz "github.com/prysmaticlabs/fastssz"
	github_com_prysmaticlabs_prysm_v3_consensus_types_primitives "github.com/prysmaticlabs/prysm/v3/consensus-types/primitives"
	v1 "github.com/prysmaticlabs/prysm/v3/proto/engine/v1"
)

// MarshalSSZ ssz marshals the BuilderPayloadAttributes object
func (b *BuilderPayloadAttributes) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(b)
}

// MarshalSSZTo ssz marshals the BuilderPayloadAttributes object to a target array
func (b *BuilderPayloadAttributes) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf

	// Field (0) 'Timestamp'
	dst = ssz.MarshalUint64(dst, b.Timestamp)

	// Field (1) 'PrevRandao'
	if size := len(b.PrevRandao); size != 32 {
		err = ssz.ErrBytesLengthFn("--.PrevRandao", size, 32)
		return
	}
	dst = append(dst, b.PrevRandao...)

	// Field (2) 'Slot'
	dst = ssz.MarshalUint64(dst, uint64(b.Slot))

	// Field (3) 'BlockHash'
	if size := len(b.BlockHash); size != 32 {
		err = ssz.ErrBytesLengthFn("--.BlockHash", size, 32)
		return
	}
	dst = append(dst, b.BlockHash...)

	return
}

// UnmarshalSSZ ssz unmarshals the BuilderPayloadAttributes object
func (b *BuilderPayloadAttributes) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 80 {
		return ssz.ErrSize
	}

	// Field (0) 'Timestamp'
	b.Timestamp = ssz.UnmarshallUint64(buf[0:8])

	// Field (1) 'PrevRandao'
	if cap(b.PrevRandao) == 0 {
		b.PrevRandao = make([]byte, 0, len(buf[8:40]))
	}
	b.PrevRandao = append(b.PrevRandao, buf[8:40]...)

	// Field (2) 'Slot'
	b.Slot = github_com_prysmaticlabs_prysm_v3_consensus_types_primitives.Slot(ssz.UnmarshallUint64(buf[40:48]))

	// Field (3) 'BlockHash'
	if cap(b.BlockHash) == 0 {
		b.BlockHash = make([]byte, 0, len(buf[48:80]))
	}
	b.BlockHash = append(b.BlockHash, buf[48:80]...)

	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the BuilderPayloadAttributes object
func (b *BuilderPayloadAttributes) SizeSSZ() (size int) {
	size = 80
	return
}

// HashTreeRoot ssz hashes the BuilderPayloadAttributes object
func (b *BuilderPayloadAttributes) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(b)
}

// HashTreeRootWith ssz hashes the BuilderPayloadAttributes object with a hasher
func (b *BuilderPayloadAttributes) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()

	// Field (0) 'Timestamp'
	hh.PutUint64(b.Timestamp)

	// Field (1) 'PrevRandao'
	if size := len(b.PrevRandao); size != 32 {
		err = ssz.ErrBytesLengthFn("--.PrevRandao", size, 32)
		return
	}
	hh.PutBytes(b.PrevRandao)

	// Field (2) 'Slot'
	hh.PutUint64(uint64(b.Slot))

	// Field (3) 'BlockHash'
	if size := len(b.BlockHash); size != 32 {
		err = ssz.ErrBytesLengthFn("--.BlockHash", size, 32)
		return
	}
	hh.PutBytes(b.BlockHash)

	if ssz.EnableVectorizedHTR {
		hh.MerkleizeVectorizedHTR(indx)
	} else {
		hh.Merkleize(indx)
	}
	return
}

// MarshalSSZ ssz marshals the BuilderPayloadAttributesV2 object
func (b *BuilderPayloadAttributesV2) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(b)
}

// MarshalSSZTo ssz marshals the BuilderPayloadAttributesV2 object to a target array
func (b *BuilderPayloadAttributesV2) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(84)

	// Field (0) 'Timestamp'
	dst = ssz.MarshalUint64(dst, b.Timestamp)

	// Field (1) 'PrevRandao'
	if size := len(b.PrevRandao); size != 32 {
		err = ssz.ErrBytesLengthFn("--.PrevRandao", size, 32)
		return
	}
	dst = append(dst, b.PrevRandao...)

	// Field (2) 'Slot'
	dst = ssz.MarshalUint64(dst, uint64(b.Slot))

	// Field (3) 'BlockHash'
	if size := len(b.BlockHash); size != 32 {
		err = ssz.ErrBytesLengthFn("--.BlockHash", size, 32)
		return
	}
	dst = append(dst, b.BlockHash...)

	// Offset (4) 'Withdrawals'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(b.Withdrawals) * 44

	// Field (4) 'Withdrawals'
	if size := len(b.Withdrawals); size > 16 {
		err = ssz.ErrListTooBigFn("--.Withdrawals", size, 16)
		return
	}
	for ii := 0; ii < len(b.Withdrawals); ii++ {
		if dst, err = b.Withdrawals[ii].MarshalSSZTo(dst); err != nil {
			return
		}
	}

	return
}

// UnmarshalSSZ ssz unmarshals the BuilderPayloadAttributesV2 object
func (b *BuilderPayloadAttributesV2) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 84 {
		return ssz.ErrSize
	}

	tail := buf
	var o4 uint64

	// Field (0) 'Timestamp'
	b.Timestamp = ssz.UnmarshallUint64(buf[0:8])

	// Field (1) 'PrevRandao'
	if cap(b.PrevRandao) == 0 {
		b.PrevRandao = make([]byte, 0, len(buf[8:40]))
	}
	b.PrevRandao = append(b.PrevRandao, buf[8:40]...)

	// Field (2) 'Slot'
	b.Slot = github_com_prysmaticlabs_prysm_v3_consensus_types_primitives.Slot(ssz.UnmarshallUint64(buf[40:48]))

	// Field (3) 'BlockHash'
	if cap(b.BlockHash) == 0 {
		b.BlockHash = make([]byte, 0, len(buf[48:80]))
	}
	b.BlockHash = append(b.BlockHash, buf[48:80]...)

	// Offset (4) 'Withdrawals'
	if o4 = ssz.ReadOffset(buf[80:84]); o4 > size {
		return ssz.ErrOffset
	}

	if o4 < 84 {
		return ssz.ErrInvalidVariableOffset
	}

	// Field (4) 'Withdrawals'
	{
		buf = tail[o4:]
		num, err := ssz.DivideInt2(len(buf), 44, 16)
		if err != nil {
			return err
		}
		b.Withdrawals = make([]*v1.Withdrawal, num)
		for ii := 0; ii < num; ii++ {
			if b.Withdrawals[ii] == nil {
				b.Withdrawals[ii] = new(v1.Withdrawal)
			}
			if err = b.Withdrawals[ii].UnmarshalSSZ(buf[ii*44 : (ii+1)*44]); err != nil {
				return err
			}
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the BuilderPayloadAttributesV2 object
func (b *BuilderPayloadAttributesV2) SizeSSZ() (size int) {
	size = 84

	// Field (4) 'Withdrawals'
	size += len(b.Withdrawals) * 44

	return
}

// HashTreeRoot ssz hashes the BuilderPayloadAttributesV2 object
func (b *BuilderPayloadAttributesV2) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(b)
}

// HashTreeRootWith ssz hashes the BuilderPayloadAttributesV2 object with a hasher
func (b *BuilderPayloadAttributesV2) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()

	// Field (0) 'Timestamp'
	hh.PutUint64(b.Timestamp)

	// Field (1) 'PrevRandao'
	if size := len(b.PrevRandao); size != 32 {
		err = ssz.ErrBytesLengthFn("--.PrevRandao", size, 32)
		return
	}
	hh.PutBytes(b.PrevRandao)

	// Field (2) 'Slot'
	hh.PutUint64(uint64(b.Slot))

	// Field (3) 'BlockHash'
	if size := len(b.BlockHash); size != 32 {
		err = ssz.ErrBytesLengthFn("--.BlockHash", size, 32)
		return
	}
	hh.PutBytes(b.BlockHash)

	// Field (4) 'Withdrawals'
	{
		subIndx := hh.Index()
		num := uint64(len(b.Withdrawals))
		if num > 16 {
			err = ssz.ErrIncorrectListSize
			return
		}
		for _, elem := range b.Withdrawals {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return
			}
		}
		if ssz.EnableVectorizedHTR {
			hh.MerkleizeWithMixinVectorizedHTR(subIndx, num, 16)
		} else {
			hh.MerkleizeWithMixin(subIndx, num, 16)
		}
	}

	if ssz.EnableVectorizedHTR {
		hh.MerkleizeVectorizedHTR(indx)
	} else {
		hh.Merkleize(indx)
	}
	return
}
