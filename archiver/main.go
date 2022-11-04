package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/lannparty/kob/internal/archivers"

	_ "github.com/mattn/go-sqlite3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// Initialize SQLite
	log.Print("Initializing database...")
	const file = "/opt/kob/obituaries.db"
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err.Error())
	}
	log.Print("Success!")

	log.Print("Initializing k8s client...")
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	log.Print("Success!")

	// Watch all pods.
	log.Print("Watching...")
	watch, err := clientset.CoreV1().Pods("").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	if err != nil {
		log.Fatal(err.Error())
	}

	// Process pods.
	for event := range watch.ResultChan() {
		pod := event.Object.(*v1.Pod)
		archivers.ArchivePodManifest(pod, db)
	}
}
