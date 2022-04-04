package genesis

import (
	"bytes"
	gocrypto "crypto"

	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/crypto"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/errors"
)

var (
	ErrRootGenesisIsNil   = errors.New("root genesis is nil")
	ErrVerifierIsNil      = errors.New("verifier is nil")
	ErrPartitionsNotFound = errors.New("partitions not found")
)

func (x *RootGenesis) IsValid(verifier crypto.Verifier, hashAlgorithm gocrypto.Hash) error {
	if x == nil {
		return ErrRootGenesisIsNil
	}
	if verifier == nil {
		return ErrVerifierIsNil
	}
	pubKeyBytes, err := verifier.MarshalPublicKey()
	if err != nil {
		return err
	}
	if !bytes.Equal(pubKeyBytes, x.TrustBase) {
		return errors.Errorf("invalid trust base. expected %X, got %X", pubKeyBytes, x.TrustBase)
	}

	if len(x.Partitions) == 0 {
		return ErrPartitionsNotFound
	}
	for _, p := range x.Partitions {
		if err = p.IsValid(verifier, hashAlgorithm); err != nil {
			return err
		}
	}
	return nil
}

func (x *RootGenesis) GetRoundNumber() uint64 {
	return x.Partitions[0].Certificate.UnicitySeal.RootChainRoundNumber
}

func (x *RootGenesis) GetPreviousBlockHash() []byte {
	return x.Partitions[0].Certificate.UnicitySeal.PreviousHash
}
func (x *RootGenesis) GetBlockHash() []byte {
	return x.Partitions[0].Certificate.UnicitySeal.Hash
}

func (x *RootGenesis) GetPartitionRecords() []*PartitionRecord {
	records := make([]*PartitionRecord, len(x.Partitions))
	for i, partition := range x.Partitions {
		records[i] = &PartitionRecord{
			SystemDescriptionRecord: partition.SystemDescriptionRecord,
			Validators:              partition.Nodes,
		}
	}
	return records
}
