package tokens

import (
	"bytes"
	"crypto"
	"github.com/alphabill-org/alphabill/internal/block"
	"hash"

	"github.com/alphabill-org/alphabill/internal/errors"
	"github.com/alphabill-org/alphabill/internal/txsystem"
	"github.com/alphabill-org/alphabill/internal/util"
	"github.com/holiman/uint256"
	"google.golang.org/protobuf/proto"
)

const (
	TypeCreateNonFungibleTokenTypeAttributes = "CreateNonFungibleTokenTypeAttributes"
	TypeMintNonFungibleTokenAttributes       = "MintNonFungibleTokenAttributes"
	TypeTransferNonFungibleTokenAttributes   = "TransferNonFungibleTokenAttributes"
	TypeUpdateNonFungibleTokenAttributes     = "UpdateNonFungibleTokenAttributes"
	TypeCreateFungibleTokenTypeAttributes    = "CreateFungibleTokenTypeAttributes"
	TypeMintFungibleTokenAttributes          = "MintFungibleTokenAttributes"
	TypeTransferFungibleTokenAttributes      = "TransferFungibleTokenAttributes"
	TypeSplitFungibleTokenAttributes         = "SplitFungibleTokenAttributes"
	TypeBurnFungibleTokenAttributes          = "BurnFungibleTokenAttributes"
	TypeJoinFungibleTokenAttributes          = "JoinFungibleTokenAttributes"

	protobufTypeUrlPrefix                       = "type.googleapis.com/alphabill.tokens.v1."
	typeURLCreateNonFungibleTokenTypeAttributes = protobufTypeUrlPrefix + TypeCreateNonFungibleTokenTypeAttributes
	typeURLMintNonFungibleTokenAttributes       = protobufTypeUrlPrefix + TypeMintNonFungibleTokenAttributes
	typeURLTransferNonFungibleTokenAttributes   = protobufTypeUrlPrefix + TypeTransferNonFungibleTokenAttributes
	typeURLUpdateNonFungibleTokenAttributes     = protobufTypeUrlPrefix + TypeUpdateNonFungibleTokenAttributes
	typeURLCreateFungibleTokenTypeAttributes    = protobufTypeUrlPrefix + TypeCreateFungibleTokenTypeAttributes
	typeURLMintFungibleTokenAttributes          = protobufTypeUrlPrefix + TypeMintFungibleTokenAttributes
	typeURLTransferFungibleTokenAttributes      = protobufTypeUrlPrefix + TypeTransferFungibleTokenAttributes
	typeURLSplitFungibleTokenAttributes         = protobufTypeUrlPrefix + TypeSplitFungibleTokenAttributes
	typeURLBurnFungibleTokenAttributes          = protobufTypeUrlPrefix + TypeBurnFungibleTokenAttributes
	typeURLJoinFungibleTokenAttributes          = protobufTypeUrlPrefix + TypeJoinFungibleTokenAttributes
)

// TransactionTypes contains all transaction types supported by the user token partition.
var TransactionTypes = map[string]proto.Message{
	TypeCreateNonFungibleTokenTypeAttributes: &CreateNonFungibleTokenTypeAttributes{},
	TypeMintNonFungibleTokenAttributes:       &MintNonFungibleTokenAttributes{},
	TypeTransferNonFungibleTokenAttributes:   &TransferNonFungibleTokenAttributes{},
	TypeUpdateNonFungibleTokenAttributes:     &UpdateNonFungibleTokenAttributes{},
	TypeCreateFungibleTokenTypeAttributes:    &CreateFungibleTokenTypeAttributes{},
	TypeMintFungibleTokenAttributes:          &MintFungibleTokenAttributes{},
	TypeTransferFungibleTokenAttributes:      &TransferFungibleTokenAttributes{},
	TypeSplitFungibleTokenAttributes:         &SplitFungibleTokenAttributes{},
	TypeBurnFungibleTokenAttributes:          &BurnFungibleTokenAttributes{},
	TypeJoinFungibleTokenAttributes:          &JoinFungibleTokenAttributes{},
}

