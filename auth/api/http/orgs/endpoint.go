package orgs

import (
	"context"

	"github.com/MainfluxLabs/mainflux/auth"
	"github.com/go-kit/kit/endpoint"
)

func createOrgEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createOrgReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		org := auth.Org{
			Name:        req.Name,
			Description: req.Description,
			Metadata:    req.Metadata,
		}

		org, err := svc.CreateOrg(ctx, req.token, org)
		if err != nil {
			return nil, err
		}

		return orgRes{created: true, id: org.ID}, nil
	}
}

func viewOrgEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(orgReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		org, err := svc.ViewOrg(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		res := viewOrgRes{
			ID:          org.ID,
			Name:        org.Name,
			Description: org.Description,
			Metadata:    org.Metadata,
			OwnerID:     org.OwnerID,
			CreatedAt:   org.CreatedAt,
			UpdatedAt:   org.UpdatedAt,
		}

		return res, nil
	}
}

func updateOrgEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateOrgReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		org := auth.Org{
			ID:          req.id,
			Name:        req.Name,
			Description: req.Description,
			Metadata:    req.Metadata,
		}

		_, err := svc.UpdateOrg(ctx, req.token, org)
		if err != nil {
			return nil, err
		}

		return orgRes{created: false}, nil
	}
}

func deleteOrgEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(orgReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.RemoveOrg(ctx, req.token, req.id); err != nil {
			return nil, err
		}

		return deleteRes{}, nil
	}
}

func listOrgsEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listOrgsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		pm := auth.PageMetadata{
			Metadata: req.metadata,
			Name:     req.name,
		}

		page, err := svc.ListOrgs(ctx, req.token, pm)
		if err != nil {
			return nil, err
		}

		return buildOrgsResponse(page), nil
	}
}

func listMemberships(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listOrgMembershipsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		pm := auth.PageMetadata{
			Name:     req.name,
			Offset:   req.offset,
			Limit:    req.limit,
			Metadata: req.metadata,
		}

		page, err := svc.ListOrgMemberships(ctx, req.token, req.id, pm)
		if err != nil {
			return nil, err
		}

		return buildOrgsResponse(page), nil
	}
}

func assignMembersEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(membersReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.AssignMembers(ctx, req.token, req.orgID, req.Members...); err != nil {
			return nil, err
		}

		return assignRes{}, nil
	}
}

func unassignMembersEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(unassignMembersReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.UnassignMembers(ctx, req.token, req.orgID, req.MemberIDs...); err != nil {
			return nil, err
		}

		return unassignRes{}, nil
	}
}

func updateMembersEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(membersReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.UpdateMembers(ctx, req.token, req.orgID, req.Members...); err != nil {
			return nil, err
		}

		return assignRes{}, nil
	}
}

func listMembersEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listOrgMembersReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		pm := auth.PageMetadata{
			Offset:   req.offset,
			Limit:    req.limit,
			Metadata: req.metadata,
		}
		page, err := svc.ListOrgMembers(ctx, req.token, req.id, pm)
		if err != nil {
			return nil, err
		}

		return buildMembersResponse(page), nil
	}
}

func assignOrgGroupsEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(groupsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.AssignGroups(ctx, req.token, req.orgID, req.GroupIDs...); err != nil {
			return nil, err
		}

		return assignRes{}, nil
	}
}

func unassignOrgGroupsEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(groupsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.UnassignGroups(ctx, req.token, req.orgID, req.GroupIDs...); err != nil {
			return nil, err
		}

		return unassignRes{}, nil
	}
}

func listGroupsEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listOrgGroupsReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		pm := auth.PageMetadata{
			Offset:   req.offset,
			Limit:    req.limit,
			Metadata: req.metadata,
		}
		page, err := svc.ListOrgGroups(ctx, req.token, req.id, pm)
		if err != nil {
			return nil, err
		}

		return buildGroupsResponse(page), nil
	}
}

func backupEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(backupReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		backup, err := svc.Backup(ctx, req.token)
		if err != nil {
			return nil, err
		}

		return buildBackupResponse(backup), nil
	}
}

func restoreEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(restoreReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		backup := auth.Backup{
			Orgs:            req.Orgs,
			MemberRelations: req.MemberRelations,
			GroupRelations:  req.GroupRelations,
		}

		err := svc.Restore(ctx, req.token, backup)
		if err != nil {
			return nil, err
		}

		return restoreRes{}, nil
	}
}

func buildOrgsResponse(gp auth.OrgsPage) orgsPageRes {
	res := orgsPageRes{
		pageRes: pageRes{
			Total: gp.Total,
		},
		Orgs: []viewOrgRes{},
	}

	for _, org := range gp.Orgs {
		view := viewOrgRes{
			ID:          org.ID,
			OwnerID:     org.OwnerID,
			Name:        org.Name,
			Description: org.Description,
			Metadata:    org.Metadata,
			CreatedAt:   org.CreatedAt,
			UpdatedAt:   org.UpdatedAt,
		}
		res.Orgs = append(res.Orgs, view)
	}

	return res
}

func buildMembersResponse(mp auth.MembersPage) memberPageRes {
	res := memberPageRes{
		pageRes: pageRes{
			Total:  mp.Total,
			Offset: mp.Offset,
			Limit:  mp.Limit,
		},
		Members: []viewMemberRes{},
	}

	for _, memb := range mp.Members {
		m := viewMemberRes{
			ID:    memb.ID,
			Email: memb.Email,
			Role:  memb.Role,
		}
		res.Members = append(res.Members, m)
	}

	return res
}

func buildGroupsResponse(mp auth.GroupsPage) groupsPageRes {
	res := groupsPageRes{
		pageRes: pageRes{
			Total:  mp.Total,
			Offset: mp.Offset,
			Limit:  mp.Limit,
		},
		Groups: []viewGroupRes{},
	}

	for _, group := range mp.Groups {
		g := viewGroupRes{
			ID:          group.ID,
			OwnerID:     group.OwnerID,
			Name:        group.Name,
			Description: group.Description,
		}
		res.Groups = append(res.Groups, g)
	}

	return res
}

func buildBackupResponse(b auth.Backup) backupRes {
	res := backupRes{
		Orgs:            []viewOrgRes{},
		MemberRelations: []viewMemberRelations{},
		GroupRelations:  []viewGroupRelations{},
	}

	for _, org := range b.Orgs {
		view := viewOrgRes{
			ID:          org.ID,
			OwnerID:     org.OwnerID,
			Name:        org.Name,
			Description: org.Description,
			Metadata:    org.Metadata,
			CreatedAt:   org.CreatedAt,
			UpdatedAt:   org.UpdatedAt,
		}
		res.Orgs = append(res.Orgs, view)
	}

	for _, mRel := range b.MemberRelations {
		view := viewMemberRelations{
			OrgID:     mRel.OrgID,
			MemberID:  mRel.MemberID,
			Role:      mRel.Role,
			CreatedAt: mRel.CreatedAt,
			UpdatedAt: mRel.UpdatedAt,
		}
		res.MemberRelations = append(res.MemberRelations, view)
	}

	for _, groupRel := range b.GroupRelations {
		view := viewGroupRelations{
			GroupID:   groupRel.GroupID,
			OrgID:     groupRel.OrgID,
			CreatedAt: groupRel.CreatedAt,
			UpdatedAt: groupRel.UpdatedAt,
		}
		res.GroupRelations = append(res.GroupRelations, view)
	}

	return res
}
