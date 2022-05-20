package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	// Initialize SQLite
	log.Print("Initializing database...")
	const file = "obituary.db"
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err.Error())
	}
	log.Print("Success!")

	log.Print("Initializing tables...")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS 'terminated'(name TEXT, manifest TEXT)")
	if err != nil {
		panic(err.Error())
	}

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	log.Print("Success!")

	log.Print("Initializing Client...")
	// use the current context in kubeconfg
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	log.Print("Success!")

	log.Print("Watching...")
	// Insert manifest into database for all terminating pods.
	watch, err := clientset.CoreV1().Pods("").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	if err != nil {
		log.Fatal(err.Error())
	}
	go func() {
		for event := range watch.ResultChan() {
			pod := event.Object.(*v1.Pod)
			if pod.ObjectMeta.DeletionTimestamp != nil {
				marshalledPod, err := json.Marshal(pod)
				if err != nil {
					log.Print("Cannot convert pod object to JSON, pod name: ", pod.Name, ", error: ", err.Error())
				}
				results, err := db.Exec("INSERT INTO terminated(name, manifest) VALUES(?, ?)", pod.Name, string(marshalledPod))
				if err != nil {
					log.Print("Cannot insert pod manifest into database, pod name: ", pod.Name, ", error: ", err.Error())
				}
				fmt.Println(results)
			}
		}
	}()
	time.Sleep(600 * time.Second)
}
