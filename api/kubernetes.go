package api

import (
	"context"

	"github.com/shaardie/mondane-operator/api/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var scheme = runtime.NewScheme()

type User struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	URLs     []string `json:"urls"`
	Email    string   `json:"email"`
}

type k8sClient interface {
	Create(ctx context.Context, user *User) error
	Read(ctx context.Context, name string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, user *User) error
}

type k8sClientImpl struct {
	client    client.Client
	namespace string
}

func initK8sClient(namespace string) (k8sClient, error) {
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	cli, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}
	return &k8sClientImpl{
		client:    cli,
		namespace: namespace,
	}, nil
}

func (k *k8sClientImpl) Create(ctx context.Context, user *User) error {
	return k.client.Create(ctx, k.userToK8sUser(user))
}

func (k *k8sClientImpl) Read(ctx context.Context, name string) (*User, error) {
	user := &v1alpha1.User{}
	err := k.client.Get(ctx, types.NamespacedName{Namespace: k.namespace, Name: name}, user)
	if err != nil {
		return &User{}, err
	}
	r := k.k8sUserToUser(user)
	return &r, nil
}

func (k *k8sClientImpl) Update(ctx context.Context, user *User) error {
	// k8sUser := &v1alpha1.User{}
	// err := k.client.Get(ctx, types.NamespacedName{Name: user.ID, Namespace: k.namespace}, k8sUser)
	// if err != nil {
	// 	return err
	// }
	// k8sUser.Spec = v1alpha1.UserSpec{
	// 	Username: user.Username,
	// 	Email:    user.Email,
	// 	URLs:     user.URLs,
	// }
	// fmt.Println(k8sUser)
	return k.client.Update(ctx, k.userToK8sUser(user))
}

func (k *k8sClientImpl) Delete(ctx context.Context, user *User) error {
	return k.client.Delete(ctx, k.userToK8sUser(user))
}

func (k *k8sClientImpl) userToK8sUser(u *User) *v1alpha1.User {
	return &v1alpha1.User{
		ObjectMeta: v1.ObjectMeta{
			Name:      u.ID,
			Namespace: k.namespace,
		},
		Spec: v1alpha1.UserSpec{
			Username: u.Username,
			Email:    u.Email,
			URLs:     u.URLs,
		},
	}
}

func (k *k8sClientImpl) k8sUserToUser(u *v1alpha1.User) User {
	return User{
		ID:       u.Name,
		Username: u.Spec.Username,
		URLs:     u.Spec.URLs,
		Email:    u.Spec.Username,
	}
}
