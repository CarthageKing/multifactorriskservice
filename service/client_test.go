package service

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

type ClientSuite struct {
	suite.Suite
	Server *httptest.Server
}

func (suite *ClientSuite) SetupTest() {
	require := suite.Require()
	assert := suite.Assert()

	f, err := os.Open("../fixtures/example_records.json")
	require.NoError(err)
	suite.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert request form values are correct for REDCap API
		assert.Equal("123456789", r.FormValue("token"))
		assert.Equal("record", r.FormValue("content"))
		assert.Equal("json", r.FormValue("format"))
		assert.Equal("flat", r.FormValue("type"))
		assert.Equal("study_id, redcap_event_name, mrn, participant_information_complete, rf_date, rf_cmc_risk_cat, rf_func_risk_cat, rf_sb_risk_cat, rf_util_risk_cat, rf_risk_predicted, risk_factors_complete", r.FormValue("fields"))
		assert.Equal("json", r.FormValue("returnFormat"))

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.Copy(w, f)
	}))
}

func (suite *ClientSuite) TearDownTest() {
	if suite.Server != nil {
		suite.Server.Close()
	}
}

func (suite *ClientSuite) TestGetREDCapData() {
	assert := suite.Assert()
	require := suite.Require()

	m, err := GetREDCapData(suite.Server.URL, "123456789")
	require.NoError(err)
	require.Len(m, 2)

	s, ok := m["1"]
	require.True(ok)
	assert.Equal("1", s.ID)
	assert.Equal("1-1", s.MedicalRecordNumber)
	require.Len(s.Records, 2)

	s, ok = m["a"]
	require.True(ok)
	assert.Equal("a", s.ID)
	assert.Equal("1-a", s.MedicalRecordNumber)
	require.Len(s.Records, 1)
}