type (
	wrapper struct {
		transaction *txsystem.Transaction
		hashFunc    crypto.Hash
		hashValue   []byte
	}

	createNonFungibleTokenTypeWrapper struct {
		wrapper
		attributes *CreateNonFungibleTokenTypeAttributes
	}

	mintNonFungibleTokenWrapper struct {
		wrapper
		attributes *MintNonFungibleTokenAttributes
	}

	transferNonFungibleTokenWrapper struct {
		wrapper
		attributes *TransferNonFungibleTokenAttributes
	}

	updateNonFungibleTokenWrapper struct {
		wrapper
		attributes *UpdateNonFungibleTokenAttributes
	}

	createFungibleTokenTypeWrapper struct {
		wrapper
		attributes *CreateFungibleTokenTypeAttributes
	}

	mintFungibleTokenWrapper struct {
		wrapper
		attributes *MintFungibleTokenAttributes
	}

	transferFungibleTokenWrapper struct {
		wrapper
		attributes *TransferFungibleTokenAttributes
	}

	splitFungibleTokenWrapper struct {
		wrapper
		attributes *SplitFungibleTokenAttributes
	}

	burnFungibleTokenWrapper struct {
		wrapper
		attributes *BurnFungibleTokenAttributes
	}

	joinFungibleTokenWrapper struct {
		wrapper
		attributes       *JoinFungibleTokenAttributes
		burnTransactions []BurnFungibleToken
	}
)

func NewGenericTx(tx *txsystem.Transaction) (txsystem.GenericTransaction, error) {
	switch tx.TransactionAttributes.TypeUrl {
	case typeURLCreateNonFungibleTokenTypeAttributes:
		return convertToWrapper(
			&CreateNonFungibleTokenTypeAttributes{},
			func(a *CreateNonFungibleTokenTypeAttributes) (txsystem.GenericTransaction, error) {
				return &createNonFungibleTokenTypeWrapper{
					wrapper:    wrapper{transaction: tx},
					attributes: a,
				}, nil
			},
			tx)
	case typeURLMintNonFungibleTokenAttributes:
		return convertToWrapper(
			&MintNonFungibleTokenAttributes{},
			func(a *MintNonFungibleTokenAttributes) (txsystem.GenericTransaction, error) {
				return &mintNonFungibleTokenWrapper{
					wrapper:    wrapper{transaction: tx},
					attributes: a,
				}, nil
			},
			tx)
	case typeURLTransferNonFungibleTokenAttributes:
		return convertToWrapper(
			&TransferNonFungibleTokenAttributes{},
			func(a *TransferNonFungibleTokenAttributes) (txsystem.GenericTransaction, error) {
				return &transferNonFungibleTokenWrapper{
					wrapper:    wrapper{transaction: tx},
					attributes: a,
				}, nil
			},
			tx)
	case typeURLUpdateNonFungibleTokenAttributes:
		return convertToWrapper(
			&UpdateNonFungibleTokenAttributes{},
			func(a *UpdateNonFungibleTokenAttributes) (txsystem.GenericTransaction, error) {
				return &updateNonFungibleTokenWrapper{
					wrapper:    wrapper{transaction: tx},
					attributes: a,
				}, nil
			},
			tx)
	case typeURLCreateFungibleTokenTypeAttributes:
		return convertToWrapper(
			&CreateFungibleTokenTypeAttributes{},
			func(a *CreateFungibleTokenTypeAttributes) (txsystem.GenericTransaction, error) {
				return &createFungibleTokenTypeWrapper{
					wrapper:    wrapper{transaction: tx},
					attributes: a,
				}, nil
			},
			tx)
	case typeURLMintFungibleTokenAttributes:
		return convertToWrapper(
			&MintFungibleTokenAttributes{},
			func(a *MintFungibleTokenAttributes) (txsystem.GenericTransaction, error) {
				return &mintFungibleTokenWrapper{
					wrapper:    wrapper{transaction: tx},
					attributes: a,
				}, nil
			},
			tx)
	case typeURLTransferFungibleTokenAttributes:
		return convertToWrapper(
			&TransferFungibleTokenAttributes{},
			func(a *TransferFungibleTokenAttributes) (txsystem.GenericTransaction, error) {
				return &transferFungibleTokenWrapper{
					wrapper:    wrapper{transaction: tx},
					attributes: a,
				}, nil
			},
			tx)
	case typeURLSplitFungibleTokenAttributes:
		return convertToWrapper(
			&SplitFungibleTokenAttributes{},
			func(a *SplitFungibleTokenAttributes) (txsystem.GenericTransaction, error) {
				return &splitFungibleTokenWrapper{
					wrapper:    wrapper{transaction: tx},
					attributes: a,
				}, nil
			},
			tx)
	case typeURLBurnFungibleTokenAttributes:
		return convertToWrapper(
			&BurnFungibleTokenAttributes{},
			func(a *BurnFungibleTokenAttributes) (txsystem.GenericTransaction, error) {
				return &burnFungibleTokenWrapper{
					wrapper:    wrapper{transaction: tx},
					attributes: a,
				}, nil
			},
			tx)
	case typeURLJoinFungibleTokenAttributes:
		return convertToWrapper(
			&JoinFungibleTokenAttributes{},
			func(a *JoinFungibleTokenAttributes) (txsystem.GenericTransaction, error) {
				burnTransactions := a.BurnTransactions
				lenBTxs := len(burnTransactions)
				lenProofs := len(a.Proofs)
				if lenProofs != lenBTxs {
					return nil, errors.Errorf("invalid proofs count: expected %v, got %v", lenBTxs, lenProofs)
				}
				var bTxs = make([]BurnFungibleToken, lenBTxs)
				for i, btx := range burnTransactions {
					genericBurnTx, err := NewGenericTx(btx)
					if err != nil {
						return nil, errors.Errorf("burn transaction with index %v is invalid: %v", i, err)
					}
					bTxs[i] = genericBurnTx.(*burnFungibleTokenWrapper)
				}

				return &joinFungibleTokenWrapper{
					wrapper:          wrapper{transaction: tx},
					attributes:       a,
					burnTransactions: bTxs,
				}, nil
			},
			tx)
	default:
		return nil, errors.Errorf("unknown transaction type %s", tx.TransactionAttributes.TypeUrl)
	}
}

