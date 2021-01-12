package util

import (
	"context"
	"net/http"

	// "sync"
	claimsv1alpha1 "github.com/tmax-cloud/claim-operator/api/v1alpha1"
	clusterv1alpha1 "github.com/tmax-cloud/cluster-manager-operator/api/v1alpha1"
	client "github.com/tmax-cloud/hypercloud-server/client"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	// hyperv1 "github.com/tmax-cloud/hypercloud-server/external/hyper/v1"
	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
	"k8s.io/klog"
	"k8s.io/utils/pointer"
)

const (
	CLAIM_API_GROUP             = "claims.tmax.io"
	CLAIM_API_Kind              = "clusterclaims"
	CLAIM_API_GROUP_VERSION     = "claims.tmax.io/v1alpha1"
	CLUSTER_API_GROUP           = "cluster.tmax.io"
	CLUSTER_API_Kind            = "clustermanagers"
	CLUSTER_API_GROUP_VERSION   = "cluster.tmax.io/v1alpha1"
	HYPERCLOUD_SYSTEM_NAMESPACE = "hypercloud5-system"
)

var customClientset *client.Clientset
var k8sClientset *kubernetes.Clientset

func init() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset

	customClientset, err = client.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	k8sClientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func createSubjectAccessReview(userId string, group string, resource string, namespace string, name string, verb string) (*authv1.SubjectAccessReview, error) {
	sar := &authv1.SubjectAccessReview{
		Spec: authv1.SubjectAccessReviewSpec{
			ResourceAttributes: &authv1.ResourceAttributes{
				Group:     group,
				Resource:  resource,
				Namespace: namespace,
				Name:      name,
				Verb:      verb,
			},
			User: userId,
		},
	}

	sarResult, err := k8sClientset.AuthorizationV1().SubjectAccessReviews().Create(context.TODO(), sar, metav1.CreateOptions{})
	if err != nil {
		klog.Errorln(err)
		return nil, err
	}

	return sarResult, nil
}

func AdmitClusterClaim(userId string, clusterClaim *claimsv1alpha1.ClusterClaim, admit bool) (*claimsv1alpha1.ClusterClaim, string, int) {
	clusterClaimStatusUpdateRuleReview := authv1.SubjectAccessReview{
		Spec: authv1.SubjectAccessReviewSpec{
			ResourceAttributes: &authv1.ResourceAttributes{
				Resource: "clusterclaims/status",
				Verb:     "update",
				Group:    CLAIM_API_GROUP,
			},
			User: userId,
		},
	}
	sarResult, err := k8sClientset.AuthorizationV1().SubjectAccessReviews().Create(context.TODO(), &clusterClaimStatusUpdateRuleReview, metav1.CreateOptions{})
	if err != nil {
		klog.Errorln(err)
		return nil, err.Error(), http.StatusInternalServerError
	}
	if sarResult.Status.Allowed {
		klog.Infoln(" User [ " + userId + " ] has ClusterClaims/status Update Role, Can Update ClusterClaims")

		if admit == true {
			clusterClaim.Status.Phase = "Admitted"
			clusterClaim.Status.Reason = "Administrator approve the claim"
		} else {
			clusterClaim.Status.Phase = "Rejected"
			clusterClaim.Status.Reason = "Administrator reject the claim"
		}

		result, err := customClientset.ClaimsV1alpha1().ClusterClaims(HYPERCLOUD_SYSTEM_NAMESPACE).
			UpdateStatus(context.TODO(), clusterClaim, metav1.UpdateOptions{})
		if err != nil {
			klog.Errorln("Update ClusterClaim [ " + clusterClaim.Name + " ] Failed")
			return nil, err.Error(), http.StatusInternalServerError
		} else {
			msg := "Update ClusterClaim [ " + clusterClaim.Name + " ] Success"
			klog.Infoln(msg)
			return result, msg, http.StatusOK
		}
	} else {
		msg := " User [ " + userId + " ] has No ClusterClaims/status Update Role, Check If user has ClusterClaims/status Update Role"
		klog.Infoln(msg)
		return nil, msg, http.StatusForbidden
	}
}

