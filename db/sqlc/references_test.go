package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/require"
)

// createRandomReference adds a new random reference table
// (College, Funnel, LessonLocation, LessonSubject or PaymentMethod) to the database,
// and returns it.
func createRandomReferenceStruct(t *testing.T, key ReferenceStructName) ReferenceStruct {
	ref, ok := ReferenceStructMap[key]
	require.True(t, ok)

	name := util.RandomName()

	var err error
	switch key {
	case ReferenceCollege:
		ref, err = testQueries.CreateCollege(context.Background(), name)
	case ReferenceFunnel:
		ref, err = testQueries.CreateFunnel(context.Background(), name)
	case ReferenceLessonLocation:
		ref, err = testQueries.CreateLessonLocation(context.Background(), name)
	case ReferenceLessonSubject:
		ref, err = testQueries.CreateLessonSubject(context.Background(), name)
	case ReferencePaymentMethod:
		ref, err = testQueries.CreatePaymentMethod(context.Background(), name)
	}

	require.NoError(t, err)
	require.NotEmpty(t, ref)
	require.Equal(t, name, ref.GetName())
	require.NotZero(t, ref.GetID())
	return ref
}

func TestCreateReferences(t *testing.T) {
	for key := range ReferenceStructMap {
		t.Run(string(key), func(t *testing.T) {
			createRandomReferenceStruct(t, key)
		})
	}
}

func TestGetReferences(t *testing.T) {
	for key := range ReferenceStructMap {
		t.Run(string(key), func(t *testing.T) {
			ref1 := createRandomReferenceStruct(t, key)

			var ref2 ReferenceStruct
			var err error
			switch key {
			case ReferenceCollege:
				ref2, err = testQueries.GetCollege(context.Background(), ref1.GetID())
			case ReferenceFunnel:
				ref2, err = testQueries.GetFunnel(context.Background(), ref1.GetID())
			case ReferenceLessonLocation:
				ref2, err = testQueries.GetLessonLocation(context.Background(), ref1.GetID())
			case ReferenceLessonSubject:
				ref2, err = testQueries.GetLessonSubject(context.Background(), ref1.GetID())
			case ReferencePaymentMethod:
				ref2, err = testQueries.GetPaymentMethod(context.Background(), ref1.GetID())
			}

			require.NoError(t, err)
			require.NotEmpty(t, ref2)
			require.Equal(t, ref1.GetID(), ref2.GetID())
			require.Equal(t, ref1.GetName(), ref2.GetName())

		})
	}
}

func TestDeleteReferences(t *testing.T) {
	for key := range ReferenceStructMap {
		t.Run(string(key), func(t *testing.T) {
			ref1 := createRandomReferenceStruct(t, key)

			var ref2 ReferenceStruct
			var err error
			switch key {
			case ReferenceCollege:
				err = testQueries.DeleteCollege(context.Background(), ref1.GetID())
				require.NoError(t, err)

				ref2, err = testQueries.GetCollege(context.Background(), ref1.GetID())
			case ReferenceFunnel:
				err = testQueries.DeleteFunnel(context.Background(), ref1.GetID())
				require.NoError(t, err)

				ref2, err = testQueries.GetFunnel(context.Background(), ref1.GetID())
			case ReferenceLessonLocation:
				err = testQueries.DeleteLessonLocation(context.Background(), ref1.GetID())
				require.NoError(t, err)

				ref2, err = testQueries.GetLessonLocation(context.Background(), ref1.GetID())
			case ReferenceLessonSubject:
				err = testQueries.DeleteLessonSubject(context.Background(), ref1.GetID())
				require.NoError(t, err)

				ref2, err = testQueries.GetLessonSubject(context.Background(), ref1.GetID())
			case ReferencePaymentMethod:
				err = testQueries.DeletePaymentMethod(context.Background(), ref1.GetID())
				require.NoError(t, err)

				ref2, err = testQueries.GetPaymentMethod(context.Background(), ref1.GetID())
			}

			require.Error(t, err)
			require.EqualError(t, err, sql.ErrNoRows.Error())
			require.Empty(t, ref2)
		})
	}
}

