package main

import (
	"encoding/hex"
	"fmt"

	"github.com/Salvionied/apollo"
	"github.com/Salvionied/apollo/serialization/Address"
	"github.com/Salvionied/apollo/serialization/Key"
	"github.com/Salvionied/apollo/serialization/TransactionOutput"
	"github.com/Salvionied/apollo/serialization/UTxO"
	"github.com/Salvionied/apollo/serialization/Value"
	"github.com/Salvionied/apollo/txBuilding/Backend/FixedChainContext"
	"github.com/fxamacker/cbor/v2"
)

func main() {
	cc := FixedChainContext.InitFixedChainContext()
	apollob := apollo.New(&cc)

	utxo := UTxO.UTxO{}
	utxo.Input.Index = 0
	utxo.Input.TransactionId, _ = hex.DecodeString("f168dc43fe77ae2ecd3af9e21f0c4f29cb4b740db2b20bc7580612fa8f82062b")

	address, err := Address.DecodeAddress("addr_test1vqxvmfj9ky733sqhwv38n95hammmlrdwmzjnv7huezh8cwgxrcmtd")
	if err != nil {
		panic(err)
	}
	utxo.Output = TransactionOutput.SimpleTransactionOutput(address, Value.SimpleValue(10_000_000, nil))

	apollob.AddLoadedUTxOs(utxo)
	apollob.PayToAddressBech32("addr_test1vqxvmfj9ky733sqhwv38n95hammmlrdwmzjnv7huezh8cwgxrcmtd", 5_000_000)
	apollob.SetChangeAddressBech32("addr_test1vqxvmfj9ky733sqhwv38n95hammmlrdwmzjnv7huezh8cwgxrcmtd")
	_, err = apollob.Complete()
	if err != nil {
		panic(err)
	}

	vkey, err := Key.VerificationKeyFromCbor("58203b1d38f9a136083f3b7d332a37b320e8409b9c5278ba3716f2af7733d6eb4045")
	if err != nil {
		panic(err)
	}
	skey, err := Key.SigningKeyFromHexString("58209fc392611d00553b1e034fbe1a017b894727a055fae498cd8d4b65c315550dbf")
	if err != nil {
		panic(err)
	}

	_, err = apollob.SignWithSkey(*vkey, *skey)
	if err != nil {
		panic(err)
	}

	tx := apollob.GetTx()
	cborred, err := cbor.Marshal(tx)
	if err != nil {
		panic(err)
	}
	fmt.Println(hex.EncodeToString(cborred))
}
