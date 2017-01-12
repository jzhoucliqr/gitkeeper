// Copyright 2016 Jun Zhou <zhoujun06@gmail.com>. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

// webhook to provide api endpoints to integrate with github
// so github can push events to keeper, rather than keeper pull
// periodically from github, which may have issues due to github
// api request limits

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libkv/store"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"net/http"
	//"reflect"
)

type EventsHandler struct {
	KVStoreType string
	KVStoreUrl  string
	SecretKey   string
	port        int
}

func (h *EventsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte(h.SecretKey))
	if err != nil {
		log.Errorf("%v", err)
	}
	log.Debug(string(payload))
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Errorf("%v", err)
	}
	str, _ := json.Marshal(r.Header)
	log.Print(string(str))

	processWebHookEvents(event)
}

func processWebHookEvents(event interface{}) {
	key, err := getKeyForEvent(event)
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	value, err := json.Marshal(event)
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	store := getKVStore()
	store.Put(key, value, nil)
}

func getKeyForEvent(event interface{}) (string, error) {
	var s string
	var err error
	switch event.(type) {
	case *github.IssuesEvent:
		e := event.(*github.IssuesEvent)
		s = fmt.Sprintf("issues/%s/%d", *e.Repo.FullName, *e.Issue.Number)
	case *github.PullRequestEvent:
		e := event.(*github.PullRequestEvent)
		s = fmt.Sprintf("prs/%s/%d", *e.Repo.FullName, *e.PullRequest.Number)
	case *github.IssueCommentEvent:
		e := event.(*github.IssueCommentEvent)
		if e.Issue.PullRequestLinks != nil {
			s = fmt.Sprintf("comments/%s/%d/%d", *e.Repo.FullName, *e.Issue.Number, *e.Comment.ID)
		} else {
			err = errors.New("not supported comment for issue, only for pr: ")
		}
	default:
		err = errors.New("not supported event type: ")
		s = ""
	}

	log.Debugf("key is %s", s)
	return s, err
}

func getKVStore() store.Store {
	storeConfig := &StoreConfig{KVStoreType: "etcd", KVStoreUrl: "localhost:2379"}
	store, err := NewKVStore(storeConfig)
	if err != nil {
		log.Fatal(err)
	}
	return store
}

func startServer(h *EventsHandler) {
	log.Print(h.KVStoreType)
	log.Print(h.KVStoreUrl)
	log.Print(h.SecretKey)
	log.Print(h.port)

	http.Handle("/api", h)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", h.port), nil))

}

func addFlags(cmd *cobra.Command, handler *EventsHandler) {
	cmd.Flags().StringVarP(&handler.KVStoreType, "kvstoreType", "t", "etcd", "etcd or consul")
	cmd.Flags().StringVarP(&handler.KVStoreUrl, "kvstoreUrl", "u", "localhost:2379", "url for etcd or consul")
	cmd.Flags().StringVarP(&handler.SecretKey, "secretKey", "s", "", "webhook secret key")
	cmd.Flags().IntVarP(&handler.port, "port", "p", 8888, "port server to listen on")
}

func main() {
	log.SetLevel(log.DebugLevel)
	handler := &EventsHandler{}
	rootCmd := &cobra.Command{
		Use:   "webhook",
		Short: "webhook to receive events from github",
		Long:  "webhook to receive events from github",
		Run: func(cmd *cobra.Command, args []string) {
			startServer(handler)
		},
	}

	addFlags(rootCmd, handler)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