func GetClusterClaim(userId string, clusterClaimName string) (*claimsv1alpha1.ClusterClaim, string, int) {

	var clusterClaim = &claimsv1alpha1.ClusterClaim{}

	clusterClaimGetRuleResult, err := createSubjectAccessReview(userId, CLAIM_API_GROUP, "clusterclaims", HYPERCLOUD_SYSTEM_NAMESPACE, clusterClaimName, "get")
	if err != nil {
		klog.Errorln(err)
		return nil, err.Error(), http.StatusInternalServerError
	}

	if clusterClaimGetRuleResult.Status.Allowed {
		clusterClaim, err = customClientset.ClaimsV1alpha1().ClusterClaims(HYPERCLOUD_SYSTEM_NAMESPACE).Get(context.TODO(), clusterClaimName, metav1.GetOptions{})
		if err != nil {
			klog.Errorln(err)
			return nil, err.Error(), http.StatusInternalServerError
		}
	} else {
		msg := "User [" + userId + "] authorization is denied for clusterclaims [" + clusterClaimName + "]"
		klog.Infoln(msg)
		return nil, msg, http.StatusForbidden
	}

	return clusterClaim, "Get claim success", http.StatusOK
}

func ListAccessibleClusterClaims(userId string) (*claimsv1alpha1.ClusterClaimList, string, int) {
	var clusterClaimList = &claimsv1alpha1.ClusterClaimList{}

	clusterClaimListRuleResult, err := createSubjectAccessReview(userId, CLAIM_API_GROUP, "clusterclaims", HYPERCLOUD_SYSTEM_NAMESPACE, "", "list")
	if err != nil {
		klog.Errorln(err)
		return nil, err.Error(), http.StatusInternalServerError
	}

	clusterClaimList, err = customClientset.ClaimsV1alpha1().ClusterClaims(HYPERCLOUD_SYSTEM_NAMESPACE).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Errorln(err)
		return nil, err.Error(), http.StatusInternalServerError
	}
	clusterClaimList.Kind = "ClusterClaimList"
	clusterClaimList.APIVersion = "claims.tmax.io/v1alpha1"

	if clusterClaimListRuleResult.Status.Allowed {
		msg := "User [ " + userId + " ] has clusterClaim List Role, Can Access All clusterClaim"
		klog.Infoln(msg)
		return clusterClaimList, msg, http.StatusOK
	} else {
		klog.Infoln(" User [ " + userId + " ] has No ClusterClaim List Role, Check If user has ClusterClaim Get Role & has Owner Annotation on certain ClusterClaim")
		// 0. claim 이름을 갖고 각 claim을 대상으로 sar을 날려서 user가 get할 수 있는지 확인..
		// 1. role 리스트 가져와서 userid-clusterName-clm-role 이름을 갖는 role이 있는지 확인
		// 2. claimlist에서 creator == userId인것을 찾아서 반환 (이렇게 하면 annotaion을 특정 sa제외 못 바꾸게 설정해야함!)

		// 특정 clusterclaim에 get 권한이 있는지 확인해야 되나 ..
		// clusterClaimGetRuleReview := authv1.SubjectAccessReview{
		// 	Spec: authv1.SubjectAccessReviewSpec{
		// 		ResourceAttributes: &authv1.ResourceAttributes{
		// 			Resource: "clusterclaims",
		// 			Verb:     "get",
		// 			Group:    CLAIM_API_GROUP,
		// 		},
		// 		User: userId,
		// 	},
		// }

		// clusterClaimGetRuleResult, err := k8sClientset.AuthorizationV1().SubjectAccessReviews().Create(context.TODO(), &clusterClaimGetRuleReview, metav1.CreateOptions{})
		// if err != nil {
		// 	klog.Errorln(err)
		// 	return nil, err.Error(), http.StatusInternalServerError
		// }
		// if clusterClaimGetRuleResult.Status.Allowed {
		klog.Infoln(" User [ " + userId + " ] has ClusterClaim Get Role")

		_clusterClaimList := []claimsv1alpha1.ClusterClaim{}
		// var wg sync.WaitGroup
		// wg.Add(len(clusterClaimList.Items))
		for _, clusterClaim := range clusterClaimList.Items {
			// go func(clusterClaim claimsv1alpha1.ClusterClaim, userId string, _clusterClaimList []claimsv1alpha1.ClusterClaim) {
			// defer wg.Done()
			if clusterClaim.Annotations["creator"] == userId {
				klog.Infoln(" User [ " + userId + " ] has owner annotation in ClusterClaim [ " + clusterClaim.Name + " ]")
				_clusterClaimList = append(_clusterClaimList, clusterClaim)
			}
			// }(clusterClaim, userId, _clusterClaimList)
		}
		// wg.Wait()

		clusterClaimList.Items = _clusterClaimList

		if len(clusterClaimList.Items) == 0 {
			msg := " User [ " + userId + " ] has No ClusterClaim"
			klog.Infoln(msg)
			return nil, msg, http.StatusForbidden
		}
		// } else {
		// 	msg := "User [ " + userId + " ] has no ClusterClaim Get Role, User Cannot Access any ClusterClaim"
		// 	klog.Infoln(msg)
		// 	return nil, msg, http.StatusForbidden
		// }
	}
	msg := "Success to get ClusterClaim for User [ " + userId + " ]"
	klog.Infoln(msg)
	return clusterClaimList, msg, http.StatusOK
}

