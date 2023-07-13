package resources

import (
	"context"
	gamekruiseiov1alpha1 "github.com/openkruise/kruise-game/apis/v1alpha1"
	"github.com/openkruise/kruise-game/pkg/util"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"strconv"
)

// ResourceManager conteroller
type ResourceManager struct {
	client.Client
}

func (rm *ResourceManager) ListResources(namespaces []string, resourcesLabels map[string]string) ([]Resource, error) {
	var resources []Resource

	if len(namespaces) == 0 {
		gssList := &gamekruiseiov1alpha1.GameServerSetList{}
		err := rm.List(context.Background(), gssList, &client.ListOptions{
			LabelSelector: labels.SelectorFromSet(resourcesLabels),
		})
		if err != nil {
			return nil, NewResourceError(ApiCallError, "", err.Error())
		}

		for _, gss := range gssList.Items {
			resources = append(resources, gss.DeepCopy())
		}
	}

	for _, ns := range namespaces {
		listOptions := client.ListOptions{
			Namespace:     ns,
			LabelSelector: labels.SelectorFromSet(resourcesLabels),
		}
		gssList := &gamekruiseiov1alpha1.GameServerSetList{}
		err := rm.List(context.Background(), gssList, &listOptions)
		if err != nil {
			return nil, NewResourceError(ApiCallError, "", err.Error())
		}
		for _, gss := range gssList.Items {
			resources = append(resources, gss.DeepCopy())
		}
	}

	return resources, nil
}

func (rm *ResourceManager) GetResource(meta *ResourceMeta) (Resource, error) {
	if err := checkResourceMeta(meta, &metaNeed{ID: true, Name: true, Namespace: true}); err != nil {
		return nil, err
	}

	gs := &gamekruiseiov1alpha1.GameServer{}
	err := rm.Get(context.Background(), types.NamespacedName{
		Name:      meta.Name + "-" + meta.ID,
		Namespace: meta.Namespace,
	}, gs)
	if err != nil {
		if errors.IsNotFound(err) {
			gss := &gamekruiseiov1alpha1.GameServerSet{}
			err := rm.Get(context.Background(), types.NamespacedName{
				Name:      meta.Name,
				Namespace: meta.Namespace,
			}, gss)
			if err != nil {
				return nil, NewResourceError(ApiCallError, "", err.Error())
			}
			idInt, _ := strconv.Atoi(meta.ID)
			if util.IsNumInList(idInt, gss.Spec.ReserveGameServerIds) {
				return nil, NewResourceError(NotFoundError, PauseReason, "")
			} else {
				return nil, NewResourceError(NotFoundError, NotExistReason, "")
			}
		}
		return nil, NewResourceError(ApiCallError, "", err.Error())
	}

	return gs, nil
}

func (rm *ResourceManager) GetResourceEndpoint(meta *ResourceMeta) (string, error) {
	gs := &gamekruiseiov1alpha1.GameServer{}
	err := rm.Get(context.Background(), types.NamespacedName{
		Name:      meta.Name + "-" + meta.ID,
		Namespace: meta.Namespace,
	}, gs)
	if err != nil {
		if errors.IsNotFound(err) {
			return "", nil
		}
		return "", NewResourceError(ApiCallError, "", err.Error())
	}
	if len(gs.Status.NetworkStatus.ExternalAddresses) == 0 {
		return "", nil
	}

	return gs.Status.NetworkStatus.ExternalAddresses[0].EndPoint, nil
}

func (rm *ResourceManager) CreateResource(meta *ResourceMeta) (*ResourceMeta, error) {
	if err := checkResourceMeta(meta, &metaNeed{Name: true, Namespace: true}); err != nil {
		return nil, err
	}

	// get GameServerSet
	gss := &gamekruiseiov1alpha1.GameServerSet{}
	err := rm.Get(context.Background(), types.NamespacedName{
		Name:      meta.Name,
		Namespace: meta.Namespace,
	}, gss)
	if err != nil {
		return nil, NewResourceError(ApiCallError, "", err.Error())
	}

	newId := len(gss.Spec.ReserveGameServerIds) + int(*gss.Spec.Replicas)

	// update GameServerSet
	gss.Spec.Replicas = pointer.Int32(*gss.Spec.Replicas + 1)
	err = rm.Update(context.Background(), gss)
	if err != nil {
		return nil, NewResourceError(ApiCallError, "", err.Error())
	}

	return &ResourceMeta{
		Namespace: meta.Namespace,
		Name:      meta.Name,
		ID:        strconv.Itoa(newId),
	}, nil
}

func (rm *ResourceManager) PauseResource(meta *ResourceMeta) error {
	if err := checkResourceMeta(meta, &metaNeed{ID: true, Name: true, Namespace: true}); err != nil {
		return err
	}

	// get GameServerSet
	gss := &gamekruiseiov1alpha1.GameServerSet{}
	err := rm.Get(context.Background(), types.NamespacedName{
		Name:      meta.Name,
		Namespace: meta.Namespace,
	}, gss)
	if err != nil {
		return NewResourceError(ApiCallError, "", err.Error())
	}

	idInt, _ := strconv.Atoi(meta.ID)
	// check if already paused
	if util.IsNumInList(idInt, gss.Spec.ReserveGameServerIds) {
		return NewResourceError(NotFoundError, PauseReason, "")
	}

	// update GameServerSet
	gss.Spec.Replicas = pointer.Int32(*gss.Spec.Replicas - 1)
	gss.Spec.ReserveGameServerIds = append(gss.Spec.ReserveGameServerIds, []int{idInt}...)
	err = rm.Update(context.Background(), gss)
	if err != nil {
		return err
	}

	return nil
}