// convertToWrapper converts given tx to a generic transaction. attrType is the type of the tx attributes. createGenericTxFunc creates an instance of given generic transaction.
func convertToWrapper[A proto.Message, G txsystem.GenericTransaction](attrType A, createGenericTxFunc func(a A) (G, error), tx *txsystem.Transaction) (g G, err error) {
	err = tx.TransactionAttributes.UnmarshalTo(attrType)
	if err != nil {
		return g, errors.Wrapf(err, "invalid tx attributes")
	}
	return createGenericTxFunc(attrType)
}

func (w *wrapper) UnitID() *uint256.Int              { return uint256.NewInt(0).SetBytes(w.transaction.UnitId) }
func (w *wrapper) Timeout() uint64                   { return w.transaction.Timeout }
func (w *wrapper) SystemID() []byte                  { return w.transaction.SystemId }
func (w *wrapper) OwnerProof() []byte                { return w.transaction.OwnerProof }
func (w *wrapper) ToProtoBuf() *txsystem.Transaction { return w.transaction }

func (w *wrapper) sigBytes(b *bytes.Buffer) {
	b.Write(w.transaction.SystemId)
	b.Write(w.transaction.UnitId)
	b.Write(util.Uint64ToBytes(w.transaction.Timeout))
}

func (w *wrapper) hashComputed(hashFunc crypto.Hash) bool {
	return w.hashFunc == hashFunc && w.hashValue != nil
}

func (w *wrapper) addTransactionFieldsToHasher(hasher hash.Hash) {
	hasher.Write(w.transaction.SystemId)
	hasher.Write(w.transaction.UnitId)
	hasher.Write(w.transaction.OwnerProof)
	hasher.Write(util.Uint64ToBytes(w.transaction.Timeout))
}

