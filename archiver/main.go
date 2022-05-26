package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

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

	log.Print("Initializing tables...")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS pods(name TEXT, uid TEXT UNIQUE, manifest TEXT)")
	if err != nil {
		panic(err.Error())
	}

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

	// Insert manifest into database for all terminating pods.
	log.Print("Watching...")
	watch, err := clientset.CoreV1().Pods("").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	if err != nil {
		log.Fatal(err.Error())
	}
	for event := range watch.ResultChan() {
		pod := event.Object.(*v1.Pod)
		if pod.ObjectMeta.DeletionTimestamp != nil {
			marshalledPod, err := json.Marshal(pod)
			if err != nil {
				log.Print("Cannot convert pod object to JSON, pod name: ", pod.Name, ", error: ", err.Error())
			}
			_, err = db.Exec("INSERT INTO pods(name, uid, manifest) VALUES(?, ?, ?)", pod.Name, pod.UID, string(marshalledPod))
			if err != nil {
				if err.Error() != "UNIQUE constraint failed: pods.uid" {
					log.Print("Cannot insert pod manifest into database, pod name: ", pod.Name, ", error: ", err.Error())
				}
			} else {
				log.Print("Created entry for ", pod.Name, " ", pod.UID)
			}
		}
	}
}
