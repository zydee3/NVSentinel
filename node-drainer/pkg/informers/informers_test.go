// Copyright (c) 2025, NVIDIA CORPORATION.  All rights reserved.
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

package informers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestFilterPodsWithGPURequests(t *testing.T) {
	podWithGPU := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "gpu-pod", Namespace: "default"},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name: "container",
					Resources: v1.ResourceRequirements{
						Limits: v1.ResourceList{
							"nvidia.com/gpu": resource.MustParse("1"),
						},
					},
				},
			},
		},
	}

	podWithPGPU := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "pgpu-pod", Namespace: "default"},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name: "container",
					Resources: v1.ResourceRequirements{
						Limits: v1.ResourceList{
							"nvidia.com/pgpu": resource.MustParse("1"),
						},
					},
				},
			},
		},
	}

	podWithoutGPU := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "cpu-pod", Namespace: "default"},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name: "container",
					Resources: v1.ResourceRequirements{
						Limits: v1.ResourceList{
							v1.ResourceCPU: resource.MustParse("1"),
						},
					},
				},
			},
		},
	}

	drainGPUPodsFlag := true
	informersInstance := &Informers{
		drainGPUPods: &drainGPUPodsFlag,
	}

	pods := []*v1.Pod{podWithGPU, podWithoutGPU, podWithPGPU}
	filteredPods := informersInstance.filterPodsWithGPURequests(pods)

	assert.Len(t, filteredPods, 2)
	assert.Equal(t, "gpu-pod", filteredPods[0].Name)
	assert.Equal(t, "pgpu-pod", filteredPods[1].Name)
}

func TestFilterPodsWithGPURequests_InitContainers(t *testing.T) {
	podWithGPUInInit := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "gpu-init-pod", Namespace: "default"},
		Spec: v1.PodSpec{
			InitContainers: []v1.Container{
				{
					Name: "init-container",
					Resources: v1.ResourceRequirements{
						Limits: v1.ResourceList{
							"nvidia.com/gpu": resource.MustParse("1"),
						},
					},
				},
			},
			Containers: []v1.Container{
				{
					Name: "container",
					Resources: v1.ResourceRequirements{
						Limits: v1.ResourceList{
							v1.ResourceCPU: resource.MustParse("1"),
						},
					},
				},
			},
		},
	}

	podWithoutGPU := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "cpu-pod", Namespace: "default"},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name: "container",
					Resources: v1.ResourceRequirements{
						Limits: v1.ResourceList{
							v1.ResourceCPU: resource.MustParse("1"),
						},
					},
				},
			},
		},
	}

	drainGPUPodsFlag := true
	informersInstance := &Informers{
		drainGPUPods: &drainGPUPodsFlag,
	}

	pods := []*v1.Pod{podWithGPUInInit, podWithoutGPU}
	filteredPods := informersInstance.filterPodsWithGPURequests(pods)

	assert.Len(t, filteredPods, 1)
	assert.Equal(t, "gpu-init-pod", filteredPods[0].Name)
}