func (w *wrapper) IsPrimary() bool {
	return true
}

func (c *createNonFungibleTokenTypeWrapper) parentTypeIdInt() *uint256.Int {
	return uint256.NewInt(0).SetBytes(c.attributes.ParentTypeId)
}

func (c *createNonFungibleTokenTypeWrapper) ParentTypeId() []byte {
	return c.attributes.ParentTypeId
}

func (c *createNonFungibleTokenTypeWrapper) Hash(hashFunc crypto.Hash) []byte {
	if c.wrapper.hashComputed(hashFunc) {
		return c.wrapper.hashValue
	}
	hasher := hashFunc.New()
	c.AddToHasher(hasher)
	c.wrapper.hashValue = hasher.Sum(nil)
	c.wrapper.hashFunc = hashFunc
	return c.wrapper.hashValue
}

func (c *createNonFungibleTokenTypeWrapper) SigBytes() []byte {
	var b bytes.Buffer
	c.wrapper.sigBytes(&b)
	b.Write([]byte(c.Symbol()))
	b.Write(c.ParentTypeId())
	b.Write(c.SubTypeCreationPredicate())
	b.Write(c.TokenCreationPredicate())
	b.Write(c.InvariantPredicate())
	b.Write(c.DataUpdatePredicate())
	return b.Bytes()
}

func (c *createNonFungibleTokenTypeWrapper) AddToHasher(hasher hash.Hash) {
	c.wrapper.addTransactionFieldsToHasher(hasher)
	hasher.Write([]byte(c.Symbol()))
	hasher.Write(c.ParentTypeId())
	hasher.Write(c.SubTypeCreationPredicate())
	hasher.Write(c.TokenCreationPredicate())
	hasher.Write(c.InvariantPredicate())
	hasher.Write(c.DataUpdatePredicate())
	hasher.Write(c.SubTypeCreationPredicateSignature())
}

func (c *createNonFungibleTokenTypeWrapper) Symbol() string {
	return c.attributes.Symbol
}

func (c *createNonFungibleTokenTypeWrapper) SubTypeCreationPredicate() []byte {
	return c.attributes.SubTypeCreationPredicate
}

func (c *createNonFungibleTokenTypeWrapper) TokenCreationPredicate() []byte {
	return c.attributes.TokenCreationPredicate
}

func (c *createNonFungibleTokenTypeWrapper) InvariantPredicate() []byte {
	return c.attributes.InvariantPredicate
}

func (c *createNonFungibleTokenTypeWrapper) DataUpdatePredicate() []byte {
	return c.attributes.DataUpdatePredicate
}

func (c *createNonFungibleTokenTypeWrapper) SubTypeCreationPredicateSignature() []byte {
	return c.attributes.SubTypeCreationPredicateSignature
}

func (c *mintNonFungibleTokenWrapper) NFTTypeIdInt() *uint256.Int {
	return uint256.NewInt(0).SetBytes(c.NFTTypeId())
}

func (c *mintNonFungibleTokenWrapper) NFTTypeId() []byte {
	return c.attributes.NftType
}

func (c *mintNonFungibleTokenWrapper) Hash(hashFunc crypto.Hash) []byte {
	if c.wrapper.hashComputed(hashFunc) {
		return c.wrapper.hashValue
	}
	hasher := hashFunc.New()
	c.AddToHasher(hasher)
	c.wrapper.hashValue = hasher.Sum(nil)
	c.wrapper.hashFunc = hashFunc
	return c.wrapper.hashValue
}

func (c *mintNonFungibleTokenWrapper) SigBytes() []byte {
	var b bytes.Buffer
	c.wrapper.sigBytes(&b)
	b.Write(c.Bearer())
	b.Write(c.NFTTypeId())
	b.Write([]byte(c.URI()))
	b.Write(c.Data())
	b.Write(c.DataUpdatePredicate())
	return b.Bytes()
}

