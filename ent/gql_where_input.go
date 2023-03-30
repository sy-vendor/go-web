// Code generated by ent, DO NOT EDIT.

package ent

import (
	"errors"
	"fmt"
	"go-web/ent/predicate"
	"go-web/ent/user"
)

// UserWhereInput represents a where input for filtering User queries.
type UserWhereInput struct {
	Predicates []predicate.User  `json:"-"`
	Not        *UserWhereInput   `json:"not,omitempty"`
	Or         []*UserWhereInput `json:"or,omitempty"`
	And        []*UserWhereInput `json:"and,omitempty"`

	// "id" field predicates.
	ID      *uint64  `json:"id,omitempty"`
	IDNEQ   *uint64  `json:"idNEQ,omitempty"`
	IDIn    []uint64 `json:"idIn,omitempty"`
	IDNotIn []uint64 `json:"idNotIn,omitempty"`
	IDGT    *uint64  `json:"idGT,omitempty"`
	IDGTE   *uint64  `json:"idGTE,omitempty"`
	IDLT    *uint64  `json:"idLT,omitempty"`
	IDLTE   *uint64  `json:"idLTE,omitempty"`

	// "name" field predicates.
	Name             *string  `json:"name,omitempty"`
	NameNEQ          *string  `json:"nameNEQ,omitempty"`
	NameIn           []string `json:"nameIn,omitempty"`
	NameNotIn        []string `json:"nameNotIn,omitempty"`
	NameGT           *string  `json:"nameGT,omitempty"`
	NameGTE          *string  `json:"nameGTE,omitempty"`
	NameLT           *string  `json:"nameLT,omitempty"`
	NameLTE          *string  `json:"nameLTE,omitempty"`
	NameContains     *string  `json:"nameContains,omitempty"`
	NameHasPrefix    *string  `json:"nameHasPrefix,omitempty"`
	NameHasSuffix    *string  `json:"nameHasSuffix,omitempty"`
	NameEqualFold    *string  `json:"nameEqualFold,omitempty"`
	NameContainsFold *string  `json:"nameContainsFold,omitempty"`

	// "sex" field predicates.
	Sex    *bool `json:"sex,omitempty"`
	SexNEQ *bool `json:"sexNEQ,omitempty"`

	// "age" field predicates.
	Age      *int  `json:"age,omitempty"`
	AgeNEQ   *int  `json:"ageNEQ,omitempty"`
	AgeIn    []int `json:"ageIn,omitempty"`
	AgeNotIn []int `json:"ageNotIn,omitempty"`
	AgeGT    *int  `json:"ageGT,omitempty"`
	AgeGTE   *int  `json:"ageGTE,omitempty"`
	AgeLT    *int  `json:"ageLT,omitempty"`
	AgeLTE   *int  `json:"ageLTE,omitempty"`

	// "account" field predicates.
	Account             *string  `json:"account,omitempty"`
	AccountNEQ          *string  `json:"accountNEQ,omitempty"`
	AccountIn           []string `json:"accountIn,omitempty"`
	AccountNotIn        []string `json:"accountNotIn,omitempty"`
	AccountGT           *string  `json:"accountGT,omitempty"`
	AccountGTE          *string  `json:"accountGTE,omitempty"`
	AccountLT           *string  `json:"accountLT,omitempty"`
	AccountLTE          *string  `json:"accountLTE,omitempty"`
	AccountContains     *string  `json:"accountContains,omitempty"`
	AccountHasPrefix    *string  `json:"accountHasPrefix,omitempty"`
	AccountHasSuffix    *string  `json:"accountHasSuffix,omitempty"`
	AccountEqualFold    *string  `json:"accountEqualFold,omitempty"`
	AccountContainsFold *string  `json:"accountContainsFold,omitempty"`

	// "password" field predicates.
	Password             *string  `json:"password,omitempty"`
	PasswordNEQ          *string  `json:"passwordNEQ,omitempty"`
	PasswordIn           []string `json:"passwordIn,omitempty"`
	PasswordNotIn        []string `json:"passwordNotIn,omitempty"`
	PasswordGT           *string  `json:"passwordGT,omitempty"`
	PasswordGTE          *string  `json:"passwordGTE,omitempty"`
	PasswordLT           *string  `json:"passwordLT,omitempty"`
	PasswordLTE          *string  `json:"passwordLTE,omitempty"`
	PasswordContains     *string  `json:"passwordContains,omitempty"`
	PasswordHasPrefix    *string  `json:"passwordHasPrefix,omitempty"`
	PasswordHasSuffix    *string  `json:"passwordHasSuffix,omitempty"`
	PasswordEqualFold    *string  `json:"passwordEqualFold,omitempty"`
	PasswordContainsFold *string  `json:"passwordContainsFold,omitempty"`
}

