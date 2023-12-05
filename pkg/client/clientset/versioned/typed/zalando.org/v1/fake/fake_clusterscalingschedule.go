/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	zalandoorgv1 "github.com/zalando-incubator/kube-metrics-adapter/pkg/apis/zalando.org/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeClusterScalingSchedules implements ClusterScalingScheduleInterface
type FakeClusterScalingSchedules struct {
	Fake *FakeZalandoV1
}

var clusterscalingschedulesResource = schema.GroupVersionResource{Group: "zalando.org", Version: "v1", Resource: "clusterscalingschedules"}

var clusterscalingschedulesKind = schema.GroupVersionKind{Group: "zalando.org", Version: "v1", Kind: "ClusterScalingSchedule"}

// Get takes name of the clusterScalingSchedule, and returns the corresponding clusterScalingSchedule object, and an error if there is any.
func (c *FakeClusterScalingSchedules) Get(ctx context.Context, name string, options v1.GetOptions) (result *zalandoorgv1.ClusterScalingSchedule, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(clusterscalingschedulesResource, name), &zalandoorgv1.ClusterScalingSchedule{})
	if obj == nil {
		return nil, err
	}
	return obj.(*zalandoorgv1.ClusterScalingSchedule), err
}

// List takes label and field selectors, and returns the list of ClusterScalingSchedules that match those selectors.
func (c *FakeClusterScalingSchedules) List(ctx context.Context, opts v1.ListOptions) (result *zalandoorgv1.ClusterScalingScheduleList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(clusterscalingschedulesResource, clusterscalingschedulesKind, opts), &zalandoorgv1.ClusterScalingScheduleList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &zalandoorgv1.ClusterScalingScheduleList{ListMeta: obj.(*zalandoorgv1.ClusterScalingScheduleList).ListMeta}
	for _, item := range obj.(*zalandoorgv1.ClusterScalingScheduleList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterScalingSchedules.
func (c *FakeClusterScalingSchedules) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(clusterscalingschedulesResource, opts))
}

// Create takes the representation of a clusterScalingSchedule and creates it.  Returns the server's representation of the clusterScalingSchedule, and an error, if there is any.
func (c *FakeClusterScalingSchedules) Create(ctx context.Context, clusterScalingSchedule *zalandoorgv1.ClusterScalingSchedule, opts v1.CreateOptions) (result *zalandoorgv1.ClusterScalingSchedule, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(clusterscalingschedulesResource, clusterScalingSchedule), &zalandoorgv1.ClusterScalingSchedule{})
	if obj == nil {
		return nil, err
	}
	return obj.(*zalandoorgv1.ClusterScalingSchedule), err
}

// Update takes the representation of a clusterScalingSchedule and updates it. Returns the server's representation of the clusterScalingSchedule, and an error, if there is any.
func (c *FakeClusterScalingSchedules) Update(ctx context.Context, clusterScalingSchedule *zalandoorgv1.ClusterScalingSchedule, opts v1.UpdateOptions) (result *zalandoorgv1.ClusterScalingSchedule, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(clusterscalingschedulesResource, clusterScalingSchedule), &zalandoorgv1.ClusterScalingSchedule{})
	if obj == nil {
		return nil, err
	}
	return obj.(*zalandoorgv1.ClusterScalingSchedule), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeClusterScalingSchedules) UpdateStatus(ctx context.Context, clusterScalingSchedule *zalandoorgv1.ClusterScalingSchedule, opts v1.UpdateOptions) (*zalandoorgv1.ClusterScalingSchedule, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(clusterscalingschedulesResource, "status", clusterScalingSchedule), &zalandoorgv1.ClusterScalingSchedule{})
	if obj == nil {
		return nil, err
	}
	return obj.(*zalandoorgv1.ClusterScalingSchedule), err
}

// Delete takes name of the clusterScalingSchedule and deletes it. Returns an error if one occurs.
func (c *FakeClusterScalingSchedules) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(clusterscalingschedulesResource, name, opts), &zalandoorgv1.ClusterScalingSchedule{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClusterScalingSchedules) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(clusterscalingschedulesResource, listOpts)

	_, err := c.Fake.Invokes(action, &zalandoorgv1.ClusterScalingScheduleList{})
	return err
}

// Patch applies the patch and returns the patched clusterScalingSchedule.
func (c *FakeClusterScalingSchedules) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *zalandoorgv1.ClusterScalingSchedule, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(clusterscalingschedulesResource, name, pt, data, subresources...), &zalandoorgv1.ClusterScalingSchedule{})
	if obj == nil {
		return nil, err
	}
	return obj.(*zalandoorgv1.ClusterScalingSchedule), err
}
