package kubernetes

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/yaml"

	ctrl "sigs.k8s.io/controller-runtime/pkg/client"
)

// LoadResourceFromYaml loads a k8s resource from a yaml definition
func LoadResourceFromYaml(scheme *runtime.Scheme, data string) (ctrl.Object, error) {
	source := []byte(data)
	jsonSource, err := yaml.ToJSON(source)
	if err != nil {
		return nil, err
	}
	u := unstructured.Unstructured{}
	err = u.UnmarshalJSON(jsonSource)
	if err != nil {
		return nil, err
	}
	ro, err := runtimeObjectFromUnstructured(scheme, &u)
	if err != nil {
		return nil, err
	}
	if o, ok := ro.(ctrl.Object); !ok {
		return nil, err
	} else {
		return o, nil
	}
}

// LoadRawResourceFromYaml loads a k8s resource from a yaml definition without making assumptions on the underlying type
func LoadRawResourceFromYaml(data string) (runtime.Object, error) {
	source := []byte(data)
	jsonSource, err := yaml.ToJSON(source)
	if err != nil {
		return nil, err
	}
	var objmap map[string]interface{}
	if err = json.Unmarshal(jsonSource, &objmap); err != nil {
		return nil, err
	}
	return &unstructured.Unstructured{
		Object: objmap,
	}, nil
}

func runtimeObjectFromUnstructured(scheme *runtime.Scheme, u *unstructured.Unstructured) (runtime.Object, error) {
	gvk := u.GroupVersionKind()
	codecs := serializer.NewCodecFactory(scheme)
	decoder := codecs.UniversalDecoder(gvk.GroupVersion())

	b, err := u.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("error running MarshalJSON on unstructured object: %v", err)
	}
	ro, _, err := decoder.Decode(b, &gvk, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json data with gvk(%v): %v", gvk.String(), err)
	}
	return ro, nil
}
