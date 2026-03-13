package tag

import (
	"context"
	"vcx/agent/internal/consts/tagtype"
	"vcx/agent/internal/domains"
	db "vcx/agent/internal/infra/db/store/tag"
	"vcx/agent/internal/session"
	"vcx/pkg/toolkit/mapkit"
)

const Domain = "Tag"

type Tag struct {
	domains.Meta
	AccountID   string
	Name        string
	Description string
	TagType     tagtype.TagType
	FileID      string
	BranchID    string
	ProjectID   string
	ChangeID    string
}

func mapToStruct(data map[string]any) *Tag {
	return &Tag{
		Meta: domains.Meta{
			ID:           mapkit.GetString(data, db.COL_ID),
			CreationDate: mapkit.GetString(data, db.COL_CREATIONDATE),
			LMU:          mapkit.GetString(data, db.COL_LMU),
			LMD:          mapkit.GetString(data, db.COL_LMD),
			GUID:         mapkit.GetString(data, db.COL_GUID),
		},
		AccountID:   mapkit.GetString(data, db.COL_ACCOUNTID),
		Name:        mapkit.GetString(data, db.COL_NAME),
		Description: mapkit.GetString(data, db.COL_DESCRIPTION),
		TagType:     tagtype.FromString(mapkit.GetString(data, db.COL_TAGTYPE)),
		FileID:      mapkit.GetString(data, db.COL_FILEID),
		BranchID:    mapkit.GetString(data, db.COL_BRANCHID),
		ProjectID:   mapkit.GetString(data, db.COL_PROJECTID),
		ChangeID:    mapkit.GetString(data, db.COL_CHANGEID),
	}
}


func create(ctx context.Context, tt tagtype.TagType, name, description, fileID string) (*Tag, error) {
	data := map[string]any{
		db.COL_ACCOUNTID: session.GetAccountID(ctx),
		db.COL_TAGTYPE:   tt.ToString(),
		db.COL_CHANGEID:  session.GetChangeID(ctx),
	}

	if name != "" {
		data[db.COL_NAME] = name
	}
	if description != "" {
		data[db.COL_DESCRIPTION] = description
	}
	if fileID != "" {
		data[db.COL_FILEID] = fileID
	}

	switch tt {
	case tagtype.SYSTEM_FILE, tagtype.USER_FILE, tagtype.SYSTEM_BRANCH, tagtype.USER_BRANCH:
		data[db.COL_BRANCHID] = session.GetBranchID(ctx)
		fallthrough
	case tagtype.SYSTEM_PROJECT, tagtype.USER_PROJECT:
		data[db.COL_PROJECTID] = session.GetProjectID(ctx)
	}

	result, err := db.Create(ctx, data)
	if err != nil {
		domains.LogError(Domain, "Creation", err)
		return nil, err
	}

	return mapToStruct(result), nil
}


func NewSystemProjectTag(ctx context.Context) (*Tag, error) {
	return create(ctx, tagtype.SYSTEM_PROJECT, "", "", "")
}


func NewSystemBranchTag(ctx context.Context) (*Tag, error) {
	return create(ctx, tagtype.SYSTEM_BRANCH, "", "", "")
}


func NewSystemFileTag(ctx context.Context, fileID string) (*Tag, error) {
	return create(ctx, tagtype.SYSTEM_FILE, "", "", fileID)
}


func NewUserTag(ctx context.Context, name, description string) (*Tag, error) {
	return create(ctx, tagtype.USER, name, description, "")
}


func NewUserProjectTag(ctx context.Context, name, description string) (*Tag, error) {
	return create(ctx, tagtype.USER_PROJECT, name, description, "")
}


func NewUserBranchTag(ctx context.Context, name, description string) (*Tag, error) {
	return create(ctx, tagtype.USER_BRANCH, name, description, "")
}


func NewUserFileTag(ctx context.Context, name, description, fileID string) (*Tag, error) {
	return create(ctx, tagtype.USER_FILE, name, description, fileID)
}


func (t *Tag) Update(ctx context.Context) error {
	data := map[string]any{
		db.COL_ACCOUNTID:   t.AccountID,
		db.COL_NAME:        t.Name,
		db.COL_DESCRIPTION: t.Description,
		db.COL_TAGTYPE:     t.TagType.ToString(),
		db.COL_FILEID:      t.FileID,
		db.COL_BRANCHID:    t.BranchID,
		db.COL_PROJECTID:   t.ProjectID,
		db.COL_CHANGEID:    t.ChangeID,
	}
	_, err := db.Update(ctx, t.ID, data)
	if err != nil {
		domains.LogError(Domain, "Update", err)
	}

	return err
}

func GetByID(ctx context.Context, id string) (*Tag, error) {
	data, err := db.GetByID(ctx, id)
	if err != nil {
		domains.LogError(Domain, "Retrieval", err)
		return nil, err
	}

	return mapToStruct(data), nil
}
