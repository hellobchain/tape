package basic

import (
	"crypto/rand"
	"encoding/asn1"
	"github.com/wsw365904/newcryptosm"
	"github.com/wsw365904/newcryptosm/ecdsa"
	"github.com/wsw365904/newcryptosm/x509"
	"math/big"

	"github.com/wsw365904/tape/internal/fabric/bccsp/utils"
	"github.com/wsw365904/tape/internal/fabric/common/crypto"

	"github.com/hyperledger/fabric-protos-go/common"
)

type CryptoConfig struct {
	MSPID      string
	PrivKey    string
	SignCert   string
	TLSCACerts []string
}

type ECDSASignature struct {
	R, S *big.Int
}

type CryptoImpl struct {
	Creator  []byte
	PrivKey  *ecdsa.PrivateKey
	SignCert *x509.Certificate
	HashType newcryptosm.Hash
}

func (s *CryptoImpl) Sign(msg []byte) ([]byte, error) {
	ri, si, err := ecdsa.Sign(rand.Reader, s.PrivKey, digest(msg, s.Hash()))
	if err != nil {
		return nil, err
	}

	si, _, err = utils.ToLowS(&s.PrivKey.PublicKey, si)
	if err != nil {
		return nil, err
	}

	return asn1.Marshal(ECDSASignature{ri, si})
}

func (s *CryptoImpl) Serialize() ([]byte, error) {
	return s.Creator, nil
}

func (s *CryptoImpl) Hash() newcryptosm.Hash {
	return s.HashType
}

func (s *CryptoImpl) NewSignatureHeader() (*common.SignatureHeader, error) {
	creator, err := s.Serialize()
	if err != nil {
		return nil, err
	}
	nonce, err := crypto.GetRandomNonce()
	if err != nil {
		return nil, err
	}

	return &common.SignatureHeader{
		Creator: creator,
		Nonce:   nonce,
	}, nil
}

func digest(in []byte, hashType newcryptosm.Hash) []byte {
	h := hashType.New()
	h.Write(in)
	return h.Sum(nil)
}
