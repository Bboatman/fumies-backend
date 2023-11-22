package handlers

import (
	"cmp"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/drewlanenga/govector"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"fumies/api/models"
	"fumies/api/queries"
)

func GetPerfume(ctx *gin.Context) {
	userId, _ := ctx.Get("user_id")
	pool, _ := ctx.Get("db_pool")
	var perfumes []models.PerfumeResponse
	err := pgxscan.Select(ctx, pool.(*pgxpool.Pool), &perfumes, queries.SelectPerfumesForUser, userId)
	if err != nil {
		log.Default().Println(err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, perfumes)
}

func CreateOrUpdatePerfume(ctx *gin.Context) {
	userId, _ := ctx.Get("user_id")
	perfumeId, updatePerfume := ctx.Params.Get("id")
	var requestBody models.ModifyPerfumeBody
	if err := ctx.BindJSON(&requestBody); err != nil {
		log.Default().Printf("Failed to bind request: %+v", err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	pool, _ := ctx.Get("db_pool")

	var result models.Perfume
	var query string
	var modifierId interface{}

	if updatePerfume {
		query = queries.UpdatePerfumeForUser
		modifierId = perfumeId
	} else {
		query = queries.CreatePerfumeForUser
		modifierId = userId
	}

	rows, err := pool.(*pgxpool.Pool).Query(
		ctx,
		query,
		modifierId,
		requestBody.Name,
		requestBody.House,
		requestBody.Url,
		requestBody.Description,
		requestBody.IsEmpty,
	)

	if err != nil {
		log.Default().Printf("Failed to query request: %+v", err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = pgxscan.ScanOne(&result, rows)

	if err != nil {
		log.Default().Printf("Failed to bind request: %+v", err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = SetMetricsForPerfume(pool.(*pgxpool.Pool), ctx, result.Id, requestBody.Notes)
	if err != nil {
		log.Default().Printf("Failed to save metrics: %+v", err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if updatePerfume {
		ctx.JSON(http.StatusOK, result)
	} else {
		ctx.JSON(http.StatusCreated, result)
	}
}

func SetMetricsForPerfume(pool *pgxpool.Pool, ctx *gin.Context, perfumeId uuid.UUID, noteIds *[]uuid.UUID) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, queries.DeleteMetricsForPerfume, perfumeId)

	if err != nil {
		return err
	}

	if noteIds == nil {
		err = tx.Commit(ctx)
		if err != nil {
			return err
		}
		return nil
	}

	perfumeMetrics := [][]any{}
	for _, x := range *noteIds {
		perfumeMetrics = append(perfumeMetrics, []any{perfumeId, x})
	}
	count, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"perfume_metric"},
		[]string{"perfume_id", "note_id"},
		pgx.CopyFromRows(perfumeMetrics),
	)

	if err != nil {
		return err
	}
	log.Default().Printf("updated %v", count)

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func RecommendPerfume(ctx *gin.Context) {
	userId, _ := ctx.Get("user_id")
	var requestBody models.RecommendationRequest
	pool, _ := ctx.Get("db_pool")

	if err := ctx.BindJSON(&requestBody); err != nil {
		log.Default().Printf("Failed to bind request: %+v", err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if len(requestBody.Notes) < 1 {
		ctx.JSON(http.StatusBadRequest, "invalid note data")
		return
	}

	var perfumeVectors []models.PerfumeVector
	err := pgxscan.Select(ctx, pool.(*pgxpool.Pool), &perfumeVectors, queries.SelectPerfumeVectorsForUser, userId)

	if err != nil {
		log.Default().Printf("Failed to get perfume vectors for user: %+v", err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	log.Default().Printf("perfume vectors %+v", perfumeVectors)

	var moodVectorScan []models.PerfumeVector
	moodIds := fmt.Sprintf("'%s'", strings.Join(requestBody.Notes, "','"))
	err = pgxscan.Select(
		ctx,
		pool.(*pgxpool.Pool),
		&moodVectorScan,
		fmt.Sprintf(queries.SelectMoodVectorForUser, moodIds),
	)
	if err != nil {
		log.Default().Printf("Failed to get mood vector for user: %+v", err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	moodVector, _ := govector.AsVector(moodVectorScan[0].Vector)
	var retVal []models.PerfumeVector
	for ind, x := range perfumeVectors {
		xVector, _ := govector.AsVector(x.Vector)
		sim, _ := govector.Cosine(xVector, moodVector)

		inclusionRate := float64(0)
		if x.Epoch != nil {
			now := time.Now()
			t := float64(now.Unix()) - *x.Epoch

			// lim -> 1 at 0,0 (e.g epoch) and 0 at inf,inf (e.g approaching eternity)
			inclusionRate = 1 - (math.Log(1+t) / t)
		}

		include := rand.Float64() > inclusionRate
		perfumeVectors[ind].CosineSim = &sim
		perfumeVectors[ind].Include = &include

		if include {
			retVal = append(retVal, perfumeVectors[ind])
		}
	}

	slices.SortFunc(retVal,
		func(a, b models.PerfumeVector) int {
			return cmp.Compare(*b.CosineSim, *a.CosineSim)
		})

	ctx.JSON(http.StatusOK, retVal)
}

func WearPerfume(ctx *gin.Context) {
	userId, _ := ctx.Get("user_id")
	perfumeId, _ := ctx.Params.Get("id")
	pool, _ := ctx.Get("db_pool")

	_, err := pool.(*pgxpool.Pool).Exec(
		ctx,
		queries.CreatePerfumeWearForUser,
		userId,
		perfumeId,
	)

	if err != nil {
		log.Default().Printf("Failed to save wear")
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, nil)
}