func (rm *ResourceManager) RecoverResource(meta *ResourceMeta) (Resource, error) {
	if err := checkResourceMeta(meta, &metaNeed{ID: true, Name: true, Namespace: true}); err != nil {
		return nil, err
	}

	// get GameServer
	gs := &gamekruiseiov1alpha1.GameServer{}
	err := rm.Get(context.Background(), types.NamespacedName{
		Name:      meta.Name + "-" + meta.ID,
		Namespace: meta.Namespace,
	}, gs)
	if err != nil {
		if errors.IsNotFound(err) {
			// get GameServerSet
			gss := &gamekruiseiov1alpha1.GameServerSet{}
			err = rm.Get(context.Background(), types.NamespacedName{
				Name:      meta.Name,
				Namespace: meta.Namespace,
			}, gss)
			if err != nil {
				return nil, NewResourceError(ApiCallError, "", err.Error())
			}

			idInt, _ := strconv.Atoi(meta.ID)
			if util.IsNumInList(idInt, gss.Spec.ReserveGameServerIds) {
				// update GameServerSet
				gss.Spec.Replicas = pointer.Int32(*gss.Spec.Replicas + 1)
				idInt, _ := strconv.Atoi(meta.ID)
				gss.Spec.ReserveGameServerIds = util.GetSliceInANotInB(gss.Spec.ReserveGameServerIds, []int{idInt})
				err = rm.Update(context.Background(), gss)
				if err != nil {
					return nil, NewResourceError(ApiCallError, "", err.Error())
				}
				return nil, NewResourceError(NotFoundError, PauseReason, "")
			}

			return nil, NewResourceError(NotFoundError, NotExistReason, "")
		}

		return nil, NewResourceError(ApiCallError, "", err.Error())
	}

	return gs, nil
}

func (rm *ResourceManager) DeleteResource(meta *ResourceMeta) error {
	if err := checkResourceMeta(meta, &metaNeed{ID: true, Name: true, Namespace: true}); err != nil {
		return err
	}

	// get GameServerSet
	gss := &gamekruiseiov1alpha1.GameServerSet{}
	err := rm.Get(context.Background(), types.NamespacedName{
		Name:      meta.Name,
		Namespace: meta.Namespace,
	}, gss)
	if err != nil {
		return NewResourceError(ApiCallError, "", err.Error())
	}

	idInt, _ := strconv.Atoi(meta.ID)
	// check if gs exist or not
	if !util.IsNumInList(idInt, gss.Spec.ReserveGameServerIds) {
		// update GameServerSet to delete gs
		gss.Spec.Replicas = pointer.Int32(*gss.Spec.Replicas - 1)
		gss.Spec.ReserveGameServerIds = append(gss.Spec.ReserveGameServerIds, []int{idInt}...)
		err = rm.Update(context.Background(), gss)
		if err != nil {
			return err
		}
	}

	// delete pvcs related to gss
	for _, vct := range gss.Spec.GameServerTemplate.VolumeClaimTemplates {
		pvc := &v1.PersistentVolumeClaim{}
		err = rm.Get(context.Background(), types.NamespacedName{
			Name:      vct.GetName() + "-" + meta.Name + "-" + meta.ID,
			Namespace: meta.Namespace,
		}, pvc)
		if err != nil {
			if errors.IsNotFound(err) {
				return nil
			}
			return NewResourceError(ApiCallError, "", err.Error())
		}
		err = rm.Delete(context.Background(), pvc)
		if err != nil && !errors.IsNotFound(err) {
			return NewResourceError(ApiCallError, "", err.Error())
		}
	}

	return nil
}

func (rm *ResourceManager) RestartResource(meta *ResourceMeta) error {
	if err := checkResourceMeta(meta, &metaNeed{ID: true, Name: true, Namespace: true}); err != nil {
		return err
	}

	// delete pod
	pod := &v1.Pod{}
	err := rm.Get(context.Background(), types.NamespacedName{
		Name:      meta.Name + "-" + meta.ID,
		Namespace: meta.Namespace,
	}, pod)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return NewResourceError(ApiCallError, "", err.Error())
	}
	err = rm.Delete(context.Background(), pod)
	if err != nil && !errors.IsNotFound(err) {
		return NewResourceError(ApiCallError, "", err.Error())
	}

	return nil
}

type metaNeed struct {
	ID        bool
	Name      bool
	Namespace bool
}

func checkResourceMeta(meta *ResourceMeta, metaNeed *metaNeed) error {
	if meta == nil {
		return NewResourceError(ParameterError, ResourceMetaNullReason, "%s: %s", ParameterError, ResourceMetaNullReason)
	}

	if metaNeed.Namespace && meta.Namespace == "" {
		return NewResourceError(ParameterError, NamespaceNullReason, "%s: %s", ParameterError, NamespaceNullReason)
	}

	if metaNeed.Name && meta.Name == "" {
		return NewResourceError(ParameterError, NameNullReason, "%s: %s", ParameterError, NameNullReason)
	}

	if metaNeed.ID {
		if meta.ID == "" {
			return NewResourceError(ParameterError, IdNullReason, "%s: %s", ParameterError, IdNullReason)
		}
		_, err := strconv.Atoi(meta.ID)
		if err != nil {
			return NewResourceError(ParameterError, IdNotIntegerReason, "%s: %s", ParameterError, IdNotIntegerReason)
		}
	}

	return nil
}

var rm *ResourceManager

func NewResourceManager() *ResourceManager {
	return rm
}

func init() {
	cfg := config.GetConfigOrDie()

	c, err := client.New(cfg, client.Options{})
	if err != nil {
		panic(err)
	}
	rm = &ResourceManager{
		Client: c,
	}
	gamekruiseiov1alpha1.AddToScheme(scheme.Scheme)
}
