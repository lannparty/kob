package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"path/filepath"

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
	const file = "/opt/kob/obituaries.db"
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err.Error())
	}
	log.Print("Success!")

	// --
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// --

	// --
	//log.Print("Initializing k8s client...")
	//// creates the in-cluster config
	//config, err := rest.InClusterConfig()
	//if err != nil {
	//	panic(err.Error())
	//}
	// --

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

			podEvents, err := clientset.CoreV1().Events(pod.Namespace).List(context.TODO(), metav1.ListOptions{FieldSelector: "involvedObject.name=" + pod.Name, TypeMeta: metav1.TypeMeta{Kind: "Pod"}})
			if err != nil {
				log.Print("Error retrieving pod events: ", err)
			}

			marshalledPodEvents, err := json.Marshal(podEvents)
			if err != nil {
				log.Print("Error marshalling pod events", err)
			}

			node, err := clientset.CoreV1().Nodes().Get(context.TODO(), pod.Spec.NodeName, metav1.GetOptions{})
			if err != nil {
				log.Print("Error retrieving node manifest: ", err)
			}

			marshalledNode, err := json.Marshal(node)
			if err != nil {
				log.Print("Error retrieving node snapshot: ", err)
			}

			nodeEvents, err := clientset.CoreV1().Events(pod.Namespace).List(context.TODO(), metav1.ListOptions{FieldSelector: "involvedObject.name=" + pod.Spec.NodeName, TypeMeta: metav1.TypeMeta{Kind: "Node"}})
			if err != nil {
				log.Print("Error retrieving node events: ", err)
			}

			marshalledNodeEvents, err := json.Marshal(nodeEvents)
			if err != nil {
				log.Print("Error marshalling node events: ", err)
			}

			_, err = db.Exec("INSERT INTO pods(name, uid, manifest) VALUES(?, ?, ?)", pod.Name, pod.UID, string(marshalledPod))
			if err != nil {
				if err.Error() != "UNIQUE constraint failed: pods.uid" {
					log.Print("Cannot insert pod manifest into database, pod name: ", pod.Name, ", error: ", err.Error())
				}
			} else {
				log.Print("Created pod entry for ", pod.Name, " ", pod.UID)
			}

			_, err = db.Exec("INSERT INTO pod_events(uid, pod_events) VALUES(?, ?)", pod.UID, string(marshalledPodEvents))
			if err != nil {
				if err.Error() != "UNIQUE constraint failed: pods.uid" {
					log.Print("Cannot insert pod manifest into database, pod name: ", pod.Name, ", error: ", err.Error())
				}
			} else {
				log.Print("Created node entry for ", pod.Name, " ", pod.UID)
			}

			_, err = db.Exec("INSERT INTO nodes(uid, manifest) VALUES(?, ?)", pod.UID, string(marshalledNode))
			if err != nil {
				if err.Error() != "UNIQUE constraint failed: pods.uid" {
					log.Print("Cannot insert pod manifest into database, pod name: ", pod.Name, ", error: ", err.Error())
				}
			} else {
				log.Print("Created node entry for ", pod.Name, " ", pod.UID)
			}

			_, err = db.Exec("INSERT INTO node_events(uid, manifest) VALUES(?, ?)", pod.UID, string(marshalledNodeEvents))
			if err != nil {
				if err.Error() != "UNIQUE constraint failed: pods.uid" {
					log.Print("Cannot insert pod manifest into database, pod name: ", pod.Name, ", error: ", err.Error())
				}
			} else {
				log.Print("Created node entry for ", pod.Name, " ", pod.UID)
			}

		}
	}
}