func (c *mintNonFungibleTokenWrapper) AddToHasher(hasher hash.Hash) {
	c.wrapper.addTransactionFieldsToHasher(hasher)
	hasher.Write(c.Bearer())
	hasher.Write(c.NFTTypeId())
	hasher.Write([]byte(c.URI()))
	hasher.Write(c.Data())
	hasher.Write(c.DataUpdatePredicate())
}

func (c *mintNonFungibleTokenWrapper) Bearer() []byte {
	return c.attributes.Bearer
}

func (c *mintNonFungibleTokenWrapper) URI() string {
	return c.attributes.Uri
}

func (c *mintNonFungibleTokenWrapper) Data() []byte {
	return c.attributes.Data
}

func (c *mintNonFungibleTokenWrapper) DataUpdatePredicate() []byte {
	return c.attributes.DataUpdatePredicate
}

func (c *mintNonFungibleTokenWrapper) TokenCreationPredicateSignature() []byte {
	return c.attributes.TokenCreationPredicateSignature
}

func (t *transferNonFungibleTokenWrapper) Hash(hashFunc crypto.Hash) []byte {
	if t.wrapper.hashComputed(hashFunc) {
		return t.wrapper.hashValue
	}
	hasher := hashFunc.New()
	t.AddToHasher(hasher)
	t.wrapper.hashValue = hasher.Sum(nil)
	t.wrapper.hashFunc = hashFunc
	return t.wrapper.hashValue
}

func (t *transferNonFungibleTokenWrapper) SigBytes() []byte {
	var b bytes.Buffer
	t.wrapper.sigBytes(&b)
	b.Write(t.NewBearer())
	b.Write(t.Nonce())
	b.Write(t.Backlink())
	return b.Bytes()
}

func (t *transferNonFungibleTokenWrapper) AddToHasher(hasher hash.Hash) {
	t.wrapper.addTransactionFieldsToHasher(hasher)
	hasher.Write(t.NewBearer())
	hasher.Write(t.Nonce())
	hasher.Write(t.Backlink())
	hasher.Write(t.InvariantPredicateSignature())
}

func (t *transferNonFungibleTokenWrapper) NewBearer() []byte {
	return t.attributes.NewBearer
}

func (t *transferNonFungibleTokenWrapper) Nonce() []byte {
	return t.attributes.Nonce
}

func (t *transferNonFungibleTokenWrapper) Backlink() []byte {
	return t.attributes.Backlink
}

func (t *transferNonFungibleTokenWrapper) InvariantPredicateSignature() []byte {
	return t.attributes.InvariantPredicateSignature
}

func (u *updateNonFungibleTokenWrapper) Hash(hashFunc crypto.Hash) []byte {
	if u.wrapper.hashComputed(hashFunc) {
		return u.wrapper.hashValue
	}
	hasher := hashFunc.New()
	u.AddToHasher(hasher)
	u.wrapper.hashValue = hasher.Sum(nil)
	u.wrapper.hashFunc = hashFunc
	return u.wrapper.hashValue
}

func (u *updateNonFungibleTokenWrapper) SigBytes() []byte {
	var b bytes.Buffer
	u.wrapper.sigBytes(&b)
	b.Write(u.Data())
	b.Write(u.Backlink())
	return b.Bytes()
}

func (u *updateNonFungibleTokenWrapper) AddToHasher(hasher hash.Hash) {
	u.wrapper.addTransactionFieldsToHasher(hasher)
	hasher.Write(u.Data())
	hasher.Write(u.Backlink())
	hasher.Write(u.DataUpdateSignature())
}

func (u *updateNonFungibleTokenWrapper) Data() []byte {
	return u.attributes.Data
}

func (u *updateNonFungibleTokenWrapper) Backlink() []byte {
	return u.attributes.Backlink
}

func (u *updateNonFungibleTokenWrapper) DataUpdateSignature() []byte {
	return u.attributes.DataUpdateSignature
}

