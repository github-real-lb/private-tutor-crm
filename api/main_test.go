package api

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// requireBodyMatchStruct asserts that a JSON httptest.ResponseRecorder.Body
// equal to a Struct object.
func requireBodyMatchStruct(t *testing.T, body *bytes.Buffer, obj interface{}) {
	jsonBodyData, err := io.ReadAll(body)
	require.NoError(t, err)

	jsonObjData, err := json.Marshal(obj)
	require.NoError(t, err)

	require.Equal(t, string(jsonObjData), string(jsonBodyData))
}
