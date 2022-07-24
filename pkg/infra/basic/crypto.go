package basic

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"github.com/wsw365904/cryptosm/ecdsa"
	"github.com/wsw365904/cryptosm/x509"
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
}

func (s *CryptoImpl) Sign(msg []byte) ([]byte, error) {
	ri, si, err := ecdsa.Sign(rand.Reader, s.PrivKey, digest(msg))
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

func digest(in []byte) []byte {
	h := sha256.New()
	h.Write(in)
	return h.Sum(nil)
}
