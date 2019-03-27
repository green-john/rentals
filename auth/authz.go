package auth

type Permission uint

const (
	Create Permission = iota
	Read
	Update
	Delete
)

type AuthzService struct {
	// This is a map of roles to permissions per resource
	//
	//		role1 -> {
	//			resource1 -> [Read, Write]
	//			resource2 -> [Read, Write]
	//		},
	//		role2 -> {
	//			resource1 -> [Create, Delete]
	//			resource3 -> [Read, Update]
	//		}
	//
	perm map[string]map[string][]Permission
}

func (a *AuthzService) AddPermission(role string, resource string, permissions ...Permission) {
	if _, ok := a.perm[role]; !ok {
		a.perm[role] = make(map[string][]Permission)
	}

	if _, ok := a.perm[role][resource]; !ok {
		a.perm[role][resource] = make([]Permission, 0)
	}

	for _, p := range permissions {
		if !containsPerm(a.perm[role][resource], p) {
			a.perm[role][resource] = append(a.perm[role][resource], p)
		}
	}
}

func (a *AuthzService) Allowed(role string, resource string, permission Permission) bool {
	// If role does not exist, return false
	if _, ok := a.perm[role]; !ok {
		return false
	}

	// Same for resource
	if _, ok := a.perm[role][resource]; !ok {
		return false
	}

	return containsPerm(a.perm[role][resource], permission)
}

func NewAuthzService() *AuthzService {
	p := make(map[string]map[string][]Permission)
	return &AuthzService{p}
}

func containsPerm(permissions []Permission, perm Permission) bool {
	for _, elt := range permissions {
		if elt == perm {
			return true
		}
	}

	return false
}
