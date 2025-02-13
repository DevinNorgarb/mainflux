package mocks

import (
	"context"
	"sync"

	"github.com/MainfluxLabs/mainflux/auth"
	"github.com/MainfluxLabs/mainflux/pkg/errors"
)

var _ auth.OrgRepository = (*orgRepositoryMock)(nil)

type orgRepositoryMock struct {
	mu      sync.Mutex
	orgs    map[string]auth.Org
	members map[string]auth.Member
	groups  map[string]auth.Group
}

// NewOrgRepository returns mock of org repository
func NewOrgRepository() auth.OrgRepository {
	return &orgRepositoryMock{
		orgs:    make(map[string]auth.Org),
		members: make(map[string]auth.Member),
		groups:  make(map[string]auth.Group),
	}
}

func (orm *orgRepositoryMock) Save(ctx context.Context, orgs ...auth.Org) error {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	for _, org := range orgs {
		if _, ok := orm.orgs[org.ID]; ok {
			return errors.ErrConflict
		}

		orm.orgs[org.ID] = org
	}

	return nil
}

func (orm *orgRepositoryMock) Update(ctx context.Context, org auth.Org) error {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	if _, ok := orm.orgs[org.ID]; !ok {
		return errors.ErrNotFound
	}

	orm.orgs[org.ID] = org

	return nil
}

func (orm *orgRepositoryMock) Delete(ctx context.Context, owner, id string) error {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	if _, ok := orm.orgs[id]; !ok && orm.orgs[id].OwnerID != owner {
		return errors.ErrNotFound
	}
	delete(orm.orgs, id)

	return nil
}

func (orm *orgRepositoryMock) RetrieveByID(ctx context.Context, id string) (auth.Org, error) {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	org, ok := orm.orgs[id]
	if !ok {
		return auth.Org{}, errors.ErrNotFound
	}

	return org, nil
}

