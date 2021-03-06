package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/code.google.com/p/go-uuid/uuid"
	kstore "github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/eris-ltd/eris-keys/crypto/key_store"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/ed25519"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/account"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/tendermint/tendermint/wire"
)

func Pubkeyer(k *kstore.Key) ([]byte, error) {
	priv := k.PrivateKey
	privKeyBytes := new([64]byte)
	copy(privKeyBytes[:32], priv)
	pubKeyBytes := ed25519.MakePublicKey(privKeyBytes)
	return pubKeyBytes[:], nil
}

func init() {
	kstore.SetPubkeyer(Pubkeyer)
}

type PrivValidator struct {
	Address    []byte                 `json:"address"`
	PubKey     account.PubKeyEd25519  `json:"pub_key"`
	PrivKey    account.PrivKeyEd25519 `json:"priv_key"`
	LastHeight int                    `json:"last_height"`
	LastRound  int                    `json:"last_round"`
	LastStep   int                    `json:"last_step"`
}

func cliConvertErisKeyToPrivValidator(cmd *cobra.Command, args []string) {
	cmd.ParseFlags(args)
	if len(args) == 0 {
		exit(fmt.Errorf("Please enter the address of your key"))
	}

	addr := args[0]
	addrBytes, err := hex.DecodeString(addr)
	ifExit(err)

	privVal, err := coreConvertErisKeyToPrivValidator(addrBytes)
	ifExit(err)

	fmt.Println(string(wire.JSONBytes(privVal)))
}

func coreConvertErisKeyToPrivValidator(addrBytes []byte) (*PrivValidator, error) {
	keyStore := kstore.NewKeyStorePlain(DefaultKeyStore)
	key, err := keyStore.GetKey(addrBytes, "")
	ifExit(err)

	pub, err := key.Pubkey()
	if err != nil {
		return nil, err
	}

	var pubKey account.PubKeyEd25519
	copy(pubKey[:], pub)

	var privKey account.PrivKeyEd25519
	copy(privKey[:], key.PrivateKey)

	return &PrivValidator{
		Address: addrBytes,
		PubKey:  pubKey,
		PrivKey: privKey,
	}, nil
}

func cliConvertPrivValidatorToErisKey(cmd *cobra.Command, args []string) {
	cmd.ParseFlags(args)
	if len(args) == 0 {
		exit(fmt.Errorf("Please enter the path to the priv_validator.json"))
	}

	pvf := args[0]
	b, err := ioutil.ReadFile(pvf)
	ifExit(err)

	key, err := coreConvertPrivValidatorToErisKey(b)
	ifExit(err)

	fmt.Printf("%X\n", key.Address)
}

func coreConvertPrivValidatorToErisKey(b []byte) (key *kstore.Key, err error) {

	pv := new(PrivValidator)
	wire.ReadJSON(pv, b, &err)
	if err != nil {
		return nil, err
	}

	keyStore := kstore.NewKeyStorePlain(DefaultKeyStore)

	key = &kstore.Key{
		Id:         uuid.NewRandom(),
		Type:       kstore.KeyType{kstore.CurveTypeEd25519, kstore.AddrTypeRipemd160},
		Address:    pv.Address,
		PrivateKey: pv.PrivKey[:],
	}

	err = keyStore.StoreKey(key, "")
	return key, err
}
