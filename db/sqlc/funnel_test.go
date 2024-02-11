package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

// createRandomFunnel tests adding a new random funnel to the database, and returns the Funnel data type.
func createRandomFunnel(t *testing.T) Funnel {
	name := util.RandomName()
	funnel, err := testQueries.CreateFunnel(context.Background(), name)
	require.NoError(t, err)
	require.NotEmpty(t, funnel)
	require.Equal(t, name, funnel.Name)
	require.NotZero(t, funnel.FunnelID)

	return funnel
}

func TestCreateFunnel(t *testing.T) {
	createRandomFunnel(t)
}

func TestGetFunnel(t *testing.T) {
	funnel1 := createRandomFunnel(t)
	funnel2, err := testQueries.GetFunnel(context.Background(), funnel1.FunnelID)
	require.NoError(t, err)
	require.NotEmpty(t, funnel2)

	require.Equal(t, funnel1.FunnelID, funnel2.FunnelID)
	require.Equal(t, funnel1.Name, funnel2.Name)
}

func TestUpdateFunnel(t *testing.T) {
	funnel1 := createRandomFunnel(t)

	arg := UpdateFunnelParams{
		FunnelID: funnel1.FunnelID,
		Name:     util.RandomName(),
	}
	err := testQueries.UpdateFunnel(context.Background(), arg)
	require.NoError(t, err)

	funnel2, err := testQueries.GetFunnel(context.Background(), arg.FunnelID)
	require.NoError(t, err)
	require.NotEmpty(t, funnel2)

	require.Equal(t, funnel1.FunnelID, funnel2.FunnelID)
	require.Equal(t, arg.Name, funnel2.Name)
}

func TestDeleteFunnel(t *testing.T) {
	funnel1 := createRandomFunnel(t)

	err := testQueries.DeleteFunnel(context.Background(), funnel1.FunnelID)
	require.NoError(t, err)

	funnel2, err := testQueries.GetFunnel(context.Background(), funnel1.FunnelID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, funnel2)
}

func TestListFunnels(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomFunnel(t)
	}

	arg := ListFunnelsParams{
		Limit:  5,
		Offset: 5,
	}

	funnels, err := testQueries.ListFunnels(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, funnels, 5)

	for _, funnel := range funnels {
		require.NotEmpty(t, funnel)
	}
}