func (c *createFungibleTokenTypeWrapper) Hash(hashFunc crypto.Hash) []byte {
	if c.wrapper.hashComputed(hashFunc) {
		return c.wrapper.hashValue
	}
	hasher := hashFunc.New()
	c.AddToHasher(hasher)
	c.wrapper.hashValue = hasher.Sum(nil)
	c.wrapper.hashFunc = hashFunc
	return c.wrapper.hashValue
}

func (c *createFungibleTokenTypeWrapper) SigBytes() []byte {
	var b bytes.Buffer
	c.wrapper.sigBytes(&b)
	b.Write([]byte(c.Symbol()))
	b.Write(c.ParentTypeId())
	b.Write(util.Uint32ToBytes(c.DecimalPlaces()))
	b.Write(c.SubTypeCreationPredicate())
	b.Write(c.TokenCreationPredicate())
	b.Write(c.InvariantPredicate())
	return b.Bytes()
}

func (c *createFungibleTokenTypeWrapper) ParentTypeIdInt() *uint256.Int {
	return uint256.NewInt(0).SetBytes(c.attributes.ParentTypeId)
}

func (c *createFungibleTokenTypeWrapper) AddToHasher(hasher hash.Hash) {
	c.wrapper.addTransactionFieldsToHasher(hasher)
	hasher.Write([]byte(c.Symbol()))
	hasher.Write(c.ParentTypeId())
	hasher.Write(util.Uint32ToBytes(c.DecimalPlaces()))
	hasher.Write(c.SubTypeCreationPredicate())
	hasher.Write(c.TokenCreationPredicate())
	hasher.Write(c.InvariantPredicate())
	hasher.Write(c.SubTypeCreationPredicateSignature())
}

func (c *createFungibleTokenTypeWrapper) ParentTypeId() []byte {
	return c.attributes.ParentTypeId
}

func (c *createFungibleTokenTypeWrapper) DecimalPlaces() uint32 {
	return c.attributes.DecimalPlaces
}

func (c *createFungibleTokenTypeWrapper) Symbol() string {
	return c.attributes.Symbol
}

func (c *createFungibleTokenTypeWrapper) SubTypeCreationPredicate() []byte {
	return c.attributes.SubTypeCreationPredicate
}

func (c *createFungibleTokenTypeWrapper) TokenCreationPredicate() []byte {
	return c.attributes.TokenCreationPredicate
}

func (c *createFungibleTokenTypeWrapper) InvariantPredicate() []byte {
	return c.attributes.InvariantPredicate
}

func (c *createFungibleTokenTypeWrapper) SubTypeCreationPredicateSignature() []byte {
	return c.attributes.SubTypeCreationPredicateSignature
}

func (m *mintFungibleTokenWrapper) Hash(hashFunc crypto.Hash) []byte {
	if m.wrapper.hashComputed(hashFunc) {
		return m.wrapper.hashValue
	}
	hasher := hashFunc.New()
	m.AddToHasher(hasher)
	m.wrapper.hashValue = hasher.Sum(nil)
	m.wrapper.hashFunc = hashFunc
	return m.wrapper.hashValue
}

func (m *mintFungibleTokenWrapper) TypeIdInt() *uint256.Int {
	return uint256.NewInt(0).SetBytes(m.TypeId())
}

func (m *mintFungibleTokenWrapper) TypeId() []byte {
	return m.attributes.Type
}

func (m *mintFungibleTokenWrapper) Value() uint64 {
	return m.attributes.Value
}

func (m *mintFungibleTokenWrapper) Bearer() []byte {
	return m.attributes.Bearer
}

func (m *mintFungibleTokenWrapper) TokenCreationPredicateSignature() []byte {
	return m.attributes.TokenCreationPredicateSignature
}

func (m *mintFungibleTokenWrapper) SigBytes() []byte {
	var b bytes.Buffer
	m.wrapper.sigBytes(&b)
	b.Write(m.Bearer())
	b.Write(m.TypeId())
	b.Write(util.Uint64ToBytes(m.Value()))
	return b.Bytes()
}

