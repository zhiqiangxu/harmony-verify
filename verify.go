package main

import (
	"encoding/binary"
	"github.com/ethereum/go-ethereum/common"
	bls_core "github.com/harmony-one/bls/ffi/go/bls"
	"github.com/harmony-one/harmony/block"
	"github.com/harmony-one/harmony/consensus/quorum"
	"github.com/harmony-one/harmony/crypto/bls"
	"github.com/harmony-one/harmony/shard"
	"github.com/pkg/errors"
	"math/big"
)

func verify(parent, header *block.Header) error {
	pas := payloadArgsFromHeader(parent)
	sas := sigArgs{
		sig:    header.LastCommitSignature(),
		bitmap: header.LastCommitBitmap(),
	}
	return verifySignature(nil, pas, sas)
}

type sigArgs struct {
	sig    bls.SerializedSignature
	bitmap []byte
}

type payloadArgs struct {
	blockHash common.Hash
	shardID   uint32
	epoch     *big.Int
	number    uint64
	viewID    uint64
}

func (args payloadArgs) constructPayload() []byte {
	return ConstructCommitPayload(args.epoch, args.blockHash, args.number, args.viewID)
}

// ConstructCommitPayload returns the commit payload for consensus signatures.
func ConstructCommitPayload(epoch *big.Int, blockHash common.Hash, blockNum, viewID uint64,
) []byte {
	blockNumBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(blockNumBytes, blockNum)
	commitPayload := append(blockNumBytes, blockHash.Bytes()...)

	viewIDBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(viewIDBytes, viewID)
	return append(commitPayload, viewIDBytes...)
}

func payloadArgsFromHeader(header *block.Header) payloadArgs {
	return payloadArgs{
		blockHash: header.Hash(),
		shardID:   header.ShardID(),
		epoch:     header.Epoch(),
		number:    header.Number().Uint64(),
		viewID:    header.ViewID().Uint64(),
	}
}

func getEpochCtxC(ss *shard.State, shardID uint32, epoch *big.Int) (ec epochCtx, err error) {
	shardComm, err := ss.FindCommitteeByID(shardID)
	if err != nil {
		return
	}
	pubKeys, err := shardComm.BLSPublicKeys()
	if err != nil {
		return
	}
	qrVerifier, err := quorum.NewVerifier(shardComm, epoch, true)
	if err != nil {
		return
	}
	ec = epochCtx{
		qrVerifier: qrVerifier,
		pubKeys:    pubKeys,
	}
	return
}
func verifySignature(ss *shard.State, pas payloadArgs, sas sigArgs) (err error) {

	ec, err := getEpochCtxC(ss, pas.shardID, pas.epoch)
	if err != nil {
		return
	}

	var (
		pubKeys      = ec.pubKeys
		qrVerifier   = ec.qrVerifier
		commitSig    = sas.sig
		commitBitmap = sas.bitmap
	)
	aggSig, mask, err := DecodeSigBitmap(commitSig, commitBitmap, pubKeys)
	if err != nil {
		return errors.Wrap(err, "deserialize signature and bitmap")
	}
	if !qrVerifier.IsQuorumAchievedByMask(mask) {
		return errors.New("not enough signature collected")
	}
	commitPayload := pas.constructPayload()
	if !aggSig.VerifyHash(mask.AggregatePublic, commitPayload) {
		return errors.New("Unable to verify aggregated signature for block")
	}

	return nil
}

// DecodeSigBitmap decode and parse the signature, bitmap with the given public keys
func DecodeSigBitmap(sigBytes bls.SerializedSignature, bitmap []byte, pubKeys []bls.PublicKeyWrapper) (*bls_core.Sign, *bls.Mask, error) {
	aggSig := bls_core.Sign{}
	err := aggSig.Deserialize(sigBytes[:])
	if err != nil {
		return nil, nil, errors.New("unable to deserialize multi-signature from payload")
	}
	mask, err := bls.NewMask(pubKeys, nil)
	if err != nil {
		return nil, nil, errors.New("unable to setup mask from payload")
	}
	if err := mask.SetMask(bitmap); err != nil {
		return nil, nil, errors.New("mask.SetMask failed")
	}
	return &aggSig, mask, nil
}

type epochCtx struct {
	qrVerifier quorum.Verifier
	pubKeys    []bls.PublicKeyWrapper
}
