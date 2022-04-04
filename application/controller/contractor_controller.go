package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"service_admin_contractor/application/cerrors"
	"service_admin_contractor/application/cvalidator"
	"service_admin_contractor/application/dto"
	"service_admin_contractor/application/respond"
	"service_admin_contractor/application/service"
	"strconv"
)

type ContractorController struct {
	s service.ContractorService
}

func NewContractorController(s service.ContractorService) *ContractorController {
	return &ContractorController{s}
}

func (c *ContractorController) HandleRoutes(r *mux.Router) {
	r.HandleFunc("/contractors", c.GetAllContractors).Methods(http.MethodOptions, http.MethodGet)
	r.HandleFunc("/contractors", c.CreateContractor).Methods(http.MethodOptions, http.MethodPost)
	r.HandleFunc("/contractors/{id}", c.GetContractor).Methods(http.MethodOptions, http.MethodGet)
	r.HandleFunc("/contractors/{id}", c.UpdateContractor).Methods(http.MethodOptions, http.MethodPut)
	r.HandleFunc("/contractors/{id}", c.DeleteContractor).Methods(http.MethodOptions, http.MethodDelete)

	r.HandleFunc("/contractors/{id}/employee", c.CreateContractorEmployee).Methods(http.MethodOptions, http.MethodPost)
	r.HandleFunc("/contractors/{id}/employee/{employeeId}", c.UpdateContractorEmployee).Methods(http.MethodOptions, http.MethodPut)
	r.HandleFunc("/contractors/{id}/employee/{employeeId}", c.DeleteContractorEmployee).Methods(http.MethodOptions, http.MethodDelete)

}

func (c *ContractorController) GetAllContractors(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	searchParameters, err := dto.ParseContractorSearchParameters(r.Form)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	res, total, err := c.s.FindContractors(r.Context(), *searchParameters)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	respond.WithPagination(w, r, dto.ConvertContractors(res), total)
}

func (c *ContractorController) GetContractor(w http.ResponseWriter, r *http.Request) {
	rid := mux.Vars(r)["id"]
	err := cvalidator.Validate.Var(rid, "required,numeric")
	if err != nil {
		respond.WithError(w, r, cerrors.ErrBadRequestVar(err, "id"))
		return
	}

	id, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	data, err := c.s.GetContractor(r.Context(), id)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	respond.With(w, r, dto.ConvertContractor(data))
}

func (c *ContractorController) CreateContractor(w http.ResponseWriter, r *http.Request) {
	requestDto := &dto.ContractorDto{}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&requestDto)
	if err != nil {
		respond.WithError(w, r, cerrors.ErrCouldNotDecodeBody(err))
		return
	}

	err = cvalidator.ValidateStruct(requestDto)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	contractor := dto.ConvertContractorDtoToEntity(requestDto)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	ctx := r.Context()

	err = c.s.CreateContractor(ctx, contractor)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	respond.With(w, r, dto.ConvertContractor(*contractor))
}

func (c *ContractorController) UpdateContractor(w http.ResponseWriter, r *http.Request) {
	rid := mux.Vars(r)["id"]
	err := cvalidator.Validate.Var(rid, "required,numeric")
	if err != nil {
		respond.WithError(w, r, cerrors.ErrBadRequestVar(err, "id"))
		return
	}

	id, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	requestDto := &dto.ContractorDto{}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&requestDto)
	if err != nil {
		respond.WithError(w, r, cerrors.ErrCouldNotDecodeBody(err))
		return
	}

	contractor := dto.ConvertContractorDtoToEntity(requestDto)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	ctx := r.Context()
	err = c.s.UpdateContractor(ctx, id, contractor)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	respond.With(w, r, true)
}

func (c *ContractorController) DeleteContractor(w http.ResponseWriter, r *http.Request) {
	rid := mux.Vars(r)["id"]
	err := cvalidator.Validate.Var(rid, "required,numeric")
	if err != nil {
		respond.WithError(w, r, cerrors.ErrBadRequestVar(err, "id"))
		return
	}

	id, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}
	err = c.s.DeleteContractor(id)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	respond.With(w, r, true)
}

func (c *ContractorController) CreateContractorEmployee(w http.ResponseWriter, r *http.Request) {
	rid := mux.Vars(r)["id"]
	err := cvalidator.Validate.Var(rid, "required,numeric")
	if err != nil {
		respond.WithError(w, r, cerrors.ErrBadRequestVar(err, "id"))
		return
	}

	contractorId, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	requestDto := &dto.EmployeeDto{}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&requestDto)
	if err != nil {
		respond.WithError(w, r, cerrors.ErrCouldNotDecodeBody(err))
		return
	}

	err = cvalidator.Validate.Struct(requestDto)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	employee := dto.ConvertEmployeeDtoToEntity(requestDto)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	ctx := r.Context()
	err = c.s.CreateContractorEmployee(ctx, contractorId, employee)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	respond.With(w, r, dto.ConvertContractorEmployee(*employee))
}

func (c *ContractorController) UpdateContractorEmployee(w http.ResponseWriter, r *http.Request) {
	rid := mux.Vars(r)["employeeId"]
	err := cvalidator.Validate.Var(rid, "required,numeric")
	if err != nil {
		respond.WithError(w, r, cerrors.ErrBadRequestVar(err, "employeeId"))
		return
	}

	employeeId, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	requestDto := &dto.EmployeeDto{}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&requestDto)
	if err != nil {
		respond.WithError(w, r, cerrors.ErrCouldNotDecodeBody(err))
		return
	}

	employee := dto.ConvertEmployeeDtoToEntity(requestDto)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	ctx := r.Context()
	err = c.s.UpdateContractorEmployee(ctx, employeeId, employee)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	respond.With(w, r, true)
}

func (c *ContractorController) DeleteContractorEmployee(w http.ResponseWriter, r *http.Request) {
	rid := mux.Vars(r)["employeeId"]
	err := cvalidator.Validate.Var(rid, "required,numeric")
	if err != nil {
		respond.WithError(w, r, cerrors.ErrBadRequestVar(err, "employeeId"))
		return
	}

	id, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}
	err = c.s.DeleteContractorEmployee(id)
	if err != nil {
		respond.WithError(w, r, err)
		return
	}

	respond.With(w, r, true)
}
