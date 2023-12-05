package collector

import (
	"errors"
	"fmt"
	"math"
	"time"

	v1 "github.com/zalando-incubator/kube-metrics-adapter/pkg/apis/zalando.org/v1"
	scheduledscaling "github.com/zalando-incubator/kube-metrics-adapter/pkg/controller/scheduledscaling"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/custom_metrics"
)

var (
	// ErrScalingScheduleNotFound is returned when a item referenced in
	// the HPA config is not in the ScalingScheduleCollectorPlugin.store.
	ErrScalingScheduleNotFound = errors.New("referenced ScalingSchedule not found")
	// ErrNotScalingScheduleFound is returned when a item returned from
	// the ScalingScheduleCollectorPlugin.store was expected to
	// be an ScalingSchedule but the type assertion failed.
	ErrNotScalingScheduleFound = errors.New("error converting returned object to ScalingSchedule")
	// ErrClusterScalingScheduleNotFound is returned when a item referenced in
	// the HPA config is not in the ClusterScalingScheduleCollectorPlugin.store.
	ErrClusterScalingScheduleNotFound = errors.New("referenced ClusterScalingSchedule not found")
	// ErrNotClusterScalingScheduleFound is returned when a item returned from
	// the ClusterScalingScheduleCollectorPlugin.store was expected to
	// be an ClusterScalingSchedule but the type assertion failed. When
	// returned the type assertion to ScalingSchedule failed too.
	ErrNotClusterScalingScheduleFound = errors.New("error converting returned object to ClusterScalingSchedule")
)

// Now is the function that returns a time.Time object representing the
// current moment. Its main implementation is the time.Now func in the
// std lib. It's used mainly for test/mock purposes.
type Now func() time.Time

// Store represent an in memory Store for the [Cluster]ScalingSchedule
// objects. Its main implementation is the [cache.cache][0] struct
// returned by the [cache.NewStore][1] function. Here it's used mainly
// for tests/mock purposes.
//
// [1]: https://pkg.go.dev/k8s.io/client-go/tools/cache#NewStore
// [0]: https://github.com/kubernetes/client-go/blob/v0.21.1/tools/cache/Store.go#L132-L140
type Store interface {
	GetByKey(key string) (item interface{}, exists bool, err error)
}

// ScalingScheduleCollectorPlugin is a collector plugin for initializing metrics
// collectors for getting ScalingSchedule configured metrics.
type ScalingScheduleCollectorPlugin struct {
	store                Store
	now                  Now
	defaultScalingWindow time.Duration
	defaultTimeZone      string
	rampSteps            int
}

// ClusterScalingScheduleCollectorPlugin is a collector plugin for initializing metrics
// collectors for getting ClusterScalingSchedule configured metrics.
type ClusterScalingScheduleCollectorPlugin struct {
	store                Store
	now                  Now
	defaultScalingWindow time.Duration
	defaultTimeZone      string
	rampSteps            int
}

// NewScalingScheduleCollectorPlugin initializes a new ScalingScheduleCollectorPlugin.
func NewScalingScheduleCollectorPlugin(store Store, now Now, defaultScalingWindow time.Duration, defaultTimeZone string, rampSteps int) (*ScalingScheduleCollectorPlugin, error) {
	return &ScalingScheduleCollectorPlugin{
		store:                store,
		now:                  now,
		defaultScalingWindow: defaultScalingWindow,
		defaultTimeZone:      defaultTimeZone,
		rampSteps:            rampSteps,
	}, nil
}

// NewClusterScalingScheduleCollectorPlugin initializes a new ClusterScalingScheduleCollectorPlugin.
func NewClusterScalingScheduleCollectorPlugin(store Store, now Now, defaultScalingWindow time.Duration, defaultTimeZone string, rampSteps int) (*ClusterScalingScheduleCollectorPlugin, error) {
	return &ClusterScalingScheduleCollectorPlugin{
		store:                store,
		now:                  now,
		defaultScalingWindow: defaultScalingWindow,
		defaultTimeZone:      defaultTimeZone,
		rampSteps:            rampSteps,
	}, nil
}

// NewCollector initializes a new scaling schedule collector from the
// specified HPA. It's the only required method to implement the
// collector.CollectorPlugin interface.
func (c *ScalingScheduleCollectorPlugin) NewCollector(hpa *autoscalingv2.HorizontalPodAutoscaler, config *MetricConfig, interval time.Duration) (Collector, error) {
	return NewScalingScheduleCollector(c.store, c.defaultScalingWindow, c.defaultTimeZone, c.rampSteps, c.now, hpa, config, interval)
}

