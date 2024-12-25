/*
Copyright 2024 The CodeFuture Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package podanalyzer

import (
	"time"

	api_v1 "k8s.io/api/core/v1"

	ctlutil "sigs.k8s.io/kubefed/pkg/controller/utils"
)

type PodAnalysisResult struct {
	// Total number of pods created.
	Total int
	// Number of pods that are running and ready.
	RunningAndReady int
	// Number of pods that have been in unschedulable state for UnshedulableThreshold seconds.
	Unschedulable int

	// TODO: Handle other scenarios like pod waiting too long for scheduler etc.
}

const (
	// TODO: make it configurable
	UnschedulableThreshold = 60 * time.Second
)

// AnalyzePods calculates how many pods from the list are in one of
// the meaningful (from the replica set perspective) states. This function is
// a temporary workaround against the current lack of ownerRef in pods.
func AnalyzePods(podList *api_v1.PodList, currentTime time.Time) (PodAnalysisResult, ctlutil.ReconciliationStatus) {
	result := PodAnalysisResult{}
	unschedulableRightNow := 0
	for _, pod := range podList.Items {
		result.Total++
		for _, condition := range pod.Status.Conditions {
			if pod.Status.Phase == api_v1.PodRunning {
				if condition.Type == api_v1.PodReady {
					result.RunningAndReady++
				}
			} else if condition.Type == api_v1.PodScheduled &&
				condition.Status == api_v1.ConditionFalse &&
				condition.Reason == api_v1.PodReasonUnschedulable {
				unschedulableRightNow++
				if condition.LastTransitionTime.Add(UnschedulableThreshold).Before(currentTime) {
					result.Unschedulable++
				}
			}
		}
	}
	if unschedulableRightNow != result.Unschedulable {
		// We get the reconcile event almost immediately after  the status of a
		// pod changes, however we will not consider the unschedulable pods as
		// unschedulable immediately (until 60 secs), because we don't  want to
		// change state frequently (it can lead to continuously moving replicas
		// around). We need to reconcile again after a timeout. We use the return
		// status to indicate retry for reconcile.
		return result, ctlutil.StatusNeedsRecheck
	}

	return result, ctlutil.StatusAllOK
}
