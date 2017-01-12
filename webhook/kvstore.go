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

import (
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
	"time"
)

type StoreConfig struct {
	KVStoreType string //etcd, consul
	KVStoreUrl  string //localhost:2379
}

func NewKVStore(config *StoreConfig) (store.Store, error) {
	etcd.Register()
	backend := store.ETCD
	backendUrl := config.KVStoreUrl

	kvstore, err := libkv.NewStore(
		backend,
		[]string{backendUrl},
		&store.Config{ConnectionTimeout: 10 * time.Second},
	)

	return kvstore, err
}
