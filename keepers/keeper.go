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

package keepers

import (
	"github.com/jzhoucliqr/gitkeeper/github"
)

// keeper defines the interface for each keeper
type Keeper interface {
	Name() string

	// Keeper receive each github object, like if an PR, try merge it
	Receive(obj *github.GitObject)

	// keep doing
	KeepDoing(token string)
}

var keepers []Keeper

func GetEnabledKeepers() []Keeper {
	if keepers == nil {
		keepers = []Keeper{
			&ClaKeeper{Ch: make(chan *github.GitObject)},
			&MergerKeeper{Ch: make(chan *github.GitObject)},
		}
	}

	return keepers
}

func RegisterKeeper() {
}
