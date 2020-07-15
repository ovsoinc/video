package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"storj.io/uplink"
)

func main() {
	fmt.Printf("HellO!")

	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	access, err := uplink.ParseAccess("13P8jzLQy7NaRdY3M5vgrupZ9YJg1nxttMajudFjTeQVtPKTyBAaYnjAD6yT9KQA4ahwzoFKdJE5WAY8NZrpqQqRxABaRi5ZrjEGckkL9dLRbfxVw94K6J74bBWX8hg1whgF4XJ6XUb7H2Fkr1M8NjmirVycAAcc4bxCJYSBX3Fq6CGkeboA483mjEdk9RmPnh71nYhzCBVmncjvo2y8deez4a5MPkYMyamuxt8aE3upUgR518ZfVyGXGUCTHZmw7uuzF8Kh7eaVN19BFuRqQGhsSRWFw6qVPtDq4xgrdmFHMCZVRmat6xL")

	bucket := "video"

	if err != nil {
		fmt.Errorf("could not request access grant: %v", err)
		panic(err)
	}

	project, err := uplink.OpenProject(context.Background(), access)

	if err != nil {
		fmt.Errorf("could not open project: %v", err)
		panic(err)
	}

	bytes := r.Header.Get("Range")

	fmt.Printf("range: %v", bytes)

	filename := r.URL.Path[1:]

	w.Header().Set("Content-Type", "video/mp4")

	object, err := project.StatObject(context.Background(), bucket, filename)
	orPanic(err)

	w.Header().Set("Content-Length", strconv.FormatInt(object.System.ContentLength, 10))

	download, err := project.DownloadObject(context.Background(), bucket, filename, nil)

	defer download.Close()

	_, err = io.Copy(w, download)
}
