package crypto

import (
	"crypto/elliptic"

	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/domain/canonicalizer"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/errors"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/errors/errstr"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/hash"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

type (
	verifierSecp256k1 struct {
		pubKey []byte
	}
)

// CompressedSecp256K1PublicKeySize is size of public key in compressed format
const CompressedSecp256K1PublicKeySize = 33

// NewVerifierSecp256k1 creates new verifier from an existing Secp256k1 compressed public key.
func NewVerifierSecp256k1(compressedPubKey []byte) (Verifier, error) {
	if len(compressedPubKey) != CompressedSecp256K1PublicKeySize {
		return nil, errors.Wrapf(errors.ErrInvalidArgument, "pubkey must be %d bytes long, but is %d", CompressedSecp256K1PublicKeySize, len(compressedPubKey))
	}
	x, y := secp256k1.DecompressPubkey(compressedPubKey)
	pubkey := elliptic.Marshal(secp256k1.S256(), x, y)
	return &verifierSecp256k1{pubkey}, nil
}

// VerifyBytes hashes the data with SHA256 and verifies it using the public key of the verifier.
func (v *verifierSecp256k1) VerifyBytes(sig []byte, data []byte) error {
	if v == nil || v.pubKey == nil || sig == nil || data == nil {
		return errors.Wrap(errors.ErrInvalidArgument, errstr.NilArgument)
	}
	if len(sig) == crypto.SignatureLength {
		// If signature contains recovery ID, then remove it.
		sig = sig[:len(sig)-1]
	}
	if len(sig) != crypto.RecoveryIDOffset {
		return errors.Wrapf(errors.ErrInvalidState, "signature length is %d b (expected %d b)", len(sig), crypto.RecoveryIDOffset)
	}
	if secp256k1.VerifySignature(v.pubKey, hash.Sum256(data), sig) {
		return nil
	}
	return errors.Wrap(errors.ErrVerificationFailed, "signature verify failed")
}

// VerifyObject verifies the signature of canonicalizable object with public key.
func (v *verifierSecp256k1) VerifyObject(sig []byte, obj canonicalizer.Canonicalizer, opts ...canonicalizer.Option) error {
	data, err := canonicalizer.Canonicalize(obj, opts...)
	if err != nil {
		return errors.Wrap(err, "could not canonicalize the object")
	}
	return v.VerifyBytes(sig, data)
}

// MarshalPublicKey returns compressed public key, 33 bytes
func (v *verifierSecp256k1) MarshalPublicKey() ([]byte, error) {
	if v == nil || v.pubKey == nil {
		return nil, errors.Wrap(errors.ErrInvalidArgument, errstr.NilArgument)
	}
	pubkey, err := crypto.UnmarshalPubkey(v.pubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert public key bytes to ECDSA public key")
	}
	compressed := secp256k1.CompressPubkey(pubkey.X, pubkey.Y)
	return compressed, nil
}