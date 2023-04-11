package blockproposal

import (
	"bytes"
	gocrypto "crypto"
	"errors"
	"fmt"

	"github.com/alphabill-org/alphabill/internal/crypto"
)

var (
	ErrBlockProposalIsNil      = errors.New("block proposal is nil")
	ErrTrustBaseIsNil          = errors.New("trust base is nil")
	ErrSignerIsNil             = errors.New("signer is nil")
	ErrNodeVerifierIsNil       = errors.New("node signature verifier is nil")
	ErrInvalidSystemIdentifier = errors.New("invalid system identifier")
	errBlockProposerIDMissing  = errors.New("block proposer id is missing")
)

func (x *BlockProposal) IsValid(nodeSignatureVerifier crypto.Verifier, ucTrustBase map[string]crypto.Verifier, algorithm gocrypto.Hash, systemIdentifier []byte, systemDescriptionHash []byte) error {
	if x == nil {
		return ErrBlockProposalIsNil
	}
	if nodeSignatureVerifier == nil {
		return ErrNodeVerifierIsNil
	}
	if len(x.NodeIdentifier) == 0 {
		return errBlockProposerIDMissing
	}
	if ucTrustBase == nil {
		return ErrTrustBaseIsNil
	}
	if !bytes.Equal(systemIdentifier, x.SystemIdentifier) {
		return fmt.Errorf("%w, expected %X, got %X", ErrInvalidSystemIdentifier, systemIdentifier, x.SystemIdentifier)
	}
	if err := x.UnicityCertificate.IsValid(ucTrustBase, algorithm, systemIdentifier, systemDescriptionHash); err != nil {
		return err
	}
	return x.Verify(algorithm, nodeSignatureVerifier)
}

func (x *BlockProposal) Hash(algorithm gocrypto.Hash) ([]byte, error) {
	hasher := algorithm.New()
	hasher.Write(x.SystemIdentifier)
	hasher.Write([]byte(x.NodeIdentifier))
	x.UnicityCertificate.AddToHasher(hasher)
	for _, tx := range x.Transactions {
		hasher.Write(tx.Bytes())
	}
	return hasher.Sum(nil), nil
}

func (x *BlockProposal) Sign(algorithm gocrypto.Hash, signer crypto.Signer) error {
	if signer == nil {
		return ErrSignerIsNil
	}
	hash, err := x.Hash(algorithm)
	if err != nil {
		return err
	}
	x.Signature, err = signer.SignHash(hash)
	if err != nil {
		return err
	}
	return nil
}

func (x *BlockProposal) Verify(algorithm gocrypto.Hash, nodeSignatureVerifier crypto.Verifier) error {
	if nodeSignatureVerifier == nil {
		return ErrNodeVerifierIsNil
	}
	hash, err := x.Hash(algorithm)
	if err != nil {
		return err
	}
	return nodeSignatureVerifier.VerifyHash(x.Signature, hash)
}
