package tokens

import (
	"crypto"

	"github.com/alphabill-org/alphabill-go-base/txsystem/tokens"
	"github.com/alphabill-org/alphabill/predicates"
	"github.com/alphabill-org/alphabill/state"
	txtypes "github.com/alphabill-org/alphabill/txsystem/types"
)

var _ txtypes.Module = (*NonFungibleTokensModule)(nil)

type NonFungibleTokensModule struct {
	state         *state.State
	hashAlgorithm crypto.Hash
	execPredicate predicates.PredicateRunner
}

func NewNonFungibleTokensModule(options *Options) (*NonFungibleTokensModule, error) {
	return &NonFungibleTokensModule{
		state:         options.state,
		hashAlgorithm: options.hashAlgorithm,
		execPredicate: predicates.NewPredicateRunner(options.exec),
	}, nil
}

func (n *NonFungibleTokensModule) TxHandlers() map[string]txtypes.TxExecutor {
	return map[string]txtypes.TxExecutor{
		tokens.PayloadTypeDefineNFT:   txtypes.NewTxHandler[tokens.DefineNonFungibleTokenAttributes, tokens.DefineNonFungibleTokenAuthProof](n.validateDefineNFT, n.executeDefineNFT),
		tokens.PayloadTypeMintNFT:     txtypes.NewTxHandler[tokens.MintNonFungibleTokenAttributes, tokens.MintNonFungibleTokenAuthProof](n.validateMintNFT, n.executeMintNFT),
		tokens.PayloadTypeTransferNFT: txtypes.NewTxHandler[tokens.TransferNonFungibleTokenAttributes, tokens.TransferNonFungibleTokenAuthProof](n.validateTransferNFT, n.executeTransferNFT),
		tokens.PayloadTypeUpdateNFT:   txtypes.NewTxHandler[tokens.UpdateNonFungibleTokenAttributes, tokens.UpdateNonFungibleTokenAuthProof](n.validateUpdateNFT, n.executeUpdateNFT),
	}
}