func (m *mintFungibleTokenWrapper) AddToHasher(hasher hash.Hash) {
	m.wrapper.addTransactionFieldsToHasher(hasher)
	hasher.Write(m.Bearer())
	hasher.Write(m.TypeId())
	hasher.Write(util.Uint64ToBytes(m.Value()))
	hasher.Write(m.TokenCreationPredicateSignature())
}

func (t *transferFungibleTokenWrapper) NewBearer() []byte {
	return t.attributes.NewBearer
}

func (t *transferFungibleTokenWrapper) Value() uint64 {
	return t.attributes.Value
}

func (t *transferFungibleTokenWrapper) Nonce() []byte {
	return t.attributes.Nonce
}

func (t *transferFungibleTokenWrapper) Backlink() []byte {
	return t.attributes.Backlink
}

func (t *transferFungibleTokenWrapper) InvariantPredicateSignature() []byte {
	return t.attributes.InvariantPredicateSignature
}

func (t *transferFungibleTokenWrapper) Hash(hashFunc crypto.Hash) []byte {
	if t.wrapper.hashComputed(hashFunc) {
		return t.wrapper.hashValue
	}
	hasher := hashFunc.New()
	t.AddToHasher(hasher)
	t.wrapper.hashValue = hasher.Sum(nil)
	t.wrapper.hashFunc = hashFunc
	return t.wrapper.hashValue
}

func (t *transferFungibleTokenWrapper) SigBytes() []byte {
	var b bytes.Buffer
	t.wrapper.sigBytes(&b)
	b.Write(t.NewBearer())
	b.Write(util.Uint64ToBytes(t.Value()))
	b.Write(t.Nonce())
	b.Write(t.Backlink())
	return b.Bytes()
}

func (t *transferFungibleTokenWrapper) AddToHasher(hasher hash.Hash) {
	t.wrapper.addTransactionFieldsToHasher(hasher)
	hasher.Write(t.NewBearer())
	hasher.Write(util.Uint64ToBytes(t.Value()))
	hasher.Write(t.Nonce())
	hasher.Write(t.Backlink())
	hasher.Write(t.InvariantPredicateSignature())
}

func (s *splitFungibleTokenWrapper) Hash(hashFunc crypto.Hash) []byte {
	if s.wrapper.hashComputed(hashFunc) {
		return s.wrapper.hashValue
	}
	hasher := hashFunc.New()
	s.AddToHasher(hasher)
	s.wrapper.hashValue = hasher.Sum(nil)
	s.wrapper.hashFunc = hashFunc
	return s.wrapper.hashValue
}

func (s *splitFungibleTokenWrapper) HashForIdCalculation(hashFunc crypto.Hash) []byte {
	hasher := hashFunc.New()
	idBytes := s.UnitID().Bytes32()
	hasher.Write(idBytes[:])
	s.addAttributesToHasher(hasher)
	hasher.Write(util.Uint64ToBytes(s.Timeout()))
	return hasher.Sum(nil)
}

func (s *splitFungibleTokenWrapper) addAttributesToHasher(hasher hash.Hash) {
	hasher.Write(s.NewBearer())
	hasher.Write(util.Uint64ToBytes(s.TargetValue()))
	hasher.Write(s.Nonce())
	hasher.Write(s.Backlink())
	hasher.Write(s.InvariantPredicateSignature())
}

func (s *splitFungibleTokenWrapper) SigBytes() []byte {
	var b bytes.Buffer
	s.wrapper.sigBytes(&b)
	b.Write(s.NewBearer())
	b.Write(util.Uint64ToBytes(s.TargetValue()))
	b.Write(s.Nonce())
	b.Write(s.Backlink())
	return b.Bytes()
}

func (s *splitFungibleTokenWrapper) AddToHasher(hasher hash.Hash) {
	s.wrapper.addTransactionFieldsToHasher(hasher)
	s.addAttributesToHasher(hasher)
}

func (s *splitFungibleTokenWrapper) NewBearer() []byte {
	return s.attributes.NewBearer
}

func (s *splitFungibleTokenWrapper) TargetValue() uint64 {
	return s.attributes.TargetValue
}

