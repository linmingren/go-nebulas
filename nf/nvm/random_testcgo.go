package nvm

import (
	"C"

	"github.com/nebulasio/go-nebulas/core"
)
import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/nebulasio/go-nebulas/consensus/dpos"
	"github.com/nebulasio/go-nebulas/core/state"
	"github.com/nebulasio/go-nebulas/crypto"
	"github.com/nebulasio/go-nebulas/crypto/keystore"
	"github.com/nebulasio/go-nebulas/crypto/keystore/secp256k1"
	"github.com/nebulasio/go-nebulas/storage"
	"github.com/nebulasio/go-nebulas/util"
	"github.com/nebulasio/go-nebulas/util/byteutils"
	"github.com/stretchr/testify/assert"
)

func newUint128FromIntWrapper2(a int64) *util.Uint128 {
	b, _ := util.NewUint128FromInt(a)
	return b
}

func mockNormalTransaction2(from, to, value string) *core.Transaction {

	fromAddr, _ := core.AddressParse(from)
	toAddr, _ := core.AddressParse(to)
	payload, _ := core.NewBinaryPayload(nil).ToBytes()
	gasPrice, _ := util.NewUint128FromString("1000000")
	gasLimit, _ := util.NewUint128FromString("2000000")
	v, _ := util.NewUint128FromString(value)
	tx, _ := core.NewTransaction(1, fromAddr, toAddr, v, 1, core.TxPayloadBinaryType, payload, gasPrice, gasLimit)

	priv1 := secp256k1.GeneratePrivateKey()
	signature, _ := crypto.NewSignature(keystore.SECP256K1)
	signature.InitSign(priv1)
	tx.Sign(signature)
	return tx
}

type testBlock2 struct {
}

// Coinbase mock
func (block *testBlock2) Coinbase() *core.Address {
	addr, _ := core.AddressParse("n1FkntVUMPAsESuCAAPK711omQk19JotBjM")
	return addr
}

// Hash mock
func (block *testBlock2) Hash() byteutils.Hash {
	return []byte("59fc526072b09af8a8ca9732dae17132c4e9127e43cf2232")
}

// Height mock
func (block *testBlock2) Height() uint64 {
	return core.NvmMemoryLimitWithoutInjectHeight
}

// RandomSeed mock
func (block *testBlock2) RandomSeed() string {
	return "59fc526072b09af8a8ca9732dae17132c4e9127e43cf2232"
}

// RandomAvailable mock
func (block *testBlock2) RandomAvailable() bool {
	return true
}

// DateAvailable
func (block *testBlock2) DateAvailable() bool {
	return true
}

// GetTransaction mock
func (block *testBlock2) GetTransaction(hash byteutils.Hash) (*core.Transaction, error) {
	return nil, nil
}

// RecordEvent mock
func (block *testBlock2) RecordEvent(txHash byteutils.Hash, topic, data string) error {
	return nil
}

func (block *testBlock2) Timestamp() int64 {
	return int64(0)
}

func mockBlock2() Block {
	block := &testBlock2{}
	return block
}

func testRandomFunc(t *testing.T) {
	mem, _ := storage.NewMemoryStorage()
	context, _ := state.NewWorldState(dpos.NewDpos(), mem)
	contractAddr, err := core.AddressParse("n1FkntVUMPAsESuCAAPK711omQk19JotBjM")
	assert.Nil(t, err)
	contract, _ := context.CreateContractAccount(contractAddr.Bytes(), nil)
	contract.AddBalance(newUint128FromIntWrapper2(5))

	tx := mockNormalTransaction2("n1FkntVUMPAsESuCAAPK711omQk19JotBjM", "n1TV3sU6jyzR4rJ1D7jCAmtVGSntJagXZHC", "0")
	ctx, err := NewContext(mockBlock2(), tx, contract, context)

	// execute.
	engine := NewV8Engine(ctx)
	assert.Nil(t, engine.ctx.rand)

	r1 := GetTxRandomFunc(unsafe.Pointer(uintptr(engine.lcsHandler)))
	assert.NotNil(t, r1)
	assert.NotNil(t, engine.ctx.rand)
	rs1 := C.GoString(r1)

	r2 := GetTxRandomFunc(unsafe.Pointer(uintptr(engine.lcsHandler)))
	assert.NotNil(t, r2)
	rs2 := C.GoString(r2)

	assert.NotEqual(t, rs1, rs2)

	fmt.Println(rs1, rs2)

	engine.Dispose()
}