func TestUpdateReferences(t *testing.T) {
	for key := range ReferenceStructMap {
		t.Run(string(key), func(t *testing.T) {
			ref1 := createRandomReferenceStruct(t, key)

			name := util.RandomName()

			var ref2 ReferenceStruct
			var err error
			switch key {
			case ReferenceCollege:
				arg := UpdateCollegeParams{
					CollegeID: ref1.GetID(),
					Name:      name,
				}

				err = testQueries.UpdateCollege(context.Background(), arg)
				require.NoError(t, err)

				ref2, err = testQueries.GetCollege(context.Background(), ref1.GetID())
			case ReferenceFunnel:
				arg := UpdateFunnelParams{
					FunnelID: ref1.GetID(),
					Name:     name,
				}

				err = testQueries.UpdateFunnel(context.Background(), arg)
				require.NoError(t, err)

				ref2, err = testQueries.GetFunnel(context.Background(), ref1.GetID())
			case ReferenceLessonLocation:
				arg := UpdateLessonLocationParams{
					LocationID: ref1.GetID(),
					Name:       name,
				}

				err = testQueries.UpdateLessonLocation(context.Background(), arg)
				require.NoError(t, err)

				ref2, err = testQueries.GetLessonLocation(context.Background(), ref1.GetID())
			case ReferenceLessonSubject:
				arg := UpdateLessonSubjectParams{
					SubjectID: ref1.GetID(),
					Name:      name,
				}

				err = testQueries.UpdateLessonSubject(context.Background(), arg)
				require.NoError(t, err)

				ref2, err = testQueries.GetLessonSubject(context.Background(), ref1.GetID())
			case ReferencePaymentMethod:
				arg := UpdatePaymentMethodParams{
					PaymentMethodID: ref1.GetID(),
					Name:            name,
				}

				err = testQueries.UpdatePaymentMethod(context.Background(), arg)
				require.NoError(t, err)

				ref2, err = testQueries.GetPaymentMethod(context.Background(), ref1.GetID())
			}

			require.NoError(t, err)
			require.NotEmpty(t, ref2)

			require.Equal(t, ref1.GetID(), ref2.GetID())
			require.Equal(t, name, ref2.GetName())

		})
	}
}

func TestListReferences(t *testing.T) {
	for key := range ReferenceStructMap {
		t.Run(string(key), func(t *testing.T) {
			for i := 0; i < 10; i++ {
				createRandomReferenceStruct(t, key)
			}

			limit := 5
			offset := 5

			switch key {
			case ReferenceCollege:
				list, err := testQueries.ListColleges(context.Background(),
					ListCollegesParams{
						Limit:  int32(limit),
						Offset: int32(offset),
					})
				require.NoError(t, err)
				require.Len(t, list, limit)

				for _, v := range list {
					require.NotEmpty(t, v)
				}
			case ReferenceFunnel:
				list, err := testQueries.ListFunnels(context.Background(),
					ListFunnelsParams{
						Limit:  int32(limit),
						Offset: int32(offset),
					})
				require.NoError(t, err)
				require.Len(t, list, limit)

				for _, v := range list {
					require.NotEmpty(t, v)
				}
			case ReferenceLessonLocation:
				list, err := testQueries.ListLessonLocations(context.Background(),
					ListLessonLocationsParams{
						Limit:  int32(limit),
						Offset: int32(offset),
					})
				require.NoError(t, err)
				require.Len(t, list, limit)

				for _, v := range list {
					require.NotEmpty(t, v)
				}
			case ReferenceLessonSubject:
				list, err := testQueries.ListLessonSubjects(context.Background(),
					ListLessonSubjectsParams{Limit: int32(limit), Offset: int32(offset)})
				require.NoError(t, err)
				require.Len(t, list, limit)

				for _, v := range list {
					require.NotEmpty(t, v)
				}
			case ReferencePaymentMethod:
				list, err := testQueries.ListPaymentMethods(context.Background(),
					ListPaymentMethodsParams{
						Limit:  int32(limit),
						Offset: int32(offset),
					})
				require.NoError(t, err)
				require.Len(t, list, limit)

				for _, v := range list {
					require.NotEmpty(t, v)
				}
			}
		})
	}
}