func ListCluster(userId string) (*clusterv1alpha1.ClusterManagerList, string, int) {

	var clmList = &clusterv1alpha1.ClusterManagerList{}

	clmListRuleResult, err := createSubjectAccessReview(userId, CLUSTER_API_GROUP, "clusterclaims", HYPERCLOUD_SYSTEM_NAMESPACE, "", "list")
	if err != nil {
		klog.Errorln(err)
		return nil, err.Error(), http.StatusInternalServerError
	}

	clmList, err = customClientset.ClusterV1alpha1().ClusterManagers(HYPERCLOUD_SYSTEM_NAMESPACE).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Errorln(err)
		return nil, err.Error(), http.StatusInternalServerError
	}
	clmList.Kind = "ClusterManagerList"
	clmList.APIVersion = "cluster.tmax.io/v1alpha1"

	if clmListRuleResult.Status.Allowed {
		msg := "User [ " + userId + " ] has ClusterManager List Role, Can Access All ClusterManager"
		klog.Infoln(msg)
		return clmList, msg, http.StatusOK
	} else {
		msg := "User [ " + userId + " ] has No ClusterManager List Role"
		klog.Infoln(msg)
		return nil, msg, http.StatusForbidden
	}
}

func ListOwnerCluster(userId string) (*clusterv1alpha1.ClusterManagerList, string, int) {

	var clmList = &clusterv1alpha1.ClusterManagerList{}

	clmList, err := customClientset.ClusterV1alpha1().ClusterManagers(HYPERCLOUD_SYSTEM_NAMESPACE).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Errorln(err)
		return nil, err.Error(), http.StatusInternalServerError
	}
	clmList.Kind = "ClusterManagerList"
	clmList.APIVersion = "cluster.tmax.io/v1alpha1"

	_clmList := []clusterv1alpha1.ClusterManager{}
	for _, clm := range clmList.Items {
		if clm.Status.Owner == userId {
			_clmList = append(_clmList, clm)
		}
	}
	clmList.Items = _clmList

	if len(clmList.Items) == 0 {
		msg := " User [ " + userId + " ] has No own Cluster"
		klog.Infoln(msg)
		return nil, msg, http.StatusOK
	}
	msg := " User [ " + userId + " ] has own Cluster"
	klog.Infoln(msg)
	return clmList, msg, http.StatusOK
}

func ListMemberCluster(userId string) (*clusterv1alpha1.ClusterManagerList, string, int) {

	var clmList = &clusterv1alpha1.ClusterManagerList{}

	clmList, err := customClientset.ClusterV1alpha1().ClusterManagers(HYPERCLOUD_SYSTEM_NAMESPACE).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Errorln(err)
		return nil, err.Error(), http.StatusInternalServerError
	}
	clmList.Kind = "ClusterManagerList"
	clmList.APIVersion = "cluster.tmax.io/v1alpha1"

	_clmList := []clusterv1alpha1.ClusterManager{}
	for _, clm := range clmList.Items {
		if Contains(clm.Status.Members, userId) {
			_clmList = append(_clmList, clm)
		}
	}
	clmList.Items = _clmList

	if len(clmList.Items) == 0 {
		msg := " User [ " + userId + " ] has No belonging Cluster"
		klog.Infoln(msg)
		return nil, msg, http.StatusOK
	}

	msg := " User [ " + userId + " ] has belonging Clusters"
	klog.Infoln(msg)
	return clmList, msg, http.StatusOK
}

