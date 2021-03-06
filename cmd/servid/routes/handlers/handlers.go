package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

func jsonResponse(entity interface{}, w io.Writer) {
	jsonResponse, _ := json.Marshal(entity)

	_, err := fmt.Fprintf(w, string(jsonResponse))

	if err != nil {
		log.Println("Can't write response in response writer", err)
	}
}

func jsonCollectionResponse(key string, entity interface{}, w io.Writer) {
	collectionMap := map[string]interface{}{
		key: entity,
	}

	jsonResponse(collectionMap, w)
}
