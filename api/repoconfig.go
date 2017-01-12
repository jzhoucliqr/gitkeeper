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
package api

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/jzhoucliqr/gitkeeper/config"
	"github.com/jzhoucliqr/gitkeeper/github"
	"github.com/jzhoucliqr/gitkeeper/kvstore"
	"io/ioutil"
	"net/http"
)

type ApiHandler struct {
	KeeperConfig *config.KeeperConfig
}

func (h *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Debug("failed to read request")
		return
	}

	conf := &config.RepoConfig{}
	err = json.Unmarshal(payload, conf)
	if err != nil {
		log.Debug("failed to unmarshal repoconfig")
		return
	}

	h.addRepoConfig(conf)

	h.createWebHook(conf)
}

func (h *ApiHandler) addRepoConfig(conf *config.RepoConfig) {
	if !conf.IsValid() {
		log.Debug("invalid repo config")
		return
	}

	h.persistRepoConfig(conf)
}

func (h *ApiHandler) createWebHook(conf *config.RepoConfig) {
	hook, err := github.CreateWebHook(h.KeeperConfig, conf)
	if err != nil {
		log.Debug("failed to create web hook")
	}

	if hook != nil {
		log.Debugf("hook: %v\n", *hook)
	}
}

func (h *ApiHandler) persistRepoConfig(conf *config.RepoConfig) {
	//storeConf := &config.KeeperConfig{KVStoreType: "etcd", KVStoreUrl: "localhost:2379"}
	key := fmt.Sprintf("repo/%s", conf.Repo.FullName)
	payload, err := json.Marshal(*conf)
	if err != nil {
		log.Debug("failed to marshal repoconfig")
		return
	}

	store, _ := kvstore.NewKVStore(h.KeeperConfig)

	store.Put(key, payload, nil)
}

func StartApiServer(h *ApiHandler) {
	log.Debug("start api server")
	http.Handle("/api", h)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 9999), nil))
	log.Debug("api server started")
}
