package wallet

import (
	"bytes"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/crypto"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/errors"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/hash"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/rpc/transaction"
	"gitdc.ee.guardtime.com/alphabill/alphabill/internal/script"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"sort"
)

func createTransaction(pubKey []byte, k *accountKey, amount uint64, b *bill, timeout uint64) (*transaction.Transaction, error) {
	var err error
	var tx *transaction.Transaction
	if b.Value == amount {
		tx, err = createTransferTx(pubKey, k, b, timeout)
	} else {
		tx, err = createSplitTx(amount, pubKey, k, b, timeout)
	}
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func createTransferTx(pubKey []byte, k *accountKey, bill *bill, timeout uint64) (*transaction.Transaction, error) {
	tx := &transaction.Transaction{
		UnitId:                bill.getId(),
		TransactionAttributes: new(anypb.Any),
		Timeout:               timeout,
	}
	err := anypb.MarshalFrom(tx.TransactionAttributes, &transaction.BillTransfer{
		NewBearer:   script.PredicatePayToPublicKeyHashDefault(hash.Sum256(pubKey)),
		TargetValue: bill.Value,
		Backlink:    bill.TxHash,
	}, proto.MarshalOptions{})
	if err != nil {
		return nil, err
	}
	err = signTx(tx, k)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func createSplitTx(amount uint64, pubKey []byte, k *accountKey, bill *bill, timeout uint64) (*transaction.Transaction, error) {
	tx := &transaction.Transaction{
		UnitId:                bill.getId(),
		TransactionAttributes: new(anypb.Any),
		Timeout:               timeout,
	}
	err := anypb.MarshalFrom(tx.TransactionAttributes, &transaction.BillSplit{
		Amount:         bill.Value,
		TargetBearer:   script.PredicatePayToPublicKeyHashDefault(hash.Sum256(pubKey)),
		RemainingValue: bill.Value - amount,
		Backlink:       bill.TxHash,
	}, proto.MarshalOptions{})
	if err != nil {
		return nil, err
	}
	err = signTx(tx, k)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func createDustTx(k *accountKey, bill *bill, nonce []byte, timeout uint64) (*transaction.Transaction, error) {
	tx := &transaction.Transaction{
		UnitId:                bill.getId(),
		TransactionAttributes: new(anypb.Any),
		Timeout:               timeout,
	}
	err := anypb.MarshalFrom(tx.TransactionAttributes, &transaction.TransferDC{
		TargetValue:  bill.Value,
		TargetBearer: script.PredicatePayToPublicKeyHashDefault(k.PubKeyHashSha256),
		Backlink:     bill.TxHash,
		Nonce:        nonce,
	}, proto.MarshalOptions{})
	if err != nil {
		return nil, err
	}
	err = signTx(tx, k)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func createSwapTx(k *accountKey, dcBills []*bill, dcNonce []byte, timeout uint64) (*transaction.Transaction, error) {
	if len(dcBills) == 0 {
		return nil, errors.New("cannot create swap transaction as no dust bills exist")
	}
	// sort bills by ids in ascending order
	sort.Slice(dcBills, func(i, j int) bool {
		return bytes.Compare(dcBills[i].getId(), dcBills[j].getId()) < 0
	})

	var billIds [][]byte
	var dustTransferProofs [][]byte
	var dustTransferOrders []*transaction.Transaction
	var billValueSum uint64
	for _, b := range dcBills {
		billIds = append(billIds, b.getId())
		dustTransferOrders = append(dustTransferOrders, b.DcTx)
		// TODO add DC proofs: https://guardtime.atlassian.net/browse/AB-99
		dustTransferProofs = append(dustTransferProofs, nil)
		billValueSum += b.Value
	}

	swapTx := &transaction.Transaction{
		UnitId:                dcNonce,
		TransactionAttributes: new(anypb.Any),
		Timeout:               timeout,
	}
	err := anypb.MarshalFrom(swapTx.TransactionAttributes, &transaction.Swap{
		OwnerCondition:  script.PredicatePayToPublicKeyHashDefault(k.PubKeyHashSha256),
		BillIdentifiers: billIds,
		DcTransfers:     dustTransferOrders,
		Proofs:          dustTransferProofs,
		TargetValue:     billValueSum,
	}, proto.MarshalOptions{})
	if err != nil {
		return nil, err
	}
	err = signTx(swapTx, k)
	if err != nil {
		return nil, err
	}
	return swapTx, nil
}

func signTx(tx *transaction.Transaction, ac *accountKey) error {
	signer, err := crypto.NewInMemorySecp256K1SignerFromKey(ac.PrivKey)
	if err != nil {
		return err
	}
	gtx, err := transaction.NewMoneyTx(tx)
	if err != nil {
		return err
	}
	sig, err := signer.SignBytes(gtx.SigBytes())
	if err != nil {
		return err
	}
	tx.OwnerProof = script.PredicateArgumentPayToPublicKeyHashDefault(sig, ac.PubKey)
	return nil
}