func (orm *orgRepositoryMock) RetrieveByOwner(ctx context.Context, ownerID string, pm auth.PageMetadata) (auth.OrgsPage, error) {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	i := uint64(0)
	orgs := make([]auth.Org, 0)
	for _, org := range orm.orgs {
		if i >= pm.Offset && i < pm.Offset+pm.Limit {
			if org.OwnerID == ownerID {
				orgs = append(orgs, org)
			}
		}
		i++
	}

	return auth.OrgsPage{
		Orgs: orgs,
		PageMetadata: auth.PageMetadata{
			Total:  uint64(len(orm.orgs)),
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}, nil
}

func (orm *orgRepositoryMock) RetrieveMemberships(ctx context.Context, memberID string, pm auth.PageMetadata) (auth.OrgsPage, error) {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	i := uint64(0)
	orgs := make([]auth.Org, 0)
	for _, org := range orm.orgs {
		if i >= pm.Offset && i < pm.Offset+pm.Limit {
			if _, ok := orm.members[memberID]; ok {
				orgs = append(orgs, org)
			}
		}
		i++
	}

	return auth.OrgsPage{
		Orgs: orgs,
		PageMetadata: auth.PageMetadata{
			Total:  uint64(len(orm.orgs)),
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}, nil
}

func (orm *orgRepositoryMock) AssignMembers(ctx context.Context, mrs ...auth.MemberRelation) error {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	for _, mr := range mrs {
		if _, ok := orm.orgs[mr.OrgID]; !ok {
			return errors.ErrNotFound
		}
		orm.members[mr.MemberID] = auth.Member{
			ID:   mr.MemberID,
			Role: mr.Role,
		}
	}

	return nil
}

func (orm *orgRepositoryMock) UnassignMembers(ctx context.Context, orgID string, memberIDs ...string) error {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	for _, memberID := range memberIDs {
		if _, ok := orm.members[memberID]; !ok || orm.orgs[orgID].ID != orgID {
			return errors.ErrNotFound
		}
		delete(orm.members, memberID)
	}

	return nil
}

func (orm *orgRepositoryMock) UpdateMembers(ctx context.Context, mrs ...auth.MemberRelation) error {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	for _, mr := range mrs {
		if _, ok := orm.members[mr.MemberID]; !ok || orm.orgs[mr.OrgID].ID != mr.OrgID {
			return errors.ErrNotFound
		}
		orm.members[mr.MemberID] = auth.Member{
			ID:   mr.MemberID,
			Role: mr.Role,
		}
	}

	return nil
}

func (orm *orgRepositoryMock) RetrieveRole(ctx context.Context, memberID, orgID string) (string, error) {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	if _, ok := orm.members[memberID]; !ok {
		return "", errors.ErrNotFound
	}

	return orm.members[memberID].Role, nil
}

func (orm *orgRepositoryMock) RetrieveMembers(ctx context.Context, orgID string, pm auth.PageMetadata) (auth.OrgMembersPage, error) {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	i := uint64(0)
	members := []auth.Member{}
	for _, member := range orm.members {
		if i >= pm.Offset && i < pm.Offset+pm.Limit {
			if _, ok := orm.orgs[orgID]; ok {
				members = append(members, member)
			}
		}
		i++
	}

	return auth.OrgMembersPage{
		Members: members,
		PageMetadata: auth.PageMetadata{
			Total:  uint64(len(orm.members)),
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}, nil
}

func (orm *orgRepositoryMock) AssignGroups(ctx context.Context, grs ...auth.GroupRelation) error {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	for _, gr := range grs {
		if _, ok := orm.orgs[gr.OrgID]; !ok {
			return errors.ErrNotFound
		}
		orm.groups[gr.GroupID] = auth.Group{
			ID: gr.GroupID,
		}
	}

	return nil
}

func (orm *orgRepositoryMock) UnassignGroups(ctx context.Context, orgID string, groupIDs ...string) error {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	for _, groupID := range groupIDs {
		if _, ok := orm.groups[groupID]; !ok || orm.orgs[orgID].ID != orgID {
			return errors.ErrNotFound
		}
		delete(orm.groups, groupID)
	}

	return nil
}

func (orm *orgRepositoryMock) RetrieveGroups(ctx context.Context, orgID string, pm auth.PageMetadata) (auth.GroupRelationsPage, error) {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	i := uint64(0)
	grs := []auth.GroupRelation{}
	for _, group := range orm.groups {
		if i >= pm.Offset && i < pm.Offset+pm.Limit {
			if _, ok := orm.orgs[orgID]; ok {
				grs = append(grs, auth.GroupRelation{
					OrgID:   orgID,
					GroupID: group.ID,
				})
			}
		}
		i++
	}

	return auth.GroupRelationsPage{
		GroupRelations: grs,
		PageMetadata: auth.PageMetadata{
			Total:  uint64(len(orm.groups)),
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}, nil
}

func (orm *orgRepositoryMock) RetrieveByGroupID(ctx context.Context, groupID string) (auth.Org, error) {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	org, ok := orm.groups[groupID]
	if !ok {
		return auth.Org{}, errors.ErrNotFound
	}

	return orm.orgs[org.ID], nil
}

func (orm *orgRepositoryMock) RetrieveAll(ctx context.Context) ([]auth.Org, error) {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	var orgs []auth.Org
	for _, org := range orm.orgs {
		orgs = append(orgs, org)
	}

	return orgs, nil
}

func (orm *orgRepositoryMock) RetrieveByAdmin(ctx context.Context, pm auth.PageMetadata) (auth.OrgsPage, error) {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	i := uint64(0)
	orgs := []auth.Org{}
	for _, org := range orm.orgs {
		if i >= pm.Offset && i < pm.Offset+pm.Limit {
			orgs = append(orgs, org)
		}
		i++
	}

	return auth.OrgsPage{
		Orgs: orgs,
		PageMetadata: auth.PageMetadata{
			Total:  uint64(len(orm.orgs)),
			Offset: pm.Offset,
			Limit:  pm.Limit,
		},
	}, nil
}

func (orm *orgRepositoryMock) RetrieveAllMemberRelations(ctx context.Context) ([]auth.MemberRelation, error) {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	var mrs []auth.MemberRelation
	for _, org := range orm.orgs {
		for _, member := range orm.members {
			mrs = append(mrs, auth.MemberRelation{
				OrgID:    org.ID,
				MemberID: member.ID,
			})
		}
	}

	return mrs, nil
}

func (orm *orgRepositoryMock) RetrieveAllGroupRelations(ctx context.Context) ([]auth.GroupRelation, error) {
	orm.mu.Lock()
	defer orm.mu.Unlock()

	var grs []auth.GroupRelation
	for _, org := range orm.orgs {
		for _, group := range orm.groups {
			grs = append(grs, auth.GroupRelation{
				OrgID:   org.ID,
				GroupID: group.ID,
			})
		}
	}

	return grs, nil
}
