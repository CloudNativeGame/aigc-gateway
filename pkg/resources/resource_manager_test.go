package resources

import (
	"context"
	gameKruiseV1alpha1 "github.com/openkruise/kruise-game/apis/v1alpha1"
	"github.com/openkruise/kruise-game/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/utils/pointer"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

var (
	schemeTest = runtime.NewScheme()
)

func init() {
	utilruntime.Must(gameKruiseV1alpha1.AddToScheme(schemeTest))
}

func TestResourceManager_checkResourceMeta(t *testing.T) {
	tests := []struct {
		meta     *ResourceMeta
		metaNeed *metaNeed
		err      error
	}{
		{
			meta: &ResourceMeta{
				Name:      "gss-name",
				Namespace: "gss-ns",
				ID:        "6",
			},
			metaNeed: &metaNeed{
				Name:      true,
				Namespace: true,
				ID:        true,
			},
			err: nil,
		},
		{
			meta: &ResourceMeta{
				Name:      "gss-name",
				Namespace: "gss-ns",
				ID:        "f",
			},
			metaNeed: &metaNeed{
				Name:      true,
				Namespace: true,
				ID:        true,
			},
			err: NewResourceError(ParameterError, IdNotIntegerReason, "%s: %s", ParameterError, IdNotIntegerReason),
		},
		{
			meta: &ResourceMeta{
				Name:      "gss-name",
				Namespace: "gss-ns",
			},
			metaNeed: &metaNeed{
				Name:      true,
				Namespace: true,
			},
			err: nil,
		},
	}

	for i, test := range tests {
		err := checkResourceMeta(test.meta, test.metaNeed)
		if !reflect.DeepEqual(test.err, err) {
			t.Errorf("case %d: expect err is: %s, but actual err is: %s", i, test.err.Error(), err.Error())
		}
	}
}

func TestResourceManager_ListResources(t *testing.T) {
	tests := []struct {
		namespaces      []string
		resourcesLabels map[string]string
		gssList         []*gameKruiseV1alpha1.GameServerSet
		selectedNum     []int
	}{
		{
			gssList: []*gameKruiseV1alpha1.GameServerSet{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:       "xxx",
						Name:            "gss-0",
						ResourceVersion: "999",
						Labels:          map[string]string{"types": "cv"},
					},
					Spec: gameKruiseV1alpha1.GameServerSetSpec{
						Replicas: pointer.Int32(3),
					},
				},
			},
			selectedNum: []int{0},
		},
		{
			namespaces: []string{"xx"},
			gssList: []*gameKruiseV1alpha1.GameServerSet{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:       "xxx",
						Name:            "gss-0",
						ResourceVersion: "999",
						Labels:          map[string]string{"types": "cv"},
					},
					Spec: gameKruiseV1alpha1.GameServerSetSpec{
						Replicas: pointer.Int32(3),
					},
				},
			},
			selectedNum: nil,
		},
		{
			namespaces:      []string{"xxx"},
			resourcesLabels: map[string]string{"types": "cv"},
			gssList: []*gameKruiseV1alpha1.GameServerSet{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:       "xxx",
						Name:            "gss-0",
						ResourceVersion: "999",
						Labels:          map[string]string{"types": "cv"},
					},
					Spec: gameKruiseV1alpha1.GameServerSetSpec{
						Replicas: pointer.Int32(3),
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:       "xxx",
						Name:            "gss-1",
						ResourceVersion: "999",
						Labels:          map[string]string{"types": "nlp"},
					},
					Spec: gameKruiseV1alpha1.GameServerSetSpec{
						Replicas: pointer.Int32(3),
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:       "xxx",
						Name:            "gss-2",
						ResourceVersion: "999",
						Labels:          map[string]string{"types": "cv"},
					},
					Spec: gameKruiseV1alpha1.GameServerSetSpec{
						Replicas: pointer.Int32(3),
					},
				},
			},
			selectedNum: []int{0, 2},
		},
		{
			resourcesLabels: map[string]string{"types": "cv"},
			gssList: []*gameKruiseV1alpha1.GameServerSet{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:       "xxx",
						Name:            "gss-0",
						ResourceVersion: "999",
						Labels:          map[string]string{"types": "cv"},
					},
					Spec: gameKruiseV1alpha1.GameServerSetSpec{
						Replicas: pointer.Int32(3),
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:       "xx",
						Name:            "gss-1",
						ResourceVersion: "999",
						Labels:          map[string]string{"types": "nlp"},
					},
					Spec: gameKruiseV1alpha1.GameServerSetSpec{
						Replicas: pointer.Int32(3),
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:       "xx",
						Name:            "gss-2",
						ResourceVersion: "999",
						Labels:          map[string]string{"types": "cv"},
					},
					Spec: gameKruiseV1alpha1.GameServerSetSpec{
						Replicas: pointer.Int32(3),
					},
				},
			},
			selectedNum: []int{0, 2},
		},
	}

	for caseNum, test := range tests {
		var objs []client.Object
		for _, gss := range test.gssList {
			objs = append(objs, gss)
		}
		c := fake.NewClientBuilder().WithScheme(schemeTest).WithObjects(objs...).Build()
		rm := &ResourceManager{Client: c}
		actualResources, _ := rm.ListResources(test.namespaces, test.resourcesLabels)
		var actualNum []int
		for _, resource := range actualResources {
			num := util.GetIndexFromGsName(resource.GetName())
			actualNum = append(actualNum, num)
		}
		if !util.IsSliceEqual(actualNum, test.selectedNum) {
			t.Errorf("case %d: expect resources: %v, but actual resources: %v", caseNum, test.selectedNum, actualNum)
		}
	}
}

