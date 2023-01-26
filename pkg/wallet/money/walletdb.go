package money

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/alphabill-org/alphabill/internal/util"
	"github.com/alphabill-org/alphabill/pkg/wallet/account"
	"github.com/alphabill-org/alphabill/pkg/wallet/log"
	"github.com/holiman/uint256"
	bolt "go.etcd.io/bbolt"
)

var (
	accountsBucket      = []byte("accounts")
	accountBillsBucket  = []byte("accountBills")
	accountDcMetaBucket = []byte("accountDcMeta")
	metaBucket          = []byte("meta")
)

var (
	blockHeightKeyName = []byte("blockHeightKey")
)

var (
	errWalletDbAlreadyExists = errors.New("wallet db already exists")
	errWalletDbDoesNotExists = errors.New("cannot open wallet db, file does not exist")
	errBillNotFound          = errors.New("bill does not exist")
	errAccountNotFound       = errors.New("account does not exist")
)

const WalletFileName = "wallet.db"

type Db interface {
	Do() TxContext
	WithTransaction(func(tx TxContext) error) error
	Close()
	DeleteDb()
}

type TxContext interface {
	GetBlockNumber() (uint64, error)
	SetBlockNumber(blockNumber uint64) error

	GetBill(accountIndex uint64, id []byte) (*Bill, error)
	SetBill(accountIndex uint64, bill *Bill) error
	ContainsBill(accountIndex uint64, id *uint256.Int) (bool, error)
	RemoveBill(accountIndex uint64, id *uint256.Int) error
	GetBills(accountIndex uint64) ([]*Bill, error)
	GetAllBills(am account.Manager) ([][]*Bill, error)
	GetBalance(cmd GetBalanceCmd) (uint64, error)
	GetBalances(cmd GetBalanceCmd) ([]uint64, error)

	GetDcMetadataMap(accountIndex uint64) (map[uint256.Int]*dcMetadata, error)
	GetDcMetadata(accountIndex uint64, nonce []byte) (*dcMetadata, error)
	SetDcMetadata(accountIndex uint64, nonce []byte, dcMetadata *dcMetadata) error
}

type wdb struct {
	db         *bolt.DB
	dbFilePath string
}

type wdbtx struct {
	wdb *wdb
	tx  *bolt.Tx
}

func OpenDb(config WalletConfig) (*wdb, error) {
	walletDir, err := config.GetWalletDir()
	if err != nil {
		return nil, err
	}
	dbFilePath := path.Join(walletDir, WalletFileName)
	return openDb(dbFilePath, false)
}

func (w *wdbtx) GetBill(accountIndex uint64, billId []byte) (*Bill, error) {
	var b *Bill
	err := w.withTx(w.tx, func(tx *bolt.Tx) error {
		bkt, err := getAccountBucket(tx, util.Uint64ToBytes(accountIndex))
		if err != nil {
			return err
		}
		billBytes := bkt.Bucket(accountBillsBucket).Get(billId)
		if billBytes == nil {
			return errBillNotFound
		}
		b, err = parseBill(billBytes)
		if err != nil {
			return err
		}
		return nil
	}, false)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (w *wdbtx) SetBill(accountIndex uint64, bill *Bill) error {
	return w.withTx(w.tx, func(tx *bolt.Tx) error {
		val, err := json.Marshal(bill)
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("adding bill: value=%d id=%s, for account=%d", bill.Value, bill.Id.String(), accountIndex))
		bkt, err := getAccountBucket(tx, util.Uint64ToBytes(accountIndex))
		if err != nil {
			return err
		}
		return bkt.Bucket(accountBillsBucket).Put(bill.GetID(), val)
	}, true)
}

func (w *wdbtx) ContainsBill(accountIndex uint64, id *uint256.Int) (bool, error) {
	var res bool
	err := w.withTx(w.tx, func(tx *bolt.Tx) error {
		billId := id.Bytes32()
		bkt, err := getAccountBucket(tx, util.Uint64ToBytes(accountIndex))
		if err != nil {
			return err
		}
		res = bkt.Bucket(accountBillsBucket).Get(billId[:]) != nil
		return nil
	}, false)
	if err != nil {
		return false, err
	}
	return res, nil
}

