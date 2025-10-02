package tag

import (
	"context"
	"fmt"
	db "vcx/agent/internal/infra/db/store/tag"
	"vcx/agent/internal/domains"
	"vcx/pkg/toolkit/mapkit"

	"vcx/pkg/logging"
)

var log = logging.GetLogger()
const Domain = "Tag"

type Tag struct {
	domains.Meta
	AccountID   string
	Name        string
	Description string
	TagType     string
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
		TagType:     mapkit.GetString(data, db.COL_TAGTYPE),
		FileID:      mapkit.GetString(data, db.COL_FILEID),
		BranchID:    mapkit.GetString(data, db.COL_BRANCHID),
		ProjectID:   mapkit.GetString(data, db.COL_PROJECTID),
		ChangeID:    mapkit.GetString(data, db.COL_CHANGEID),
	}
}

func New(ctx context.Context, accountID, name, description, tagType string) (*Tag, error) {
	data := map[string]any{
		db.COL_ACCOUNTID:   accountID,
		db.COL_NAME:        name,
		db.COL_DESCRIPTION: description,
		db.COL_TAGTYPE:     tagType,
	}
	result, err := db.Create(ctx, data)
	if err != nil {
		log.Error(fmt.Sprintf("%s Creation Failed: %v", Domain, err))
		return nil, err
	}

	return mapToStruct(result), nil
}

func (t *Tag) Update(ctx context.Context) error {
	data := map[string]any{
		db.COL_ACCOUNTID:   t.AccountID,
		db.COL_NAME:        t.Name,
		db.COL_DESCRIPTION: t.Description,
		db.COL_TAGTYPE:     t.TagType,
		db.COL_FILEID:      t.FileID,
		db.COL_BRANCHID:    t.BranchID,
		db.COL_PROJECTID:   t.ProjectID,
		db.COL_CHANGEID:    t.ChangeID,
	}
	_, err := db.Update(ctx, t.ID, data)
	if err != nil {
		log.Error(fmt.Sprintf("%s Update Failed: %v", Domain, err))
	}

	return err
}

func GetByID(ctx context.Context, id string) (*Tag, error) {
	data, err := db.GetByID(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("%s Retrieval Failed: %v", Domain, err))
		return nil, err
	}

	return mapToStruct(data), nil
}