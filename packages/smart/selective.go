// Copyright (C) 2017, 2018, 2019 EGAAS S.A.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or (at
// your option) any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.

package smart

import (
	"fmt"
	"strings"

	"github.com/AplaProject/go-apla/packages/consts"
	"github.com/AplaProject/go-apla/packages/model"
	"github.com/AplaProject/go-apla/packages/model/querycost"
	"github.com/AplaProject/go-apla/packages/types"

	qb "github.com/AplaProject/go-apla/packages/smart/queryBuilder"
	log "github.com/sirupsen/logrus"
)

func addRollback(sc *SmartContract, table, tableID, rollbackInfoStr string) error {
	rollbackTx := &model.RollbackTx{
		BlockID:   sc.BlockData.BlockID,
		TxHash:    sc.TxHash,
		NameTable: table,
		TableID:   tableID,
		Data:      rollbackInfoStr,
	}

	err := rollbackTx.Create(sc.DbTransaction)
	if err != nil {
		return logErrorDB(err, "creating rollback tx")
	}
	return nil
}

func (sc *SmartContract) selectiveLoggingAndUpd(fields []string, ivalues []interface{},
	table string, inWhere *types.Map, generalRollback bool, exists bool) (int64, string, error) {

	var (
		cost            int64
		rollbackInfoStr string
		logData         map[string]string
	)

	logger := sc.GetLogger()
	if generalRollback && sc.BlockData == nil {
		logger.WithFields(log.Fields{"type": consts.EmptyObject}).Error("Block is undefined")
		return 0, ``, fmt.Errorf(`It is impossible to write to DB when Block is undefined`)
	}
	for i, field := range fields {
		fields[i] = strings.ToLower(field)
	}
	sqlBuilder := &qb.SQLQueryBuilder{
		Entry:        logger,
		Table:        table,
		Fields:       fields,
		FieldValues:  ivalues,
		Where:        inWhere,
		KeyTableChkr: model.KeyTableChecker{},
	}

	queryCoster := querycost.GetQueryCoster(querycost.FormulaQueryCosterType)
	if exists {
		selectQuery, err := sqlBuilder.GetSelectExpr()
		if err != nil {
			logger.WithFields(log.Fields{"error": err}).Error("on getting sql select statement")
			return 0, "", err
		}

		selectCost, err := queryCoster.QueryCost(sc.DbTransaction, selectQuery)
		if err != nil {
			logger.WithFields(log.Fields{"type": consts.DBError, "error": err, "table": table, "query": selectQuery, "fields": fields, "values": ivalues, "where": inWhere}).Error("getting query total cost")
			return 0, "", err
		}

		logData, err = model.GetOneRowTransaction(sc.DbTransaction, selectQuery).String()
		if err != nil {
			logger.WithFields(log.Fields{"type": consts.DBError, "error": err, "query": selectQuery}).Error("getting one row transaction")
			return 0, "", err
		}
		cost += selectCost
		if len(logData) == 0 {
			logger.WithFields(log.Fields{"type": consts.NotFound, "err": errUpdNotExistRecord, "table": table, "fields": fields, "values": shortString(fmt.Sprintf("%+v", ivalues), 100), "where": inWhere, "query": shortString(selectQuery, 100)}).Error("updating for not existing record")
			return 0, "", errUpdNotExistRecord
		}
		if sqlBuilder.IsEmptyWhere() {
			logger.WithFields(log.Fields{"type": consts.NotFound,
				"error": errWhereUpdate}).Error("update without where")
			return 0, "", errWhereUpdate
		}
	}

	if !sqlBuilder.Where.IsEmpty() && len(logData) > 0 {
		var err error
		rollbackInfoStr, err = sqlBuilder.GenerateRollBackInfoString(logData)
		if err != nil {
			return 0, "", err
		}

		updateExpr, err := sqlBuilder.GetSQLUpdateExpr(logData)
		if err != nil {
			return 0, "", err
		}

		whereExpr, err := sqlBuilder.GetSQLWhereExpr()
		if err != nil {
			logger.WithFields(log.Fields{"error": err}).Error("on getting where expression for update")
			return 0, "", err
		}
		if !sc.OBS {
			updateQuery := `UPDATE "` + sqlBuilder.Table + `" SET ` + updateExpr + " " + whereExpr
			updateCost, err := queryCoster.QueryCost(sc.DbTransaction, updateQuery)
			if err != nil {
				logger.WithFields(log.Fields{"type": consts.DBError, "error": err, "query": updateQuery}).Error("getting query total cost for update query")
				return 0, "", err
			}
			cost += updateCost
		}

		err = model.Update(sc.DbTransaction, sqlBuilder.Table, updateExpr, whereExpr)
		if err != nil {
			logger.WithFields(log.Fields{"type": consts.DBError, "error": err, "sql": updateExpr}).Error("getting update query")
			return 0, "", err
		}
		sqlBuilder.SetTableID(logData[`id`])
	} else {

		insertQuery, err := sqlBuilder.GetSQLInsertQuery(model.NextIDGetter{Tx: sc.DbTransaction})
		if err != nil {
			logger.WithFields(log.Fields{"error": err}).Error("on build insert qwery")
			return 0, "", err
		}

		insertCost, err := queryCoster.QueryCost(sc.DbTransaction, insertQuery)
		if err != nil {
			logger.WithFields(log.Fields{"type": consts.DBError, "error": err, "query": insertQuery}).Error("getting total query cost for insert query")
			return 0, "", err
		}

		cost += insertCost
		err = model.GetDB(sc.DbTransaction).Exec(insertQuery).Error
		if err != nil {
			logger.WithFields(log.Fields{"type": consts.DBError, "error": err, "query": insertQuery}).Error("executing insert query")
			return 0, "", err
		}
	}

	if generalRollback {
		if off := strings.IndexByte(sqlBuilder.Table, '_'); off > 0 {
			name := sqlBuilder.Table[off+1:]
			if sqlBuilder.KeyTableChkr.IsKeyTable(name) {
				sqlBuilder.Table = fmt.Sprintf(`%s_%s`, sqlBuilder.GetEcosystem(), name)
			}
		}
		if err := addRollback(sc, sqlBuilder.Table, sqlBuilder.TableID(), rollbackInfoStr); err != nil {
			return 0, sqlBuilder.TableID(), err
		}
	}
	return cost, sqlBuilder.TableID(), nil
}

func (sc *SmartContract) insert(fields []string, ivalues []interface{},
	table string) (int64, string, error) {
	return sc.selectiveLoggingAndUpd(fields, ivalues, table, nil, !sc.OBS && sc.Rollback, false)
}

func (sc *SmartContract) updateWhere(fields []string, values []interface{},
	table string, where *types.Map) (int64, string, error) {
	return sc.selectiveLoggingAndUpd(fields, values, table, where, !sc.OBS && sc.Rollback, true)
}

func (sc *SmartContract) update(fields []string, values []interface{},
	table string, whereField string, whereValue interface{}) (int64, string, error) {
	return sc.updateWhere(fields, values, table, types.LoadMap(map[string]interface{}{
		whereField: fmt.Sprint(whereValue)}))
}

func shortString(raw string, length int) string {
	if len(raw) > length {
		return raw[:length]
	}

	return raw
}
