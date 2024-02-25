package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

// createRandomReference adds a new random reference table
// (Funnel, Funnel, LessonLocation, LessonSubject or PaymentMethod) to the database,
// and returns it.
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
	require.Equal(t, funnel1.FunnelID, funnel2.FunnelID)

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

func TestUpdateFunnel(t *testing.T) {
	funnel1 := createRandomFunnel(t)
	name := util.RandomName()

	arg := UpdateFunnelParams{
		FunnelID: funnel1.FunnelID,
		Name:     name,
	}
	err := testQueries.UpdateFunnel(context.Background(), arg)
	require.NoError(t, err)

	funnel2, err := testQueries.GetFunnel(context.Background(), funnel1.FunnelID)
	require.NoError(t, err)
	require.NotEmpty(t, funnel2)

	require.Equal(t, funnel1.FunnelID, funnel2.FunnelID)
	require.Equal(t, name, funnel2.Name)
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
	require.Len(t, funnels, int(arg.Limit))

	for _, v := range funnels {
		require.NotEmpty(t, v)
	}
}