func GetCluster(userId string, clusterName string) (*clusterv1alpha1.ClusterManager, string, int) {

	var clm = &clusterv1alpha1.ClusterManager{}
	clusterGetRuleResult, err := createSubjectAccessReview(userId, CLUSTER_API_GROUP, "clustermanagers", HYPERCLOUD_SYSTEM_NAMESPACE, clusterName, "get")
	if err != nil {
		klog.Errorln(err)
		return nil, err.Error(), http.StatusInternalServerError
	}

	if clusterGetRuleResult.Status.Allowed {
		clm, err = customClientset.ClusterV1alpha1().ClusterManagers(HYPERCLOUD_SYSTEM_NAMESPACE).Get(context.TODO(), clusterName, metav1.GetOptions{})
		if err != nil {
			klog.Errorln(err)
			return nil, err.Error(), http.StatusInternalServerError
		}
	} else {
		msg := "User [" + userId + "] authorization is denied for cluster [" + clusterName + "]"
		klog.Infoln(msg)
		return nil, msg, http.StatusForbidden
	}
	return clm, "Get cluster success", http.StatusOK
}

func AddMembers(userId string, clm *clusterv1alpha1.ClusterManager, memberList []string) (*clusterv1alpha1.ClusterManager, string, int) {

	clmUpdateRuleResult, err := createSubjectAccessReview(userId, CLUSTER_API_GROUP, "clustermanagers", HYPERCLOUD_SYSTEM_NAMESPACE, clm.Name, "update")
	if err != nil {
		klog.Errorln(err)
		return nil, err.Error(), http.StatusInternalServerError
	}

	if clmUpdateRuleResult.Status.Allowed {
		for _, member := range memberList {
			clm.Status.Members = append(clm.Status.Members, member)
		}
		result, err := customClientset.ClusterV1alpha1().ClusterManagers(HYPERCLOUD_SYSTEM_NAMESPACE).UpdateStatus(context.TODO(), clm, metav1.UpdateOptions{})
		if err != nil {
			klog.Errorln("Update member list in cluster [ " + clm.Name + " ] Failed")
			return nil, err.Error(), http.StatusInternalServerError
		} else {
			msg := "Update member list in cluster [ " + clm.Name + " ] Success"
			klog.Infoln(msg)
			return result, msg, http.StatusOK
		}
	} else {
		msg := " User [ " + userId + " ] is not a cluster admin, Cannot invite members"
		klog.Infoln(msg)
		return nil, msg, http.StatusForbidden
	}
}

func DeleteMembers(userId string, clm *clusterv1alpha1.ClusterManager, memberList []string) (*clusterv1alpha1.ClusterManager, string, int) {

	hcrUpdateRuleResult, err := createSubjectAccessReview(userId, CLUSTER_API_GROUP, "clustermanagers", HYPERCLOUD_SYSTEM_NAMESPACE, clm.Name, "update")
	if err != nil {
		klog.Errorln(err)
		return nil, err.Error(), http.StatusInternalServerError
	}

	if hcrUpdateRuleResult.Status.Allowed {
		clm.Status.Members = Remove(clm.Status.Members, memberList)
		result, err := customClientset.ClusterV1alpha1().ClusterManagers(HYPERCLOUD_SYSTEM_NAMESPACE).UpdateStatus(context.TODO(), clm, metav1.UpdateOptions{})
		if err != nil {
			klog.Errorln("Update member list in cluster [ " + clm.Name + " ] Failed")
			return nil, err.Error(), http.StatusInternalServerError
		} else {
			msg := "Update member list in cluster [ " + clm.Name + " ] Success"
			klog.Infoln(msg)
			return result, msg, http.StatusOK
		}
	} else {
		msg := " User [ " + userId + " ] is not a cluster admin, Cannot invite members"
		klog.Infoln(msg)
		return nil, msg, http.StatusForbidden
	}
}

