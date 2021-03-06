// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package broker

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc"
)

var log = logger.GetLogger("coordinator", "BrokerStorageClusterState")

const dummy = ""

type StorageClusterState struct {
	state             *models.StorageState
	taskStreams       map[string]string
	taskClientFactory rpc.TaskClientFactory
}

func newStorageClusterState(taskClientFactory rpc.TaskClientFactory) *StorageClusterState {
	return &StorageClusterState{
		taskClientFactory: taskClientFactory,
		taskStreams:       make(map[string]string),
	}
}

func (s *StorageClusterState) SetState(state *models.StorageState) {
	log.Info("set new storage cluster state")
	var needDelete []string
	for nodeID := range s.taskStreams {
		_, ok := state.ActiveNodes[nodeID]
		if !ok {
			needDelete = append(needDelete, nodeID)
		}
	}

	for _, nodeID := range needDelete {
		s.taskClientFactory.CloseTaskClient(nodeID)
		delete(s.taskStreams, nodeID)
	}

	for nodeID, node := range state.ActiveNodes {
		// create a new client stream
		if err := s.taskClientFactory.CreateTaskClient(node.Node); err != nil {
			log.Error("create task client stream",
				logger.String("target", (&node.Node).Indicator()), logger.Error(err))
			s.taskClientFactory.CloseTaskClient(nodeID)
			delete(s.taskStreams, nodeID)
			continue
		}
		s.taskStreams[nodeID] = dummy
	}

	s.state = state
	log.Info("set new storage cluster successfully")
}

func (s *StorageClusterState) close() {
	log.Info("start close storage cluster state")
	for nodeID := range s.taskStreams {
		s.taskClientFactory.CloseTaskClient(nodeID)
		delete(s.taskStreams, nodeID)
	}
	log.Info("close storage cluster state successfully")
}
