//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha2

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClientIntents) DeepCopyInto(out *ClientIntents) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(IntentsSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(IntentsStatus)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClientIntents.
func (in *ClientIntents) DeepCopy() *ClientIntents {
	if in == nil {
		return nil
	}
	out := new(ClientIntents)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClientIntents) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClientIntentsList) DeepCopyInto(out *ClientIntentsList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClientIntents, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClientIntentsList.
func (in *ClientIntentsList) DeepCopy() *ClientIntentsList {
	if in == nil {
		return nil
	}
	out := new(ClientIntentsList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClientIntentsList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HTTPResource) DeepCopyInto(out *HTTPResource) {
	*out = *in
	if in.Methods != nil {
		in, out := &in.Methods, &out.Methods
		*out = make([]HTTPMethod, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HTTPResource.
func (in *HTTPResource) DeepCopy() *HTTPResource {
	if in == nil {
		return nil
	}
	out := new(HTTPResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Intent) DeepCopyInto(out *Intent) {
	*out = *in
	if in.Topics != nil {
		in, out := &in.Topics, &out.Topics
		*out = make([]KafkaTopic, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.HTTPResources != nil {
		in, out := &in.HTTPResources, &out.HTTPResources
		*out = make([]HTTPResource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Intent.
func (in *Intent) DeepCopy() *Intent {
	if in == nil {
		return nil
	}
	out := new(Intent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IntentsSpec) DeepCopyInto(out *IntentsSpec) {
	*out = *in
	out.Service = in.Service
	if in.Calls != nil {
		in, out := &in.Calls, &out.Calls
		*out = make([]Intent, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IntentsSpec.
func (in *IntentsSpec) DeepCopy() *IntentsSpec {
	if in == nil {
		return nil
	}
	out := new(IntentsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IntentsStatus) DeepCopyInto(out *IntentsStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IntentsStatus.
func (in *IntentsStatus) DeepCopy() *IntentsStatus {
	if in == nil {
		return nil
	}
	out := new(IntentsStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaServerConfig) DeepCopyInto(out *KafkaServerConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaServerConfig.
func (in *KafkaServerConfig) DeepCopy() *KafkaServerConfig {
	if in == nil {
		return nil
	}
	out := new(KafkaServerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KafkaServerConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaServerConfigList) DeepCopyInto(out *KafkaServerConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]KafkaServerConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaServerConfigList.
func (in *KafkaServerConfigList) DeepCopy() *KafkaServerConfigList {
	if in == nil {
		return nil
	}
	out := new(KafkaServerConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KafkaServerConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaServerConfigSpec) DeepCopyInto(out *KafkaServerConfigSpec) {
	*out = *in
	out.Service = in.Service
	out.TLS = in.TLS
	if in.Topics != nil {
		in, out := &in.Topics, &out.Topics
		*out = make([]TopicConfig, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaServerConfigSpec.
func (in *KafkaServerConfigSpec) DeepCopy() *KafkaServerConfigSpec {
	if in == nil {
		return nil
	}
	out := new(KafkaServerConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaServerConfigStatus) DeepCopyInto(out *KafkaServerConfigStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaServerConfigStatus.
func (in *KafkaServerConfigStatus) DeepCopy() *KafkaServerConfigStatus {
	if in == nil {
		return nil
	}
	out := new(KafkaServerConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaTopic) DeepCopyInto(out *KafkaTopic) {
	*out = *in
	if in.Operations != nil {
		in, out := &in.Operations, &out.Operations
		*out = make([]KafkaOperation, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaTopic.
func (in *KafkaTopic) DeepCopy() *KafkaTopic {
	if in == nil {
		return nil
	}
	out := new(KafkaTopic)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Service) DeepCopyInto(out *Service) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Service.
func (in *Service) DeepCopy() *Service {
	if in == nil {
		return nil
	}
	out := new(Service)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TLSSource) DeepCopyInto(out *TLSSource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TLSSource.
func (in *TLSSource) DeepCopy() *TLSSource {
	if in == nil {
		return nil
	}
	out := new(TLSSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopicConfig) DeepCopyInto(out *TopicConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopicConfig.
func (in *TopicConfig) DeepCopy() *TopicConfig {
	if in == nil {
		return nil
	}
	out := new(TopicConfig)
	in.DeepCopyInto(out)
	return out
}