// Copyright 2017 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fakestorage

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
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
	var genericData interface{}
	encoder := json.NewEncoder(w)

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := newErrorResponse(http.StatusBadRequest, "Failed to read body", nil)
		encoder.Encode(err)
		return
	}
	err = json.Unmarshal(bytes, &genericData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := newErrorResponse(http.StatusBadRequest, "Failed to unmarshal json", nil)
		encoder.Encode(err)
		return
	}

	data := genericData.(map[string]interface{})
	bucketName := data["name"].(string)

	s.CreateBucket(data["name"].(string))

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

func (s *Server) deleteBucket(w http.ResponseWriter, r *http.Request) {
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