func (w *wdbtx) GetBills(accountIndex uint64) ([]*Bill, error) {
	var res []*Bill
	err := w.withTx(w.tx, func(tx *bolt.Tx) error {
		bkt, err := getAccountBucket(tx, util.Uint64ToBytes(accountIndex))
		if err != nil {
			return err
		}
		c := bkt.Bucket(accountBillsBucket).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			b, err := parseBill(v)
			if err != nil {
				return err
			}
			res = append(res, b)
		}
		return nil
	}, false)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (w *wdbtx) GetAllBills(am account.Manager) ([][]*Bill, error) {
	var res [][]*Bill
	err := w.withTx(w.tx, func(tx *bolt.Tx) error {
		maxAccountIndex, err := am.GetMaxAccountIndex()
		if err != nil {
			return err
		}
		for accountIndex := uint64(0); accountIndex <= maxAccountIndex; accountIndex++ {
			accountBills, err := w.GetBills(accountIndex)
			if err != nil {
				return err
			}
			res = append(res, accountBills)
		}
		return nil
	}, false)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (w *wdbtx) RemoveBill(accountIndex uint64, id *uint256.Int) error {
	return w.withTx(w.tx, func(tx *bolt.Tx) error {
		bytes32 := id.Bytes32()
		bkt, err := getAccountBucket(tx, util.Uint64ToBytes(accountIndex))
		if err != nil {
			return err
		}
		return bkt.Bucket(accountBillsBucket).Delete(bytes32[:])
	}, true)
}

func (w *wdbtx) GetBalance(cmd GetBalanceCmd) (uint64, error) {
	sum := uint64(0)
	err := w.withTx(w.tx, func(tx *bolt.Tx) error {
		bkt, err := getAccountBucket(tx, util.Uint64ToBytes(cmd.AccountIndex))
		if err != nil {
			return err
		}
		return bkt.Bucket(accountBillsBucket).ForEach(func(k, v []byte) error {
			var b *Bill
			err := json.Unmarshal(v, &b)
			if err != nil {
				return err
			}
			if b.IsDcBill && !cmd.CountDCBills {
				return nil
			}
			sum += b.Value
			return nil
		})
	}, false)
	if err != nil {
		return 0, err
	}
	return sum, nil
}

func (w *wdbtx) GetBalances(cmd GetBalanceCmd) ([]uint64, error) {
	res := make(map[uint64]uint64)
	err := w.withTx(w.tx, func(tx *bolt.Tx) error {
		return tx.Bucket(accountsBucket).ForEach(func(accIdx, v []byte) error {
			if v != nil { // value is nil if entry is a bucket
				return nil
			}
			sum := uint64(0)
			accountBucket, err := getAccountBucket(tx, accIdx)
			if err != nil {
				return err
			}
			accBillsBucket := accountBucket.Bucket(accountBillsBucket)
			err = accBillsBucket.ForEach(func(billId, billValue []byte) error {
				var b *Bill
				err := json.Unmarshal(billValue, &b)
				if err != nil {
					return err
				}
				if b.IsDcBill && !cmd.CountDCBills {
					return nil
				}
				sum += b.Value
				return nil
			})
			if err != nil {
				return err
			}
			res[util.BytesToUint64(accIdx)] = sum
			return nil
		})
	}, false)
	if err != nil {
		return nil, err
	}
	balances := make([]uint64, len(res))
	for accIdx, sum := range res {
		balances[accIdx] = sum
	}
	return balances, nil
}

func (w *wdbtx) GetBlockNumber() (uint64, error) {
	var res uint64
	err := w.withTx(w.tx, func(tx *bolt.Tx) error {
		blockHeightBytes := tx.Bucket(metaBucket).Get(blockHeightKeyName)
		if blockHeightBytes == nil {
			return nil
		}
		res = util.BytesToUint64(blockHeightBytes)
		return nil
	}, false)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (w *wdbtx) SetBlockNumber(blockHeight uint64) error {
	return w.withTx(w.tx, func(tx *bolt.Tx) error {
		return tx.Bucket(metaBucket).Put(blockHeightKeyName, util.Uint64ToBytes(blockHeight))
	}, true)
}

func (w *wdbtx) GetDcMetadataMap(accountIndex uint64) (map[uint256.Int]*dcMetadata, error) {
	res := map[uint256.Int]*dcMetadata{}
	err := w.withTx(w.tx, func(tx *bolt.Tx) error {
		accountBucket, err := getAccountBucket(tx, util.Uint64ToBytes(accountIndex))
		if err != nil {
			return err
		}
		return accountBucket.Bucket(accountDcMetaBucket).ForEach(func(k, v []byte) error {
			var m *dcMetadata
			err := json.Unmarshal(v, &m)
			if err != nil {
				return err
			}
			res[*uint256.NewInt(0).SetBytes(k)] = m
			return nil
		})
	}, false)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (w *wdbtx) GetDcMetadata(accountIndex uint64, nonce []byte) (*dcMetadata, error) {
	var res *dcMetadata
	err := w.withTx(w.tx, func(tx *bolt.Tx) error {
		accountBucket, err := getAccountBucket(tx, util.Uint64ToBytes(accountIndex))
		if err != nil {
			return err
		}
		m := accountBucket.Bucket(accountDcMetaBucket).Get(nonce)
		if m != nil {
			return json.Unmarshal(m, &res)
		}
		return nil
	}, false)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (w *wdbtx) SetDcMetadata(accountIndex uint64, dcNonce []byte, dcMetadata *dcMetadata) error {
	return w.withTx(w.tx, func(tx *bolt.Tx) error {
		accountBucket, err := getAccountBucket(tx, util.Uint64ToBytes(accountIndex))
		if err != nil {
			return err
		}
		if dcMetadata != nil {
			val, err := json.Marshal(dcMetadata)
			if err != nil {
				return err
			}
			return accountBucket.Bucket(accountDcMetaBucket).Put(dcNonce, val)
		}
		return accountBucket.Bucket(accountDcMetaBucket).Delete(dcNonce)
	}, true)
}

func (w *wdb) DeleteDb() {
	if w.db == nil {
		return
	}
	errClose := w.db.Close()
	if errClose != nil {
		log.Warning("error closing db: ", errClose)
	}
	errRemove := os.Remove(w.dbFilePath)
	if errRemove != nil {
		log.Warning("error removing db: ", errRemove)
	}
}

func (w *wdb) WithTransaction(fn func(txc TxContext) error) error {
	return w.db.Update(func(tx *bolt.Tx) error {
		return fn(&wdbtx{wdb: w, tx: tx})
	})
}

func (w *wdb) Do() TxContext {
	return &wdbtx{wdb: w, tx: nil}
}

func (w *wdb) Path() string {
	return w.dbFilePath
}

func (w *wdb) Close() {
	if w.db == nil {
		return
	}
	log.Info("closing wallet db")
	err := w.db.Close()
	if err != nil {
		log.Warning("error closing db: ", err)
	}
}

func (w *wdbtx) withTx(dbTx *bolt.Tx, myFunc func(tx *bolt.Tx) error, writeTx bool) error {
	if dbTx != nil {
		return myFunc(dbTx)
	} else if writeTx {
		return w.wdb.db.Update(myFunc)
	} else {
		return w.wdb.db.View(myFunc)
	}
}

func (w *wdb) createBuckets() error {
	return w.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(accountsBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(metaBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(accountDcMetaBucket)
		if err != nil {
			return err
		}
		return nil
	})
}

func openDb(dbFilePath string, create bool) (*wdb, error) {
	exists := util.FileExists(dbFilePath)
	if create && exists {
		return nil, errWalletDbAlreadyExists
	} else if !create && !exists {
		return nil, errWalletDbDoesNotExists
	}

	db, err := bolt.Open(dbFilePath, 0600, nil) // -rw-------
	if err != nil {
		return nil, err
	}

	w := &wdb{db, dbFilePath}
	err = w.createBuckets()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func createNewDb(config WalletConfig) (*wdb, error) {
	walletDir, err := config.GetWalletDir()
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(walletDir, 0700) // -rwx------
	if err != nil {
		return nil, err
	}

	dbFilePath := path.Join(walletDir, WalletFileName)
	return openDb(dbFilePath, true)
}

func parseBill(v []byte) (*Bill, error) {
	var b *Bill
	err := json.Unmarshal(v, &b)
	return b, err
}

func getAccountBucket(tx *bolt.Tx, accountIndex []byte) (*bolt.Bucket, error) {
	bkt := tx.Bucket(accountsBucket).Bucket(accountIndex)
	if bkt == nil {
		return nil, errAccountNotFound
	}
	return bkt, nil
}