// AddPredicates adds custom predicates to the where input to be used during the filtering phase.
func (i *UserWhereInput) AddPredicates(predicates ...predicate.User) {
	i.Predicates = append(i.Predicates, predicates...)
}

// Filter applies the UserWhereInput filter on the UserQuery builder.
func (i *UserWhereInput) Filter(q *UserQuery) (*UserQuery, error) {
	if i == nil {
		return q, nil
	}
	p, err := i.P()
	if err != nil {
		if err == ErrEmptyUserWhereInput {
			return q, nil
		}
		return nil, err
	}
	return q.Where(p), nil
}

// ErrEmptyUserWhereInput is returned in case the UserWhereInput is empty.
var ErrEmptyUserWhereInput = errors.New("ent: empty predicate UserWhereInput")

// P returns a predicate for filtering users.
// An error is returned if the input is empty or invalid.
func (i *UserWhereInput) P() (predicate.User, error) {
	var predicates []predicate.User
	if i.Not != nil {
		p, err := i.Not.P()
		if err != nil {
			return nil, fmt.Errorf("%w: field 'not'", err)
		}
		predicates = append(predicates, user.Not(p))
	}
	switch n := len(i.Or); {
	case n == 1:
		p, err := i.Or[0].P()
		if err != nil {
			return nil, fmt.Errorf("%w: field 'or'", err)
		}
		predicates = append(predicates, p)
	case n > 1:
		or := make([]predicate.User, 0, n)
		for _, w := range i.Or {
			p, err := w.P()
			if err != nil {
				return nil, fmt.Errorf("%w: field 'or'", err)
			}
			or = append(or, p)
		}
		predicates = append(predicates, user.Or(or...))
	}
	switch n := len(i.And); {
	case n == 1:
		p, err := i.And[0].P()
		if err != nil {
			return nil, fmt.Errorf("%w: field 'and'", err)
		}
		predicates = append(predicates, p)
	case n > 1:
		and := make([]predicate.User, 0, n)
		for _, w := range i.And {
			p, err := w.P()
			if err != nil {
				return nil, fmt.Errorf("%w: field 'and'", err)
			}
			and = append(and, p)
		}
		predicates = append(predicates, user.And(and...))
	}
	predicates = append(predicates, i.Predicates...)
	if i.ID != nil {
		predicates = append(predicates, user.IDEQ(*i.ID))
	}
	if i.IDNEQ != nil {
		predicates = append(predicates, user.IDNEQ(*i.IDNEQ))
	}
	if len(i.IDIn) > 0 {
		predicates = append(predicates, user.IDIn(i.IDIn...))
	}
	if len(i.IDNotIn) > 0 {
		predicates = append(predicates, user.IDNotIn(i.IDNotIn...))
	}
	if i.IDGT != nil {
		predicates = append(predicates, user.IDGT(*i.IDGT))
	}
	if i.IDGTE != nil {
		predicates = append(predicates, user.IDGTE(*i.IDGTE))
	}
	if i.IDLT != nil {
		predicates = append(predicates, user.IDLT(*i.IDLT))
	}
	if i.IDLTE != nil {
		predicates = append(predicates, user.IDLTE(*i.IDLTE))
	}
	if i.Name != nil {
		predicates = append(predicates, user.NameEQ(*i.Name))
	}
	if i.NameNEQ != nil {
		predicates = append(predicates, user.NameNEQ(*i.NameNEQ))
	}
	if len(i.NameIn) > 0 {
		predicates = append(predicates, user.NameIn(i.NameIn...))
	}
	if len(i.NameNotIn) > 0 {
		predicates = append(predicates, user.NameNotIn(i.NameNotIn...))
	}
	if i.NameGT != nil {
		predicates = append(predicates, user.NameGT(*i.NameGT))
	}
	if i.NameGTE != nil {
		predicates = append(predicates, user.NameGTE(*i.NameGTE))
	}
	if i.NameLT != nil {
		predicates = append(predicates, user.NameLT(*i.NameLT))
	}
	if i.NameLTE != nil {
		predicates = append(predicates, user.NameLTE(*i.NameLTE))
	}
	if i.NameContains != nil {
		predicates = append(predicates, user.NameContains(*i.NameContains))
	}
	if i.NameHasPrefix != nil {
		predicates = append(predicates, user.NameHasPrefix(*i.NameHasPrefix))
	}
	if i.NameHasSuffix != nil {
		predicates = append(predicates, user.NameHasSuffix(*i.NameHasSuffix))
	}
	if i.NameEqualFold != nil {
		predicates = append(predicates, user.NameEqualFold(*i.NameEqualFold))
	}
	if i.NameContainsFold != nil {
		predicates = append(predicates, user.NameContainsFold(*i.NameContainsFold))
	}
	if i.Sex != nil {
		predicates = append(predicates, user.SexEQ(*i.Sex))
	}
	if i.SexNEQ != nil {
		predicates = append(predicates, user.SexNEQ(*i.SexNEQ))
	}
	if i.Age != nil {
		predicates = append(predicates, user.AgeEQ(*i.Age))
	}
	if i.AgeNEQ != nil {
		predicates = append(predicates, user.AgeNEQ(*i.AgeNEQ))
	}
	if len(i.AgeIn) > 0 {
		predicates = append(predicates, user.AgeIn(i.AgeIn...))
	}
	if len(i.AgeNotIn) > 0 {
		predicates = append(predicates, user.AgeNotIn(i.AgeNotIn...))
	}
	if i.AgeGT != nil {
		predicates = append(predicates, user.AgeGT(*i.AgeGT))
	}
	if i.AgeGTE != nil {
		predicates = append(predicates, user.AgeGTE(*i.AgeGTE))
	}
	if i.AgeLT != nil {
		predicates = append(predicates, user.AgeLT(*i.AgeLT))
	}
	if i.AgeLTE != nil {
		predicates = append(predicates, user.AgeLTE(*i.AgeLTE))
	}
	if i.Account != nil {
		predicates = append(predicates, user.AccountEQ(*i.Account))
	}
	if i.AccountNEQ != nil {
		predicates = append(predicates, user.AccountNEQ(*i.AccountNEQ))
	}
	if len(i.AccountIn) > 0 {
		predicates = append(predicates, user.AccountIn(i.AccountIn...))
	}
	if len(i.AccountNotIn) > 0 {
		predicates = append(predicates, user.AccountNotIn(i.AccountNotIn...))
	}
	if i.AccountGT != nil {
		predicates = append(predicates, user.AccountGT(*i.AccountGT))
	}
	if i.AccountGTE != nil {
		predicates = append(predicates, user.AccountGTE(*i.AccountGTE))
	}
	if i.AccountLT != nil {
		predicates = append(predicates, user.AccountLT(*i.AccountLT))
	}
	if i.AccountLTE != nil {
		predicates = append(predicates, user.AccountLTE(*i.AccountLTE))
	}
	if i.AccountContains != nil {
		predicates = append(predicates, user.AccountContains(*i.AccountContains))
	}
	if i.AccountHasPrefix != nil {
		predicates = append(predicates, user.AccountHasPrefix(*i.AccountHasPrefix))
	}
	if i.AccountHasSuffix != nil {
		predicates = append(predicates, user.AccountHasSuffix(*i.AccountHasSuffix))
	}
	if i.AccountEqualFold != nil {
		predicates = append(predicates, user.AccountEqualFold(*i.AccountEqualFold))
	}
	if i.AccountContainsFold != nil {
		predicates = append(predicates, user.AccountContainsFold(*i.AccountContainsFold))
	}
	if i.Password != nil {
		predicates = append(predicates, user.PasswordEQ(*i.Password))
	}
	if i.PasswordNEQ != nil {
		predicates = append(predicates, user.PasswordNEQ(*i.PasswordNEQ))
	}
	if len(i.PasswordIn) > 0 {
		predicates = append(predicates, user.PasswordIn(i.PasswordIn...))
	}
	if len(i.PasswordNotIn) > 0 {
		predicates = append(predicates, user.PasswordNotIn(i.PasswordNotIn...))
	}
	if i.PasswordGT != nil {
		predicates = append(predicates, user.PasswordGT(*i.PasswordGT))
	}
	if i.PasswordGTE != nil {
		predicates = append(predicates, user.PasswordGTE(*i.PasswordGTE))
	}
	if i.PasswordLT != nil {
		predicates = append(predicates, user.PasswordLT(*i.PasswordLT))
	}
	if i.PasswordLTE != nil {
		predicates = append(predicates, user.PasswordLTE(*i.PasswordLTE))
	}
	if i.PasswordContains != nil {
		predicates = append(predicates, user.PasswordContains(*i.PasswordContains))
	}
	if i.PasswordHasPrefix != nil {
		predicates = append(predicates, user.PasswordHasPrefix(*i.PasswordHasPrefix))
	}
	if i.PasswordHasSuffix != nil {
		predicates = append(predicates, user.PasswordHasSuffix(*i.PasswordHasSuffix))
	}
	if i.PasswordEqualFold != nil {
		predicates = append(predicates, user.PasswordEqualFold(*i.PasswordEqualFold))
	}
	if i.PasswordContainsFold != nil {
		predicates = append(predicates, user.PasswordContainsFold(*i.PasswordContainsFold))
	}

	switch len(predicates) {
	case 0:
		return nil, ErrEmptyUserWhereInput
	case 1:
		return predicates[0], nil
	default:
		return user.And(predicates...), nil
	}
}