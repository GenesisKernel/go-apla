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

package api

import (
	"net/http"

	"github.com/AplaProject/go-apla/packages/consts"
	"github.com/AplaProject/go-apla/packages/converter"
	"github.com/AplaProject/go-apla/packages/script"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type contractField struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Optional bool   `json:"optional"`
}

type getContractResult struct {
	ID       uint32          `json:"id"`
	StateID  uint32          `json:"state"`
	TableID  string          `json:"tableid"`
	WalletID string          `json:"walletid"`
	TokenID  string          `json:"tokenid"`
	Address  string          `json:"address"`
	Fields   []contractField `json:"fields"`
	Name     string          `json:"name"`
}

func getContractInfoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	logger := getLogger(r)

	contract := getContract(r, params["name"])
	if contract == nil {
		logger.WithFields(log.Fields{"type": consts.ContractError, "contract_name": params["contract"]}).Error("contract name")
		errorResponse(w, errContract.Errorf(params["name"]))
		return
	}

	var result getContractResult
	info := getContractInfo(contract)
	fields := make([]contractField, 0)
	result = getContractResult{
		ID:       uint32(info.Owner.TableID + consts.ShiftContractID),
		TableID:  converter.Int64ToStr(info.Owner.TableID),
		Name:     info.Name,
		StateID:  info.Owner.StateID,
		WalletID: converter.Int64ToStr(info.Owner.WalletID),
		TokenID:  converter.Int64ToStr(info.Owner.TokenID),
		Address:  converter.AddressToString(info.Owner.WalletID),
	}

	if info.Tx != nil {
		for _, fitem := range *info.Tx {
			fields = append(fields, contractField{
				Name:     fitem.Name,
				Type:     script.OriginalToString(fitem.Original),
				Optional: fitem.ContainsTag(script.TagOptional),
			})
		}
	}
	result.Fields = fields

	jsonResponse(w, result)
}
