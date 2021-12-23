package main

import (
	"blog-gunk/cms/handler"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	cgp "blog-gunk/gunk/v1/category"
	pgp "blog-gunk/gunk/v1/post"
)

func main() {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("cms/env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Printf("error loading configuration: %v", err)
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	store := sessions.NewCookieStore([]byte(config.GetString("session.secret")))

	host, port := config.GetString("grpc.host"),config.GetString("grpc.port")
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil{
		log.Fatal(err)
	}

	csc := cgp.NewCategoryServiceClient(conn)
	psc := pgp.NewPostServiceClient(conn)
	r := handler.New(decoder, store, csc, psc)

	host, port = config.GetString("server.host"),config.GetString("server.port")
	log.Printf("Server starting on: http://%s:%s", host, port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r); err != nil {
		log.Fatal(err)
	}
}