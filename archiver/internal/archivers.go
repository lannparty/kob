package archivers

import (
	"database/sql"
	"encoding/json"
	"log"

	v1 "k8s.io/api/core/v1"
)

func ArchivePodManifest(pod *v1.Pod, db *sql.DB) {
	if pod.ObjectMeta.DeletionTimestamp != nil {
		marshalledPod, err := json.Marshal(pod)
		if err != nil {
			log.Print("Cannot convert pod object to JSON, pod name: ", pod.Name, ", error: ", err.Error())
		}
		_, err = db.Exec("INSERT INTO pods(name, uid, manifest) VALUES(?, ?, ?) ON CONFLICT(name) DO UPDATE SET manifest = ?", pod.Name, pod.UID, string(marshalledPod), string(marshalledPod))
		if err != nil {
			if err.Error() != "UNIQUE constraint failed: pods.uid" {
				log.Print("Cannot insert pod manifest into database, pod name: ", pod.Name, ", error: ", err.Error())
			}
		} else {
			log.Print("Created entry for ", pod.Name, " ", pod.UID)
		}
	}
}

func ArchivePodLogs(pod *v1.Pod, db *sql.DB) {

}
