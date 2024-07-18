package test

import (
	"testing"

	"github.com/sebastianmontero/bennyfi-go-client/ibc"
	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/test"
	"gotest.tools/assert"
)

type TestUtil struct {
	ibcClient *ibc.IBCContract
	*test.TestUtil
}

func NewTestUtil(t *testing.T, ibcClient *ibc.IBCContract) *TestUtil {
	return &TestUtil{
		ibcClient: ibcClient,
		TestUtil:  test.NewTestUtil(t, ibcClient.EOS),
	}
}

func (m *TestUtil) AssertLightProofExists(exists bool) {
	actual, err := m.ibcClient.GetLightProof()
	assert.NilError(m.T, err)
	assert.Equal(m.T, actual != nil, exists)
}

func (m *TestUtil) AssertHeavyProofExists(exists bool) {
	actual, err := m.ibcClient.GetHeavyProof()
	assert.NilError(m.T, err)
	assert.Equal(m.T, actual != nil, exists)
}

func (m *TestUtil) AssertLightProofValidations(bridgeContract eos.AccountName, numValidations int) {
	m.AssertNumActions(bridgeContract, "checkproofc", numValidations)
}

func (m *TestUtil) AssertHeavyProofValidations(bridgeContract eos.AccountName, numValidations int) {
	m.AssertNumActions(bridgeContract, "checkproofb", numValidations)
}

func (m *TestUtil) AssertEmitXferAction(emitXferParams *ibc.EmitXferParams) map[string]interface{} {

	actionData := map[string]interface{}{
		"xfer": map[string]interface{}{
			"owner":       emitXferParams.Owner,
			"quantity":    emitXferParams.Quantity,
			"beneficiary": emitXferParams.Beneficiary,
		},
	}
	return m.AssertAction(eos.AccountName(m.ibcClient.ContractName), "emitxfer", actionData, 0)
}
