package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/mattn/go-sqlite3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiWatch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	w "k8s.io/client-go/tools/watch"
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

	timeOut := int64(120)
	watchFunc := func(options metav1.ListOptions) (apiWatch.Interface, error) {
		return clientset.CoreV1().Pods("").Watch(context.Background(), metav1.ListOptions{TimeoutSeconds: &timeOut})
	}

	w, err := w.NewRetryWatcher("1", &cache.ListWatch{WatchFunc: watchFunc})
	if err != nil {
		panic(err)
	}

	// Insert manifest into database for all terminating pods
	for event := range w.ResultChan() {
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
