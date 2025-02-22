// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handlers

import (
	"context"

	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/wire"
)

// MsgCollStats returns a set of statistics for a collection.
func (h *Handler) MsgCollStats(ctx context.Context, msg *wire.OpMsg) (*wire.OpMsg, error) {
	document, err := msg.Document()
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	m := document.Map()
	collection := m["collStats"].(string)
	db, ok := m["$db"].(string)
	if !ok {
		return nil, lazyerrors.New("no db")
	}

	stats, err := h.pgPool.TableStats(ctx, db, collection)
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	var reply wire.OpMsg
	err = reply.SetSections(wire.OpMsgSection{
		Documents: []*types.Document{types.MustNewDocument(
			"ns", db+"."+collection,
			"count", stats.Rows,
			"size", stats.SizeTotal,
			"storageSize", stats.SizeTable,
			"totalIndexSize", stats.SizeIndexes,
			"totalSize", stats.SizeTotal,
			"scaleFactor", int32(1),
			"ok", float64(1),
		)},
	})
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	return &reply, nil
}
