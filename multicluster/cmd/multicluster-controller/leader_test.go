// Copyright 2022 Antrea Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main under directory cmd parses and validates user input,
// instantiates and initializes objects imported from pkg, and runs
// the process.

package main

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/config/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	k8smcsv1alpha1 "sigs.k8s.io/mcs-api/pkg/apis/v1alpha1"

	mcsv1alpha1 "antrea.io/antrea/multicluster/apis/multicluster/v1alpha1"
	"antrea.io/antrea/multicluster/test/mocks"
)

func initMockManager(mockManager *mocks.MockManager) {
	newScheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(newScheme))
	utilruntime.Must(k8smcsv1alpha1.AddToScheme(newScheme))
	utilruntime.Must(mcsv1alpha1.AddToScheme(newScheme))
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects().Build()

	mockManager.EXPECT().GetWebhookServer().Return(&webhook.Server{}).AnyTimes()
	mockManager.EXPECT().GetWebhookServer().Return(&webhook.Server{}).AnyTimes()
	mockManager.EXPECT().GetClient().Return(fakeClient).AnyTimes()
	mockManager.EXPECT().GetScheme().Return(newScheme).AnyTimes()
	mockManager.EXPECT().GetControllerOptions().Return(v1alpha1.ControllerConfigurationSpec{}).AnyTimes()
	mockManager.EXPECT().GetLogger().Return(klog.NewKlogr()).AnyTimes()
	mockManager.EXPECT().SetFields(gomock.Any()).Return(nil).AnyTimes()
	mockManager.EXPECT().Add(gomock.Any()).Return(nil).AnyTimes()
	mockManager.EXPECT().Start(gomock.Any()).Return(nil).AnyTimes()
	mockManager.EXPECT().GetConfig().Return(&rest.Config{}).AnyTimes()
	mockManager.EXPECT().GetRESTMapper().Return(&meta.DefaultRESTMapper{}).AnyTimes()
}

func TestRunLeader(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockLeaderManager := mocks.NewMockManager(mockCtrl)
	initMockManager(mockLeaderManager)

	testCase := struct {
		name      string
		setupFunc func(o *Options) (ctrl.Manager, error)
	}{
		name: "Start leader controller successfully",
		setupFunc: func(o *Options) (ctrl.Manager, error) {
			return mockLeaderManager, nil
		},
	}

	t.Run(testCase.name, func(t *testing.T) {
		if testCase.setupFunc != nil {
			setupManagerAndCertControllerFunc = testCase.setupFunc
		}
		ctrl.SetupSignalHandler = mockSetupSignalHandler
		err := runLeader(&Options{})
		assert.NoError(t, err, "got error when running runLeader")
	})
}