func TestResourceManager_GetResource(t *testing.T) {
	tests := []struct {
		meta    *ResourceMeta
		gs      *gameKruiseV1alpha1.GameServer
		isError bool
	}{
		{
			meta: &ResourceMeta{
				Namespace: "xxx",
				Name:      "xxx",
				ID:        "1",
			},
			gs: &gameKruiseV1alpha1.GameServer{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "xxx",
					Name:      "xxx-0",
					Labels: map[string]string{
						gameKruiseV1alpha1.GameServerOwnerGssKey: "xxx",
					},
				},
			},
			isError: true,
		},
		{
			meta: &ResourceMeta{
				Namespace: "xxx",
				Name:      "xxx",
			},
			gs: &gameKruiseV1alpha1.GameServer{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "xxx",
					Name:      "xxx-0",
					Labels: map[string]string{
						gameKruiseV1alpha1.GameServerOwnerGssKey: "xxx",
					},
				},
			},
			isError: true,
		},
		{
			meta: &ResourceMeta{
				Namespace: "xxx",
				Name:      "xxx",
				ID:        "0",
			},
			gs: &gameKruiseV1alpha1.GameServer{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "xxx",
					Name:      "xxx-0",
					Labels: map[string]string{
						gameKruiseV1alpha1.GameServerOwnerGssKey: "xxx",
					},
				},
			},
			isError: false,
		},
	}

	for caseNum, test := range tests {
		objs := []client.Object{test.gs}
		c := fake.NewClientBuilder().WithScheme(schemeTest).WithObjects(objs...).Build()
		rm := &ResourceManager{Client: c}
		_, err := rm.GetResource(test.meta)
		if !reflect.DeepEqual(test.isError, err != nil) {
			t.Errorf("case %d: errs not equal", caseNum)
		}
	}
}

