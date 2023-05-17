package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alphabill-org/alphabill/internal/network/protocol/genesis"
	"github.com/alphabill-org/alphabill/internal/script"
	testpartition "github.com/alphabill-org/alphabill/internal/testutils/partition"
	moneytx "github.com/alphabill-org/alphabill/internal/txsystem/money"
	"github.com/alphabill-org/alphabill/internal/txsystem/tokens"
	"github.com/alphabill-org/alphabill/internal/util"
	wlog "github.com/alphabill-org/alphabill/pkg/wallet/log"
	"github.com/alphabill-org/alphabill/pkg/wallet/money/backend"
	moneyclient "github.com/alphabill-org/alphabill/pkg/wallet/money/backend/client"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

var defaultTokenSDR = &genesis.SystemDescriptionRecord{
	SystemIdentifier: tokens.DefaultTokenTxSystemIdentifier,
	T2Timeout:        defaultT2Timeout,
	FeeCreditBill: &genesis.FeeCreditBill{
		UnitId:         util.Uint256ToBytes(uint256.NewInt(4)),
		OwnerPredicate: script.PredicateAlwaysTrue(),
	},
}

func TestWalletFeesCmds_MoneyPartition(t *testing.T) {
	homedir, abPartition := setupMoneyInfraAndWallet(t, []*testpartition.NodePartition{})
	// get money
	moneyPartition, err := abPartition.GetNodePartition(defaultABMoneySystemIdentifier)
	require.NoError(t, err)
	abNodeAddrFlag := "-u " + moneyPartition.Nodes[0].AddrGRPC

	// list fees
	stdout, err := execFeesCommand(homedir, "list")
	require.NoError(t, err)
	require.Equal(t, "Partition: money", stdout.lines[0])
	require.Equal(t, "Account #1 0.00000000", stdout.lines[1])

	// add fee credits
	amount := uint64(150)
	stdout, err = execFeesCommand(homedir, fmt.Sprintf("add --amount=%d %s", amount, abNodeAddrFlag))
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("Successfully created %d fee credits on money partition.", amount), stdout.lines[0])
	time.Sleep(2 * time.Second) // TODO waitForConf should use backend and not node for tx confirmations

	// verify fee credits
	expectedFees := amount*1e8 - 1
	stdout, err = execFeesCommand(homedir, "list")
	require.NoError(t, err)
	require.Equal(t, "Partition: money", stdout.lines[0])
	require.Equal(t, fmt.Sprintf("Account #1 %s", amountToString(expectedFees, 8)), stdout.lines[1])

	// add more fee credits
	stdout, err = execFeesCommand(homedir, fmt.Sprintf("add --amount=%d %s", amount, abNodeAddrFlag))
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("Successfully created %d fee credits on money partition.", amount), stdout.lines[0])
	time.Sleep(2 * time.Second) // TODO waitForConf should use backend and not node for tx confirmations

	// verify fee credits
	expectedFees = amount*2*1e8 - 2
	stdout, err = execFeesCommand(homedir, "list")
	require.NoError(t, err)
	require.Equal(t, "Partition: money", stdout.lines[0])
	require.Equal(t, fmt.Sprintf("Account #1 %s", amountToString(expectedFees, 8)), stdout.lines[1])

	// reclaim fees
	stdout, err = execFeesCommand(homedir, "reclaim "+abNodeAddrFlag)
	require.NoError(t, err)
	require.Equal(t, "Successfully reclaimed fee credits on money partition.", stdout.lines[0])

	// list fees
	stdout, err = execFeesCommand(homedir, "list")
	require.NoError(t, err)
	require.Equal(t, "Partition: money", stdout.lines[0])
	require.Equal(t, "Account #1 0.00000000", stdout.lines[1])
}

