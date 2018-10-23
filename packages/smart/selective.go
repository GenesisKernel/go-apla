// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package smart

import (
	"fmt"
	"strings"

	"github.com/GenesisKernel/go-genesis/packages/consts"
	"github.com/GenesisKernel/go-genesis/packages/converter"
	"github.com/GenesisKernel/go-genesis/packages/model"
	"github.com/GenesisKernel/go-genesis/packages/model/querycost"

	log "github.com/sirupsen/logrus"
)

const (
	prefTimestamp      = "timestamp"
	prefTimestampSpace = "timestamp "
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

func getFieldIndex(fields []string, name string) int {
	for i, v := range fields {
		if strings.ToLower(v) == name {
			return i
		}
	}
	return -1
}

func (sc *SmartContract) selectiveLoggingAndUpd(fields []string, ivalues []interface{},
	table string, whereFields, whereValues []string, generalRollback bool, exists bool) (int64, string, error) {

	var (
		cost            int64
		err             error
		rollbackInfoStr string
	)

	logger := sc.GetLogger()
	if generalRollback && sc.BlockData == nil {
		logger.WithFields(log.Fields{"type": consts.EmptyObject}).Error("Block is undefined")
		return 0, ``, fmt.Errorf(`It is impossible to write to DB when Block is undefined`)
	}

	sqlBuilder := CreateQueryBuilder(sc, table, fields, whereFields, whereValues, ivalues)
	queryCoster := querycost.GetQueryCoster(querycost.FormulaQueryCosterType)

	if err := sqlBuilder.normalizeValues(fields, ivalues); err != nil {
		return 0, "", err
	}

	values, err := converter.InterfaceSliceToStr(ivalues)
	if err != nil {
		return 0, ``, err
	}

	fieldsExpr := sqlBuilder.getSQLSelectFieldsExpr(fields)
	whereExpr := sqlBuilder.getSQLWhereExpr(whereFields, whereValues)
	selectQuery := `SELECT ` + fieldsExpr + ` FROM "` + table + `" ` + whereExpr

	selectCost, err := queryCoster.QueryCost(sc.DbTransaction, selectQuery)
	if err != nil {
		logger.WithFields(log.Fields{"type": consts.DBError, "error": err, "query": selectQuery}).Error("getting query total cost")
		return 0, "", err
	}

	logData, err := model.GetOneRowTransaction(sc.DbTransaction, selectQuery).String()
	if err != nil {
		logger.WithFields(log.Fields{"type": consts.DBError, "error": err, "query": selectQuery}).Error("getting one row transaction")
		return 0, "", err
	}

	cost += selectCost
	if exists && len(logData) == 0 {
		logger.WithFields(log.Fields{"type": consts.NotFound, "err": errUpdNotExistRecord, "query": selectQuery}).Error("updating for not existing record")
		return 0, "", errUpdNotExistRecord
	}

	if whereFields != nil && len(logData) > 0 {
		// if isKeyTable {
		// 	keyEcosystem = logData["ecosystem"]
		// }
		var err error
		rollbackInfoStr, err = sqlBuilder.generateRollBackInfoString(logData)
		if err != nil {
			return 0, "", err
		}

		updateExpr, err := sqlBuilder.getSQLUpdateExpr(fields, values, logData)
		if err != nil {
			return 0, "", err
		}

		if !sc.VDE {
			updateQuery := `UPDATE "` + sqlBuilder.table + `" SET ` + updateExpr + " " + whereExpr
			updateCost, err := queryCoster.QueryCost(sc.DbTransaction, updateQuery)
			if err != nil {
				logger.WithFields(log.Fields{"type": consts.DBError, "error": err, "query": updateQuery}).Error("getting query total cost for update query")
				return 0, "", err
			}
			cost += updateCost
		}

		err = model.Update(sc.DbTransaction, sqlBuilder.table, updateExpr, whereExpr)
		if err != nil {
			logger.WithFields(log.Fields{"type": consts.DBError, "error": err, "sql": updateExpr}).Error("getting update query")
			return 0, "", err
		}
		sqlBuilder.tableID = logData[`id`]
	} else {

		insertQuery, err := sqlBuilder.getSQLInsertQuery(fields, values, whereFields, whereValues)
		if err != nil {

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
		if err := addRollback(sc, sqlBuilder.table, sqlBuilder.tableID, rollbackInfoStr); err != nil {
			return 0, sqlBuilder.tableID, err
		}
	}
	return cost, sqlBuilder.tableID, nil
}

func escapeSingleQuotes(val string) string {
	return strings.Replace(val, `'`, `''`, -1)
}

func (sc *SmartContract) insert(fields []string, ivalues []interface{},
	table string) (int64, string, error) {
	return sc.selectiveLoggingAndUpd(fields, ivalues, table, nil, nil, !sc.VDE && sc.Rollback, false)
}

func (sc *SmartContract) update(fields []string, values []interface{},
	table string, whereField string, whereValue interface{}) (int64, string, error) {
	return sc.selectiveLoggingAndUpd(fields, values, table, []string{whereField},
		[]string{fmt.Sprint(whereValue)}, !sc.VDE && sc.Rollback, true)
}