func TestResourceManager_CreateResource(t *testing.T) {
	tests := []struct {
		meta      *ResourceMeta
		gssBefore *gameKruiseV1alpha1.GameServerSet
		newID     string
		gssAfter  *gameKruiseV1alpha1.GameServerSet
	}{
		{
			meta: &ResourceMeta{
				Namespace: "xxx",
				Name:      "case-0",
			},
			gssBefore: &gameKruiseV1alpha1.GameServerSet{
				ObjectMeta: metav1.ObjectMeta{
					Namespace:       "xxx",
					Name:            "case-0",
					ResourceVersion: "999",
					Annotations:     map[string]string{gameKruiseV1alpha1.GameServerSetReserveIdsKey: "1"},
				},
				Spec: gameKruiseV1alpha1.GameServerSetSpec{
					Replicas:             pointer.Int32(3),
					ReserveGameServerIds: []int{1},
				},
			},
			newID: "4",
			gssAfter: &gameKruiseV1alpha1.GameServerSet{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "game.kruise.io/v1alpha1",
					Kind:       "GameServerSet",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace:       "xxx",
					Name:            "case-0",
					ResourceVersion: "1000",
					Annotations:     map[string]string{gameKruiseV1alpha1.GameServerSetReserveIdsKey: "1"},
				},
				Spec: gameKruiseV1alpha1.GameServerSetSpec{
					Replicas:             pointer.Int32(4),
					ReserveGameServerIds: []int{1},
				},
			},
		},
		{
			meta: &ResourceMeta{
				Namespace: "xxx",
				Name:      "case-1",
			},
			gssBefore: &gameKruiseV1alpha1.GameServerSet{
				ObjectMeta: metav1.ObjectMeta{
					Namespace:       "xxx",
					Name:            "case-1",
					ResourceVersion: "999",
					Annotations:     map[string]string{gameKruiseV1alpha1.GameServerSetReserveIdsKey: "1,2,4"},
				},
				Spec: gameKruiseV1alpha1.GameServerSetSpec{
					Replicas:             pointer.Int32(7),
					ReserveGameServerIds: []int{1, 2, 4},
				},
			},
			newID: "10",
			gssAfter: &gameKruiseV1alpha1.GameServerSet{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "game.kruise.io/v1alpha1",
					Kind:       "GameServerSet",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace:       "xxx",
					Name:            "case-1",
					ResourceVersion: "1000",
					Annotations:     map[string]string{gameKruiseV1alpha1.GameServerSetReserveIdsKey: "1,2,4"},
				},
				Spec: gameKruiseV1alpha1.GameServerSetSpec{
					Replicas:             pointer.Int32(8),
					ReserveGameServerIds: []int{1, 2, 4},
				},
			},
		},
	}

	for caseNum, test := range tests {
		objs := []client.Object{test.gssBefore}
		c := fake.NewClientBuilder().WithScheme(schemeTest).WithObjects(objs...).Build()
		rm := &ResourceManager{Client: c}
		actualMeta, err := rm.CreateResource(test.meta)
		if err != nil {
			t.Errorf("case %d: %s", caseNum, err.Error())
		}
		if actualMeta.ID != test.newID {
			t.Errorf("case %d: expect return ID is %s, but get %s", caseNum, test.newID, actualMeta.ID)
		}
		actualGss := &gameKruiseV1alpha1.GameServerSet{}
		if err := rm.Get(context.Background(), types.NamespacedName{
			Name:      test.gssBefore.GetName(),
			Namespace: test.gssBefore.GetNamespace(),
		}, actualGss); err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(test.gssAfter, actualGss) {
			t.Errorf("case %d: expect gss %v but got %v", caseNum, test.gssAfter, actualGss)
		}
	}
}

func TestResourceManager_PauseResource(t *testing.T) {
	tests := []struct {
		meta      *ResourceMeta
		gssBefore *gameKruiseV1alpha1.GameServerSet
		gssAfter  *gameKruiseV1alpha1.GameServerSet
		isError   bool
	}{
		{
			meta: &ResourceMeta{
				Namespace: "xxx",
				Name:      "case-0",
				ID:        "1",
			},
			gssBefore: &gameKruiseV1alpha1.GameServerSet{
				ObjectMeta: metav1.ObjectMeta{
					Namespace:       "xxx",
					Name:            "case-0",
					ResourceVersion: "999",
				},
				Spec: gameKruiseV1alpha1.GameServerSetSpec{
					Replicas:             pointer.Int32(3),
					ReserveGameServerIds: []int{},
				},
			},
			gssAfter: &gameKruiseV1alpha1.GameServerSet{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "game.kruise.io/v1alpha1",
					Kind:       "GameServerSet",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace:       "xxx",
					Name:            "case-0",
					ResourceVersion: "1000",
				},
				Spec: gameKruiseV1alpha1.GameServerSetSpec{
					Replicas:             pointer.Int32(2),
					ReserveGameServerIds: []int{1},
				},
			},
		},
		{
			meta: &ResourceMeta{
				Namespace: "xxx",
				Name:      "case-1",
			},
			gssBefore: &gameKruiseV1alpha1.GameServerSet{
				ObjectMeta: metav1.ObjectMeta{
					Namespace:       "xxx",
					Name:            "case-1",
					ResourceVersion: "999",
				},
				Spec: gameKruiseV1alpha1.GameServerSetSpec{
					Replicas:             pointer.Int32(3),
					ReserveGameServerIds: []int{},
				},
			},
			isError: true,
		},
	}

	for caseNum, test := range tests {
		objs := []client.Object{test.gssBefore}
		c := fake.NewClientBuilder().WithScheme(schemeTest).WithObjects(objs...).Build()
		rm := &ResourceManager{Client: c}
		err := rm.PauseResource(test.meta)

		if !reflect.DeepEqual(test.isError, err != nil) {
			t.Errorf("case %d: errs not equal", caseNum)
			continue
		}

		if !test.isError {
			actualGss := &gameKruiseV1alpha1.GameServerSet{}
			if err := rm.Get(context.Background(), types.NamespacedName{
				Name:      test.gssBefore.GetName(),
				Namespace: test.gssBefore.GetNamespace(),
			}, actualGss); err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(test.gssAfter, actualGss) {
				t.Errorf("case %d: expect gss %v but got %v", caseNum, test.gssAfter, actualGss)
			}
		}
	}
}