func (s *splitFungibleTokenWrapper) Nonce() []byte {
	return s.attributes.Nonce
}

func (s *splitFungibleTokenWrapper) Backlink() []byte {
	return s.attributes.Backlink
}

func (s *splitFungibleTokenWrapper) InvariantPredicateSignature() []byte {
	return s.attributes.InvariantPredicateSignature
}

func (bw *burnFungibleTokenWrapper) Hash(hashFunc crypto.Hash) []byte {
	if bw.wrapper.hashComputed(hashFunc) {
		return bw.wrapper.hashValue
	}
	hasher := hashFunc.New()
	bw.AddToHasher(hasher)
	bw.wrapper.hashValue = hasher.Sum(nil)
	bw.wrapper.hashFunc = hashFunc
	return bw.wrapper.hashValue
}

func (bw *burnFungibleTokenWrapper) SigBytes() []byte {
	var b bytes.Buffer
	bw.wrapper.sigBytes(&b)
	b.Write(bw.TypeId())
	b.Write(util.Uint64ToBytes(bw.Value()))
	b.Write(bw.Nonce())
	b.Write(bw.Backlink())
	return b.Bytes()
}

func (bw *burnFungibleTokenWrapper) AddToHasher(hasher hash.Hash) {
	bw.wrapper.addTransactionFieldsToHasher(hasher)
	hasher.Write(bw.TypeId())
	hasher.Write(util.Uint64ToBytes(bw.Value()))
	hasher.Write(bw.Nonce())
	hasher.Write(bw.Backlink())
	hasher.Write(bw.InvariantPredicateSignature())
}

func (bw *burnFungibleTokenWrapper) TypeId() []byte {
	return bw.attributes.Type
}

func (bw *burnFungibleTokenWrapper) Value() uint64 {
	return bw.attributes.Value
}

func (bw *burnFungibleTokenWrapper) Nonce() []byte {
	return bw.attributes.Nonce
}

func (bw *burnFungibleTokenWrapper) Backlink() []byte {
	return bw.attributes.Backlink
}

func (bw *burnFungibleTokenWrapper) InvariantPredicateSignature() []byte {
	return bw.attributes.InvariantPredicateSignature
}

func (jw *joinFungibleTokenWrapper) Hash(hashFunc crypto.Hash) []byte {
	if jw.wrapper.hashComputed(hashFunc) {
		return jw.wrapper.hashValue
	}
	hasher := hashFunc.New()
	jw.AddToHasher(hasher)
	jw.wrapper.hashValue = hasher.Sum(nil)
	jw.wrapper.hashFunc = hashFunc
	return jw.wrapper.hashValue
}

func (jw *joinFungibleTokenWrapper) AddToHasher(hasher hash.Hash) {
	for _, tx := range jw.burnTransactions {
		tx.AddToHasher(hasher)
	}
	for _, proof := range jw.BlockProofs() {
		proof.AddToHasher(hasher)
	}
	hasher.Write(jw.Backlink())
	hasher.Write(jw.InvariantPredicateSignature())
}

func (jw *joinFungibleTokenWrapper) SigBytes() []byte {
	var b bytes.Buffer
	jw.wrapper.sigBytes(&b)
	for _, tx := range jw.burnTransactions {
		b.Write(tx.SigBytes())
		b.Write(tx.OwnerProof())
	}
	for _, proof := range jw.BlockProofs() {
		b.Write(proof.Bytes())
	}
	b.Write(jw.Backlink())
	return b.Bytes()
}

func (jw *joinFungibleTokenWrapper) BurnTransactions() []BurnFungibleToken {
	return jw.burnTransactions
}

func (jw *joinFungibleTokenWrapper) BlockProofs() []*block.BlockProof {
	return jw.attributes.Proofs
}

func (jw *joinFungibleTokenWrapper) Backlink() []byte {
	return jw.attributes.Backlink
}

func (jw *joinFungibleTokenWrapper) InvariantPredicateSignature() []byte {
	return jw.attributes.InvariantPredicateSignature
}
