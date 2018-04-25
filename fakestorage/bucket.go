// Copyright 2017 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fakestorage

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"fmt"
	"io/ioutil"
)

// CreateBucket creates a bucket inside the server, so any API calls that
// require the bucket name will recognize this bucket.
//
// If the bucket already exists, this method does nothing.
func (s *Server) CreateBucket(name string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.buckets[name]; !ok {
		s.buckets[name] = nil
	}
}

func (s *Server) createBucket(w http.ResponseWriter, r *http.Request) {
	var data2 interface{}

	bytes, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(bytes, &data2)
	fmt.Println(data2)
	fmt.Println(string(bytes))
	data := data2.(map[string]interface{})
	bucketName := data["name"].(string)

	s.CreateBucket(data["name"].(string))

	encoder := json.NewEncoder(w)
	resp := newBucketResponse(bucketName)
	w.WriteHeader(http.StatusOK)
	encoder.Encode(resp)
}

func (s *Server) listBuckets(w http.ResponseWriter, r *http.Request) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	bucketNames := make([]string, 0, len(s.buckets))
	for name := range s.buckets {
		bucketNames = append(bucketNames, name)
	}
	resp := newListBucketsResponse(bucketNames)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) deleteBucket(w http.ResponseWriter, r *http.Request){
	bucketName := mux.Vars(r)["bucketName"]
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	encoder := json.NewEncoder(w)
	if _, ok := s.buckets[bucketName]; !ok {
		w.WriteHeader(http.StatusNotFound)
		err := newErrorResponse(http.StatusNotFound, "Not found", nil)
		encoder.Encode(err)
		return
	}
	delete(s.buckets, bucketName)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) getBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := mux.Vars(r)["bucketName"]
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	encoder := json.NewEncoder(w)
	if _, ok := s.buckets[bucketName]; !ok {
		w.WriteHeader(http.StatusNotFound)
		err := newErrorResponse(http.StatusNotFound, "Not found", nil)
		encoder.Encode(err)
		return
	}
	resp := newBucketResponse(bucketName)
	w.WriteHeader(http.StatusOK)
	encoder.Encode(resp)
}
