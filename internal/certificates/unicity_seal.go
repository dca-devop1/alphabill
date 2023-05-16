package certificates

import (
	"bytes"
	"errors"
	"fmt"
	"hash"

	"github.com/alphabill-org/alphabill/internal/crypto"
	"github.com/alphabill-org/alphabill/internal/util"
)

var (
	ErrUnicitySealIsNil          = errors.New("unicity seal is nil")
	ErrSignerIsNil               = errors.New("signer is nil")
	ErrUnicitySealHashIsNil      = errors.New("hash is nil")
	ErrInvalidBlockNumber        = errors.New("invalid block number")
	ErrUnicitySealSignatureIsNil = errors.New("no signatures")
	ErrRootValidatorInfoMissing  = errors.New("root node info is missing")
	ErrUnknownSigner             = errors.New("unknown signer")
	errInvalidTimestamp          = errors.New("invalid timestamp")
	errUnicitySealNoSignature    = errors.New("unicity seal is missing signature")
)

func (x *UnicitySeal) IsValid(verifiers map[string]crypto.Verifier) error {
	if x == nil {
		return ErrUnicitySealIsNil
	}
	if len(verifiers) == 0 {
		return ErrRootValidatorInfoMissing
	}
	if x.Hash == nil {
		return ErrUnicitySealHashIsNil
	}
	if x.RootChainRoundNumber < 1 {
		return ErrInvalidBlockNumber
	}
	if x.Timestamp < util.GenesisTime {
		return errInvalidTimestamp
	}
	if len(x.Signatures) == 0 {
		return ErrUnicitySealSignatureIsNil
	}
	return x.Verify(verifiers)
}

func (x *UnicitySeal) Sign(id string, signer crypto.Signer) error {
	if signer == nil {
		return ErrSignerIsNil
	}
	signature, err := signer.SignBytes(x.Bytes())
	if err != nil {
		return fmt.Errorf("sign failed, %w", err)
	}
	// initiate signatures
	if x.Signatures == nil {
		x.Signatures = make(map[string][]byte)
	}
	x.Signatures[id] = signature
	return nil
}

func (x *UnicitySeal) Bytes() []byte {
	var b bytes.Buffer
	b.Write(util.Uint64ToBytes(x.RootChainRoundNumber))
	b.Write(util.Uint64ToBytes(x.Timestamp))
	b.Write(x.PreviousHash)
	b.Write(x.Hash)
	return b.Bytes()
}

func (x *UnicitySeal) AddToHasher(hasher hash.Hash) {
	hasher.Write(x.Bytes())
}

func (x *UnicitySeal) Verify(verifiers map[string]crypto.Verifier) error {
	if verifiers == nil {
		return ErrRootValidatorInfoMissing
	}
	if len(x.Signatures) == 0 {
		return errUnicitySealNoSignature
	}
	// Verify all signatures, all must be from known origin and valid
	for id, sig := range x.Signatures {
		// Find verifier info
		ver, f := verifiers[id]
		if !f {
			return ErrUnknownSigner
		}
		err := ver.VerifyBytes(sig, x.Bytes())
		if err != nil {
			return fmt.Errorf("invalid unicity seal signature, %w", err)
		}
	}
	return nil
}
