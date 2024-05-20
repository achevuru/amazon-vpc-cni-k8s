package config

// BuildCacheOptions returns a cache.Options struct for this controller.
/*
func BuildCacheOptions() cache.Options {
	fieldSelector := fields.Set{"metadata.name": "ip-192-168-71-163.us-west-2.compute.internal"}.AsSelector().String()
	return cache.Options{
		ReaderFailOnMissingInformer: true,
		ByObject: map[client.Object]cache.ByObject{
			&cniNode.CNINode{}: {
				// We don't mutate CNI Node objects, thus we can safely disable DeepCopy to optimize memory usage.
				UnsafeDisableDeepCopy: awssdk.Bool(true),
				Field:                 fieldSelector,
			},
		},
	}
}
*/
