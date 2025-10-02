package blob

import (
	"context"
	"fmt"
	db "vcx/agent/internal/infra/db/store/blob"
	"vcx/agent/internal/domains"
	"vcx/pkg/toolkit/mapkit"

	"vcx/pkg/logging"
)

var log = logging.GetLogger()
const Domain = "Blob"

type Blob struct {
	domains.Meta
	Blob         []byte
	FilePath     string
	IsCompressed bool
	IsBinary     bool
	RefCounter   int
}

func mapToStruct(data map[string]any) *Blob {
	return &Blob{
		Meta: domains.Meta{
			ID:           mapkit.GetString(data, db.COL_ID),
			CreationDate: mapkit.GetString(data, db.COL_CREATIONDATE),
			LMU:          mapkit.GetString(data, db.COL_LMU),
			LMD:          mapkit.GetString(data, db.COL_LMD),
			GUID:         mapkit.GetString(data, db.COL_GUID),
		},
		Blob:         mapkit.GetBytes(data, db.COL_BLOB),
		FilePath:     mapkit.GetString(data, db.COL_FILEPATH),
		IsCompressed: mapkit.GetBool(data, db.COL_ISCOMPRESSED),
		IsBinary:     mapkit.GetBool(data, db.COL_ISBINARY),
		RefCounter:   mapkit.GetInt(data, db.COL_REFCOUNTER),
	}
}

func New(ctx context.Context, blob []byte, filePath string, isCompressed, isBinary bool) (*Blob, error) {
	data := map[string]any{
		db.COL_BLOB:         blob,
		db.COL_FILEPATH:     filePath,
		db.COL_ISCOMPRESSED: isCompressed,
		db.COL_ISBINARY:     isBinary,
		db.COL_REFCOUNTER:   1,
	}
	result, err := db.Create(ctx, data)
	if err != nil {
		log.Error(fmt.Sprintf("%s Creation Failed: %v", Domain, err))
		return nil, err
	}

	return mapToStruct(result), nil
}

func (b *Blob) Update(ctx context.Context) error {
	data := map[string]any{
		db.COL_BLOB:         b.Blob,
		db.COL_FILEPATH:     b.FilePath,
		db.COL_ISCOMPRESSED: b.IsCompressed,
		db.COL_ISBINARY:     b.IsBinary,
		db.COL_REFCOUNTER:   b.RefCounter,
	}
	_, err := db.Update(ctx, b.ID, data)
	if err != nil {
		log.Error(fmt.Sprintf("%s Update Failed: %v", Domain, err))
	}

	return err
}

func GetByID(ctx context.Context, id string) (*Blob, error) {
	data, err := db.GetByID(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("%s Retrieval Failed: %v", Domain, err))
		return nil, err
	}

	return mapToStruct(data), nil
}