func CreateCLMRole(clusterManager *clusterv1alpha1.ClusterManager, members []string) (string, int) {

	for _, member := range members {
		roleName := member + "-" + clusterManager.Name + "-clm-role"
		role := &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      roleName,
				Namespace: HYPERCLOUD_SYSTEM_NAMESPACE,
				OwnerReferences: []metav1.OwnerReference{
					metav1.OwnerReference{
						APIVersion:         CLUSTER_API_GROUP_VERSION,
						Kind:               CLUSTER_API_Kind,
						Name:               clusterManager.GetName(),
						UID:                clusterManager.GetUID(),
						BlockOwnerDeletion: pointer.BoolPtr(true),
						Controller:         pointer.BoolPtr(true),
					},
				},
			},
			Rules: []rbacv1.PolicyRule{
				{APIGroups: []string{CLAIM_API_GROUP}, Resources: []string{"clustermanagers"},
					ResourceNames: []string{clusterManager.Name}, Verbs: []string{"get"}},
				{APIGroups: []string{CLAIM_API_GROUP}, Resources: []string{"clustermanagers/status"},
					ResourceNames: []string{clusterManager.Name}, Verbs: []string{"get"}},
			},
		}

		if _, err := k8sClientset.RbacV1().Roles(HYPERCLOUD_SYSTEM_NAMESPACE).Create(context.TODO(), role, metav1.CreateOptions{}); err != nil {
			klog.Errorln(err)
			return err.Error(), http.StatusInternalServerError
		}

		roleBindingName := member + "-" + clusterManager.Name + "-clm-rolebinding"
		roleBinding := &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      roleBindingName,
				Namespace: HYPERCLOUD_SYSTEM_NAMESPACE,
				OwnerReferences: []metav1.OwnerReference{
					metav1.OwnerReference{
						APIVersion:         CLAIM_API_GROUP_VERSION,
						Kind:               CLAIM_API_Kind,
						Name:               clusterManager.GetName(),
						UID:                clusterManager.GetUID(),
						BlockOwnerDeletion: pointer.BoolPtr(true),
						Controller:         pointer.BoolPtr(true),
					},
				},
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
				Name:     roleName,
			},
			Subjects: []rbacv1.Subject{
				{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     "User",
					Name:     member,
					// Namespace: HYPERCLOUD_SYSTEM_NAMESPACE,
				},
			},
		}

		if _, err := k8sClientset.RbacV1().RoleBindings(HYPERCLOUD_SYSTEM_NAMESPACE).Create(context.TODO(), roleBinding, metav1.CreateOptions{}); err != nil {
			klog.Errorln(err)
			return err.Error(), http.StatusInternalServerError
		}
		msg := "ClusterMnager role [" + roleName + "] and rolebinding [ " + roleBindingName + "]  is created"
		klog.Infoln(msg)
	}
	msg := "ClusterMnager roles and rolebindings are created for all new members"
	klog.Infoln(msg)
	return msg, http.StatusOK
}

func DeleteCLMRole(clusterManager *clusterv1alpha1.ClusterManager, members []string) (string, int) {
	for _, member := range members {
		roleName := member + "-" + clusterManager.Name + "-clm-role"
		roleBindingName := member + "-" + clusterManager.Name + "-clm-rolebinding"

		_, err := k8sClientset.RbacV1().Roles(HYPERCLOUD_SYSTEM_NAMESPACE).Get(context.TODO(), roleName, metav1.GetOptions{})
		if err == nil {
			if err := k8sClientset.RbacV1().Roles(HYPERCLOUD_SYSTEM_NAMESPACE).Delete(context.TODO(), roleName, metav1.DeleteOptions{}); err != nil {
				klog.Errorln(err)
				return err.Error(), http.StatusInternalServerError
			}
		} else if errors.IsNotFound(err) {
			klog.Infoln("Role [" + roleName + "] is already deleted. pass")
		} else {
			return err.Error(), http.StatusInternalServerError
		}

		_, err = k8sClientset.RbacV1().RoleBindings(HYPERCLOUD_SYSTEM_NAMESPACE).Get(context.TODO(), roleBindingName, metav1.GetOptions{})
		if err == nil {
			if err := k8sClientset.RbacV1().RoleBindings(HYPERCLOUD_SYSTEM_NAMESPACE).Delete(context.TODO(), roleBindingName, metav1.DeleteOptions{}); err != nil {
				klog.Errorln(err)
				return err.Error(), http.StatusInternalServerError
			}
		} else if errors.IsNotFound(err) {
			klog.Infoln("Rolebinding [" + roleBindingName + "] is already deleted. pass")
		} else {
			return err.Error(), http.StatusInternalServerError
		}
		msg := "ClusterMnager role [" + roleName + "] and rolebinding [ " + roleBindingName + "]  is deleted"
		klog.Infoln(msg)

	}
	msg := "ClusterMnager roles and rolebindings are deleted for all deleted members"
	klog.Infoln(msg)
	return msg, http.StatusOK
}