// NewCollector initializes a new cluster wide scaling schedule
// collector from the specified HPA. It's the only required method to
// implement the collector.CollectorPlugin interface.
func (c *ClusterScalingScheduleCollectorPlugin) NewCollector(hpa *autoscalingv2.HorizontalPodAutoscaler, config *MetricConfig, interval time.Duration) (Collector, error) {
	return NewClusterScalingScheduleCollector(c.store, c.defaultScalingWindow, c.defaultTimeZone, c.rampSteps, c.now, hpa, config, interval)
}

// ScalingScheduleCollector is a metrics collector for time based
// scaling metrics.
type ScalingScheduleCollector struct {
	scalingScheduleCollector
}

// ClusterScalingScheduleCollector is a metrics collector for time based
// scaling metrics.
type ClusterScalingScheduleCollector struct {
	scalingScheduleCollector
}

// scalingScheduleCollector is a representation of the internal data
// struct used by both ClusterScalingScheduleCollector and the
// ScalingScheduleCollector.
type scalingScheduleCollector struct {
	store                Store
	now                  Now
	metric               autoscalingv2.MetricIdentifier
	objectReference      custom_metrics.ObjectReference
	hpa                  *autoscalingv2.HorizontalPodAutoscaler
	interval             time.Duration
	config               MetricConfig
	defaultScalingWindow time.Duration
	defaultTimeZone      string
	rampSteps            int
}

// NewScalingScheduleCollector initializes a new ScalingScheduleCollector.
func NewScalingScheduleCollector(store Store, defaultScalingWindow time.Duration, defaultTimeZone string, rampSteps int, now Now, hpa *autoscalingv2.HorizontalPodAutoscaler, config *MetricConfig, interval time.Duration) (*ScalingScheduleCollector, error) {
	return &ScalingScheduleCollector{
		scalingScheduleCollector{
			store:                store,
			now:                  now,
			objectReference:      config.ObjectReference,
			hpa:                  hpa,
			metric:               config.Metric,
			interval:             interval,
			config:               *config,
			defaultScalingWindow: defaultScalingWindow,
			defaultTimeZone:      defaultTimeZone,
			rampSteps:            rampSteps,
		},
	}, nil
}

// NewClusterScalingScheduleCollector initializes a new ScalingScheduleCollector.
func NewClusterScalingScheduleCollector(store Store, defaultScalingWindow time.Duration, defaultTimeZone string, rampSteps int, now Now, hpa *autoscalingv2.HorizontalPodAutoscaler, config *MetricConfig, interval time.Duration) (*ClusterScalingScheduleCollector, error) {
	return &ClusterScalingScheduleCollector{
		scalingScheduleCollector{
			store:                store,
			now:                  now,
			objectReference:      config.ObjectReference,
			hpa:                  hpa,
			metric:               config.Metric,
			interval:             interval,
			config:               *config,
			defaultScalingWindow: defaultScalingWindow,
			defaultTimeZone:      defaultTimeZone,
			rampSteps:            rampSteps,
		},
	}, nil
}

// GetMetrics is the main implementation for collector.Collector interface
func (c *ScalingScheduleCollector) GetMetrics() ([]CollectedMetric, error) {
	scalingScheduleInterface, exists, err := c.store.GetByKey(fmt.Sprintf("%s/%s", c.objectReference.Namespace, c.objectReference.Name))
	if !exists {
		return nil, ErrScalingScheduleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("unexpected error retrieving the ScalingSchedule: %s", err.Error())
	}

	scalingSchedule, ok := scalingScheduleInterface.(*v1.ScalingSchedule)
	if !ok {
		return nil, ErrNotScalingScheduleFound
	}
	return calculateMetrics(scalingSchedule.Spec, c.defaultScalingWindow, c.defaultTimeZone, c.rampSteps, c.now(), c.objectReference, c.metric)
}

