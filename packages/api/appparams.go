// Apla Software includes an integrated development
// environment with a multi-level system for the management
// of access rights to data, interfaces, and Smart contracts. The
// technical characteristics of the Apla Software are indicated in
// Apla Technical Paper.

// Apla Users are granted a permission to deal in the Apla
// Software without restrictions, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of Apla Software, and to permit persons
// to whom Apla Software is furnished to do so, subject to the
// following conditions:
// * the copyright notice of GenesisKernel and EGAAS S.A.
// and this permission notice shall be included in all copies or
// substantial portions of the software;
// * a result of the dealing in Apla Software cannot be
// implemented outside of the Apla Platform environment.

// THE APLA SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY
// OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED
// TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE, ERROR FREE AND NONINFRINGEMENT. IN
// NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR
// THE USE OR OTHER DEALINGS IN THE APLA SOFTWARE.

package api

import (
	"net/http"

	"github.com/AplaProject/go-apla/packages/consts"
	"github.com/AplaProject/go-apla/packages/converter"
	"github.com/AplaProject/go-apla/packages/model"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

type appParamsResult struct {
	App  string        `json:"app_id"`
	List []paramResult `json:"list"`
}

type appParamsForm struct {
	ecosystemForm
	paramsForm
}

func (f *appParamsForm) Validate(r *http.Request) error {
	return f.ecosystemForm.Validate(r)
}

func getAppParamsHandler(w http.ResponseWriter, r *http.Request) {
	form := &appParamsForm{}
	if err := parseForm(r, form); err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	logger := getLogger(r)

	ap := &model.AppParam{}
	ap.SetTablePrefix(form.EcosystemPrefix)

	list, err := ap.GetAllAppParameters(converter.StrToInt64(params["appID"]))
	if err != nil {
		logger.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("Getting all app parameters")
	}

	result := &appParamsResult{
		App:  params["appID"],
		List: make([]paramResult, 0),
	}

	acceptNames := form.AcceptNames()
	for _, item := range list {
		if !acceptNames[item.Name] {
			continue
		}
		result.List = append(result.List, paramResult{
			ID:         converter.Int64ToStr(item.ID),
			Name:       item.Name,
			Value:      item.Value,
			Conditions: item.Conditions,
		})
	}

	jsonResponse(w, result)
}