func TestWalletFeesCmds_TokenPartition(t *testing.T) {
	// start money partition and create wallet with token partition as well
	tokensPartition := createTokensPartition(t)
	homedir, abNet := setupMoneyInfraAndWallet(t, []*testpartition.NodePartition{tokensPartition})
	// get money
	moneyPartition, err := abNet.GetNodePartition(defaultABMoneySystemIdentifier)
	require.NoError(t, err)
	moneyNodeURL := moneyPartition.Nodes[0].AddrGRPC

	// start token partition
	startPartitionRPCServers(t, tokensPartition)

	tokenBackendURL, _, _ := startTokensBackend(t, tokensPartition.Nodes[0].AddrGRPC)
	args := fmt.Sprintf("--partition=token -r %s -u %s -m %s", defaultAlphabillApiURL, moneyNodeURL, tokenBackendURL)

	// list fees on token partition
	stdout, err := execFeesCommand(homedir, "list "+args)
	require.NoError(t, err)
	require.Equal(t, "Partition: token", stdout.lines[0])
	require.Equal(t, "Account #1 0.00000000", stdout.lines[1])

	// add fee credits
	amount := uint64(150)
	stdout, err = execFeesCommand(homedir, fmt.Sprintf("add --amount=%d %s", amount, args))
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("Successfully created %d fee credits on token partition.", amount), stdout.lines[0])

	// verify fee credits
	expectedFees := amount*1e8 - 1
	stdout, err = execFeesCommand(homedir, "list "+args)
	require.NoError(t, err)
	require.Equal(t, "Partition: token", stdout.lines[0])
	require.Equal(t, fmt.Sprintf("Account #1 %s", amountToString(expectedFees, 8)), stdout.lines[1])

	// add more fee credits to token partition
	stdout, err = execFeesCommand(homedir, fmt.Sprintf("add --amount=%d %s", amount, args))
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("Successfully created %d fee credits on token partition.", amount), stdout.lines[0])

	// verify fee credits to token partition
	expectedFees = amount*2*1e8 - 2
	stdout, err = execFeesCommand(homedir, "list "+args)
	require.NoError(t, err)
	require.Equal(t, "Partition: token", stdout.lines[0])
	require.Equal(t, fmt.Sprintf("Account #1 %s", amountToString(expectedFees, 8)), stdout.lines[1])

	// reclaim fees
	stdout, err = execFeesCommand(homedir, "reclaim "+args)
	require.NoError(t, err)
	require.Equal(t, "Successfully reclaimed fee credits on token partition.", stdout.lines[0])

	// list fees
	stdout, err = execFeesCommand(homedir, "list")
	require.NoError(t, err)
	require.Equal(t, "Partition: money", stdout.lines[0])
	require.Equal(t, "Account #1 0.00000000", stdout.lines[1])
}

func execFeesCommand(homeDir, command string) (*testConsoleWriter, error) {
	return execCommand(homeDir, " fees "+command)
}

// setupMoneyInfraAndWallet starts money partiton and wallet backend and sends initial bill to wallet.
// Returns wallet homedir and reference to money partition object.
func setupMoneyInfraAndWallet(t *testing.T, otherPartitions []*testpartition.NodePartition) (string, *testpartition.AlphabillNetwork) {
	initialBill := &moneytx.InitialBill{
		ID:    uint256.NewInt(1),
		Value: 1e18,
		Owner: script.PredicateAlwaysTrue(),
	}
	moneyPartition := createMoneyPartition(t, initialBill)
	nodePartitions := []*testpartition.NodePartition{moneyPartition}
	nodePartitions = append(nodePartitions, otherPartitions...)
	abNet := startAlphabill(t, nodePartitions)

	startPartitionRPCServers(t, moneyPartition)

	startMoneyBackend(t, moneyPartition, initialBill)

	// create wallet
	wlog.InitStdoutLogger(wlog.DEBUG)
	homedir := createNewTestWallet(t)

	stdout := execWalletCmd(t, "", homedir, "get-pubkeys")
	require.Len(t, stdout.lines, 1)
	pk, _ := strings.CutPrefix(stdout.lines[0], "#1 ")

	// transfer initial bill to wallet pubkey
	spendInitialBillWithFeeCredits(t, abNet, initialBill, pk)

	// wait for initial bill tx
	waitForBalanceCLI(t, homedir, defaultAlphabillApiURL, initialBill.Value-3, 0) // initial bill minus txfees

	return homedir, abNet
}

func startMoneyBackend(t *testing.T, moneyPart *testpartition.NodePartition, initialBill *moneytx.InitialBill) (string, *moneyclient.MoneyBackendClient) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)
	go func() {
		err := backend.Run(ctx,
			&backend.Config{
				ABMoneySystemIdentifier: []byte{0, 0, 0, 0},
				AlphabillUrl:            moneyPart.Nodes[0].AddrGRPC,
				ServerAddr:              defaultAlphabillApiURL, // TODO move to random port
				DbFile:                  filepath.Join(t.TempDir(), backend.BoltBillStoreFileName),
				ListBillsPageLimit:      100,
				InitialBill: backend.InitialBill{
					Id:        util.Uint256ToBytes(initialBill.ID),
					Value:     initialBill.Value,
					Predicate: script.PredicateAlwaysTrue(),
				},
				SystemDescriptionRecords: []*genesis.SystemDescriptionRecord{defaultMoneySDR, defaultTokenSDR},
			})
		require.ErrorIs(t, err, context.Canceled)
	}()

	restClient, err := moneyclient.New(defaultAlphabillApiURL)
	require.NoError(t, err)

	return defaultAlphabillApiURL, restClient
}