package main

import (
	"context"
	"log"
	"net"

	"bikraj.movie_microservice.net/gen"
	metadatatest "bikraj.movie_microservice.net/metadata/pkg/testutils"
	movietest "bikraj.movie_microservice.net/movie/pkg/testutils"
	"bikraj.movie_microservice.net/pkg/discovery"
	memory "bikraj.movie_microservice.net/pkg/discovery/memorypackage"
	ratingtest "bikraj.movie_microservice.net/rating/pkg/testutils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	metadataServiceName = "metadata"
	ratingServiceName   = "rating"
	movieServiceName    = "movie"
	metadataServiceAddr = "localhost:8081"
	ratingServiceaddr   = "localhost:8082"
	movieServiceAddr    = "localhost:8083"
)

func main() {

	log.Println("starting the intergration test")

	ctx := context.Background()

	registry := memory.NewRegistry()

	log.Println("Setting up Service Handlers and clients")

	metadatSrv := startMetadataService(ctx, registry)

	defer metadatSrv.GracefulStop()

	ratingSrv := startRatingService(ctx, registry)

	defer ratingSrv.GracefulStop()

	movieSrv := startMovieService(ctx, registry)

	defer movieSrv.GracefulStop()

	log.Println("setting up the client")

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	// Client for connecting with metadata.

	metadataConn, err := grpc.Dial(metadataServiceAddr, opts)
	if err != nil {
		panic(err)
	}
	defer metadataConn.Close()

	metadataClient := gen.NewMetadataServiceClient(metadataConn)

	// Client for connecting with  rating.

	ratingConn, err := grpc.Dial(ratingServiceaddr, opts)
	if err != nil {
		panic(err)
	}
	defer ratingConn.Close()

	ratingClient := gen.NewRatingServiceClient(ratingConn)

	// Client for connecting with movie.

	movieConn, err := grpc.Dial(movieServiceAddr, opts)
	if err != nil {
		panic(err)
	}
	defer movieConn.Close()

	movieClient := gen.NewMovieServiceClient(movieConn)

	log.Println("Starting the testing for the metadata Service")
	log.Println("Saving test metadata via metadat service")

	m := &gen.Metadata{
		Id:          "the-movie",
		Title:       "some random title",
		Description: "Some random Description",
		Director:    "some random Director",
	}
	if _, err := metadataClient.PutMetadata(ctx, &gen.PutMetadataRequest{Metadata: m}); err != nil {
		log.Fatalf("put metadata : %v", err)
	}

	log.Println("Retrieving test metdata via metadata service")

	getMetadataResp, err := metadataClient.GetMetadata(ctx, &gen.GetMetadataReqeust{MovieId: m.Id})
	if err != nil {
		log.Fatalf("get metadata : %v", err)
	}

	if diff := cmp.Diff(getMetadataResp.Metadata, m, cmpopts.IgnoreUnexported(gen.Metadata{})); diff != "" {
		log.Fatalf("get metadata after put mismatch: %v", diff)
	}

	log.Println("Starting tests on movie service")
	log.Println("Getting movie details via movie service")

	wantMovieDetails := &gen.MovieDetails{
		Metadata: m,
	}
	getMovieDetailsResp, err := movieClient.GetMovieDetails(ctx, &gen.GetMovieDetailsRequest{MovieId: m.Id})
	if err != nil {
		log.Fatalf("get movie details error: %v", err)
	}

	// Convert the expected and actual responses to JSON strings
	gotJSON, err := protojson.Marshal(getMovieDetailsResp.MovieDetails)
	if err != nil {
		log.Fatalf("Error marshaling Got message: %v", err)
	}

	wantJSON, err := protojson.Marshal(wantMovieDetails)
	if err != nil {
		log.Fatalf("Error marshaling Want message: %v", err)
	}

	// Convert JSON bytes to strings for easy comparison
	gotStr := string(gotJSON)
	wantStr := string(wantJSON)

	// Log them for debugging
	if gotStr != wantStr {
		log.Fatalf("Mismatch:\nGot: %s\nWant: %s", gotStr, wantStr)
	}

	// log.Printf("Got JSON: %s", gotJSON)
	// log.Printf("Want JSON: %s", wantJSON)
	// if !proto.Equal(getMovieDetailsResp.MovieDetails, wantMovieDetails) {
	// 	log.Fatalf("Mismatch:\nGot: %+v\nWant: %+v", getMovieDetailsResp.MovieDetails, wantMovieDetails)
	// }
	log.Println("testing the rating service")

	log.Println("testing the put rating via ratin service")

	const userId = "user0"
	const recordTypeMovie = "movie"

	firsRating := int32(5)

	if _, err := ratingClient.PutRating(ctx, &gen.PutRatingRequest{UserId: userId, RecordType: recordTypeMovie, RecordId: m.Id, RatingValue: firsRating}); err != nil {
		log.Fatalf("put meatadata error: %v", err)
	}
	log.Println("Retrieving the initial rating")

	getAggregratedRatingresp, err := ratingClient.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{RecordId: m.Id, RecordType: recordTypeMovie})
	if err != nil {
		log.Fatalf("get aggregated Ratin error: %v", err)
	}

	if got, want := getAggregratedRatingresp.RatingValue, float64(5); got != want {
		log.Fatalf("rating mismatch: got %v want %v", got, want)
	}
	log.Println("testing the put rating the second time")
	log.Println("testing the put rating via ratin service")

	secondRating := int32(4)

	if _, err := ratingClient.PutRating(ctx, &gen.PutRatingRequest{UserId: userId, RecordType: recordTypeMovie, RecordId: m.Id, RatingValue: secondRating}); err != nil {
		log.Fatalf("put meatadata error: %v", err)
	}
	log.Println("Retrieving the initial rating")

	getAggregratedRatingresp, err = ratingClient.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{RecordId: m.Id, RecordType: recordTypeMovie})
	if err != nil {
		log.Fatalf("get aggregated Rating error: %v", err)
	}
	wantRating := float64(4.5)

	if got := getAggregratedRatingresp.RatingValue; got != wantRating {
		log.Fatalf("rating mismatch: got %v want %v", got, wantRating)
	}

	log.Println("testing get movie details after update on the rating")
	getMovieDetailsResp, err = movieClient.GetMovieDetails(ctx, &gen.GetMovieDetailsRequest{MovieId: m.Id})
	if err != nil {
		log.Fatalf("get movie Details error: %v", err)
	}
	gotJSON, err = protojson.Marshal(getMovieDetailsResp.MovieDetails)
	if err != nil {
		log.Fatalf("Error marshaling Got message: %v", err)
	}
	wantMovieDetails.Rating = float32(wantRating)
	wantJSON, err = protojson.Marshal(wantMovieDetails)
	if err != nil {
		log.Fatalf("Error marshaling Want message: %v", err)
	}
	gotStr = string(gotJSON)
	wantStr = string(wantJSON)
	if diff := cmp.Diff(wantStr, gotStr); diff != "" {

		log.Fatalf("Mismatch:\nGot: %+v\nWant: %+v", gotStr, wantStr)
	}
	log.Println(" intergration test sucessfully executed")

}
func startMetadataService(ctx context.Context, registry discovery.Registry) *grpc.Server {
	log.Println("starting movie service on" + metadataServiceName)

	h := metadatatest.NewTestMetadataGRPCServer()
	l, err := net.Listen("tcp", metadataServiceAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMetadataServiceServer(srv, h)
	go func() {
		if err := srv.Serve(l); err != nil {
			panic(err)
		}

	}()
	id := discovery.GenerateInstanceID(metadataServiceName)
	if err := registry.Register(ctx, id, metadataServiceName, metadataServiceAddr); err != nil {
		panic(err)
	}

	return srv
}
func startRatingService(ctx context.Context, registry discovery.Registry) *grpc.Server {
	log.Println("starting movie service on" + ratingServiceName)

	h := ratingtest.NewTestRatingGRPCServer()
	l, err := net.Listen("tcp", ratingServiceaddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterRatingServiceServer(srv, h)
	go func() {
		if err := srv.Serve(l); err != nil {
			panic(err)
		}

	}()
	id := discovery.GenerateInstanceID(ratingServiceName)
	if err := registry.Register(ctx, id, ratingServiceName, ratingServiceaddr); err != nil {
		panic(err)
	}

	return srv
}
func startMovieService(ctx context.Context, registry discovery.Registry) *grpc.Server {
	log.Println("starting movie service on" + movieServiceName)

	h := movietest.NewTestMovieGRPCServer(registry)
	l, err := net.Listen("tcp", movieServiceAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMovieServiceServer(srv, h)
	go func() {
		if err := srv.Serve(l); err != nil {
			panic(err)
		}

	}()
	id := discovery.GenerateInstanceID(movieServiceName)
	if err := registry.Register(ctx, id, movieServiceName, movieServiceAddr); err != nil {
		panic(err)
	}

	return srv
}
