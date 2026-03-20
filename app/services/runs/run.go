package Runs

import (
	"context"
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/mongodb"
	"dungeons/app/server"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Run struct {
	validate *validator.Validate
}

func New() *Run {
	return &Run{
		validate: validator.New(),
	}
}

func (r *Run) Get(queryParams models.QueryParams) ([]models.Run, error) {
	var (
		err  error
		runs []models.Run
		run  models.Run
		cursor *mongo.Cursor
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(run.Collection())

	filter := mongodb.SelectConstructeur(queryParams)
	cursor, err = collection.Find(context.TODO(), filter)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		// A new result variable should be declared for each document.
		var run models.Run
		err = cursor.Decode(&run)
		if err != nil {
			log.Error().Err(err).Msg("")
			return nil, err
		}
		runs = append(runs, run)
	}

	err = cursor.Err()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return runs, err
}

func (r *Run) Create(in *models.Run) (*models.Run, error) {
	var run models.Run

	srv := server.GetServer()
	collection := srv.Database.Collection(run.Collection())

	err := r.validate.Struct(in)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	err = functions.ConvertInputStructToDataStruct(in, &run)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	run.CustomID = functions.NewUUID()
	run.StartedAt = time.Now()

	_, err = collection.InsertOne(context.TODO(), run)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &run, nil
}

func (r *Run) GetByID(id string) (models.Run, error) {
	var (
		err error
		run models.Run
		queryParams models.QueryParams
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(run.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)
	err = collection.FindOne(context.TODO(), filter).Decode(&run)
	if err == nil {
		if err == mongo.ErrNoDocuments {
			log.Error().Err(err).Msg("")
			return run, err
		}
	}

	return run, err
}