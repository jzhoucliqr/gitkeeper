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
	log "github.com/Sirupsen/logrus"
	"github.com/jzhoucliqr/gitkeeper/api"
	"github.com/jzhoucliqr/gitkeeper/config"
	"github.com/jzhoucliqr/gitkeeper/keepers"
	"github.com/jzhoucliqr/gitkeeper/kvstore"
	"github.com/spf13/cobra"
)

func commandInit(cmd *cobra.Command, conf *config.KeeperConfig) {
	cmd.Flags().StringVarP(&conf.GitUserToken, "gituser-token", "t", "", "github token for the user which you want keeper to be")
	cmd.Flags().StringVarP(&conf.KVStoreType, "kvstore-type", "s", "etcd", "kvstore type [etcd|consul]")
	cmd.Flags().StringVarP(&conf.KVStoreUrl, "kvstore-url", "u", "localhost:2379", "url for kvstore [localhost:2379]")
}

func startKeepers(conf *config.KeeperConfig) {
	// start AipHandler
	startApiServer(conf)

	// run all keepers
	runKeepers(conf)

	// start to read from kvstore
	kvstore.ConsumeEvents(conf)

}

func startApiServer(conf *config.KeeperConfig) {
	h := &api.ApiHandler{KeeperConfig: conf}
	go api.StartApiServer(h)
}

func runKeepers(conf *config.KeeperConfig) {
	kps := keepers.GetEnabledKeepers()
	for _, keeper := range kps {
		go keeper.KeepDoing(conf.GitUserToken)
	}
}

func AppStart() {
	log.SetLevel(log.DebugLevel)
	conf := &config.KeeperConfig{}

	rootCmd := &cobra.Command{
		Use:   "gitkeeper",
		Short: "Keeper of your Git projects",
		Long: `Keeper of your Git projects.
Can handle cla / merge PR etc`,
		Run: func(cmd *cobra.Command, args []string) {
			startKeepers(conf)
		},
	}

	commandInit(rootCmd, conf)
	rootCmd.Execute()

}