// GetMetrics is the main implementation for collector.Collector interface
func (c *ClusterScalingScheduleCollector) GetMetrics() ([]CollectedMetric, error) {
	clusterScalingScheduleInterface, exists, err := c.store.GetByKey(c.objectReference.Name)
	if !exists {
		return nil, ErrClusterScalingScheduleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("unexpected error retrieving the ClusterScalingSchedule: %s", err.Error())
	}

	// The [cache.Store][0] returns the v1.ClusterScalingSchedule items as
	// a v1.ScalingSchedule when it first lists it. Once the objects are
	// updated/patched it asserts it correctly to the
	// v1.ClusterScalingSchedule type. It means we have to handle both
	// cases.
	// TODO(jonathanbeber): Identify why it happens and fix in the upstream.
	//
	// [0]: https://github.com/kubernetes/client-go/blob/v0.21.1/tools/cache/Store.go#L132-L140
	var clusterScalingSchedule v1.ClusterScalingSchedule
	scalingSchedule, ok := clusterScalingScheduleInterface.(*v1.ScalingSchedule)
	if !ok {
		css, ok := clusterScalingScheduleInterface.(*v1.ClusterScalingSchedule)
		if !ok {
			return nil, ErrNotClusterScalingScheduleFound
		}
		clusterScalingSchedule = *css
	} else {
		clusterScalingSchedule = v1.ClusterScalingSchedule(*scalingSchedule)
	}

	return calculateMetrics(clusterScalingSchedule.Spec, c.defaultScalingWindow, c.defaultTimeZone, c.rampSteps, c.now(), c.objectReference, c.metric)
}

// Interval returns the interval at which the collector should run.
func (c *ScalingScheduleCollector) Interval() time.Duration {
	return c.interval
}

// Interval returns the interval at which the collector should run.
func (c *ClusterScalingScheduleCollector) Interval() time.Duration {
	return c.interval
}

func calculateMetrics(spec v1.ScalingScheduleSpec, defaultScalingWindow time.Duration, defaultTimeZone string, rampSteps int, now time.Time, objectReference custom_metrics.ObjectReference, metric autoscalingv2.MetricIdentifier) ([]CollectedMetric, error) {
	scalingWindowDuration := defaultScalingWindow
	if spec.ScalingWindowDurationMinutes != nil {
		scalingWindowDuration = time.Duration(*spec.ScalingWindowDurationMinutes) * time.Minute
	}
	if scalingWindowDuration < 0 {
		return nil, fmt.Errorf("scaling window duration cannot be negative")
	}

	value := int64(0)
	for _, schedule := range spec.Schedules {
		startTime, endTime, err := scheduledscaling.ScheduleStartEnd(now, schedule, defaultTimeZone)
		if err != nil {
			return nil, err
		}
		value = maxInt64(value, valueForEntry(now, startTime, endTime, scalingWindowDuration, rampSteps, schedule.Value))
	}

	return []CollectedMetric{
		{
			Type:      autoscalingv2.ObjectMetricSourceType,
			Namespace: objectReference.Namespace,
			Custom: custom_metrics.MetricValue{
				DescribedObject: objectReference,
				Timestamp:       metav1.Time{Time: now},
				Value:           *resource.NewMilliQuantity(value*1000, resource.DecimalSI),
				Metric:          custom_metrics.MetricIdentifier(metric),
			},
		},
	}, nil
}

func valueForEntry(timestamp time.Time, startTime time.Time, endTime time.Time, scalingWindowDuration time.Duration, rampSteps int, value int64) int64 {
	scaleUpStart := startTime.Add(-scalingWindowDuration)
	scaleUpEnd := endTime.Add(scalingWindowDuration)

	if scheduledscaling.Between(timestamp, startTime, endTime) {
		return value
	}
	if scheduledscaling.Between(timestamp, scaleUpStart, startTime) {
		return scaledValue(timestamp, scaleUpStart, scalingWindowDuration, rampSteps, value)
	}
	if scheduledscaling.Between(timestamp, endTime, scaleUpEnd) {
		return scaledValue(scaleUpEnd, timestamp, scalingWindowDuration, rampSteps, value)
	}
	return 0
}

// The HPA has a rule to do not scale up or down if the change in the
// metric is less than 10% (by default) of the current value. We will
// use buckets of time using the floor of each as the returned metric.
// Any config greater or equal to 10 buckets must guarantee changes
// bigger than 10%.
func scaledValue(timestamp time.Time, startTime time.Time, scalingWindowDuration time.Duration, rampSteps int, value int64) int64 {
	if scalingWindowDuration == 0 {
		return 0
	}

	steps := float64(rampSteps)

	requiredPercentage := math.Abs(float64(timestamp.Sub(startTime))) / float64(scalingWindowDuration)
	return int64(math.Floor(requiredPercentage*steps) * (float64(value) / steps))
}

func maxInt64(i1, i2 int64) int64 {
	if i1 > i2 {
		return i1
	}
	return i2
}
