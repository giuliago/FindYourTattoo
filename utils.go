package main

import("io"
	   "log" 	 
	   "encoding/json")

func DecodeJson(payload io.Reader, entity interface{}) (err error){

	decoder := json.NewDecoder(payload)

	err = decoder.Decode(entity)

	if err != nil {
		log.Printf("[ERROR] could not convert json, because: %v", err)
		return
	}

	return
}


func EncodeJson(w io.Writer, entity interface{}) (err error){

	encoder := json.NewEncoder(w)

	err = encoder.Encode(entity)

	if err != nil {
		log.Printf("[ERROR] could not convert json, because: %v", err)
		return
	}

	return
}
